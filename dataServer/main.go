package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fengjun2016/gDStorage/dataServer/heartbeat"
	"github.com/fengjun2016/gDStorage/dataServer/locate"
	"github.com/fengjun2016/gDStorage/dataServer/objects"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
