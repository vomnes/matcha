# API routes details

### Account
#### POST - /v1/account/register
```
JSON Body :
  {
    "username": string,
    "email": string,
    "lastname": string,
    "firstname": string,
    "password": string,
    "re-password": string
  }
```
This route allows to handle the user registration by using the data sent in the body  
- Body Fields can't be empty, it must be a valid username (a-zA-Z0-9.- _ \\ {6,64}), firstname
and lastname (a-zA-Z - {6,64}), password (a-zA-Z0-9 {8,100} - At least one of each) and
email address (max 254)
- Password and reentered password must be identical

If a least one of points below is not respected :  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <error details>"  
Check in our PostgreSQL database, if the Username or/and Email address are already used  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: \<details\> already used"  
Encrypt the password and insert in the database the new user  
Return HTTP Code 201 Status Created

___

#### POST - /v1/account/login
```
JSON Body :
  {
    "username": string,
    "password": string,
    "uuid": string, // Universally unique identifier from the user's web browser
  }
```
This route allows to handle the user authentication by using the data sent in the body.  
If the Username from the body is not in our PostgreSQL database  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 403 Forbidden - JSON Content "Error: User or password incorrect"  
If the Password from the body does not match with the data linked to the username in our PostgreSQL database  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 403 Forbidden - JSON Content "Error: User or password incorrect"

Generate a JSON Web Token (JWT) with payload content :
```
{
 "iss":      "matcha.com",
 "sub":      UUID, // From body
 "userId":   UserID, // From body
 "username": Username, // From body
 "iat":      now, // As time the number of seconds elapsed since January 1, 1970 UTC
 "exp":      now + 72h, // As time the number of seconds elapsed since January 1, 1970 UTC
}
```
Set in the Redis database the key `Username + "-" + UUID` with the JWT as value  
Return HTTP Code 200 Status OK - JSON Content "token": JWT

> All the routes following the login, must contain in the header :  
> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**_Authorization: Bearer \<User_JWT\>_**  
> This token will be checked by the middleware for authentication.

___

#### JSON Web Token explanation

###### Basics
JSON Web Token is a JSON-based open standard (RFC 7519) for creating access tokens that assert some number of claims.
This token is composed of :
- Header
```
{
  "alg": "HS256",
  "typ":"JWT",
}
```
- Payload that contains data such as iss, sub, iat (token issued at), exp (token expiration date) and other personal data (userId, username)
- Signature - A secret string

>token = encodeBase64Url(header) + '.' + encodeBase64Url(payload) + '.' + encodeBase64Url(signature)

JWT is then used to identify the user, it is sent through the header **_Authorization: Bearer \<User_JWT\>_** and
we can decode the token to collect data from payload (check validity, private data).

###### JWT in this project
<img alt="JSON Web Token Schema" src="../../screenshots/JWT.png" width="60%" title="JWT Schema">

___

#### POST - /v1/account/logout
This route allows to handle the user logout  
Delete in the Redis database the key `Username + "-" + UUID` allowing to validate the JWT token, using context data  
If deletion failed  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 500 Internal Server Error - JSON Content "Error: Failed to delete token"  
Return HTTP Code 202 Status Accepted

___
