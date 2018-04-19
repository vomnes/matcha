# Websocket (WS) Documentation

#### WS - /ws/{jwt}

The Websocket server listen and serve on
```
ws://localhost:8081/ws/{jwt}'
```
The parameter jwt contains the JSON Web Token of the connected user.  

#### Login

The first Websocket connection of a user create the login state.  
This will update the online status (true) and the online status update date in the database.  
Send to every websocket connection (except this just login user) the message (5 sec delay to handle page refresh) :  
```
JSON Message :
  {
    "event": "login",
    "username": string,
  }
```

#### Logout

When the last Websocket connection of a user is closer the logout state is created.  
This will update the online status (false) and the online status update date in the database.  
Send to every websocket connection the message (5 sec delay to handle page refresh) :  
```
JSON Message :
  {
    "event": "logout",
    "username": string,
  }
```

#### Message
When the server will receive the following message, the message data will be
inserted in the database using the message data.  
A message will be emitted on the target ws connection.  
```
JSON Received :
  {
    "event": "message",
    "target": string,   // Username
    "data": string,
  }
```
```
JSON Emitted :
  {
    "event": "message",
    "data": {
      "from":    string, // Username in the JWT
      "content": string,
    },
  }
```

#### Targeted events
When the server will receive a message with the event "view", "like", "match", "unmatch" or "isTyping" a message will we sent on the ws connections of the targeted user.  
```
JSON Received :
  {
    "event":  string, // "view", "like", "match", "unmatch" or "isTyping"
    "target": string, // Username
  }
```
```
JSON Emitted :
  {
    "event": string, // "view", "like", "match", "unmatch" or "isTyping"
    "from":  string, // Username in the JWT
  }
```
