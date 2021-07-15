package log

import (
	"github.com/watertreestar/go-toy/pkg/log"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
)

var logger *log.Logger

type fileLog string

var slog *stdlog.Logger

var f *os.File

// 接收POST请求，写入logger
// Write 实现了Writer接口
func (fl fileLog) Write(data []byte) (int, error) {
	defer f.Close()
	return f.Write(data)
}

// Init 初始化
func Init(destination string) {
	logger = log.NewLogger("./logs", "log", "LogService", log.DEBUG)
	slog = stdlog.New(fileLog(destination), "LogService", stdlog.LstdFlags)
	var err error
	f, err = os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(request.Body)
			if err != nil || len(msg) == 0 {
				writer.WriteHeader(http.StatusBadRequest)
			}
			Write(string(msg))
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func Write(msg string) {
	slog.Println(msg)
}
