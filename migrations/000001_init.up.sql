CREATE TABLE team
(
	name       TEXT      PRIMARY KEY NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL    DEFAULT NOW()
);

CREATE TABLE "user"
(
	id         TEXT      PRIMARY KEY NOT NULL UNIQUE,
	name       TEXT      NOT NULL,
	is_active  BOOLEAN   NOT NULL,
	team_name  TEXT      NOT NULL REFERENCES team (name),
	created_at TIMESTAMP NOT NULL    DEFAULT NOW(),
	updated_at TIMESTAMP
);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE prs
(
	id         TEXT                   PRIMARY KEY NOT NULL UNIQUE,
	title      TEXT                   NOT NULL,
	author_id  TEXT                   NOT NULL REFERENCES "user" (id),
	status     pr_status              NOT NULL    DEFAULT 'OPEN',
	created_at TIMESTAMP WITH TIME ZONE           DEFAULT NOW(),
	merged_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE pr_reviewers
(
	pr_id       TEXT NOT NULL REFERENCES prs (id) ON DELETE CASCADE,
	reviewer_id TEXT NOT NULL REFERENCES "user" (id),
	assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	PRIMARY KEY (pr_id, reviewer_id)
);