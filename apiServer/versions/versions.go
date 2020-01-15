package versions

func Handler(rw http.ResponseWriter, req *http.Request) {
	m := req.Method
	if m != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	from := 0
	size := 1000

	name := strings.Split(req.URL.EscapedPath(), "/")[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size)
		if err != nil {
			logrus.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i := range metas {
			b, _ := json.Marshal(metas[i])
			rw.Write(b)
			rw.Write([]byte("\n"))
		}

		if len(metas) != size {
			return
		}

		from += size
	}
}
