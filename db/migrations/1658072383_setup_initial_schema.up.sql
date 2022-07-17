BEGIN;

-- ATTEMPTS TABLE

CREATE TABLE IF NOT EXISTS attempts (
     channel_id varchar NOT NULL,
     user_id varchar NOT NULL,
     user_name varchar NOT NULL,
     "day" int4 NOT NULL,
     attempts int4 NULL,
     max_attempts int4 NULL,
     success bool NOT NULL,
     attempts_json json NOT NULL,
     posted_at timestamptz NOT NULL,
     message_id varchar NOT NULL,
     hard_mode bool NULL DEFAULT false,
     CONSTRAINT attempts_pk PRIMARY KEY (channel_id, user_id, day)
);

CREATE INDEX IF NOT EXISTS attempts_attempts_btree_idx ON attempts USING btree (attempts);
CREATE INDEX IF NOT EXISTS attempts_channel_id_hash_idx ON attempts USING btree (channel_id);
CREATE INDEX IF NOT EXISTS attempts_day_hash_idx ON attempts USING btree (day);
CREATE INDEX IF NOT EXISTS attempts_posted_at_btree_idx ON attempts USING btree (posted_at);
CREATE INDEX IF NOT EXISTS attempts_success_hash_idx ON attempts USING btree (success);
CREATE INDEX IF NOT EXISTS attempts_user_id_hash_idx ON attempts USING btree (user_id);
CREATE INDEX IF NOT EXISTS attempts_user_name_idx ON attempts USING btree (user_name);


-- TRACKED CHANNELS TABLE

CREATE TABLE IF NOT EXISTS tracked_channels (
     channel_id varchar NOT NULL,
     CONSTRAINT tracked_channels_pk PRIMARY KEY (channel_id)
);

COMMIT;