package filesystem

import (
	"fmt"
	"os"
	"path"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/internal/file"
	filesystemClient "github.com/rwese/archivar/internal/filesystem"
	"github.com/rwese/archivar/internal/utils/config"

	"github.com/sirupsen/logrus"
)

// Filesystem archives directly on filesystem, whatever it may be
type FileSystemConfig struct {
	OverwriteExisting bool
	Directory         string
}

type FileSystem struct {
	client    *filesystemClient.FileSystem
	logger    *logrus.Logger
	directory string
}

func init() {
	archivers.Register(New)
}

// New will return a new filesystem archiver
func New(c interface{}, logger *logrus.Logger) archivers.Archiver {
	fsystemConfig := &FileSystemConfig{}
	config.ConfigFromStruct(c, &fsystemConfig)
	fsystem := &FileSystem{
		logger:    logger,
		directory: fsystemConfig.Directory,
		client:    filesystemClient.New(logger),
	}

	return fsystem
}

// Upload takes filename, fileDirectory and fileHandle stores it on the filesystem
func (fsystem *FileSystem) Upload(f file.File) (err error) {
	uploadFilePath := path.Join(fsystem.directory, f.Directory)
	fsystem.logger.Debugf("Storing file '%s' at '%s'", f.Filename, uploadFilePath)

	return fsystem.client.Upload(f.Filename, uploadFilePath, f.Body)
}

// Connect exists only to satisfy the archiver interface
func (fsystem *FileSystem) Connect() (err error) {
	if _, err := os.Stat(fsystem.directory); os.IsNotExist(err) {
		return fmt.Errorf("filesystem directory '%s' does not exist", fsystem.directory)
	}

	fsystem.logger.Debugf("Filesystem archive directory '%s'", fsystem.directory)
	return
}
