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

Sanitize by removed the space after and before the variables and escaping characters  
If any elements in the body is not valid  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Not a valid <detail>"  
Convert string format time from body to \*time.Time  
Update the table Users in the database with the new values  
If a new field is empty then this field won't be updated  
Return HTTP Code 200 Status OK  
___

#### POST - /v1/profiles/edit/location

```
JSON Body :
  {
    "lat"     float64,
    "lng"     float64,
    "city"    string,
    "zip"     string,
    "country" string
  }
```

If any field in the body is empty  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: No field inside the body can be empty"  
If the latitude or longitude is in overflow  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <type> value is over the limit"  
Trim and escape characters of city, zip and country  
If the city, zip or country is invalid (common name)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <detail> is invalid"  
Format as title the city name and country and format as upper case the ZIP  
Update the table Users in the database with the new values and set geolocalisation_allowed as true  
Return HTTP Code 200 Status OK  

___

#### POST - /v1/profiles/edit/password

```
JSON Body :
  {
    "password"         string,
    "new_password"     string,
    "new_rePassword"   string
  }
```

If any field in the body is empty  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: No field inside the body can be empty"  
If the current or new password is invalid  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <type> password field is not a valid password"  
If the new password and re entered new password are not identical  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Both password entered must be identical"  
Check if the userId and username match with an row in the table Users of the database  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User does not exists in the database"  
Check if the current password in the body match with the current password of the user  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 403 Forbidden - JSON Content "Error: Current password incorrect"  
Encrypt the new password and update the table Users in the database  
Return HTTP Code 200 Status OK  

___

#### GET - /v1/profiles/edit

```
JSON Body :
  {
    "ip" string
  }
```

Collect the data concerning the users in the table Users of the database  
If the user doesn't exists  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"  
Collect the tags (id, name) concerning the users in database  
If geolocalisation_allowed is false we need to set or update the location of the users by using the IP in the body  
Trim and escapte characters of the IP  
If the IP is not a valid IP4  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: IP in the body is invalid"  
Collect the latitude, longitude, city, zip and country linked to this IP using ip-api.com's API  
Update the geoposition of the user using this new data, geolocalisation_allowed still false  
city and country are fomated as Title and ZIP as upper case  
Return HTTP Code 200 Status OK
```
JSON Content Response :
  {
    "username":                string,
    "email":                   string,
    "lastname":                string,
    "firstname":               string,
    "picture_url_1":           string,
    "picture_url_2":           string,
    "picture_url_3":           string,
    "picture_url_4":           string,
    "picture_url_5":           string,
    "biography":               string,
    "birthday":                string, // Format Date DD/MM/YYYY
    "genre":                   string,
    "interesting_in":          string,
    "latitude":                float64,
    "longitude":               float64,
    "city":                    string,
    "zip":                     string,
    "country":                 string,
    "geolocalisation_allowed": bool,
    "tags":                    []string,
  }
```
