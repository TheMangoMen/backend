package model

type User struct {
	UID string `json:"uid"`
}

type Ranking struct {
	UID             string `json:"uid" db:"uid"`
	JID             int    `json:"jid" db:"jid"`
	UserRanking     int    `json:"userranking" db:"userranking"`
	EmployerRanking string `json:"employerranking" db:"employerranking"`
}

type Contribution struct {
	UID            string `json:"uid" db:"uid"`
	JID            int    `json:"jid" db:"jid"`
	OA             bool   `json:"oa" db:"oa"`
	InterviewStage int    `json:"interviewstage" db:"interviewstage"`
	OfferCall      bool   `json:"offercall" db:"offercall"`
}

type Stage struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Job struct {
	Watching bool    `json:"watching"`
	JID      int     `json:"jid"`
	Title    string  `json:"title"`
	Company  string  `json:"company"`
	Location string  `json:"location"`
	Openings int     `json:"openings"`
	Stages   []Stage `json:"stages"`
}

type InterviewRow struct {
	JID        int    `db:"jid"`
	Title      string `db:"title"`
	Company    string `db:"company"`
	Season     string `db:"season"`
	Year       string `db:"year"`
	Cycle      int    `db:"cycle"`
	Location   string `db:"location"`
	Openings   int    `db:"openings"`
	Watching   bool   `db:"watching"`
	OACount    int    `db:"oacount"`
	Int1Count  int    `db:"int1count"`
	Int2Count  int    `db:"int2count"`
	Int3Count  int    `db:"int3count"`
	OfferCount int    `db:"offercount"`
}

type RankingRow struct {
	JID       int    `db:"jid"`
	Title     string `db:"title"`
	Company   string `db:"company"`
	Season    string `db:"season"`
	Year      string `db:"year"`
	Cycle     int    `db:"cycle"`
	Location  string `db:"location"`
	Openings  int    `db:"openings"`
	Watching  bool   `db:"watching"`
	Ranked    int    `db:"ranked"`
	Taking    int    `db:"taking"`
	NotTaking int    `db:"nottaking"`
}

type WatchedStatusCount struct {
	Status string `db:"status"`
	Count  int    `db:"count"`
}
