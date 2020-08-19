# Data-conductor 
This is an example application of Davic framework.
The application, data-conductor, will fetch two JSON objects from different endpoints, and "merge" them. 

You will be able to configure how the two JSON objects being merged. Of course, with the Davic framework, you will be able to change the configuration on-the-fly. 

## Build and Run Data-conductor
### Clean Build
#### (1) Ensure/Update Davic Version in glide.yaml 

#### (2) Cleanup and Install Dependencies 
``` 
rm -rf vendor
```

```
rm glide.lock go.mod go.sum
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
### Fetching Hero's Data 
#### (1) Run data-conductor. Please refer to the previous section 

#### (2) Try out the mock data service: "get-hero"
GET request to http://localhost:8080/get-hero?id=2 

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

#### (3) Try out the Davic configurable endpoint
POST request to http://localhost:8080/davic/go with the following payload 

```
{
  "hero-id": "2", 
  "power-id": "0"
}
```

You should get the following response. 

```
{
  "hero-id": "2", 
  "power-id": "0"
}
```

At this point, the response is just an echo of the request. 
It is because we haven't done any configuration yet. 

#### (4) Create the GET request URL for the hero data
POST request to http://localhost:8080/davic/set with the following payload 

```
[
  "-opt-",
  "-store-write-",
  "hero-url",
  [
      "-opt-",
      "-str-concat-",
      "http://127.0.0.1:8080/get-hero?id=",
      [
          "-opt-",
          "-store-read-",
          "hero-id"
      ]
  ]
]
``` 

Then, agin, 
POST request to http://localhost:8080/davic/go with the following payload 

```
{
  "hero-id": "2", 
  "power-id": "0"
}
```

We now get the following response: 

```
{
  "hero-id": "0",
  "hero-url": "http://127.0.0.1:8080/get-hero?id=2",
  "power-id": "1"
}
```

We see there is a new field "hero-url" fo making the GET request for the hero data. 
Davic maintains a "store" underlying. The ***store*** is nothing but a single-layer key-value map. 
You can see the store as the "memory" of the Davic system. 

Now, it is a good time to explain what are the endpoints of ***data-conductor*** do. 

***data-conductor*** maintains a list of Davic operations. 
The payload we send to endpoint ***/set*** is a single operation, and endpoint ***/set*** append the operation to the list. 

Endpoint ***/go*** does the following: 

1. Take the Http request and save the entire payload as the ***store*** (the one maintained by Davic) 
2. Execute the Davic operation list (the one appended by the ***/set*** endpoint). Based on the operations, the Davic ***store*** may be update 
3. Responds with using the entire ***store*** as the response payload 

In this step of the tutorial, the payload we sent to ***/set*** is an operation which does the following 

1. Take the value of "hero-id" in the ***store*** 
2. Concat string "http://127.0.0.1:8080/get-hero?id=" with the value extracted from the previous step 
3. Store the concat result back to ***store*** with key "hero-url" 

#### (5) Make the GET request and save the response into Davic store 
POST request to http://127.0.0.1:8080/davic/set with the following payload 

```
[
  "-opt-",
  "-store-write-",
  "hero-data",
  [
    "-opt-",
    "-http-call-",
    [
      "-opt-",
      "-obj-update-",
      {
        "method": "GET",
        "url": "",
        "headers": {},
        "body": null
      },
      ["url"],
      [
        "-opt-",
        "-store-read-",
        "hero-url"
      ]
    ]
  ]
]
```

The above is an operation which does the follows: 

**a. Get the value of "hero-url" from the store**
```
[
  "-opt-",
  "-store-read-",
  "hero-url"
]
```
Davic will evaluate it and get the result
```
"http://127.0.0.1:8080/get-hero?id=2"
```

**b. Update the "url" field of an request object template**
```
[
  "-opt-",
  "-obj-update-",
  {
    "method": "GET",
    "url": "",
    "headers": {},
    "body": null
  },
  ["url"],
  "http://127.0.0.1:8080/get-hero?id=2"
]  
```
Davic will evaluate it and get the result (a request object)
```
{
  "method": "GET",
  "url": "http://127.0.0.1:8080/get-hero?id=2",
  "headers": {},
  "body": null
}
```

**c. Once the request object is made, make a Http GET call**
```
[
  "-opt-",
  "-http-call-",
  {
    "method": "GET",
    "url": "http://127.0.0.1:8080/get-hero?id=2",
    "headers": {},
    "body": null
  }
]
```
Davic will evaluate it and get the result
```
{
  "body": {
    "age": "50",
    "email": "batman@dc.com",
    "gander": "m",
    "id": 2,
    "name": "Batman"
  },
  "headers": {
    "Content-Length": "72",
    "Content-Type": "application/json",
    "Date": "Thu, 06 Aug 2020 17:23:48 GMT"
  },
  "status": "200"
}
```

**d. Save the Http response to Davic store with key "hero-data"**
```
[
  "-opt-",
  "-store-write-",
  "hero-data",
  {
    "body": {
      "age": "50",
      "email": "batman@dc.com",
      "gander": "m",
      "id": 2,
      "name": "Batman"
    },
    "headers": {
      "Content-Length": "72",
      "Content-Type": "application/json",
      "Date": "Thu, 06 Aug 2020 17:23:48 GMT"
    },
    "status": "200"
  }
]
```

#### (5) Acquire the hero's data through the Davic configurable endpoint 
Now, again making a call to http://127.0.0.1:8080/davic/go with 
```
{
  "hero-id": "2", 
  "power-id": "0"
}
```

You should now get 
```
{
  "hero-data": {
    "body": {
      "age": "50",
      "email": "batman@dc.com",
      "gander": "m",
      "id": 2,
      "name": "Batman"
    },
    "headers": {
      "Content-Length": "72",
      "Content-Type": "application/json",
      "Date": "Thu, 06 Aug 2020 17:23:48 GMT"
    },
    "status": "200"
  },
  "hero-id": "2",
  "hero-url": "http://127.0.0.1:8080/get-hero?id=2",
  "power-id": "0"
}
```

### Fetching Power's Data 
Similar to what we did for fetching hero's data in the previous section, use the following two payloads to send two POST requests to http://127.0.0.1:8080/davic/set Please keep the order of the two payloads/requests. 

**a. Make GET request url** 
```
[
  "-opt-",
  "-store-write-",
  "power-url",
  [
    "-opt-",
    "-str-concat-",
    "http://127.0.0.1:8080/get-power?id=",
    [
      "-opt-",
      "-store-read-",
      "power-id"
    ]
  ]
]
```

**b. Fetch Power data and save it to Davic store** 
```
[
  "-opt-",
  "-store-write-",
  "power-data",
  [
    "-opt-",
    "-http-call-",
    [
      "-opt-",
      "-obj-update-",
      {
        "method": "GET",
        "url": "",
        "headers": {},
        "body": null
      },
      ["url"],
      [
        "-opt-",
        "-store-read-",
        "power-url"
      ]
    ]
  ]
]
```

Now, again making a call to http://127.0.0.1:8080/davic/go with
```
{
  "hero-id": "2", 
  "power-id": "0"
}
```
You should now get
```
{
  "hero-data": {
    "body": {
      "age": "50",
      "email": "batman@dc.com",
      "gander": "m",
      "id": 2,
      "name": "Batman"
    },
    "headers": {
      "Content-Length": "72",
      "Content-Type": "application/json",
      "Date": "Thu, 06 Aug 2020 17:52:50 GMT"
    },
    "status": "200"
  },
  "hero-id": "2",
  "hero-url": "http://127.0.0.1:8080/get-hero?id=2",
  "power-data": {
    "body": {
      "id": 0,
      "level": 1000000,
      "name": "rich"
    },
    "headers": {
      "Content-Length": "38",
      "Content-Type": "application/json",
      "Date": "Thu, 06 Aug 2020 17:52:50 GMT"
    },
    "status": "200"
  },
  "power-id": "0",
  "power-url": "http://127.0.0.1:8080/get-power?id=0"
}
```

### Cleaning Up the Response
We have demonstrated how Davic fetch data from web based on request. 
At this point, we just put everything we have in Davic store and return. 
Let's cleaning up the Davic store a bit and responding with a clean payload. 

Here is our expected response
```
{
  "name": "(hero's name)", 
  "age": "(hero's age)", 
  "power": "(hero power's name)", 
  "level": "(hero power's level)" 
}
```

This is what we are going to do: 
1. Make store entries, "name", "age", "power", and "level", by extracting data from the store. 
2. Remove entries other than the wanted four. 

Here we just take entry "name" as an example. 
We demonstrate how to extract the hero's name from the store and save it to a new entry. 
We then demonstrate removing entry "hero-id". 
Making/removing other entries will be just setting up the similar operations. 

POST request to http://127.0.0.1:8080/davic/set with the following payload for 
1. pulling hero's name from the store entry "hero-data" under ["body", "name"], and 
2. saving the name to store entry "name"
```
[
  "-opt-",
  "-store-write-",
  "name",
  [
    "-opt-",
    "-obj-read-",
    [
      "-opt-",
      "-store-read-",
      "hero-data"
    ],
    [
      "body",
      "name"
    ]
  ]
]
```

POST request to http://127.0.0.1:8080/davic/set with the following payload for removing store entry "hero-id"
```
[
  "-opt-",
  "-store-delete-",
  "hero-id"
]
```

With pulling data and removing redundant store entries, we can now build a web services dynamically fetches web data and conducts to a desired response. 