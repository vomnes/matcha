# matcha
This project is about creating a dating website.

I have designed and realised this project from scratch using ReactJS, Golang, PostgreSQL, Redis, HTML and SCSS (CSS) during my studies at 'Ecole 42'.
It was first big project using ReactJS, PostgreSQL, Redis and SCSS.

I haven't used any HTML/CSS framework, everything from scratch with a responsive web design.

## Documentation
- [Client](./client)
- [Server](./server)

## Features

### Profile
- Create an account
- Reset your password through an unique link sent by email (Using third party MailJet)
- Modified your private data
  * Firstname
  * Lastname
  * Email address
  * Biography
  * Birth date
  * Genre
  * Interesting in
  * Password
- Add/Edit
  * Profile pictures - Max 5 - 1 Profile picture
  * Location - Default IP location
  * Tag list

<img alt="Login" src="./screenshots/login.png" width="75%" title="login">
<img alt="Profile" src="./screenshots/my-profile.png" width="75%" title="profile">

### Browsing
- List suggested profiles (age, interesting in, location)
- Search and filter profiles according age, location, rating, tags
- A map with the position of the matched profiles

<img alt="Browsing" src="./screenshots/browsing.png" width="75%" title="browsing">
<img alt="Browsing responsive" src="./screenshots/browsing-responsive.png" width="75%" title="browsing responsive">

### See profile
- User data
  * Firstname
  * Lastname
  * Biography
  * Age
  * Genre
  * Interesting in
  * Shared and personal tags
  * Online/Offline for ... (live - websocket)
- See on the same page the next and previous matched profile and be able to see the profile
- Like/Dislike the profile
- Report as fake account/Remove fake account report (a reported as fake user, block notification and doesn't appear any more in the searches)

<img alt="Profile 1" src="./screenshots/profile-1.png" width="75%" title="profile 1">
<img alt="Profile 2" src="./screenshots/profile-2.png" width="75%" title="profile 2">

### Chat with the matches (both liked each other)
- Show all the messages shared
- Able to chat with the matched profiles (live - websocket)

<img alt="Matches and notifications" src="./screenshots/matches.png" width="75%" title="Matches and notifications">

### Notifications
- See when the logged user has received a view (profile), like, match, unmatch and a new message
