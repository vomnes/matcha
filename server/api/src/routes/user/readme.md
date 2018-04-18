Go back to [Table of Contents](../../../)

### Users
#### GET - /v1/chat/matches

```
JSON Encoded Base64 - Search-Parameters Header :
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
    "sort_type": string,      // age, rating, distance, common_tags
    "sort_direction": string, // normal or reverse
    "start_position": uint,   // default = 0
    "finish_position": uint,  // default = 20
  }
```

The route match will return an array with the possible match of the connected user according his data and parameters from the header (base64).  
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

---

#### GET - /v1/users/data/match/{username}

```
JSON Encoded Base64 - Search-Parameters Header :
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
    "sort_type": string,      // age, rating, distance, common_tags
    "sort_direction": string, // normal or reverse
  }
```

The route match will return an object with next and previous profile, related to username (url).  
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

Generate an map[string]interface{} array with the profiles next to (+1 && -1) the target username from url parameter, add only if exists  
Return HTTP Code 200 Status OK - JSON Content  

```
JSON Content Response :
  {
    "previous": {
      "username":    string,
      "firstname":   string,
      "lastname":    string,
      "picture_url": string,
    },
    "next": {
      "username":    string,
      "firstname":   string,
      "lastname":    string,
      "picture_url": string,
    },
  },
```

___

#### GET - /v1/users/data/me

Collect the data concerning the user in the table Users of the database, total_new_notifications and total_new_messages
If the user doesn't exists  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User[<username>] doesn't exists"  
Get list reported as fake usernames  
Create redirect array with the name of the profile variables (age, picture)  
If this array is not empty return this array

```
JSON Content Response :
  {
    "redirect":    []string,
  }
```

Return HTTP Code 200 Status OK - JSON Content User data  

```
JSON Content Response :
  {
    "username":                   string,
    "firstname":                  string,
    "lastname":                   string,
    "birthday":                   time.Time,
    "age":                        int,
    "profile_picture":            string,
    "lat":                        float64,
    "lng":                        float64,
    "total_new_notifications":    int,
    "total_new_messages":         int,
    "reported_as_fake_usernames": []string,
  }
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
    "logout_at":        time.Time,
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
Check if now the user are connected  
Handle PushNotif like and match
Return HTTP Code 200 Status OK  

```
JSON Content Response :
  {
    "users_linked":     bool,
  }
```

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
Handle PushNotif unmatch
Return HTTP Code 200 Status OK

```
JSON Content Response :
  {
    "users_were_linked":     bool,
  }
```  

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
___

#### GET - /v1/users/data/notifications

Collect the user's notifications in the the database with profile data  
If one of the notifications is mark as unread  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Update with is_read true all the notifications  
Return HTTP Code 200 Status OK  
If notifications list is empty JSON contains "data": "No notifications"  

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
      "type":             string,
      "date":             time.Time,
      "new":              bool,
      "username":         string,
      "firstname":        string,
      "lastname":         string,
      "user_picture_url": string,
    },
  ]
```
