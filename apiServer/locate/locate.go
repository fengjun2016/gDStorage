package locate

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fengjun2016/gDStorage/apiServer/pkg/rabbitmq"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	m := req.Method
	if m != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(req.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	rw.Write(b)
}

func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
