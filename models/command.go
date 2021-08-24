package models

const (
	CMD_NICK CommandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type CommandID int

type Command struct {
	Id     CommandID
	Client *Client
	Args   []string
}
