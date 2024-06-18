CREATE TABLE movies (
  id VARCHAR(50),
  name VARCHAR(50) NOT NULL,
  director VARCHAR(50) NOT NULL,
  duration_in_seconds INTEGER NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE rooms (
  id VARCHAR(50),
  number INTEGER NOT NULL,
  description VARCHAR(50) NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE room_movies (
  fk_room_id VARCHAR(50) NOT NULL,
  FOREIGN KEY (fk_room_id) REFERENCES rooms(id),
  fk_movie_id VARCHAR(50) NOT NULL,
  FOREIGN KEY (fk_movie_id) REFERENCES movies(id),
  PRIMARY KEY(fk_room_id, fk_movie_id)
);