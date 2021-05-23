package filesystem

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
)

func init() {
	archivers.RegisterArchiver(NewArchiver)
	archivers.RegisterGatherer(NewGatherer)
}
