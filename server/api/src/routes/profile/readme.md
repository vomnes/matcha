# API routes details

## Table of Contents
- [Accounts](../accounts)
- [Mails](../mails)

### Profiles
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

Sanitise by removed the space after and before the variables and escaping characters  
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
Trim and escape characters of the IP  
If the IP is not a valid IP4  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: IP in the body is invalid"  
Collect the latitude, longitude, city, zip and country linked to this IP using ip-api.com's API  
Update the geoposition of the user using this new data, geolocalisation_allowed still false  
city and country are formated as Title and ZIP as upper case  
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
    "birthday":                string, /* Format Date DD/MM/YYYY */
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

___

#### POST - /v1/profiles/picture/{number}

```
JSON Body :
  {
    "picture_base64" string
  }
```

If the url parameter number is not a number between 1 and 5  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Url parameter must be a number between 1 and 5, not <number>"  
Convert the base64 picture a file picture, support only png, jpg and jpeg files  
The file picture is stored in '/storage/pictures/profiles/<username>' on the server (specific directory for tests)  
If the file generating failed  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 500 Server Internal Error - JSON Content "Error: Failed to generate <type>"  
If the file type is not supported  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Image type <type> not accepted, support only png, jpg and jpeg images"  
Update the picture path in the database  
Remove the old file on the server by using the old path from the database  
Return HTTP Code 200 Status OK  

```
JSON Content Response :
  {
    "picture_url":                string,
  }
```

___

#### DELETE - /v1/profiles/picture/{number}

If the url parameter number is not a number between 1 and 5  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Url parameter must be a number between 1 and 5, not <number>"  
Not possible to only remove the first picture (only update)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 403 Forbidden - JSON Content "Error: Not possible to delete the 1st picture - Only upload a new one is possible"  
Update the picture path in the database with an empty string  
Remove the old file on the server by using the old path from the database  
Return HTTP Code 200 Status OK  

___

#### POST - /v1/profiles/edit/tag

```
JSON Body :
  {
    "tag_name" string
  }
```

If in the body tag name is empty  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Tag name in body can't be empty"  
Set trim, escape characters and to lower case tag name  
If tag name is not valid  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Tag name is not valid"  
If the tag name doesn't exists in the table Tags of the database we need to insert it  
Collect his tagID  
If the user own own already this tag    
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Tag name already linked to this user"  
Else link the userId with the tagID in the table Users_Tags  
Return HTTP Code 200 Status OK

```
JSON Content Response :
  {
    "tag_id":  string
  }
```

___

#### DELETE - /v1/profiles/edit/tag

```
JSON Body :
  {
    "tag_id" string
  }
```

If in the body tagID is empty  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Tag ID in body can't be empty"  
Set trim, escape characters and to lower case tagID  
If tag name is not valid  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Tag ID is not valid"  
Delete the link between the tagID and the userID in the database in the table Users_Tags  
Return HTTP Code 200 Status OK  
