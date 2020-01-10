package heartbeat

import (
	"os"
	"time"

	"github.com/fengjun2016/gDStorage/dataServer/pkg/rabbitmq"
)

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
