package google_drive

import (
	"io"
	"path"
)

func (g *GoogleDrive) Upload(fileName string, directory string, fileHandle io.Reader) (err error) {
	_, err = g.Connect()
	if err != nil {
		return err
	}

	filePath := path.Join(g.uploadDirectory, directory, fileName)
	_, err = g.drive.PutFile(filePath, fileHandle)
	g.logger.Debugf("Uploaded '%s' to: %s", fileName, filePath)
	return
}
