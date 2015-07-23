package p4

import (
	"bytes"
	"os"
	"os/exec"
)

type Connection struct {
	port     string
	username string
	password string
	client   string
}

func Connect(port string, username string, password string, client string) *Connection {
	c := new(Connection)
	c.port = port
	c.username = username
	c.password = password
	c.client = client

	return c
}

func (c *Connection) execP4(args ...string) ([]byte, error) {
	env := []string{
		"HOME=" + os.Getenv("HOME"),
		"P4CLIENT=" + c.client,
		"P4PORT=" + c.port,
		"P4USER=" + c.username,
	}

	var b bytes.Buffer
	b.Write([]byte(c.password))

	cmd := exec.Command("p4", "login")
	cmd.Env = env
	cmd.Stdin = &b

	if data, err := cmd.CombinedOutput(); err == nil {
		cmd := exec.Command("p4", args...)
		cmd.Env = env

		if data, err := cmd.CombinedOutput(); err == nil {
			return data, nil
		} else {
			return nil, P4Error{err, append([]string{"p4"}, args...), data}
		}
	} else {
		return nil, P4Error{err, []string{"p4", "login"}, data}
	}
}
