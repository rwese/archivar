package gatherer

type Gatherer interface {
	Download() (err error)
}
