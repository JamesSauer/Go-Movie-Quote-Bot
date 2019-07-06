SELECT 'movie_count' AS stat, count(*) AS value FROM movies
UNION
SELECT 'character_count' AS stat, count(*) AS value FROM characters
UNION
SELECT 'quote_count' AS stat, count(*) AS value FROM quotes;