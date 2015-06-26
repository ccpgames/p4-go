package p4

import (
	"bytes"
	"reflect"
	"testing"
)

func testZTag(t *testing.T, input []byte, expect []map[string]string) {
	var b bytes.Buffer

	b.Write(input)

	if m, err := ParseZTag(&b); err == nil {
		if !reflect.DeepEqual(m, expect) {
			t.Error("unexpected parse result, expected:", expect, "got", m)
		}
	} else {
		t.Error(err)
	}
}

func TestEmpty(t *testing.T) {
	testZTag(t, []byte(""), []map[string]string{})
}

func TestSingleEntry(t *testing.T) {
	testZTag(t, []byte("... key Value"), []map[string]string{
		map[string]string{"key": "Value"},
	})
}

func TestSingleEntryMultipleKeys(t *testing.T) {
	testZTag(t, []byte("... one 1\n... two 2\n"), []map[string]string{
		map[string]string{
			"one": "1",
			"two": "2",
		},
	})
}

func TestMultipleEntries(t *testing.T) {
	testZTag(t, []byte("... a b\n\n... b c\n... c d\n\n... d e"), []map[string]string{
		map[string]string{"a": "b"},
		map[string]string{"b": "c", "c": "d"},
		map[string]string{"d": "e"},
	})
}

func TestMultilineEntry(t *testing.T) {
	testZTag(t, []byte("... key val\nue\n\n... 2nd 3rd\n"), []map[string]string{
		map[string]string{"key": "val\nue"},
		map[string]string{"2nd": "3rd"},
	})
}
