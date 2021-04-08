package http

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

var regexHeaderResponseCode = regexp.MustCompile(`^HTTP/\d\.\d (\d+) \w+`)
var regexHeadersBody = regexp.MustCompile(`(\S+): (.*)`)

func AddCode(w http.ResponseWriter, headersBody string) (err error) {
	returnCodeMatch := regexHeaderResponseCode.FindAllStringSubmatch(headersBody, -1)[0][1]
	headerResponseCode, err := strconv.Atoi(string(returnCodeMatch))
	if err != nil {
		logrus.Error(err)
	}
	w.WriteHeader(headerResponseCode)
	return
}

func AddHeaders(w http.ResponseWriter, headersBody string) (err error) {
	headers := regexHeadersBody.FindAllStringSubmatch(headersBody, -1)
	for _, header := range headers {
		headerKey := header[1]
		headerValue := header[2]
		w.Header().Add(headerKey, headerValue)
	}

	AddCode(w, headersBody)
	return
}
