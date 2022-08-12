package main

import (
	"net"
	"time"
)

type room struct {
	name        string
	members     map[net.Addr]*client
	timeCreated time.Time
}

func newRoom(name string) *room {
	return &room{
		name: name,
		members: make(map[net.Addr]*client),
		timeCreated: time.Now(),
	}
}

func (r *room) broadcast(msg string) {
	for _, member := range r.members {
		member.msg(msg)
	}
}

func (r *room) broadcastToOthers(sender *client, msg string) {
	for addr, member := range r.members {
		if addr != sender.conn.RemoteAddr() {
			member.msg(msg)
		}
	}
}
