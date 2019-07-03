WITH added_movie AS (
	INSERT INTO movies(wikiquote_url, title)
    VALUES($1, $2)
    ON CONFLICT (wikiquote_url)
    DO UPDATE SET last_scraped = now()
    RETURNING movies
), added_characters AS (
	INSERT INTO characters(name)
    SELECT n
    FROM UNNEST($3::varchar[]) n
    ON CONFLICT (name)
    DO NOTHING
    RETURNING name
)
SELECT added_characters.name, added_movie.wikiquote_url
FROM added_characters
JOIN added_movie
ON TRUE;