package main

type commandID int

const (
	CMD_NAME commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
