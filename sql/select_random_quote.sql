SELECT quotes.body, quotes.author, movies.title, movies.wikiquote_url as wikiquoteurl
FROM quotes
JOIN movies
ON movies.wikiquote_url = quotes.movie
OFFSET floor(random()*(
	SELECT count(id) FROM quotes
))
LIMIT 1;