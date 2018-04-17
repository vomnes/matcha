\echo '------- Initialize db_matcha --------'
DROP DATABASE IF EXISTS db_matcha;
CREATE DATABASE db_matcha;

GRANT ALL PRIVILEGES ON DATABASE db_matcha TO vomnes;

\connect db_matcha;

CREATE TABLE Users (
  ID SERIAL PRIMARY KEY,
  username VARCHAR (65)     NOT NULL,
  email VARCHAR (255) NOT NULL,
  lastname VARCHAR (65)     NOT NULL,
  firstname VARCHAR (65)     NOT NULL,
  password VARCHAR (65)     NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp,
  random_token VARCHAR (254) DEFAULT '',
  picture_url_1 VARCHAR (255) DEFAULT '',
  picture_url_2 VARCHAR (255) DEFAULT '',
  picture_url_3 VARCHAR (255) DEFAULT '',
  picture_url_4 VARCHAR (255) DEFAULT '',
  picture_url_5 VARCHAR (255) DEFAULT '',
  biography VARCHAR (2048) DEFAULT '',
  birthday timestamp DEFAULT NULL,
  genre VARCHAR (64) DEFAULT 'male',
  interesting_in  VARCHAR (255) DEFAULT 'bisexual',
  city VARCHAR (65)     DEFAULT '',
  zip VARCHAR (65)     DEFAULT '',
  country VARCHAR (65)     DEFAULT '',
  latitude DECIMAL(20,6) DEFAULT 0.0,
  longitude DECIMAL(20,6) DEFAULT 0.0,
  geolocalisation_allowed BOOLEAN NOT NULL DEFAULT FALSE,
  online BOOLEAN NOT NULL DEFAULT FALSE,
  online_status_update_date timestamp with time zone DEFAULT NULL,
  rating DECIMAL(9,6) DEFAULT 2.5
);

CREATE TABLE Tags (
  ID SERIAL PRIMARY KEY,
  name VARCHAR (65) NOT NULL
);

CREATE TABLE Users_Tags (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  tagID SERIAL NOT NULL
);

CREATE TABLE Likes (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  liked_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE Visits (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  visited_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE Fake_Reports (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  target_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE Messages (
  ID SERIAL PRIMARY KEY,
  senderID SERIAL NOT NULL,
  receiverID SERIAL NOT NULL,
  content VARCHAR (255) default '',
  created_at timestamp with time zone DEFAULT current_timestamp,
  is_read BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE Notifications (
  ID SERIAL PRIMARY KEY,
  typeID SERIAL NOT NULL,
  userID SERIAL NOT NULL,
  target_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp,
  is_read BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE Notifications_Types (
  ID SERIAL PRIMARY KEY,
  name VARCHAR (65) NOT NULL,
  description VARCHAR (255) default ''
);

INSERT INTO Notifications_Types (name, description) VALUES ('view', 'A user has viewed the profile of an other user');
INSERT INTO Notifications_Types (name, description) VALUES ('like', 'A user has liked the profile of an other user');
INSERT INTO Notifications_Types (name, description) VALUES ('match', 'A liked profile has liked your profile');
INSERT INTO Notifications_Types (name, description) VALUES ('unmatch', 'One of the two users of a match has unliked the other profile');
INSERT INTO Notifications_Types (name, description) VALUES ('message', 'A user has sent a message to an other user');

-- SQL Functions Lib --

-- geodistance return the distance in kilometers between two GPS (latitude and longitude) coordinates
CREATE OR REPLACE FUNCTION geodistance(alat double precision, alng double precision, blat double precision, blng double precision)
  RETURNS double precision AS
$BODY$
SELECT (2 * 6371 *
  asin(
    sqrt(
      sin(radians($3 - $1) / 2) ^ 2 +
      cos(radians($1)) *
      cos(radians($3)) *
      sin(radians($4 - $2) / 2) ^ 2
    )
  )) AS distance;
$BODY$
LANGUAGE sql IMMUTABLE COST 100;

INSERT INTO Users (username, email, lastname, firstname, password) VALUES ('vomnes', 'valentin.omnes@gmail.com', 'Omnes', 'Valentin', '$2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua');

\echo '----- Initialize db_matcha_tests -----'

DROP DATABASE IF EXISTS db_matcha_tests;
CREATE DATABASE db_matcha_tests;

GRANT ALL PRIVILEGES ON DATABASE db_matcha_tests TO vomnes;

\connect db_matcha_tests;

CREATE TABLE Users (
  ID SERIAL PRIMARY KEY,
  username VARCHAR (65)     NOT NULL,
  email VARCHAR (255) NOT NULL,
  lastname VARCHAR (65)     NOT NULL,
  firstname VARCHAR (65)     NOT NULL,
  password VARCHAR (65)     NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp,
  random_token VARCHAR (254) DEFAULT '',
  picture_url_1 VARCHAR (255) DEFAULT '',
  picture_url_2 VARCHAR (255) DEFAULT '',
  picture_url_3 VARCHAR (255) DEFAULT '',
  picture_url_4 VARCHAR (255) DEFAULT '',
  picture_url_5 VARCHAR (255) DEFAULT '',
  biography VARCHAR (255) DEFAULT '',
  birthday timestamp DEFAULT NULL,
  genre VARCHAR (64) DEFAULT 'male',
  interesting_in  VARCHAR (255) DEFAULT 'bisexual',
  city VARCHAR (65)     DEFAULT '',
  zip VARCHAR (65)     DEFAULT '',
  country VARCHAR (65)     DEFAULT '',
  latitude DECIMAL(9,6) DEFAULT NULL,
  longitude DECIMAL(9,6) DEFAULT NULL,
  geolocalisation_allowed BOOLEAN NOT NULL DEFAULT FALSE,
  online BOOLEAN NOT NULL DEFAULT FALSE,
  online_status_update_date timestamp with time zone DEFAULT NULL,
  rating DECIMAL(9,6) DEFAULT 2.5
);

CREATE TABLE Tags (
  ID SERIAL PRIMARY KEY,
  name VARCHAR (65) NOT NULL
);

CREATE TABLE Users_Tags (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  tagID SERIAL NOT NULL
);

CREATE TABLE Likes (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  liked_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE Visits (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  visited_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE Fake_Reports (
  ID SERIAL PRIMARY KEY,
  userID SERIAL NOT NULL,
  target_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE Messages (
  ID SERIAL PRIMARY KEY,
  senderID SERIAL NOT NULL,
  receiverID SERIAL NOT NULL,
  content VARCHAR (255) default '',
  created_at timestamp with time zone DEFAULT current_timestamp,
  is_read BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE Notifications (
  ID SERIAL PRIMARY KEY,
  typeID SERIAL NOT NULL,
  userID SERIAL NOT NULL,
  target_userID SERIAL NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp,
  is_read BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE Notifications_Types (
  ID SERIAL PRIMARY KEY,
  name VARCHAR (65) NOT NULL,
  description VARCHAR (255) default ''
);

INSERT INTO Notifications_Types (name, description) VALUES ('view', 'A user has viewed the profile of an other user');
INSERT INTO Notifications_Types (name, description) VALUES ('like', 'A user has liked the profile of an other user');
INSERT INTO Notifications_Types (name, description) VALUES ('match', 'A liked profile has liked your profile');
INSERT INTO Notifications_Types (name, description) VALUES ('unmatch', 'One of the two users of a match has unliked the other profile');
INSERT INTO Notifications_Types (name, description) VALUES ('message', 'A user has sent a message to an other user');

-- SQL Functions Lib --

-- geodistance return the distance in kilometers between two GPS (latitude and longitude) coordinates
CREATE OR REPLACE FUNCTION geodistance(alat double precision, alng double precision, blat double precision, blng double precision)
  RETURNS double precision AS
$BODY$
SELECT (2 * 6371 *
  asin(
    sqrt(
      sin(radians($3 - $1) / 2) ^ 2 +
      cos(radians($1)) *
      cos(radians($3)) *
      sin(radians($4 - $2) / 2) ^ 2
    )
  )) AS distance;
$BODY$
LANGUAGE sql IMMUTABLE COST 100;
