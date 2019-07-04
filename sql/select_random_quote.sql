SELECT quotes.body, quotes.author, movies.title
FROM quotes
JOIN movies
ON movies.wikiquote_url = quotes.movie
OFFSET floor(random()*(
	SELECT count(id) FROM quotes
))
LIMIT 1;