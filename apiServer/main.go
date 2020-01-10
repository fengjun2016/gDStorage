package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fengjun2016/gDStorage/apiServer/heartbeat"
	"github.com/fengjun2016/gDStorage/apiServer/locate"
	"github.com/fengjun2016/gDStorage/apiServer/objects"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
