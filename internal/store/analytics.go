package store

import (
	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetWatchedStatusCount(uID string) ([]model.WatchedStatusCount, error) {
	var watchedStatusCount []model.WatchedStatusCount
	query := `
with statuses as (
    select 'nothing' as status, 1 as rank
    union
    select 'oa', 2
    union
    select 'interview', 3
    union
    select 'offercall', 4
),
watched_status as (
    select
    c.jid,
    case
        when BOOL_OR(c.offercall) then 'offercall'
        when max(c.interviewstage) > 0 then 'interview'
        when BOOL_OR(c.oa) then 'oa'
        else 'nothing'
    end as status
    from users u join watching w ON w.uid = u.uid left join contributions c on w.jid = c.jid
    where u.uid = $1 group by c.jid
),
status_counts as (
    select status, count(status) as count
    from watched_status
    group by status
)
select s.status, coalesce(sc.count, 0) as count
from statuses s left outer join status_counts sc on s.status = sc.status
order by s.rank;
`
	err := s.db.Select(&watchedStatusCount, query, uID)
	if err != nil {
		return []model.WatchedStatusCount{}, err
	}
	return watchedStatusCount, nil

}
