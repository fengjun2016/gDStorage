package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fengjun2016/gDStorage/apiServer/heartbeat"
	"github.com/fengjun2016/gDStorage/apiServer/locate"
	"github.com/fengjun2016/gDStorage/apiServer/objectstream"
	"github.com/sirupsen/logrus"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}

	if m == http.MethodGet {
		get(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(rw http.ResponseWriter, req *http.Request) {
	object := strings.Split(req.URL.EscapedPath(), "/")[2]
	c, err := storeObject(req.Body, object)
	if err != nil {
		log.Println(err.Error())
	}
	rw.WriteHeader(c)
}

func storeObject(r io.Reader, object string) (int, error) {
	stream, err := putStream(object)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	io.Copy(stream, r)
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServers")
	}

	return objectstream.NewPutStream(server, object), nil
}

func get(rw http.ResponseWriter, req *http.Request) {
	object := strings.Split(req.URL.EscapedPath(), "/")[2]
	stream, err := getStream(object)
	if err != nil {
		logrus.Println(err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(rw, stream)
}

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}

	return objectstream.NewGetStream(server, object)
}
