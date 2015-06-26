package p4

import (
	"bytes"
	"strconv"
)

type Review struct {
	Change int
	User   string
	Email  string
	Name   string
}

func (c *Connection) Counters() (map[string]string, error) {
	counters := map[string]string{}

	if data, err := c.execP4("-ztag", "counters"); err == nil {
		var b bytes.Buffer
		b.Write(data)

		if zcounters, err := ParseZTag(&b); err == nil {
			for _, zcounter := range zcounters {
				counters[zcounter["counter"]] = zcounter["value"]
			}

			return counters, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (c *Connection) ReviewByCounter(counter string) ([]Review, error) {
	reviews := []Review{}

	if data, err := c.execP4("-ztag", "review", "-t", counter); err == nil {
		var b bytes.Buffer
		b.Write(data)

		if zreviews, err := ParseZTag(&b); err == nil {
			for _, zreview := range zreviews {
				intChange, err := strconv.Atoi(zreview["change"])

				if err != nil {
					return nil, err
				}

				review := Review{
					Change: intChange,
					User:   zreview["user"],
					Email:  zreview["email"],
					Name:   zreview["name"],
				}

				reviews = append(reviews, review)
			}

			return reviews, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
