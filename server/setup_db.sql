DROP DATABASE IF EXISTS db_matcha;
CREATE DATABASE db_matcha;

GRANT ALL PRIVILEGES ON DATABASE db_matcha TO vomnes;

\connect db_matcha;

CREATE TABLE Users (
  ID SERIAL PRIMARY KEY,
  username VARCHAR     NOT NULL,
  created_at timestamp with time zone DEFAULT current_timestamp
);

INSERT INTO Users (username) VALUES ('Valentin');
