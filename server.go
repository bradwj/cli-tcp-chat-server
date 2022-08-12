package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

const MAX_MESSAGE_LENGTH = 1000
const MAX_NAME_LENGTH = 20

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.setName(cmd.client, cmd.args)
		case CMD_JOIN:
			s.joinRoom(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.sendMessage(cmd.client, cmd.args)
		case CMD_USERS:
			s.listUsers(cmd.client, cmd.args)
		case CMD_LEAVE:
			s.leaveRoom(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		case CMD_HELP:
			s.help(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}

	c.msg(fmt.Sprintln("welcome to the CLI TCP chat server! for a list of commands, type \"/help\""))
	c.readInput()
}

func (s *server) setName(c *client, args []string) {
	if len(args) < 2 {
		c.err(errors.New("you must specify a name. e.g. \"/name brad\""))
		return
	}

	oldName := c.name
	newName := strings.Join(args[1:], " ")
	if len(newName) > MAX_NAME_LENGTH {
		c.err(errors.New(fmt.Sprintf("name is too long! (%d / %d maximum allowed characters)",
			len(newName),
			MAX_NAME_LENGTH,
		)))
		return
	}

	c.name = newName
	c.msg(fmt.Sprintf("changed name to \"%s\"", c.name))

	// broadcast to room when user changes name
	if c.room != nil {
		c.room.broadcast(fmt.Sprintf("user \"%s\" changed name to \"%s\"", oldName, c.name))
	}
}

func (s *server) joinRoom(c *client, args []string) {
	if len(args) < 2 {
		c.err(errors.New("you must specify a room name. e.g. \"/join groupchat\""))
		return
	}
	// remove client from current room
	s.removeClientFromRoom(c)

	roomName := args[1]
	r, ok := s.rooms[roomName]
	// create new room if does not already exist
	if !ok {
		r = newRoom(roomName)
		s.rooms[roomName] = r
		log.Printf("created room: %s", roomName)
	}

	r.members[c.conn.RemoteAddr()] = c
	c.room = r

	c.msg(fmt.Sprintf("welcome to %s", r.name))
	r.broadcast(fmt.Sprintf("%s has joined the room", c.name))
}

func (s *server) listRooms(c *client, args []string) {
	if len(s.rooms) == 0 {
		c.msg("there are no available rooms. type \"/join <room name>\" to create one!")
		return
	}

	var roomNames []string
	for name := range s.rooms {
		roomNames = append(roomNames, name)
	}

	c.msg(fmt.Sprintf("available rooms are: %s", strings.Join(roomNames, ", ")))
}

func (s *server) sendMessage(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join a room before sending a message"))
		return
	}

	message := strings.Join(args[1:], " ")

	// check if message is longer than max length allowed
	if len(message) > MAX_MESSAGE_LENGTH {
		c.err(errors.New(fmt.Sprintf("message is too long! (%d / %d maximum allowed characters)",
			len(message),
			MAX_MESSAGE_LENGTH,
		)))
		return
	}

	c.room.broadcast(c.name + ": " + message)
}

func (s *server) listUsers(c *client, args []string) {

	// check if user is in a room
	if c.room == nil {
		c.err(errors.New(fmt.Sprint("you must be in a room to list the users")))
		return
	}

	var memberNames []string
	for addr, member := range c.room.members {
		name := member.name
		if addr == c.conn.RemoteAddr() {
			name += " (you)"
		}
		memberNames = append(memberNames, name)
	}

	c.msg(fmt.Sprintf("users in this room: %s", strings.Join(memberNames, ", ")))
}

func (s *server) leaveRoom(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you are not currently in a room"))
		return
	}
	roomName := c.room.name
	s.removeClientFromRoom(c)
	c.msg(fmt.Sprintf("you have left %s", roomName))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())

	s.removeClientFromRoom(c)
	c.msg("goodbye, come back soon!")
	c.conn.Close()
}

func (s *server) removeClientFromRoom(c *client) {
	if c.room != nil {
		// remove client from room member list
		delete(c.room.members, c.conn.RemoteAddr())

		// delete room if no members left
		if len(c.room.members) == 0 {
			delete(s.rooms, c.room.name)
			log.Printf("deleted room: %s", c.room.name)
		} else {
			c.room.broadcast(fmt.Sprintf("%s has left the room", c.name))
		}
		c.room = nil
	}
}

func (s *server) help(c *client, args []string) {
	helpMsg := `available commands:
	> "/name <name>" -- Set your username. Otherwise, you will remain anonymous.
	> "/join <room name>" -- Join a chat room. If the room doesn't exist, a new one will be created.
	> "/rooms" -- Show list of available rooms to join.
	> "/msg <message>" -- Broadcast message to everyone in current room.
	> "/users" -- List the users that are in the current room.
	> "/leave" -- Leave the current room.
	> "/quit" -- Disconnect from the chat server.
	> "/help" -- List available commands.`

	c.msg(helpMsg)
}
