Go back to [Table of Contents](../../../)

### Chat
#### GET - /v1/chat/matches

Collect the user's matchesIDs in the database  
If matchesIDs is empty  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 200 OK - JSON Content "data: No matches"

```
JSON Content Response :
  {
    "data":   string,
  }
```

Get in the database for each id, the username, firstname, lastname, picture url, online status, the last message (content/date) and total unread messages  
Everything is stored in a structure, sorted by last message date and unread message count  
Return HTTP Code 200 Status OK - JSON Content Structure  

```
JSON Content Response :
  [
    {
      "username":                 string,
      "firstname":                string,
      "lastname":                 string,
      "picture_url":              string,
      "last_message_content":     string,
      "last_message_date":        time.Time,
      "online":                   bool,
      "total_unread_messages":    int,
    },
  ]
```

---

#### GET - /v1/chat/messages/{username}

If parameter username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
If logged username is equal to targetUsername  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad Request - JSON Content "Error: Cannot target your own profile"  
Get targetUserID from targetUsername  
Collect the discussion between logged user and target user and the profiles data in the database, sort by asc  
Update all the messages as read in the database  
If there are no messages  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 200 OK - JSON Content "data: "No messages"  
```
JSON Content Response :
  {
    "data": string,
  },
```
Else  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return HTTP Code 200 Status OK - JSON Content Messages  

```
JSON Content Response :
  [
    {
      "username":     string,
      "lastname":     string,
      "firstname":    string,
      "picture_url":  string,
      "content":      string,
      "received_at":  time.Time,
    },
  ]
```

---

#### POST - /v1/chat/messages/{username}

If parameter username is not a valid username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"  
If logged username is equal to targetUsername  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Return an error - HTTP Code 400 Bad Request - JSON Content "Error: Cannot target your own profile"  
Get targetUserID from targetUsername  
Update all the messages as read in the discussion between logged user and target user  
Return HTTP Code 200 Status OK  
