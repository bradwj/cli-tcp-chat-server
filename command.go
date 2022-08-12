package main

type commandID int

const (
	CMD_NAME commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_USERS
	CMD_LEAVE
	CMD_QUIT
	CMD_HELP
)

type command struct {
	id     commandID
	client *client
	args   []string
}
