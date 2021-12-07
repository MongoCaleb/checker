package intersphinx

// The Go standard library comes with excellent support
// for HTTP clients and servers in the `net/http`
// package. In this example we'll use it to issue simple
// HTTP requests.

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Intersphinx() {

	// Issue an HTTP GET request to a server. `http.Get` is a
	// convenient shortcut around creating an `http.Client`
	// object and calling its `Get` method; it uses the
	// `http.DefaultClient` object which has useful default
	// settings.
	resp, err := http.Get("https://docs.mongodb.com/drivers/go/current/objects.inv")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	target := "compressed using zlib.\n"
	cut := bytes.Index(buff, []byte(target)) + len(target)
	b := bytes.NewReader(buff[cut:])

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	// io.Copy(os.Stdout, r)

	parsed, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	refMap := make(map[string]string)

	for _, line := range strings.Split(string(parsed), "\n") {
		if len(line) == 0 {
			continue
		}
		lineSplit := strings.Split(line, " ")
		refMap[lineSplit[0]] = lineSplit[3]
	}

	fmt.Println(refMap)

	r.Close()
}
