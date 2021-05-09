package webdav

import "github.com/studio-b12/gowebdav"

func (w *Webdav) Connect() (newSession bool, err error) {
	if newSession = w.client == nil; !newSession {
		return
	}

	w.client = gowebdav.NewClient(w.Server, w.UserName, w.Password)
	w.logger.Debugf("connecting to %s as %s\n", w.Server, w.UserName)
	err = w.client.Connect()
	if err != nil {
		w.logger.Fatalf("failed to connect: %s\n", err.Error())
	}

	return
}
