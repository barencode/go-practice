INSERT INTO mines (id, coal, reserves, closed)
VALUES (1, 0, 10, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO miners (mine_id, id, name, energy, production)
VALUES
    (1, 101, 'Борис', 2, 3),
    (1, 309, 'Анна', 1, 5)
ON CONFLICT (mine_id, id) DO NOTHING;