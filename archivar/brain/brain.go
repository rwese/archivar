package brain

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/brain/brains"
	"github.com/sirupsen/logrus"
)

type BrainStorageBackend struct {
	archiver archivers.Archiver
	gatherer archivers.Gatherer
}

type BrainConfig struct {
	File string // File is the absolute path to the file on the storageBackend
}

type Brain struct {
	storageBackend BrainStorageBackend
}

func New(backendType string, backendConfig interface{}, logger *logrus.Logger) brains.Brain {
	backendArchiver := archivers.GetArchiver(backendType, backendConfig, logger)
	backendGatherer := archivers.GetGatherer(backendType, backendConfig, backendArchiver, logger)

	brainBackend := brains.GetBrain(backendType, backendConfig, backendArchiver, backendGatherer, logger)

	return brainBackend
}

func (brain *Brain) Connect() {

}
