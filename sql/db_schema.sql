CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE characters (
	name VARCHAR(128),
	PRIMARY KEY(name)
);

CREATE TABLE movies (
	wikiquote_url varchar(128),
	title VARCHAR(128),
	last_scraped TIMESTAMP (0) DEFAULT now(),
	PRIMARY KEY(wikiquote_url)
);

CREATE TABLE quotes (
	id UUID DEFAULT gen_random_uuid(),
	author varchar(128),
	movie varchar(128),
	body TEXT NOT NULL,
	FOREIGN KEY(author) REFERENCES characters(name) ON DELETE CASCADE,
	FOREIGN KEY(movie) REFERENCES movies(wikiquote_url) ON DELETE CASCADE,
	PRIMARY KEY(id)
);