package intersphinx

import (
	"bytes"
	"checker/internal/parsers/rst"
	"compress/zlib"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SphinxMap map[string]rst.RefTarget

func Intersphinx(buff []byte, domain string) SphinxMap {

	markerLine := "# The remainder of this file is compressed using zlib.\n"
	cut := bytes.Index(buff, []byte(markerLine)) + len(markerLine)
	if cut < len(markerLine) {
		log.Warn("no marker line found in inv file header for intersphinx parsing")
		return nil
	}

	b := bytes.NewReader(buff[cut:])
	if b.Size() == 0 {
		log.Errorf("no data found in input for intersphinx parsing")
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

	res := make(map[string]rst.RefTarget)

	for _, line := range strings.Split(string(parsed), "\n") {
		if len(line) == 0 {
			continue
		}
		lineSplit := strings.Split(line, " ")
		res[lineSplit[0]] = rst.RefTarget{Raw: lineSplit[3], Target: domain + lineSplit[3] + "%s", Type: "intersphinx"}
	}
	return res
}

func JoinSphinxes(input []SphinxMap) SphinxMap {
	refMap := make(SphinxMap)
	for _, m := range input {
		for k, v := range m {
			refMap[k] = v
		}
	}
	return refMap
}
