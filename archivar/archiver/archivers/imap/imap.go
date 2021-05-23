package imap

import "github.com/rwese/archivar/archivar/archiver/archivers"

func init() {
	archivers.RegisterGatherer(NewGatherer)
}
