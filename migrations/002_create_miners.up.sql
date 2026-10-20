CREATE TABLE miners (
    mine_id INTEGER NOT NULL REFERENCES mines(id),
    id INTEGER NOT NULL,
    name TEXT NOT NULL,
    energy INTEGER NOT NULL DEFAULT 0,
    production INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (mine_id, id)
);