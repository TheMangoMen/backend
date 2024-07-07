package store

import (
	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetWatchedStatusCount(uID string) ([]model.WatchedStatusCount, error) {
	var watchedStatusCount []model.WatchedStatusCount
	query := `
WITH statuses AS (
    SELECT 'Nothing' AS status, 1 AS rank
    UNION
    SELECT 'OA', 2
    UNION
    SELECT 'Interview', 3
    UNION
    SELECT 'Offer Call', 4
),
watched_status AS (
    SELECT
    c.jid,
    CASE
        WHEN BOOL_OR(c.offercall) THEN 'Offer Call'
        WHEN MAX(c.interviewstage) > 0 THEN 'Interview'
        WHEN BOOL_OR(c.oa) THEN 'OA'
        ELSE 'Nothing'
    END AS status
    FROM users u JOIN watching w ON w.uid = u.uid LEFT JOIN contributions c ON w.jid = c.jid
    WHERE u.uid = $1 GROUP BY c.jid
),
status_counts AS (
    SELECT status, COUNT(status) AS count
    FROM watched_status
    GROUP BY status
)
SELECT s.status, COALESCE(sc.count, 0) AS count
FROM statuses s LEFT OUTER JOIN status_counts sc ON s.status = sc.status
ORDER BY s.rank;
`
	err := s.db.Select(&watchedStatusCount, query, uID)
	if err != nil {
		return []model.WatchedStatusCount{}, err
	}
	return watchedStatusCount, nil

}
