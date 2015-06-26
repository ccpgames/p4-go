package p4

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

var ztagRegexp = regexp.MustCompile("^... ([a-zA-Z0-9]+) (.*)$")

type ZTag []map[string]string

func ParseZTag(r io.Reader) ([]map[string]string, error) {
	var a = []map[string]string{}
	var m = map[string]string{}
	var lastTag string

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(m) > 0 {
				a = append(a, m)
			}

			m = map[string]string{}
		} else {
			if match := ztagRegexp.FindStringSubmatch(line); match != nil {
				lastTag = match[1]
				m[match[1]] = match[2]
			} else if lastTag != "" && len(m) > 0 {
				m[lastTag] += "\n" + line
			} else {
				return nil, errors.New("ztag: parse error")
			}
		}
	}

	if len(m) > 0 {
		a = append(a, m)
	}

	return a, scanner.Err()
}
