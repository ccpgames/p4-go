package p4

import (
	"fmt"
	"regexp"
	"strconv"
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

var countersRegexp = regexp.MustCompile("(?m)^(.+) = (.+)$")
var describeRegexp = regexp.MustCompile("\\AChange (\\d+) by (.+)@(.+) on (.+)((?: *pending*)?)\n\n((?:\t.*\n)*)\nAffected files ...\n\n((?:... (?:.+) (?:[\\w/]+)\n)*)\n\\z")
var describeAffectedRegexp = regexp.MustCompile("(?m)^... (.+)#(\\d+) ([\\w/]+)$")
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

func (c *Connection) ReviewByCounter(counter string) ([]Review, error) {
	reviews := []Review{}

	if data, err := c.execP4("review", "-t", counter); err == nil {
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
	if _, err := c.execP4("sync", "-s", fmt.Sprintf("%s@%d", path, clNumber)); err == nil {
		return nil
	} else {
		return err
	}
}
