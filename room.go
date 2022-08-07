package main

import "net"

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
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
