CREATE TABLE IF NOT EXISTS users(
    uid varchar(36) NOT NULL PRIMARY KEY,
    email TEXT NOT NULL,
    pass TEXT NOT NULL,
    age integer
);
CREATE UNIQUE INDEX IF NOT EXISTS email_id ON users (email);

CREATE TABLE IF NOT EXISTS books(
    bid varchar(36) NOT NULL PRIMARY KEY,
    lable TEXT NOT NULL,
    author TEXT NOT NULL,
    "desc" TEXT NOT NULL,
    age integer NOT NULL,
    count integer NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT false
);