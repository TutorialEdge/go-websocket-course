Go WebSocket course
=======================

Building a Real-Time Football Score App in Go and Vue.js.

## Frontend Overview

Frontend Web App will be built in Vue.js - when a client connects, it'll 
populate a list of games that you can click on to follow live scores.

You'll then connect to a specific game which will receive events from the
server every time a new event happens.

This will be fairly simple, and we'll mostly just be using it as a way to
demonstrate us connecting to specific channels.


## Running Locally

```
$ docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
$ go run cmd/server/main.go
```