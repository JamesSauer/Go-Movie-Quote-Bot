INSERT INTO movies(title, wikiquote_url)
VALUES($1, $2)
ON CONFLICT (wikiquote_url)
DO UPDATE SET last_scraped = now();