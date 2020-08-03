# Data-conductor 
This is an example application of Davic framework.
The application, data-conductor, will fetch two JSON objects from different endpoints, and "merge" them. 

You will be able to configure how the two JSON objects being merged. Of course, with the Davic framework, you will be able to change the configuration on-the-fly. 

## Build and Run Data-conductor
### Clean Build
#### (1) Ensure Davic Version in glide.yaml 

#### (2) Cleanup and Install Dependencies 
``` 
rm -rf vendor
```

```
glide install 
```

#### (3) Build Go Module 
```
go mod init github.com/wfchiang/davic-sample-apps
``` 

### Run Data-conductor 
Once the Go module is built (please refer to the above "clean build" section), do the following command: 

``` 
go run main.go
```

## Step-by-step Demo
#### (1) Run data-conductor. Please refer to the previous section 

#### (2) Try out the two mock data services: "get-hero" and "get-power"
There are two data services: "get-hero" and "get-power". Try out these services with the following commands. 

**get-hero**: GET request to http://localhost:8080/get-hero?id=2 

You should get the following response. 

```
{
  "age": "50",
  "email": "batman@dc.com",
  "gander": "m",
  "id": 2,
  "name": "Batman"
}
```

**get-power**: GET request to http://localhost:8080/get-power?id=0

You should get the following response. 

```
{
  "id": 0,
  "level": 1000000,
  "name": "rich"
}
```

#### (3) Try out the Davic configurable endpoint
POST request to http://localhost:8080/davic/go 
with the following payload 

```
{
  "hero-url": "http://localhost:8080/get-hero?id=2", 
  "power-url": "http://localhost:8080/get-power?id=0"
}
```

You should get the following response. 

```
{
  "hero-url": "http://localhost:8080/get-hero?id=2", 
  "power-url": "http://localhost:8080/get-power?id=0"
}
```

At this point, the response is just an echo of the request. 
It is because we haven't done any configuration yet. 

