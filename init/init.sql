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
DROP TABLE IF EXISTS forums_users CASCADE;

CREATE EXTENSION IF NOT EXISTS CITEXT;

-- Tables

---- Users
CREATE UNLOGGED TABLE customers
(
    email    CITEXT COLLATE "C" NOT NULL UNIQUE,
    nickname CITEXT COLLATE "C" PRIMARY KEY,
    fullname TEXT               NOT NULL,
    about    TEXT
);

---- Forums
CREATE UNLOGGED TABLE forums
(
    slug    CITEXT COLLATE "C" PRIMARY KEY,
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

CREATE UNIQUE INDEX unique_vote_idx on votes (author, thread);

---- Forum Users
CREATE UNLOGGED TABLE forums_users
(
    forum    CITEXT COLLATE "C" NOT NULL REFERENCES forums (slug) ON DELETE CASCADE,
    nickname CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_idx_forum_users on forums_users (nickname, forum);

---- Update path
CREATE OR REPLACE FUNCTION update_path()
    RETURNS TRIGGER AS
$BODY$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT thread, path
        FROM posts
        WHERE thread = NEW.thread AND id = NEW.parent
        INTO first_parent_thread, parent_path;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'Parent post not found in current thread' USING ERRCODE = '00404';
        END IF ;
        NEW.path := parent_path || NEW.id;
    END IF;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER path_updater
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();

CREATE OR REPLACE FUNCTION update_forum_users_by_insert_into_threads()
RETURNS TRIGGER AS
$BODY$
BEGIN
    INSERT INTO forums_users (forum, nickname) values (NEW.forum, NEW.author)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;


CREATE TRIGGER thread_insert_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_users_by_insert_into_threads()

-- Update forum threads
CREATE OR REPLACE FUNCTION update_forum_threads()
RETURNS TRIGGER AS
$BODY$
BEGIN
    UPDATE forums SET threads = (threads + 1) WHERE slug = NEW.forum;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER upd_forum_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_threads();

-- Update thread votes
CREATE OR REPLACE FUNCTION update_votes()
RETURNS TRIGGER AS
    $BODY$
    BEGIN
        IF NEW.voice = OLD.voice THEN
            RETURN NEW;
        END IF;
        IF NEW.voice > 0 THEN
            UPDATE threads SET votes = (votes + 2) WHERE id = NEW.thread;
        ELSE
            UPDATE threads SET votes = (votes - 2) WHERE id = NEW.thread;
        END IF;
        RETURN NEW;
    END;
    $BODY$ LANGUAGE plpgsql;

CREATE TRIGGER update_votes_trigger
    AFTER UPDATE
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE update_votes();

-- Insert vote
CREATE OR REPLACE FUNCTION insert_votes()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    IF NEW.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1) WHERE id = NEW.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1) WHERE id = NEW.thread;
    END IF;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;


CREATE TRIGGER insert_vote_trigger
    AFTER INSERT
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE insert_votes()