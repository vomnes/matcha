# API routes details

## Account
### POST '/v1/account/register'
Input body - JSON :
```
  {
    "username": string,
    "email": string,
    "lastname": string,
    "firstname": string,
    "password": string,
    "re-password": string
  }
```
This route will allows to handle the user registration by using the data sent in the body.
Fields can't be empty, it must be a valid username (a-zA-Z0-9.- _ \\ {6,64}), firstname
and lastname (a-zA-Z - {6,64}), password (a-zA-Z0-9 {8,100}- At least one of each) and
email address (max 254). Password and reenter password must be identical.
Username and email address must not be already used.
If everything is correct the user is inserted in the Users table of the database,
password is hashed with bcrypt, correct http status code is 201 - Created.
