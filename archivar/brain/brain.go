package brain

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/brain/brains"
	_ "github.com/rwese/archivar/archivar/brain/brains/json"
	"github.com/sirupsen/logrus"
)

func NewBrain(backendType string, backendConfig interface{}, logger *logrus.Logger) brains.BrainInterface {
	backendArchiver := archivers.GetArchiver(backendType, backendConfig, logger)
	backendGatherer := archivers.GetGatherer(backendType, backendConfig, backendArchiver, logger)

	brainBackend := brains.GetBrain(backendType, backendConfig, backendArchiver, backendGatherer, logger)

	return brainBackend
}
