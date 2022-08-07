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
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
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

	c.readInput()
}

func (s *server) setName(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("setting name to: %s", c.name))
}

func (s *server) joinRoom(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	// create new room if does not already exist
	if !ok {
		r = &room{
			name: roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	
	r.members[c.conn.RemoteAddr()] = c

	// remove client from current room
	s.removeClientFromRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.name))
	c.msg(fmt.Sprintf("welcome to %s", r.name))
}

func (s *server) listRooms(c *client, args []string) {
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

	c.room.broadcast(c, c.name + ": " + strings.Join(args[1:], " "))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())

	s.removeClientFromRoom(c)

	c.msg("goodbye, come back soon!")
	c.conn.Close()
}

func (s *server) removeClientFromRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}