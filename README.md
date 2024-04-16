# Go-Baseline
## About the Service
This is a go-baseline project that aims as a starting line for Go project. The concept of this baseline is trying to implement [Uncle Bob's Clean Code Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) combined with the experience from past projects. 
Project is using drivers/adapters from [https://github.com/dityuiri/go-adapter](https://github.com/dityuiri/go-adapter).

## Architecture
![Architecture](baseline-architecutre.png)

## Project Directory
```
go-baseline
| application
  <App starter and dependency injector>
  
| common
  <Shared functions and variables like constant, utility function, error code etc.>
  
| config
  <App configuration and environment variables>
  
| controller
  <Interface adapters a.k.a the handler like API endpoint controller, message consumer, command line runner etc.>
  
| mock
  <Mock for all the interfaces in the project. Unit-testing purpose>
  
| model
  <Entities layer that can consist of DAO and DTO>
  
| proxy
  <Proxy client to external services>
  
| repository
  <Repository layer to interact with data storage such as db, redis, or even kafka>
  
| service
  <Use cases layer. Business logic goes here>
  
| Dockerfile
| docker-compose.yml
| main.go
| Makefile
| entrypoint.sh
```
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
