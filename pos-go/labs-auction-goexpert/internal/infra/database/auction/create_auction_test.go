package auction

import (
	"context"
	"testing"
	"time"
)

func TestScheduleCloseAuction(t *testing.T) {
	ctx := context.Background()

	// Use a short interval so test runs fast
	repo := &AuctionRepository{
		interval: 50 * time.Millisecond,
	}
	// Capture calls
	called := make(chan string, 1)
	repo.closeFunc = func(ctx context.Context, id string) error {
		called <- id
		return nil
	}

	auctionID := "test-auction-123"
	repo.scheduleCloseAuction(ctx, auctionID)

	select {
	case got := <-called:
		if got != auctionID {
			t.Errorf("closeFunc called with %q; want %q", got, auctionID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected closeFunc to be called, but it wasn't")
	}
}

func TestScheduleCloseAuction_CancelledContext(t *testing.T) {
	// If context is cancelled before interval, closeFunc must NOT be called
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := &AuctionRepository{
		interval: 50 * time.Millisecond,
	}
	called := false
	repo.closeFunc = func(ctx context.Context, id string) error {
		called = true
		return nil
	}

	repo.scheduleCloseAuction(ctx, "any")
	cancel() // cancel immediately

	time.Sleep(75 * time.Millisecond)
	if called {
		t.Error("closeFunc should not have been called after context cancellation")
	}
}
