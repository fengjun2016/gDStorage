package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

func getMetadata(name string, versionId int) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		os.Getenv("ES_SERVER"), name, versionId)
	r, err := http.Get(url)
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}

	return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash": "%s"}`, name, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	r, err := client.Do(request)
	if err != nil {
		return err
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}

	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatuCode, string(result))
	}

	return nil
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_serach?q=name:%s&size=1&sort=version:desc", os.Getenv("ES_SERVER"), url.PathEscape(name))
	r, err := http.Get(url)
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
}
