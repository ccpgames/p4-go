package p4

import (
	"bytes"
	"fmt"
)

type P4Error struct {
	Status error
	Output []byte
}

func (err P4Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Status, string(bytes.TrimRight(err.Output, "\n")))
}
