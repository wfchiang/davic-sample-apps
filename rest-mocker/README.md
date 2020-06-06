# Davic Driven REST mocker 

This sample app, REST mocker, demonstrates the power of Davic framework. 
It is a generic mock service which can be used for testing your RESTful APIs. 
When an endpoint hit, it replies with a JSON data for the testing purpose. 
Since it is a mock service, you do not expect it completely replicates the full functionalities of the original RESTful service. 
But, for the testing aspect, you may not want the mock service be stupid which always replies with a fixed JSON data. 
With Davic, you can specify some simple behaviors for the mock service. 
In addition, the most important thing, you can change the mock service behaviors **at runtime**! 
This is really cool. Trust me. 

## Installation 
I assume that you already have **Davic** installed. If not, please install Davic first. 
In addition to Davic, you will also need **Gorilla/mux** framework. 
Executing the python script of this directory should help installing it: 
```
python myinstall.py
```

## Run
In command-line:  
```
go run main.go
```

## Demo 
You will need a REST API client software (such as Postman, Insomnia, etc.). 
Please use your favorite. 

### 1. Run the application.

Please refer to the **Run** section above. 
I assume that you are taking the default settings: running it locally, using port 8080, and adding no extra context-root. 

### 2. Test the app is really running. 

Send a GET request to http://127.0.0.1:8080/ 
You should get a "200 OK" as the response. 

### 3. Test the "Davic Configurable Endpoint". 

Send a POST request to http://127.0.0.1:8080/davic/go with the following payload. 
```json
{
    "id": "0123456789", 
    "name": "Jenny"
}
```
No additional HTTP headers are required. 

You should get a "200 OK" as the response. 
However, the response payload is "null" -- nothing. 
It is a sure thing since you have not configured the "Davic Configurable Endpoint". 

### 4. Test the "Hi2You" Endpoint. 

Before we configure the "Davic" endpoint, let's try another "Hi2You" endpoint. 

Send a POST request to http://127.0.0.1:8080/hi2you with the following payload (the same as the previous one).  
```json
{
    "id": "0123456789", 
    "name": "Jenny"
}
```
No additional HTTP headers are required. 

You should get a "200 OK" with the following payload: 
```json
{
    "message": "Hi!",
    "name": "Jenny"
}
```

If you change the "name" field in the request payload, says sending the following payload, 
```json
{
    "id": "13579", 
    "name": "foo-bar"
}
```
you should get a different response as follows: 
```json
{
    "message": "Hi!",
    "name": "foo-bar"
}
```

"Hi2You" is a pre-configured endpoint. Let's see how it is configured. 
Send a GET request to http://127.0.0.1:8080/getopt/hi2you

You should get a "200 OK" with the following JSON payload: 
```json
[
  [
    "-opt-",
    "-store-write-",
    "http-resp",
    {
      "message": "Hi!",
      "name": null
    }
  ],
  [
    "-opt-",
    "-store-write-",
    "http-resp",
    [
      "-opt-",
      "-obj-update-",
      [
        "-opt-",
        "-store-read-",
        "http-resp"
      ],
      [
        "name"
      ],
      [
        "-opt-",
        "-obj-read-",
        [
          "-opt-",
          "-store-read-",
          "http-reqt"
        ],
        [
          "name"
        ]
      ]
    ]
  ]
]
```

Let's ignore what the payload JSON means for now. Just save it. 

### 5. Configure the Davic Endpoint 

Did you save the JSON payload got from http://127.0.0.1:8080/getopt/hi2you
Send it to http://127.0.0.1:8080/davic/set in a POST request. (No additional headers are needed.) 

You should get a "200 OK" with an empty response body. 

### 6. Try the Davic Endpoint again. 

Let's repeat Step 3 -- send a POST request to http://127.0.0.1:8080/davic/go with the following payload: 
```json
{ 
    "id": "0123456789", 
    "name": "Jenny"
}
``` 

You should now get the following response: 
```json
{
    "message": "Hi!", 
    "name": "Jenny"
}
```
The Davic endpoint is now having the exactly same behavior as the Hi2You endpoint! 

How that is possible!? 

### How It Works? 

In a sense, the two endpoints in our demo, Hi2You and Davic, have the same behavior as follows: 
1. Put the request payload (a JSON object) to the "http-reqt" field in Davic framework's "**Store**" -- you can image the **Store** as the "memory" of the Davic framework. 
2. Run the Davic program -- the JSON payload we got from the "getopt" endpoint. Remember what we got from Step 4 above? 
3. Get the object stored in the "http-resp" field of Davic framework's Store, and send it back as the response payload. 

The only different between the Hi2You and the Davic endpoints is that Hi2You has the fixed, pre-configured, Davic program. This Davic program is nothing but a piece of JSON data. 
On the other hand, the Davic endpoint accepts the program from the "/davic/set" endpoint -- the endpoint we hit in Step 5 of the above demo. 

In this demo, you can see how we define computation at runtime! 