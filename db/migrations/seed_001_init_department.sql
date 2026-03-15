
INSERT INTO departments (name)
VALUES ('Dev')
RETURNING *;

INSERT INTO departments (name)
VALUES ('QA')
RETURNING *;

INSERT INTO departments (name)
VALUES ('R&D')
RETURNING *;

INSERT INTO departments (name)
VALUES ('Support')
RETURNING *;