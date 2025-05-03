# FullCycle Auction Go Project

This project is a Go-based auction system that supports creating auctions, placing bids, and automatically closing auctions after a specified duration using goroutines.

## Prerequisites

- Docker and Docker Compose installed
- Go 1.20 installed

## Running the Project

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd labs-auction-goexpert
   ```

2. Build and start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

   This will start the application on `http://localhost:8080` and MongoDB on `localhost:27017`.

3. Set environment variables in `cmd/auction/.env` to configure the auction duration and batch settings:
   ```env
   AUCTION_INTERVAL=20s
   BATCH_INSERT_INTERVAL=20s
   MAX_BATCH_SIZE=4
   ```

## Testing the Goroutine for Automatic Auction Closure

The project includes a goroutine that automatically closes auctions after a specified duration. To test this functionality:

1. Run the tests using the following command:
   ```bash
   go test ./internal/infra/database/auction -v
   ```

   This will execute the test cases, including those for the automatic auction closure routine.

2. Check the test output to ensure that auctions are being closed automatically after the specified interval.

## Key Features

- **Automatic Auction Closure**: Auctions are automatically closed after the duration specified in the `AUCTION_INTERVAL` environment variable.
- **Batch Bid Processing**: Bids are processed in batches to optimize database operations.
- **REST API**: Exposes endpoints for managing auctions, bids, and users.

## API Endpoints

- `POST /auction`: Create a new auction
- `GET /auction`: List all auctions
- `GET /auction/:auctionId`: Get details of a specific auction
- `POST /bid`: Place a bid on an auction
- `GET /bid/:auctionId`: List all bids for a specific auction
- `GET /user/:userId`: Get details of a user

## Additional Notes

- Ensure that the MongoDB service is running and accessible at the URL specified in the `.env` file.
- Modify the `AUCTION_INTERVAL` in `.env` to test different durations for automatic auction closure.