package filterResult

type Results int

const (
	Allow Results = iota
	Reject
	NoAction
)
