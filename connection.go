package p4

import (
	"bytes"
	"os"
	"os/exec"
)

type Connection struct {
	host     string
	username string
	password string
}

func Connect(host string, username string, password string) *Connection {
	c := new(Connection)
	c.host = host
	c.username = username
	c.password = password

	return c
}

func (c *Connection) execP4(args ...string) ([]byte, error) {
	env := []string{
		"HOME=" + os.Getenv("HOME"),
		"P4HOST=" + c.host,
		"P4USER=" + c.username,
	}

	var b bytes.Buffer
	b.Write([]byte(c.password))

	cmd := exec.Command("p4", "login")
	cmd.Env = env
	cmd.Stdin = &b

	if err := cmd.Run(); err == nil {
		cmd := exec.Command("p4", args...)
		cmd.Env = env

		return cmd.CombinedOutput()
	} else {
		return nil, err
	}
}
