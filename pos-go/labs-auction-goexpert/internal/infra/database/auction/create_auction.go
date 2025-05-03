package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	interval   time.Duration
	closeFunc  func(ctx context.Context, id string) error
}

// NewAuctionRepository reads AUCTION_DURATION and returns a repo with scheduling
func NewAuctionRepository(db *mongo.Database) *AuctionRepository {
	// read auction duration from env, default 5m
	dur := getAuctionDuration()

	repo := &AuctionRepository{
		Collection: db.Collection("auctions"),
		interval:   dur,
	}
	// default closeFunc
	repo.closeFunc = repo.updateAuctionStatusToClosed
	return repo
}

// getAuctionDuration reads AUCTION_DURATION env var (e.g. "1h", "30m", "10s").
func getAuctionDuration() time.Duration {
	s := os.Getenv("AUCTION_DURATION")
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	return 5 * time.Minute
}

// CreateAuction inserts the auction and schedules its automatic closing
func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction,
) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// schedule automatic closing
	ar.scheduleCloseAuction(ctx, auctionEntity.Id)
	return nil
}

// scheduleCloseAuction launches a goroutine that waits `interval` then closes the auction
func (ar *AuctionRepository) scheduleCloseAuction(ctx context.Context, auctionId string) {
	go func() {
		select {
		case <-time.After(ar.interval):
			if err := ar.closeFunc(ctx, auctionId); err != nil {
				logger.Error("Error closing auction automatically", err)
			}
		case <-ctx.Done():
		}
	}()
}

// updateAuctionStatusToClosed marks the auction as Completed in MongoDB
func (ar *AuctionRepository) updateAuctionStatusToClosed(ctx context.Context, auctionId string) error {
	filter := bson.M{"_id": auctionId}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	return err
}
