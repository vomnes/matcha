Go back to [Table of Contents](../../../)

### Mails
#### POST - /v1/mails/forgotpassword
This route allows to send forgot password email.
```
JSON Body :
  {
    "email": string,
    "test": bool, // Avoid to send real email during the tests
  }
```
If email address from the body is empty or not a valid email  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad Request - JSON Content "Error: Email address \<details\>"  
If email address from the body doesn't match with any user in the database  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad Request - JSON Content "Error: Email address does not exists in the database"  
Generate a unique token using user firstname and current time  
Insert the unique token in the user row of the table Users in the database  
Send mail (http://localhost:3000/resetpassword/:token) :
- If in the body test is true, then the route with return HTTP Code 200 StatusOK with a JSON containing the email, fullname and forgotPasswordUrl, this is used for tests.
- Else send 'Forgot password' email to the email addres from body with variables firstname and forgotPasswordUrl used in the mailjet template, return HTTP Code 202 Status Accepted
