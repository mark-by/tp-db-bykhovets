ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';
ALTER SYSTEM SET random_page_cost = '0.1';
ALTER SYSTEM SET seq_page_cost = '0.1';


CREATE EXTENSION IF NOT EXISTS CITEXT;

-- Tables

---- Users
DROP TABLE IF EXISTS customers CASCADE;
CREATE UNLOGGED TABLE customers
(
    id       SERIAL PRIMARY KEY,
    email    CITEXT COLLATE "C" NOT NULL UNIQUE,
    nickname CITEXT COLLATE "C" NOT NULL UNIQUE,
    fullname TEXT               NOT NULL,
    about    TEXT
);

create index index_users_all on customers (nickname, fullname, email, about);
cluster customers using index_users_all;


---- Forums
DROP TABLE IF EXISTS forums CASCADE;
CREATE UNLOGGED TABLE forums
(
    id      SERIAL PRIMARY KEY,
    slug    CITEXT COLLATE "C" NOT NULL UNIQUE,
    author  CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE,
    title   TEXT               NOT NULL,
    threads BIGINT DEFAULT 0,
    posts   BIGINT DEFAULT 0
);

create index index_forum_slug_hash on forums using hash (slug);
create index index_users_fk on forums (author);
create index index_forum_all on forums (slug, title, author, posts, threads);

---- Threads
DROP TABLE IF EXISTS threads CASCADE;
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

create index index_thread_forum_created on threads (forum, created);
create index index_thread_slug on threads (slug);
create index index_thread_slug_hash on threads using hash (slug);
create index index_thread_all on threads (title, message, created, slug, author, forum, votes);
create index index_thread_users_fk on threads (author);
create index index_thread_forum_fk on threads (forum);

---- Posts
DROP TABLE IF EXISTS posts CASCADE;
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

create index index_post_thread_id on posts (thread, id);
create index index_post_thread_path on posts (thread, path);
create index index_post_thread_parent_path on posts (thread, parent, path);
create index index_post_path1_path on posts ((path[1]), path);
create index index_post_thread_created_id on posts (thread, created, id);
create index index_post_users_fk on posts (author);
create index index_post_forum_fk on posts (forum);

---- Votes
DROP TABLE IF EXISTS votes CASCADE;
CREATE UNLOGGED TABLE votes
(
    id     SERIAL PRIMARY KEY,
    voice  SMALLINT           NOT NULL,
    author CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE,
    thread INTEGER            NOT NULL REFERENCES threads ON DELETE CASCADE
);

create index index_vote_thread on votes (thread);
CREATE UNIQUE INDEX unique_vote_idx on votes (author, thread);

---- Forum Users
DROP TABLE IF EXISTS forums_users CASCADE;
CREATE UNLOGGED TABLE forums_users
(
    forum    CITEXT COLLATE "C" NOT NULL REFERENCES forums (slug) ON DELETE CASCADE,
    nickname CITEXT COLLATE "C" NOT NULL REFERENCES customers (nickname) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_idx_forum_users on forums_users (forum, nickname);
cluster forums_users using unique_idx_forum_users;

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
        IF NOT FOUND THEN
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
EXECUTE PROCEDURE update_forum_users_by_insert_into_threads();

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
    EXECUTE PROCEDURE insert_votes();

CREATE OR REPLACE FUNCTION update_forum_posts()
    RETURNS TRIGGER AS
$update_forum_posts$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = new.forum;
    RETURN new;
END;
$update_forum_posts$ LANGUAGE plpgsql;

CREATE TRIGGER update_forum_posts
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_posts();