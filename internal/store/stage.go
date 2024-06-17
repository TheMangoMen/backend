package store

func (s *Store) GetIsRankingStage() (isRankingStage bool, err error) {
	query := "SELECT * FROM Stage;"
	err = s.db.Get(&isRankingStage, query)
	return isRankingStage, err
}
