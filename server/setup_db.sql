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
  created_at timestamp with time zone DEFAULT current_timestamp
);

INSERT INTO Users (username, email, lastname, firstname, password) VALUES ('vomnes', 'valentin.omnes@gmail.com', 'Omnes', 'Valentin', 'abc');
