package db

type CarEntity struct {
	CarID         string
	Owner         string
	BodyStyle     string
	OriginalColor string
	CurrentColor  string
	Painted       bool
}
