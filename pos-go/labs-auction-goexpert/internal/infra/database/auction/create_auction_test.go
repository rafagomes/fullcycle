package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAuctionClosureRoutine(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	ctx := context.Background()
	database := mt.Client.Database("testdb")
	repo := NewAuctionRepository(database)

	// Insert a mock auction with an expired timestamp
	expiredAuction := &auction_entity.Auction{
		Id:          "expired-auction",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-2 * time.Hour),
	}
	repo.CreateAuction(ctx, expiredAuction)

	// Start the closure routine
	go startAuctionClosureRoutine(ctx, repo)

	// Wait for the routine to process
	time.Sleep(3 * time.Second)

	// Verify the auction status is updated to closed
	updatedAuction, err := repo.FindAuctionById(ctx, "expired-auction")
	if err != nil {
		t.Fatalf("Error fetching auction: %v", err)
	}

	if updatedAuction.Status != auction_entity.Completed {
		t.Errorf("Expected auction status to be 'Completed', got '%v'", updatedAuction.Status)
	}
}
