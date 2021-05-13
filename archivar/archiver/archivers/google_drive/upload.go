package google_drive

import (
	"path"

	"github.com/rwese/archivar/internal/file"
)

func (g *GoogleDrive) Upload(f file.File) (err error) {
	_, err = g.Connect()
	if err != nil {
		return err
	}

	filePath := path.Join(g.uploadDirectory, f.Directory, f.Filename)
	_, err = g.drive.PutFile(filePath, f.Body)
	g.logger.Debugf("Uploaded '%s' to: %s", f.Filename, filePath)
	return
}
