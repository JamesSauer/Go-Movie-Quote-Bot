INSERT INTO characters(name)
VALUES($1)
ON CONFLICT (name)
DO NOTHING;