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
  latitude DECIMAL(9,6) DEFAULT NULL,
  longitude DECIMAL(9,6) DEFAULT NULL,
  geolocalisation_allowed BOOLEAN NOT NULL DEFAULT FALSE
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
  latitude DECIMAL(9,6) DEFAULT NULL,
  longitude DECIMAL(9,6) DEFAULT NULL,
  geolocalisation_allowed BOOLEAN NOT NULL DEFAULT FALSE
);
