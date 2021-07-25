package log

import (
	"bytes"
	"fmt"
	stdlog "log"
	"net/http"
)

func SetClientLogger(serviceUrl string, clientService string) {
	stdlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stdlog.SetFlags(0)
	stdlog.SetOutput(&clientLogger{
		url: serviceUrl,
	})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message,Service responsed with code %v", res.StatusCode)
	}
	return len(data), nil
}
