# OHLC with Volume & Value (Stockbit Challenge)
## About the Service
1. This is a project for Stockbit's Go Take Home Challenge. 
2. Running based on environment variables
3. It has several functions:
   1. HTTP Server running on :8080 to serve `POST /publish/transaction` to populate accepted `.ndjson` file
   2. Kafka for Stockbit's Transaction 
      1. Produce Transaction
      2. Consume Transaction -> Save Stock Summary to Redis
   3. GRPC Server running on :50051 to server `GetSummary` for getting summary info of a stock based on code
      1. GRPC client file to test provided under `grpc-client` directory
4. Assumption based on problem statement on PDF
   1. "A" type transaction only be calculated for Previous Price when the Quantity = 0
   2. Transaction won't be calculated (meaning it will be skipped) if it's not started with Previous Price
   3. Systemic error during process of TransactionRecorded will send the transaction to DLQ in Kafka for later ops

### Preparations
1. Make sure your redis already running on port 6379 and please adjust the config on `/config/local.env`
2. Make sure your kafka already running on port 9092 and please adjust the config on `/config/local.env`

## Setting Up and Run
1. Set up local environment variables and run locally (based on `/config/local.env`)
    ```sh
    $ make run
    ```
2. Run test
    ```sh
    $ make test
    ```
3. Run test with coverage
    ```sh
    $ make test-coverage
    ```
4. Run golangci-lint
    ```sh
    $ make lint
    ```

## End to End Run
1. Run the service locally with your environment variables
    ```sh
    $ make run
    ```
   
    or run on docker container
   ```sh
   $ docker compose up
   ```
2. Request to POST API `POST /publish/transaction with your file
   ```sh
   $ curl --request POST \
   --url http://localhost:8080/publish/transaction \
   --header 'Content-Type: multipart/form-data' \
   --form file=@/Users/dityap/Repository/stockbit-challenge/test.ndjson
   ```
3. Wait the process until all transactions consumed by the service
4. Check the stock summary by hitting the GRPC `GetSummary`. By running this on your local
   ```sh
   $ go run grpc-client/main.go (your stock code)
   ```
