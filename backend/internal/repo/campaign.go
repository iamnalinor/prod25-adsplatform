package repo

import (
	"backend/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CampaignRepo struct {
	db *sqlx.DB
}

func (r *CampaignRepo) Add(campaign model.Campaign) (err error) {
	_, err = r.db.Exec(
		`INSERT INTO campaigns (id, advertiser_id, ad_title, ad_text, start_date, end_date, targeting_gender,
                       targeting_age_from, targeting_age_to, targeting_location, cost_per_impression, 
                       impressions_limit, cost_per_click, clicks_limit, image_path, moderation_task_id) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		campaign.Id, campaign.AdvertiserId, campaign.AdTitle, campaign.AdText, campaign.StartDate,
		campaign.EndDate, campaign.CampaignTargeting.Gender, campaign.CampaignTargeting.AgeFrom,
		campaign.CampaignTargeting.AgeTo, campaign.CampaignTargeting.Location, campaign.CostPerImpression,
		campaign.ImpressionsLimit, campaign.CostPerClick, campaign.ClicksLimit, campaign.ImagePath,
		campaign.ModerationTaskId,
	)
	return
}

func (r *CampaignRepo) GetById(id uuid.UUID) (res model.Campaign, err error) {
	err = r.db.Get(&res, `SELECT * FROM campaigns_moderation WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return
}

func (r *CampaignRepo) GetList(advertiserId uuid.UUID, size int, page int) ([]model.Campaign, error) {
	offset := (page - 1) * size
	campaigns := make([]model.Campaign, 0)
	err := r.db.Select(&campaigns,
		`SELECT * FROM campaigns_moderation WHERE advertiser_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`,
		advertiserId, size, offset)
	return campaigns, err
}

func (r *CampaignRepo) Update(campaign model.Campaign) error {
	res, err := r.db.Exec(
		`UPDATE campaigns SET ad_title = $1, ad_text = $2, start_date = $3, end_date = $4, 
					   targeting_gender = $5, targeting_age_from = $6, targeting_age_to = $7, 
					   targeting_location = $8, cost_per_impression = $9, impressions_limit = $10, 
					   cost_per_click = $11, clicks_limit = $12, image_path = $13, moderation_task_id = $14
				WHERE id = $15`,
		campaign.AdTitle, campaign.AdText, campaign.StartDate, campaign.EndDate,
		campaign.CampaignTargeting.Gender, campaign.CampaignTargeting.AgeFrom, campaign.CampaignTargeting.AgeTo,
		campaign.CampaignTargeting.Location, campaign.CostPerImpression, campaign.ImpressionsLimit,
		campaign.CostPerClick, campaign.ClicksLimit, campaign.ImagePath, campaign.ModerationTaskId, campaign.Id,
	)
	if err != nil {
		return fmt.Errorf("run query: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("fetch affected rows: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CampaignRepo) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM campaigns WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("run query: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("fetch affected rows: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// GetStats aggregates impressions and clicks over all time.
// If campaignId is uuid.Nil, it is omitted from the query.
func (r *CampaignRepo) GetStats(advertiserId uuid.UUID, campaignId uuid.UUID) (stats model.CampaignStats, err error) {
	query := `SELECT
    COUNT(ai.spent) AS impressions_count,
    COUNT(ac.spent) AS clicks_count,
    COALESCE(SUM(ai.spent), 0) AS spent_impressions,
    COALESCE(SUM(ac.spent), 0) AS spent_clicks,
    COALESCE(SUM(ai.spent), 0) + COALESCE(SUM(ac.spent), 0) AS spent_total
FROM ad_impressions ai
JOIN campaigns c on ai.campaign_id = c.id
LEFT JOIN ad_clicks ac
    ON ai.client_id = ac.client_id AND ai.campaign_id = ac.campaign_id
WHERE c.advertiser_id = $1`
	args := []any{advertiserId}

	// I wanted to omit advertiserId when it equals uuid.Nil, but because their ids are set
	// by the API caller (and it's possible to set an ID that equals to uuid.Nil), there is a rare case
	// when the query will be incorrect.
	if campaignId != uuid.Nil {
		query += ` AND c.id = $2`
		args = append(args, campaignId)
	}

	err = r.db.Get(&stats, query, args...)
	if err == nil && stats.ImpressionsCount > 0 {
		stats.Conversion = float64(stats.ClicksCount) / float64(stats.ImpressionsCount) * 100
	}
	return
}

// GetStatsDaily aggregates impressions and clicks, grouped by each day.
// If campaignId is uuid.Nil, it is omitted from the query.
func (r *CampaignRepo) GetStatsDaily(advertiserId uuid.UUID, campaignId uuid.UUID) ([]model.CampaignStats, error) {
	where := "c.advertiser_id = $1"
	args := []any{advertiserId}
	if campaignId != uuid.Nil {
		where += ` AND c.id = $2`
		args = append(args, campaignId)
	}

	var impressions []model.CampaignStats
	err := r.db.Select(&impressions, fmt.Sprintf(`
SELECT
	COUNT(spent) AS impressions_count,
	COALESCE(SUM(spent), 0) AS spent_impressions,
	date
FROM ad_impressions
JOIN campaigns c on campaign_id = c.id
WHERE %s
GROUP BY date
ORDER BY date ASC
`, where), args...)
	if err != nil {
		return nil, fmt.Errorf("get impressions stats: %w", err)
	}

	var clicks []model.CampaignStats
	err = r.db.Select(&clicks, fmt.Sprintf(`
SELECT
	COUNT(spent) AS clicks_count,
	COALESCE(SUM(spent), 0) AS spent_clicks,
	date
FROM ad_clicks
JOIN campaigns c on campaign_id = c.id
WHERE %s
GROUP BY date
ORDER BY date ASC
`, where), args...)
	if err != nil {
		return nil, fmt.Errorf("get clicks stats: %w", err)
	}

	// Two-pointer approach to merge impressions and clicks

	stats := make([]model.CampaignStats, 0, len(impressions))
	for i, j := 0, 0; i < len(impressions) || j < len(clicks); {
		if j == len(clicks) || *impressions[i].Date < *clicks[j].Date {
			// clicks exhausted or ai.date < ac.date
			// => N impressions, 0 clicks for this date

			stats = append(stats, model.CampaignStats{
				ImpressionsCount: impressions[i].ImpressionsCount,
				ClicksCount:      0,
				Conversion:       0,
				SpentImpressions: impressions[i].SpentImpressions,
				SpentClicks:      0,
				SpentTotal:       impressions[i].SpentImpressions,
				Date:             impressions[i].Date,
			})
			i++
		} else if i == len(impressions) || *clicks[j].Date < *impressions[i].Date {
			// impressions exhausted or ai.date > ac.date
			// => 0 impressions, N clicks for this date

			stats = append(stats, model.CampaignStats{
				ImpressionsCount: 0,
				ClicksCount:      clicks[j].ClicksCount,
				Conversion:       0,
				SpentImpressions: 0,
				SpentClicks:      clicks[j].SpentClicks,
				SpentTotal:       clicks[j].SpentClicks,
				Date:             clicks[j].Date,
			})
		} else {
			// ai.date == ac.date
			if *impressions[i].Date != *clicks[j].Date {
				panic("assert failed: impression.Date != click.Date")
			}

			conversion := 0.0
			if impressions[i].ImpressionsCount != 0 {
				conversion = float64(clicks[j].ClicksCount) / float64(impressions[i].ImpressionsCount) * 100
			}

			stats = append(stats, model.CampaignStats{
				ImpressionsCount: impressions[i].ImpressionsCount,
				ClicksCount:      clicks[j].ClicksCount,
				Conversion:       conversion,
				SpentImpressions: impressions[i].SpentImpressions,
				SpentClicks:      clicks[j].SpentClicks,
				SpentTotal:       impressions[i].SpentImpressions + clicks[j].SpentClicks,
				Date:             impressions[i].Date,
			})
			i++
			j++
		}
	}

	return stats, nil
}

func (r *CampaignRepo) GetModerationFailed(size int, page int) ([]model.Campaign, error) {
	offset := (page - 1) * size
	campaigns := make([]model.Campaign, 0)
	err := r.db.Select(&campaigns,
		`SELECT * FROM campaigns_moderation WHERE moderation_result->>'acceptable' = 'false'
            ORDER BY created_at LIMIT $1 OFFSET $2`, size, offset)
	return campaigns, err
}

// GetAdCandidates fetches campaigns that could be a candidate for an ad for the specified client.
// The following criteria are applied:
// 1. The current date must be between campaign's start_date and end_date, inclusive.
// 2. The campaign must not exceed impressions_limit or clicks_limit.
// 3. If campaign has targeting by gender, the gender should match.
// 4. If campaign has targeting by age_from, the client age must be greater than or equal to it.
// 5. If campaign has targeting by age_to, the client age must be less than or equal to it.
// 6. If campaign has targeting by location, the client location must match.
// It includes data from ml_scores, ad_impressions and ad_clicks tables as described in model.AdCandidate.
// The result is ordered by the date of creation in ascending order.
func (r *CampaignRepo) GetAdCandidates(clientId uuid.UUID, limitsThreshold float64) ([]model.AdCandidate, error) {
	query := `
SELECT
    ad_id, ad_title, ad_text, ms.advertiser_id, image_path,
    cost_per_impression, impressions_count, impressions_limit, viewed,
    cost_per_click, clicks_count, clicks_limit, clicked,
    COALESCE(ms.score, 0) ml_score
FROM (
    SELECT
        c.id AS ad_id,
        c.created_at,
        c.ad_title,
        c.ad_text,
        c.advertiser_id,
        c.image_path,
        c.cost_per_impression,
        COUNT(ai.spent) AS impressions_count,
        c.impressions_limit,
        MAX(CASE WHEN ai.client_id = $1 THEN 1 ELSE 0 END) AS viewed,
        c.cost_per_click,
        COUNT(ac.spent) AS clicks_count,
        c.clicks_limit,
        MAX(CASE WHEN ac.client_id = $1 THEN 1 ELSE 0 END) AS clicked
    FROM
        campaigns c
            CROSS JOIN (SELECT "current_date" FROM settings) s
            CROSS JOIN (SELECT gender, age, location FROM clients WHERE id = $1) cl
            LEFT JOIN ad_impressions ai ON c.id = ai.campaign_id
            LEFT JOIN ad_clicks ac ON c.id = ac.campaign_id AND ai.client_id = ac.client_id
    WHERE
        (c.start_date <= s."current_date" AND s."current_date" <= c.end_date)
      AND (c.targeting_gender = 'ALL' OR c.targeting_gender IS NULL OR c.targeting_gender::TEXT = cl.gender::TEXT)
      AND (c.targeting_age_from IS NULL OR cl.age >= c.targeting_age_from)
      AND (c.targeting_age_to IS NULL OR cl.age <= c.targeting_age_to)
      AND (c.targeting_location IS NULL OR cl.location = c.targeting_location)
    GROUP BY c.id
    ) as cc
LEFT JOIN ml_scores ms ON cc.advertiser_id = ms.advertiser_id AND ms.client_id = $1
WHERE
	(cc.impressions_limit > 0 AND ((cc.impressions_count::float + 1) / cc.impressions_limit::float) <= $2)
ORDER BY cc.created_at
`

	campaigns := make([]model.AdCandidate, 0)
	err := r.db.Select(&campaigns, query, clientId, limitsThreshold)
	return campaigns, err
}

// AddAdImpression adds a record that the ad was viewed.
// This function is idempotent - if the record already exists, no error is returned.
func (r *CampaignRepo) AddAdImpression(impression model.AdImpression) error {
	_, err := r.db.Exec(`INSERT INTO ad_impressions (client_id, campaign_id, spent, date)
						VALUES ($1, $2, $3, $4) ON CONFLICT (client_id, campaign_id) DO NOTHING`,
		impression.ClientId, impression.CampaignId, impression.Spent, impression.Date)
	return err
}

func (r *CampaignRepo) GetAdImpression(clientId uuid.UUID, campaignId uuid.UUID) (res model.AdImpression, err error) {
	err = r.db.Get(&res, `SELECT * FROM ad_impressions WHERE client_id = $1 AND campaign_id = $2`, clientId, campaignId)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return
}

// AddAdClick adds a record that the ad was clicked.
// This function is idempotent - if the record already exists, no error is returned.
func (r *CampaignRepo) AddAdClick(click model.AdClick) error {
	_, err := r.db.Exec(`INSERT INTO ad_clicks (client_id, campaign_id, spent, date)
						VALUES ($1, $2, $3, $4) ON CONFLICT (client_id, campaign_id) DO NOTHING`,
		click.ClientId, click.CampaignId, click.Spent, click.Date)
	return err
}
