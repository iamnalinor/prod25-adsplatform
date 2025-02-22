CREATE TABLE api_requests (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    endpoint TEXT NOT NULL,
    duration_ms FLOAT NOT NULL
);

CREATE TABLE settings (
    id SERIAL PRIMARY KEY,
    "current_date" INT NOT NULL,
    moderation_enabled BOOL NOT NULL
);
INSERT INTO settings ("current_date", moderation_enabled) VALUES (0, false);

CREATE TABLE ai_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    type TEXT NOT NULL,
    prompt TEXT NOT NULL,
    "format" JSONB NOT NULL
);

CREATE TABLE ai_task_results (
    task_id UUID PRIMARY KEY REFERENCES ai_tasks(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    answer JSONB NOT NULL
);

CREATE TYPE gender AS ENUM ('MALE', 'FEMALE');
CREATE TYPE targeting_gender AS ENUM ('MALE', 'FEMALE', 'ALL');

CREATE TABLE clients (
    id UUID PRIMARY KEY,
    login TEXT NOT NULL,
    age INT NOT NULL,
    location TEXT NOT NULL,
    gender gender NOT NULL
);

CREATE TABLE advertisers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE ml_scores (
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    advertiser_id UUID NOT NULL REFERENCES advertisers(id) ON DELETE CASCADE,
    score INT NOT NULL,
    PRIMARY KEY (client_id, advertiser_id)
);

CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    advertiser_id UUID NOT NULL REFERENCES advertisers(id) ON DELETE CASCADE,
    impressions_limit INT NOT NULL,
    clicks_limit INT NOT NULL,
    cost_per_impression FLOAT NOT NULL,
    cost_per_click FLOAT NOT NULL,
    ad_title TEXT NOT NULL,
    ad_text TEXT NOT NULL,
    start_date INT NOT NULL,
    end_date INT NOT NULL,
    targeting_gender targeting_gender,
    targeting_age_from INT,
    targeting_age_to INT,
    targeting_location TEXT,
    image_path TEXT NOT NULL,
    moderation_task_id UUID REFERENCES ai_tasks(id) ON DELETE RESTRICT
);

CREATE INDEX campaigns_start_date_end_date_index ON campaigns(start_date, end_date);
CREATE INDEX campaigns_targeting_gender_index ON campaigns(targeting_gender);
CREATE INDEX campaigns_targeting_age_from_index ON campaigns(targeting_age_from);
CREATE INDEX campaigns_targeting_age_to_index ON campaigns(targeting_age_to);
CREATE INDEX campaigns_targeting_location_index ON campaigns(targeting_location);

CREATE VIEW campaigns_moderation AS
    SELECT
        c.*,
        r.answer AS moderation_result
    FROM campaigns c
    JOIN advertisers a ON c.advertiser_id = a.id
    LEFT JOIN ai_task_results r on c.moderation_task_id = r.task_id;

CREATE TABLE ad_impressions (
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    spent FLOAT NOT NULL,
    date INT NOT NULL,
    PRIMARY KEY (client_id, campaign_id)
);

CREATE TABLE ad_clicks (
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    spent FLOAT NOT NULL,
    date INT NOT NULL,
    PRIMARY KEY (client_id, campaign_id),
    FOREIGN KEY (client_id, campaign_id) REFERENCES ad_impressions (client_id, campaign_id) ON DELETE CASCADE
);
