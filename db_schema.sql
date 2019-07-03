CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE characters (
	id UUID DEFAULT gen_random_uuid(),
	name VARCHAR(128),
	PRIMARY KEY(id)
);

CREATE TABLE movies (
	id UUID DEFAULT gen_random_uuid(),
	title VARCHAR(128),
	wikiquote_url varchar(128),
	last_scraped TIMESTAMP (0) DEFAULT now(),
	PRIMARY KEY(id)
);

CREATE TABLE quotes (
	id UUID DEFAULT gen_random_uuid(),
	author UUID,
	movie UUID,
	body TEXT NOT NULL,
	FOREIGN KEY(author) REFERENCES characters(id) ON DELETE CASCADE,
	FOREIGN KEY(movie) REFERENCES movies(id) ON DELETE CASCADE,
	PRIMARY KEY(id)
);