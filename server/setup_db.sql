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
  online BOOLEAN NOT NULL DEFAULT FALSE
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
  online BOOLEAN NOT NULL DEFAULT FALSE
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
