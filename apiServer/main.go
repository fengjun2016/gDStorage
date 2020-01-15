package main

import (
	"log"
	"net/http"
	"networkStorage/apiServer/versions"
	"os"

	"github.com/fengjun2016/gDStorage/apiServer/heartbeat"
	"github.com/fengjun2016/gDStorage/apiServer/locate"
	"github.com/fengjun2016/gDStorage/apiServer/objects"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
