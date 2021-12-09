package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func Intersphinx(url string) RefMap {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Errorf("Error creating request to url %s: %v", url, err)
		return nil
	}
	resp, err := Client.Do(req)
	if err != nil {
		log.Errorf("Error getting response from url %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil
	}

	markerLine := "# The remainder of this file is compressed using zlib.\n"
	cut := bytes.Index(buff, []byte(markerLine)) + len(markerLine)
	if cut < len(markerLine) {
		log.Warn("no marker line found in inv file header for url: ", url)
		return nil
	}

	b := bytes.NewReader(buff[cut:])
	if b.Size() == 0 {
		log.Errorf("no data found in file from url: %s", url)
		return nil
	}

	r, err := zlib.NewReader(b)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil
	}
	defer r.Close()

	parsed, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil
	}

	refMap := make(map[string]string)

	for _, line := range strings.Split(string(parsed), "\n") {
		if len(line) == 0 {
			continue
		}
		lineSplit := strings.Split(line, " ")
		refMap[lineSplit[0]] = url[:len(url)-len("objects.inv")] + lineSplit[3]
	}

	return refMap
}
