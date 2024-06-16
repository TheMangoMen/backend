package store

import (
	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetJobs(uID string) ([]model.Job, error) {
	var rows []model.JobRow
	err := s.db.Select(&rows, `
		with
goodjobs as (
    select j.*, coalesce(oaCount, 0) oaCount, coalesce(int1Count, 0) int1Count, coalesce(int2Count, 0) int2Count, coalesce(int3Count, 0) int3Count, coalesce(offerCount, 0) offerCount from jobs j
    left join (select jid, count(*) as oaCount from Contributions where oa='t' group by jid) oa on j.jid = oa.jid
    left join (select jid, count(*) as int1Count from Contributions where InterviewStage>=1 group by jid) int1 on j.jid = int1.jid
    left join (select jid, count(*) as int2Count from Contributions where InterviewStage>=2 group by jid) int2 on j.jid = int2.jid
    left join (select jid, count(*) as int3Count from Contributions where InterviewStage>=3 group by jid) int3 on j.jid = int3.jid
    left join (select jid, count(*) as offerCount from Contributions where offercall='t' group by jid) offer on j.jid = offer.jid
),
watches as (
    select *, TRUE as watch from Watching where UID = $1
)
select j.jid, j.title, j.company, coalesce(j.location, 'N/A') location, j.openings, j.season, j.year, j.cycle, j.oaCount, j.int1Count, j.int2Count, j.int3Count, j.offerCount, coalesce(w.watch, false) watching  from goodjobs j left join watches w on w.JID = j.JID
order by j.jid;`, uID)
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
			{Name: "Offer", Count: row.OfferCount},
		}
		// Truncate the trailing stages that have 0 counts
		for i := len(stages) - 1; i >= 0; i-- {
			if stages[i].Count > 0 {
				stages = stages[:i+1]
				break
			}
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

func (s *Store) CreateWatching(uID string, jIDs []string) error {
	for _, jID := range jIDs {
		_, err := s.db.Exec("INSERT INTO Watching VALUES ($1, $2)", uID, jID)
		if err != nil {
			return err
		}
	}
	return nil
}
