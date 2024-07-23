package store

import (
	"github.com/TheMangoMen/backend/internal/model"
)

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
),
tags AS (
    WITH MostCommonOADifficulty AS (
        SELECT 
            JID,
            OADifficulty,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE OADifficulty IS NOT NULL and OADifficulty <> ''
        GROUP BY JID, OADifficulty
    ),
    MostCommonOALength AS (
        SELECT 
            JID,
            OALength,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE OALength IS NOT NULL and OALength <> ''
        GROUP BY JID, OALength
    ),
    MostCommonInterviewVibe AS (
        SELECT 
            JID,
            InterviewVibe,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE InterviewVibe IS NOT NULL and InterviewVibe <> ''
        GROUP BY JID, InterviewVibe
    ),
    MostCommonInterviewTechnical AS (
        SELECT 
            JID,
            InterviewTechnical,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE InterviewTechnical IS NOT NULL and InterviewTechnical <> ''
        GROUP BY JID, InterviewTechnical
    ),
    MostCommonOfferComp AS (
        SELECT 
            JID,
            OfferComp,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE OfferComp IS NOT NULL
        GROUP BY JID, OfferComp
    )
    SELECT 
        d.JID,
        d.OADifficulty,
        l.OALength,
        v.InterviewVibe,
        t.InterviewTechnical,
        c.OfferComp
    FROM MostCommonOADifficulty d
    LEFT JOIN MostCommonOALength l ON d.JID = l.JID AND l.RowNum = 1
    LEFT JOIN MostCommonInterviewVibe v ON d.JID = v.JID AND v.RowNum = 1
    LEFT JOIN MostCommonInterviewTechnical t ON d.JID = t.JID AND t.RowNum = 1
    LEFT JOIN MostCommonOfferComp c ON d.JID = c.JID AND c.RowNum = 1
    WHERE d.RowNum = 1
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
    t.OADifficulty,
    t.OALength,
    t.InterviewVibe,
    t.InterviewTechnical,
    t.OfferComp,
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
LEFT JOIN tags t
    ON j.JID = t.JID
WHERE j.season = (SELECT * FROM season) AND j.year = (SELECT * FROM year) AND j.cycle = (SELECT * FROM cycle)
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

		tags := model.JobTags{}

		if row.OADifficulty.Valid == false {
			tags.OADifficulty = ""
		} else {
			tags.OADifficulty = row.OADifficulty.String
		}

		if row.OALength.Valid == false {
			tags.OALength = ""
		} else {
			tags.OALength = row.OALength.String
		}

		if row.InterviewVibe.Valid == false {
			tags.InterviewVibe = ""
		} else {
			tags.InterviewVibe = row.InterviewVibe.String
		}

		if row.InterviewTechnical.Valid == false {
			tags.InterviewTech = ""
		} else {
			tags.InterviewTech = row.InterviewTechnical.String
		}

		if row.OfferComp.Valid == false {
			tags.OfferComp = 0
		} else {
			tags.OfferComp = row.OfferComp.Float64
		}

		job := model.Job{
			Watching: row.Watching,
			JID:      row.JID,
			Title:    row.Title,
			Company:  row.Company,
			Location: row.Location,
			Openings: row.Openings,
			Tags:     tags,
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
),
tags AS (
    WITH MostCommonOADifficulty AS (
        SELECT 
            JID,
            OADifficulty,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE OADifficulty IS NOT NULL and OADifficulty <> ''
        GROUP BY JID, OADifficulty
    ),
    MostCommonOALength AS (
        SELECT 
            JID,
            OALength,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE OALength IS NOT NULL and OALength <> ''
        GROUP BY JID, OALength
    ),
    MostCommonInterviewVibe AS (
        SELECT 
            JID,
            InterviewVibe,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE InterviewVibe IS NOT NULL and InterviewVibe <> ''
        GROUP BY JID, InterviewVibe
    ),
    MostCommonInterviewTechnical AS (
        SELECT 
            JID,
            InterviewTechnical,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE InterviewTechnical IS NOT NULL and InterviewTechnical <> ''
        GROUP BY JID, InterviewTechnical
    ),
    MostCommonOfferComp AS (
        SELECT 
            JID,
            OfferComp,
            COUNT(*) AS Occurrence,
            ROW_NUMBER() OVER (PARTITION BY JID ORDER BY COUNT(*) DESC) AS RowNum
        FROM Tags
        WHERE OfferComp IS NOT NULL
        GROUP BY JID, OfferComp
    )
    SELECT 
        d.JID,
        d.OADifficulty,
        l.OALength,
        v.InterviewVibe,
        t.InterviewTechnical,
        c.OfferComp
    FROM MostCommonOADifficulty d
    LEFT JOIN MostCommonOALength l ON d.JID = l.JID AND l.RowNum = 1
    LEFT JOIN MostCommonInterviewVibe v ON d.JID = v.JID AND v.RowNum = 1
    LEFT JOIN MostCommonInterviewTechnical t ON d.JID = t.JID AND t.RowNum = 1
    LEFT JOIN MostCommonOfferComp c ON d.JID = c.JID AND c.RowNum = 1
    WHERE d.RowNum = 1
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
    t.OADifficulty,
    t.OALength,
    t.InterviewVibe,
    t.InterviewTechnical,
    t.OfferComp,
    COALESCE(r.ranked, 0) ranked,
    COALESCE(r.taking, 0) taking,
    COALESCE(r.nottaking, 0) nottaking,
    COALESCE(w.Watch, FALSE) watching
FROM jobs j
LEFT JOIN rankingcounts r
    ON j.JID = r.JID
LEFT JOIN watches w
    ON j.JID = w.JID
LEFT JOIN tags t
    ON j.JID = t.JID
WHERE j.season = (SELECT * FROM season) AND j.year = (SELECT * FROM year) AND j.cycle = (SELECT * FROM cycle)
ORDER BY j.company;
`
	err := s.db.Select(&rows, query, uID)
	if err != nil {
		return nil, err
	}

	jobs := make([]model.Job, 0, len(rows))
	for _, row := range rows {
		var stages []model.Stage
		if !(row.Ranked == 0 && row.Taking == 0 && row.NotTaking == 0) {
			stages = []model.Stage{
				{Name: "Ranked", Count: row.Ranked},
				{Name: "Taking", Count: row.Taking},
				{Name: "Not Taking", Count: row.NotTaking},
			}
		}

		tags := model.JobTags{}

		if row.OADifficulty.Valid == false {
			tags.OADifficulty = ""
		} else {
			tags.OADifficulty = row.OADifficulty.String
		}

		if row.OALength.Valid == false {
			tags.OALength = ""
		} else {
			tags.OALength = row.OALength.String
		}

		if row.InterviewVibe.Valid == false {
			tags.InterviewVibe = ""
		} else {
			tags.InterviewVibe = row.InterviewVibe.String
		}

		if row.InterviewTechnical.Valid == false {
			tags.InterviewTech = ""
		} else {
			tags.InterviewTech = row.InterviewTechnical.String
		}

		if row.OfferComp.Valid == false {
			tags.OfferComp = 0
		} else {
			tags.OfferComp = row.OfferComp.Float64
		}

		job := model.Job{
			Watching: row.Watching,
			JID:      row.JID,
			Title:    row.Title,
			Company:  row.Company,
			Location: row.Location,
			Openings: row.Openings,
			Tags:     tags,
			Stages:   stages,
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *Store) DeleteWatching(uID string, jID int) error {

	_, err := s.db.Exec("DELETE FROM Watching WHERE uID = $1 AND jID = $2;", uID, jID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateWatching(uID string, jID int) error {
	_, err := s.db.Exec("INSERT INTO Watching VALUES ($1, $2) ON CONFLICT DO NOTHING;", uID, jID)
	if err != nil {
		return err
	}
	return nil
}
