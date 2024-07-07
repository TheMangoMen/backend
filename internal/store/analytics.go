package store

import (
	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetWatchedJobsStatusCounts(uID string) ([]model.StatusCount, error) {
	var statusCounts []model.StatusCount
	query := `
WITH watched_status AS (
    SELECT
    c.jid,
    CASE
        WHEN BOOL_OR(c.offercall) THEN 'offercall'
        WHEN MAX(c.interviewstage) > 0 THEN 'interview'
        WHEN BOOL_OR(c.oa) THEN 'oa'
        ELSE 'nothing'
    END AS status
    FROM users u JOIN watching w ON w.uid = u.uid LEFT JOIN contributions c ON w.jid = c.jid
    WHERE u.uid = $1 GROUP BY c.jid
),
status_counts AS (
    SELECT status, COUNT(status) AS count
    FROM watched_status
    GROUP BY status
)
SELECT s.statusname AS status, COALESCE(sc.count, 0) AS count
FROM AnalyticsStatuses s LEFT OUTER JOIN status_counts sc ON s.statusid = sc.status
ORDER BY s.ranking;
`
	err := s.db.Select(&statusCounts, query, uID)
	if err != nil {
		return []model.StatusCount{}, err
	}
	return statusCounts, nil
}

func (s *Store) GetWatchedCompaniesStatusCounts(uID string) ([]model.StatusCount, error) {
	var statusCounts []model.StatusCount
	query := `
WITH watched_status AS (
    SELECT
    j2.company,
    CASE
        WHEN BOOL_OR(c.offercall) THEN 'offercall'
        WHEN MAX(c.interviewstage) > 0 THEN 'interview'
        WHEN BOOL_OR(c.oa) THEN 'oa'
        ELSE 'nothing'
    END AS status
    FROM users u JOIN watching w ON w.uid = u.uid
        JOIN jobs j1 ON w.jid = j1.jid
        JOIN jobs j2 ON j1.company = j2.company
        LEFT JOIN contributions c ON j2.jid = c.jid
    WHERE u.uid = $1
    GROUP BY j2.company
),
status_counts AS (
    SELECT status, COUNT(status) AS count
    FROM watched_status
    GROUP BY status
)
SELECT s.statusname AS status, COALESCE(sc.count, 0) AS count
FROM AnalyticsStatuses s LEFT OUTER JOIN status_counts sc ON s.statusid = sc.status
ORDER BY s.ranking;
`
	err := s.db.Select(&statusCounts, query, uID)
	if err != nil {
		return []model.StatusCount{}, err

	}
	return statusCounts, nil
}
