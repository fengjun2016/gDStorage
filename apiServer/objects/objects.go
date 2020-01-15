package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"networkStorage/pkg/es"
	"strings"

	"github.com/fengjun2016/gDStorage/apiServer/heartbeat"
	"github.com/fengjun2016/gDStorage/apiServer/locate"
	"github.com/fengjun2016/gDStorage/apiServer/objectstream"
	"github.com/fengjun2016/gDStorage/apiServer/pkg/utils"
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

	if m == http.MethodDelete {
		del(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(rw http.ResponseWriter, req *http.Request) {
	// object := strings.Split(req.URL.EscapedPath(), "/")[2]
	// c, err := storeObject(req.Body, object)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// rw.WriteHeader(c)
	hash := utils.GetHashFromHeader(req.Header)
	if hash == "" {
		logrus.Println("missing object hash in digest header")
		rw.WriteHeader(http.StatusRequest)
		return
	}

	c, err := storeObject(req.Body, url.PathEscape(hash))
	if err != nil {
		lgorus.Println(err)
		rw.WriteHeader(c)
		return
	}

	if c != http.StatusOK {
		rw.WriteHeader(c)
		return
	}

	name := strings.Split(req.URL.EscapedPath(), "/")[2]
	size := utils.GetSizeFromHeader(req.Header)
	err = es.AddVersion(name, hash, size)
	if err != nil {
		logrus.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
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
	// object := strings.Split(req.URL.EscapedPath(), "/")[2]
	// stream, err := getStream(object)
	// if err != nil {
	// 	logrus.Println(err.Error())
	// 	rw.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// io.Copy(rw, stream)
	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	versionId := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			logrus.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	meta, err := es.GetMetadata(name, version)
	if err != nil {
		logrus.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if meta.Hash == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	object := url.PathEscape(meta.Hash)
	stream, err := getStream(object)
	if err != nil {
		logrus.Println(err)
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

func del(rw http.ResponseWriter, req *http.Request) {
	name := strings.Split(req.URL.EscapedPath(), "/")[2]
	version, err := es.SearchLatestVersion(name)
	if err != nil {
		logrus.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = es.PutMetadata(name, version.Version+1, 0, "") // size 为0 hash为空字符串 表示这是一个删除标记
	if err != nil {
		logrus.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
