package objects

import (
	"io"
	"net/http"
	"os"
	"strings"

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

func put(w http.ResponseWriter, r *http.Request) {
	f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		logrus.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}

// func GetHashFromHeader(h http.Header) string {
// 	digest := h.Get("digest")
// 	if len(digest) < 9 {
// 		return ""
// 	}

// 	if digest[:8] != "SHA-256=" {
// 		return ""
// 	}

// 	return digest[8:]
// }

// func GetSizeFromHeader(h http.Header) int64 {
// 	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
// 	return size
// }

// func storeObject(r io.Reader, object string) (int, error) {
// 	stream, e := putStream(object)
// 	if e != nil {
// 		return http.StatusServiceUnvaliable, e
// 	}

// 	io.Copy(stream, r)
// 	e = stream.Close()
// 	if e != nil {
// 		return http.StatusInternalServerError, e
// 	}
// 	return http.StatusOK, nil
// }

// func putStream(object string) (*objectstream.PutStream, error) {
// 	server := heartbeat.ChooseRandomDataServer()
// 	if server == "" {
// 		return nil, fmt.Errorf("cannot find any dataServer")
// 	}
// 	return objectstream.NewPutStream(server, object), nil
// }

func get(rw http.ResponseWriter, req *http.Request) {
	f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(req.URL.EscapedPath(), "/")[2])
	if err != nil {
		logrus.Println(err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(rw, f)
}

// func getStream(object string) (io.Reader, error) {
// 	server := locate.Locate(object)
// 	if server == "" {
// 		return nil, fmt.Errorf("object %s locate fail", object)
// 	}

// 	return objectstream.NewGetStream(Server, object)

// }

// func del(w http.ResponseWriter, r *http.Request) {
// 	name := strings.Split(r.URL.EscapedPath(), "/")[2]
// 	version, e := es.SearchLatestVersion(name)
// 	if e != nil {
// 		logrus.Println(e)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	e = es.PutMetadata(name, version.Version+1, 0, "")
// 	if e != nil {
// 		logrus.Println(e)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// }
