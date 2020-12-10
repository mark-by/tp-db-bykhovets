ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET random_page_cost = '1.1';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET seq_page_cost = '0.1';
ALTER SYSTEM SET random_page_cost = '0.1';
ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';

DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE EXTENSION IF NOT EXISTS CITEXT;

-- Tables

---- Users
CREATE UNLOGGED TABLE customers
(
    id       SERIAL PRIMARY KEY,
    email    CITEXT COLLATE "C" NOT NULL UNIQUE,
    nickname CITEXT COLLATE "C" NOT NULL UNIQUE,
    fullname TEXT               NOT NULL,
    about    TEXT
);

---- Forums
CREATE UNLOGGED TABLE forums
(
    id      SERIAL PRIMARY KEY,
    slug    CITEXT COLLATE "C" NOT NULL UNIQUE,
    author  CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE,
    title   TEXT               NOT NULL,
    threads BIGINT DEFAULT 0,
    posts   BIGINT DEFAULT 0
);

---- Threads
CREATE UNLOGGED TABLE threads
(
    id      SERIAL PRIMARY KEY,
    slug    CITEXT COLLATE "C" unique,
    title   TEXT               NOT NULL,
    message TEXT               NOT NULL,
    created TIMESTAMP WITH TIME ZONE,
    votes   INT DEFAULT 0,
    author  CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE,
    forum   CITEXT COLLATE "C" NOT NULL REFERENCES forums (slug) ON DELETE CASCADE
);

---- Posts
CREATE UNLOGGED TABLE posts
(
    id        BIGSERIAL PRIMARY KEY,
    message   TEXT                  NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE NOT NULL,
    parent    INTEGER DEFAULT 0,
    created   TIMESTAMP WITH TIME ZONE,
    author    CITEXT COLLATE "C"    NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE,
    thread    INTEGER               NOT NULL REFERENCES threads ON DELETE CASCADE,
    forum     CITEXT COLLATE "C"    NOT NULL REFERENCES forums (slug) ON DELETE CASCADE,
    path      BIGINT[]
);

---- Votes
CREATE UNLOGGED TABLE votes
(
    id     SERIAL PRIMARY KEY,
    voice  SMALLINT           NOT NULL,
    author CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE,
    thread INTEGER            NOT NULL REFERENCES threads ON DELETE CASCADE
);

---- Forum Users
CREATE UNLOGGED TABLE forums_users
(
    forum    CITEXT COLLATE "C" NOT NULL REFERENCES forums (slug) ON DELETE CASCADE,
    nickname CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE
);
