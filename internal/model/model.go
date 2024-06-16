package model

type User struct {
	UID string `json:"uid"`
}

type Ranking struct {
	// TODO
}

type Contribution struct {
	UID            string
	JID            string
	OA             bool
	InterviewStage int
	OfferCall      bool
}
