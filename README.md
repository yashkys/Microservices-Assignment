# Microservices Project
## How to run locally
*__Note:__*
- Add a .env inside goService package. Add the Alchemy apiKey(API_KEY)
- Add a .env inside scalaService\com.kys.bs2 package. Add the the following keys as per you relational database ->  DB_URL, DB_USER, DB_PASSWORD, DB_DRIVER 
### Service 2 ( Scala Grpc Server )
```bash
sbt compile
sbt run
```
### Service 1 ( Golang Grpc Client )
```bash
go mod tidy
protoc --go_out=. --go-grpc_out=. transactionReceipt.proto
go run .
```
