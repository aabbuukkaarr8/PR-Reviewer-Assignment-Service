-- Insert realistic test data

-- Teams
INSERT INTO teams (team_name) VALUES
    ('backend'),
    ('frontend'),pullrequests
    ('devops'),
    ('qa')
ON CONFLICT (team_name) DO NOTHING;

-- Backend team members
INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('user_backend_001', 'alex.petrov', 'backend', true),
    ('user_backend_002', 'maria.ivanova', 'backend', true),
    ('user_backend_003', 'dmitry.sidorov', 'backend', true),
    ('user_backend_004', 'anna.kuznetsova', 'backend', true),
    ('user_backend_005', 'sergey.volkov', 'backend', false)
ON CONFLICT (user_id) DO NOTHING;

-- Frontend team members
INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('user_frontend_001', 'elena.smirnova', 'frontend', true),
    ('user_frontend_002', 'pavel.kozlov', 'frontend', true),
    ('user_frontend_003', 'olga.lebedeva', 'frontend', true),
    ('user_frontend_004', 'ivan.popov', 'frontend', true)
ON CONFLICT (user_id) DO NOTHING;

-- DevOps team members
INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('user_devops_001', 'maxim.orlov', 'devops', true),
    ('user_devops_002', 'svetlana.novikova', 'devops', true),
    ('user_devops_003', 'andrey.morozov', 'devops', true)
ON CONFLICT (user_id) DO NOTHING;

-- QA team members
INSERT INTO users (user_id, username, team_name, is_active) VALUES
    ('user_qa_001', 'tatiana.romanova', 'qa', true),
    ('user_qa_002', 'nikolay.sokolov', 'qa', true),
    ('user_qa_003', 'ekaterina.vasilieva', 'qa', true),
    ('user_qa_004', 'vladimir.fedorov', 'qa', false)
ON CONFLICT (user_id) DO NOTHING;

-- Pull Requests
-- Open PRs
INSERT INTO pullrequests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at) VALUES
    ('pr_backend_001', 'Implement user authentication service', 'user_backend_001', 'OPEN', ARRAY['user_backend_002', 'user_backend_003'], NOW() - INTERVAL '2 days'),
    ('pr_backend_002', 'Add database connection pooling', 'user_backend_002', 'OPEN', ARRAY['user_backend_001', 'user_backend_004'], NOW() - INTERVAL '1 day'),
    ('pr_frontend_001', 'Create login page component', 'user_frontend_001', 'OPEN', ARRAY['user_frontend_002', 'user_frontend_003'], NOW() - INTERVAL '3 days'),
    ('pr_frontend_002', 'Implement responsive navigation menu', 'user_frontend_003', 'OPEN', ARRAY['user_frontend_001'], NOW() - INTERVAL '5 hours'),
    ('pr_devops_001', 'Setup CI/CD pipeline for staging', 'user_devops_001', 'OPEN', ARRAY['user_devops_002'], NOW() - INTERVAL '1 day'),
    ('pr_qa_001', 'Add integration tests for API endpoints', 'user_qa_001', 'OPEN', ARRAY['user_qa_002'], NOW() - INTERVAL '4 hours')
ON CONFLICT (pull_request_id) DO NOTHING;

-- Merged PRs
INSERT INTO pullrequests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at) VALUES
    ('pr_backend_003', 'Refactor error handling middleware', 'user_backend_003', 'MERGED', ARRAY['user_backend_001', 'user_backend_002'], NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days'),
    ('pr_frontend_003', 'Fix memory leak in dashboard component', 'user_frontend_002', 'MERGED', ARRAY['user_frontend_001', 'user_frontend_004'], NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days'),
    ('pr_devops_002', 'Configure monitoring alerts', 'user_devops_002', 'MERGED', ARRAY['user_devops_001'], NOW() - INTERVAL '14 days', NOW() - INTERVAL '12 days'),
    ('pr_qa_002', 'Update test coverage documentation', 'user_qa_002', 'MERGED', ARRAY['user_qa_001'], NOW() - INTERVAL '6 days', NOW() - INTERVAL '4 days')
ON CONFLICT (pull_request_id) DO NOTHING;

