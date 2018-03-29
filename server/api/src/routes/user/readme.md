Go back to [Table of Contents](../../../)

### Users
#### GET - /v1/users/match

```
JSON Body :
  {
    "age": {
     "min": int,
     "max": int
    },
    "rating": {
     "min": float,
     "max": float
    },
    "distance": {
     "max": int
    },
    "tags": []string,
    "lat": float,
    "lng": float,
    "sort_type": string,
    "sort_direction": string,
  }
```

The route match will return an array with the possible match of the connected user according his data and parameters in the body.  
Check input :  
- Age must be a float greater than 1  
- Rating must be a float between 0.1 and 5.0  
- If min > max, automatic swap  
- Distance is an integer with default value 50 (km)  
- Sort type available are age, common_tags (when there are no selected tags) distance, rating (default)  
- Sort direction available are reverse and normal  
- Finish position is an unsigned integer, default value 20  

Collect the logged in user data (users, tags)  
Handle genre by creating an array with the possible match  
Create the request according the logged in user data and options from the body that will the matching users  
- Default range age is between -3 and +3 the age of the logged in user  

Generate an map[string]interface{} array with the users from the SQL request output between StartPosition and FinishPosition  
Return HTTP Code 200 Status OK  
If the array is empty return JSON Content "data": "No (more) users"  

```
JSON Content Response :
  {
    "data":   string,
  }
```

Else JSON Content Array  

```
JSON Content Response :
  [
    {
      "username":    string,
      "firstname":   string,
      "lastname":    string,
      "picture_url": string,
      "age":         int,
      "rating":      float64,
      "latitude":    float64,
      "longitude":   float64,
      "distance":    float64 /* Round about 0.1 */
    },
  ]
```

___

#### GET - /v1/users/{username}

The url contains the parameter username  
If username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
Collect the data concerning the user in the table Users of the database  
If the user doesn't exists  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"  
Collect the tags (id, name) concerning the users (target/connected) in database  
Split the shared and not shared tags  
Check if the connected user  
- has liked the target user and so if they have liked each other  
- has reported the user as fake  
If the targetUser is not the connectedUser  
- Add a profile visit in the table Visits in the database  
- Update target user rating  

Return HTTP Code 200 Status OK - JSON Content User data  

```
JSON Content Response :
  {
    "username":         string,
    "firstname":        string,
    "lastname":         string,
    "biography":        string,
    "genre":            string,
    "interesting_in":   string,
    "location":         string,
    "age":              int,
    "pictures":         []string,
    "rating":           float64,
    "liked":            bool,
    "users_linked":     bool,
    "reported_as_fake": bool,
    "online":           bool,
    "tags": []interface{
      "shared":   sharedTags,
      "personal": notSharedTags,
    },
    "isMe": bool,
  }
```

___

#### POST - /v1/users/{username}/like

The url contains the parameter username  
If username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
If parameter username and logged in username identical  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad request - JSON Content "Error: Cannot like your own profile"  
Collect the userId corresponding to the username in the database  
If the username doesn't match with any data  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"  
If the profile is already liked by the connected user  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Profile already liked by the user"  
Insert like in the table Likes in the database  
Update target user rating  
Return HTTP Code 200 Status OK  

___

#### DELETE - /v1/users/{username}/like

The url contains the parameter username  
If username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
If parameter username and logged in username identical  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad request - JSON Content "Error: Cannot like your own profile"  
Collect the userId corresponding to the username in the database  
If the username doesn't match with any data  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"  
Remove the like from the table Likes in the database
Update target user rating
Return HTTP Code 200 Status OK  

___

#### POST - /v1/users/{username}/fake

The url contains the parameter username  
If username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
If parameter username and logged in username identical  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad request - JSON Content "Error: Cannot like your own profile"  
Collect the userId corresponding to the username in the database  
If the username doesn't match with any data  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"  
If the profile is already liked by the connected user  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Profile already reported as fake by the user"  
Insert fake in the table Fake_Reports in the database  
Update target user rating  
Return HTTP Code 200 Status OK  

___

#### DELETE - /v1/users/{username}/fake

The url contains the parameter username  
If username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
If parameter username and logged in username identical  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad request - JSON Content "Error: Cannot like your own profile"  
Collect the userId corresponding to the username in the database  
If the username doesn't match with any data  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"  
Remove the like from the table Fake_Reports in the database
Update target user rating
Return HTTP Code 200 Status OK  
