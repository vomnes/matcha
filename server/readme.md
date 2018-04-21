# Server

This is the backend of the project written in Golang, everything is unit tested.

## Documentation
- [API](./api)
- [Websocket](./websocket)

**Note to myself for my next web project:**
> Store useful variables (db, userid, username, errors has occurred) inside a structure

> Create a lib to handle easily and in standardised way the database requests


## Launch on UNIX system

### Setup
You must have Golang, PostgreSQL and Redis installed/started.

#### Initialise database
```
  psql -f setup_db.sql
```

#### Export environment variables
```
  export DB_NAME=db_matcha
  export DB_NAME_TEST=db_matcha_tests
  export MJ_APIKEY_PUBLIC=<your_mailjet_api_key_public>
  export MJ_APIKEY_PRIVATE=<your_mailjet_api_key_private>
  export JWT_SECRET=ItsTheJWTSecret
```

### API
```
   cd server/api
```
#### Run
```
   make vendor_get
   make run
```
#### Test
```
   make test
   # OR
   make test-verbose
```

### WebSocket
```
   cd server/websocket
```
#### Run
```
   make vendor_get
   make run
```
#### Test
```
   make test
   # OR
   make test-verbose
```
