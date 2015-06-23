package p4

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
