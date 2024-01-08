CREATE TABLE passage (
  id serial PRIMARY KEY,
  book TEXT NOT NULL,
  start_chapter INT NOT NULL,
  start_verse INT NOT NULL,
  end_chapter INT NOT NULL,
  end_verse INT NOT NULL
);