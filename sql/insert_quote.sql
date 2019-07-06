INSERT INTO quotes(body, author, movie)
VALUES($1, $2, $3)
ON CONFLICT (body)
DO NOTHING;