package p4

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
)

var tokenRegexp = regexp.MustCompile("([0-9A-Z]{32})")

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
		"P4PORT=" + c.port,
		"P4CLIENT=" + c.client,
		"P4USER=" + c.username,
	}

	var password bytes.Buffer
	var token bytes.Buffer
	var errors bytes.Buffer

	password.Write([]byte(c.password))

	cmd := exec.Command("p4", "login", "-p")
	cmd.Env = env
	cmd.Stdin = &password
	cmd.Stdout = &token
	cmd.Stderr = &errors

	if err := cmd.Run(); err == nil {
		env = append(env, "P4PASSWD="+tokenRegexp.FindString(token.String()))

		cmd := exec.Command("p4", args...)
		cmd.Env = env

		if data, err := cmd.CombinedOutput(); err == nil {
			return data, nil
		} else {
			return nil, P4Error{err, append([]string{"p4"}, args...), data}
		}
	} else {
		return nil, P4Error{err, []string{"p4", "login"}, errors.Bytes()}
	}
}
