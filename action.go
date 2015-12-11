package p4

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Describe struct {
	Change      int
	User        string
	Client      string
	Time        string
	Description string
	Status      string
	Files       []DescribeFile
}

type DescribeFile struct {
	Path    string
	Version int
	Action  string
}

type Review struct {
	Change int
	User   string
	Email  string
	Name   string
}

var newlineRegexp = regexp.MustCompile("\r\n|\r|\n")

var countersRegexp = regexp.MustCompile("(?m)^(.+) = (.+)$")
var describeRegexp = regexp.MustCompile("\\AChange (\\d+) by (.+)@(.+) on (.+)((?: *pending*)?)\n\n((?:\t.*\n)*)\nAffected files ...\n\n((?:... (?:.+) (?:[\\w/]+)\n)*)\n\\z")
var describeAffectedRegexp = regexp.MustCompile("(?m)^... (.+)#(\\d+) ([\\w/]+)$")
var printRegexp = regexp.MustCompile("(?m)\\A(.+)(@|#)(\\d+)(?: - | )(.+)$")
var reviewRegexp = regexp.MustCompile("(?m)^Change (\\d+) (.+) <(.+)> \\((.+)\\)$")

func (c *Connection) Counters() (map[string]string, error) {
	counters := map[string]string{}

	if data, err := c.execP4("counters"); err == nil {
		submatch := countersRegexp.FindAllSubmatch(data, 1000000)

		for _, counter := range submatch {
			counters[string(counter[1])] = string(counter[2])
		}

		return counters, nil
	} else {
		return nil, err
	}
}

func (c *Connection) Describe(change int) (Describe, error) {
	var describe Describe

	if data, err := c.execP4("describe", "-s", strconv.Itoa(change)); err == nil {
		submatch := describeRegexp.FindSubmatch(data)
		intChange, err := strconv.Atoi(string(submatch[1]))

		if err != nil {
			return describe, err
		}

		status := "submitted"

		if string(submatch[5]) == " *pending" {
			status = "pending"
		}

		describe = Describe{
			Change:      intChange,
			User:        string(submatch[2]),
			Client:      string(submatch[3]),
			Time:        string(submatch[4]),
			Description: string(submatch[6]),
			Status:      status,
		}

		affectedSubmatch := describeAffectedRegexp.FindAllSubmatch(submatch[7], 10000000)

		for _, m := range affectedSubmatch {
			intVersion, err := strconv.Atoi(string(m[2]))

			if err != nil {
				return describe, err
			}

			describe.Files = append(describe.Files, DescribeFile{
				Path:    string(m[1]),
				Version: intVersion,
				Action:  string(m[3]),
			})
		}

		return describe, nil
	} else {
		return describe, err
	}
}

func (c *Connection) GetCounter(counter string) (string, error) {
	if data, err := c.execP4("counter", counter); err == nil {
		return strings.TrimRight(string(data), "\n"), nil
	} else {
		return "", err
	}
}

func (c *Connection) Print(path string, clNumber int) ([]byte, error) {
	// We can not use p4's -q flag here, as that leaves us with no method of
	// distinguishing an actual error from a file happens to contain an error
	// message. The process exits with a status code of 0 in both cases.
	//
	// The first line of output differs slightly between a successful request
	// and an error:
	//
	//  * On error, the line begins with the path followed by `@` and the
	//    changelist number.
	//  * On success, the line begins with the path followed by `#` and the
	//    file revision.
	//
	// Another limitation of p4's print is automatic line-ending conversions
	// on text files. This can not be disabled. It also can not be determined
	// if the file is treated as text or binary by Perforce.
	//
	// No attempt is made to correct this anomaly. Whatever p4 gives us, we
	// give you.

	url := fmt.Sprintf("%s@%d", path, clNumber)

	if data, err := c.execP4("print", url); err == nil {
		lines := newlineRegexp.Split(string(data), 2)

		if len(lines) != 2 {
			return nil, errors.New("no newlines found in p4's output")
		}

		submatch := printRegexp.FindSubmatch([]byte(lines[0]))

		if len(submatch) == 0 {
			return nil, errors.New("first line from p4 print of invalid format")
		}

		if submatch[2][0] != '#' {
			return nil, P4Error{
				errors.New(string(submatch[4])),
				[]string{"p4", "print", url},
				data,
			}
		}

		return []byte(lines[1]), nil
	} else {
		return nil, err
	}
}

func (c *Connection) ReviewByChangelist(clNumber int) ([]Review, error) {
	return c.review("review", "-c", strconv.Itoa(clNumber))
}

func (c *Connection) ReviewByCounter(counter string) ([]Review, error) {
	return c.review("review", "-t", counter)
}

func (c *Connection) review(arguments ...string) ([]Review, error) {
	reviews := []Review{}

	if data, err := c.execP4(arguments...); err == nil {
		submatch := reviewRegexp.FindAllSubmatch(data, 10000000)

		for _, review := range submatch {
			intChange, err := strconv.Atoi(string(review[1]))

			if err != nil {
				return nil, err
			}

			reviews = append(reviews, Review{
				Change: intChange,
				User:   string(review[2]),
				Email:  string(review[3]),
				Name:   string(review[4]),
			})
		}

		return reviews, nil
	} else {
		return nil, err
	}
}

func (c *Connection) SetCounter(counter string, value string) error {
	if _, err := c.execP4("counter", counter, value); err == nil {
		return nil
	} else {
		return err
	}
}

func (c *Connection) Sync(path string, clNumber int) error {
	if _, err := c.execP4("sync", "-f", fmt.Sprintf("%s@%d", path, clNumber)); err == nil {
		return nil
	} else {
		return err
	}
}
