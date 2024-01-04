# Go-Baseline
## About the Service
TBA

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
