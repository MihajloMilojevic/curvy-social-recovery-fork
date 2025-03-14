package keyrecovery

type Share struct {
	Point        string
	SpendingEval string
	ViewingEval  string
}

func (s *Share) Serialize() ([]byte, error) {
	return []byte(s.Point + s.SpendingEval + s.ViewingEval), nil
}
