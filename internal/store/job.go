package store

import "github.com/TheMangoMen/backend/internal/model"

func (s *Store) GetJobInterviews(uID string) ([]model.Job, error) {
	var rows []model.InterviewRow
	query := `
WITH interviewcounts AS (
    SELECT
        JID,
        SUM(CASE WHEN OA = true THEN 1 ELSE 0 END) oaCount,
        SUM(CASE WHEN InterviewStage >= 1 THEN 1 ELSE 0 END) int1Count,
        SUM(CASE WHEN InterviewStage >= 2 THEN 1 ELSE 0 END) int2Count,
        SUM(CASE WHEN InterviewStage >= 3 THEN 1 ELSE 0 END) int3Count,
        SUM(CASE WHEN OfferCall = true THEN 1 ELSE 0 END) offerCount
    FROM Contributions
    GROUP BY JID
),
watches AS (
    SELECT *, TRUE as watch from Watching where UID = $1
)
SELECT
    j.jid,
    j.title,
    j.company,
    coalesce(j.location, 'N/A') location,
    j.openings,
    j.season,
    j.year,
    j.cycle,
    COALESCE(i.oaCount, 0) oaCount,
    COALESCE(i.int1Count, 0) int1Count,
    COALESCE(i.int2Count, 0) int2Count,
    COALESCE(i.int3Count, 0) int3Count,
    COALESCE(i.OfferCount, 0) OfferCount,
    COALESCE(w.watch, FALSE) watching
FROM jobs j
LEFT JOIN interviewcounts i
    ON j.JID = i.JID
LEFT JOIN watches w
    ON j.JID = w.JID
ORDER BY j.company;
`
	err := s.db.Select(&rows, query, uID)
	if err != nil {
		return nil, err
	}

	jobs := make([]model.Job, 0, len(rows))
	for _, row := range rows {
		stages := []model.Stage{
			{Name: "OA", Count: row.OACount},
			{Name: "Interview 1", Count: row.Int1Count},
			{Name: "Interview 2", Count: row.Int2Count},
			{Name: "Interview 3", Count: row.Int3Count},
			{Name: "Offer Call", Count: row.OfferCount},
		}

		var filteredStages []model.Stage
		for _, stage := range stages {
			if stage.Count > 0 {
				filteredStages = append(filteredStages, stage)
			}
		}

		job := model.Job{
			Watching: row.Watching,
			JID:      row.JID,
			Title:    row.Title,
			Company:  row.Company,
			Location: row.Location,
			Openings: row.Openings,
			Stages:   filteredStages,
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *Store) GetJobRankings(uID string) ([]model.Job, error) {
	var rows []model.RankingRow
	query := `
WITH rankingcounts AS (
    SELECT
        JID,
        COUNT(*) ranked,
        SUM(CASE WHEN EmployerRanking = 'Offer' and UserRanking = 1 THEN 1 ELSE 0 END) taking,
        SUM(CASE WHEN EmployerRanking = 'Offer' and UserRanking = -1 THEN 1 ELSE 0 END) nottaking
    FROM Rankings
    GROUP BY JID
),
watches AS (
    SELECT *, TRUE AS watch FROM Watching WHERE UID = $1
)
SELECT
    j.jid,
    j.title,
    j.company,
    coalesce(j.location, 'N/A') location,
    j.openings,
    j.season,
    j.year,
    j.cycle,
    COALESCE(r.ranked, 0) ranked,
    COALESCE(r.taking, 0) taking,
    COALESCE(r.nottaking, 0) nottaking,
    COALESCE(w.Watch, FALSE) watching
FROM jobs j
LEFT JOIN rankingcounts r
    ON j.JID = r.JID
LEFT JOIN watches w
    ON j.JID = w.JID
ORDER BY j.company;
`
	err := s.db.Select(&rows, query, uID)
	if err != nil {
		return nil, err
	}

	jobs := make([]model.Job, 0, len(rows))
	for _, row := range rows {
		stages := []model.Stage{
			{Name: "Ranked", Count: row.Ranked},
			{Name: "Taking", Count: row.Taking},
			{Name: "Not Taking", Count: row.NotTaking},
		}

		job := model.Job{
			Watching: row.Watching,
			JID:      row.JID,
			Title:    row.Title,
			Company:  row.Company,
			Location: row.Location,
			Openings: row.Openings,
			Stages:   stages,
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
