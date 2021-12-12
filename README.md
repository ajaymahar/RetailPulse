# RetailPulse
Interview assignment to implement an api

Lang: Golang 1.17.2

## Description
Simple Job service to process images collected from stores.

## Assumptions
 1. DB/Queue and some other services implementation is missing. Handled them programmetically.
 2. Used Store Master StoreID data in memory.
 3. This service is not persistent, means stopping process will wipes all the jobs in queue and all the running go routines.
 4. I was not clear on submit service request,  will it contains diff store_ids in single req. 

## Installing
Clone the git project

>```git clone https://github.com/ajaymahar/RetailPulse.git```

>```cd RetailPulse```

>```mkdir images```

>```go mod tidy```
>
If you want to quickly test the service

>```go run cmd/main.go```
>
If you want to build the service

>```go build -o jobservice cmd/main.go```
>
> ```./jobservice```
If you want to install it to your local machine as executable binary. (make sure you have added GOPATH to your PATH env var)

>```go install cmd/main.go```

> ```main <ENTER>```

TO Test the endpoints you can clone the postman collections and import it.
## Postman via workspace
https://www.postman.com/ajaymaharYT/workspace/retailpulse-pub/collection/4820437-f9cb3e1f-8f90-4c92-90fa-95d336084522

## Work env used
OS: MacOS
System: iMac
IDE: vim
Library: gorilla mux (https://github.com/gorilla/mux)


## Improvements

* error handling can be improved via implement error wrapping.
* can be implemented using context package for default timeout or task cancellations.
* codebase can be organized in more simple format.
* could have implement API key/JWT functionality for authontication.
* Returned error to the end user could have been more descriptive and user friendly.
* code refactor could have done 
* request validation can be done in more robot way 
* api access limit could have been implemented (how frequently user/actor can req the api) 
* logging is missing in this service, can be implement logger for monitoring.
* Unit test/Benchmarking are missing.
* and some other missing parts...
* Can make more configurable or access the required fields via ENV Var.
