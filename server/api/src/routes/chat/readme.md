Go back to [Table of Contents](../../../)

### Chat
#### GET - /v1/users/data/match

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
