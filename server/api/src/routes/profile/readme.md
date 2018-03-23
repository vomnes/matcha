# API routes details

## Table of Contents
- [Accounts](../#accounts)
- [Mails](../#mails)

### Profiles
#### POST - /v1/account/register
```
JSON Body :
  {
    "username": string,
    "email": string,
    "lastname": string,
    "firstname": string,
    "password": string,
    "rePassword": string
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


```
JSON Body :
  {
    "username": string,
    "email": string,
    "lastname": string,
    "firstname": string,
    "password": string,
    "rePassword": string
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

#### POST - /v1/profiles/edit/data

```
JSON Body :
  {
    "lastname"       string,
    "firstname"      string,
    "email"          string,
    "biography"      string,
    "birthday"       string,
    "genre"          string,
    "interesting_in" string
  }
```

The body contains the lastname, firstname, email, biography, birthday, genre and interesting_in  
Sanitize by removed the space after and before the variables and escaping characters  
If any elements in the body is not valid  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Not a valid <detail>"  
Convert string format time from body to *time.Time  
Update the table Users in the database with the new values  
If a new field is empty then this field won't be updated  
Return HTTP Code 200 Status OK  
