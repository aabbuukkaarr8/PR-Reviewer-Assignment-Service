CREATE TABLE pullrequests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    assigned_reviewers TEXT[] DEFAULT '{}',
    "createdAt" TIMESTAMP,
    "mergedAt" TIMESTAMP
);

CREATE INDEX idx_pullrequests_author_id ON pullrequests(author_id);
CREATE INDEX idx_pullrequests_status ON pullrequests(status);
CREATE INDEX idx_pullrequests_assigned_reviewers ON pullrequests USING GIN(assigned_reviewers);

