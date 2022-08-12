# CLI TCP Chat Server
A command-line interface where you can send and receive messages through a TCP chat server built with Go.

# How to Run
Clone the repo and run the server
```bash
$ git clone github.com/bradwj/cli-tcp-chat-server
$ cd cli-tcp-chat-server
$ go run .
```
In a separate terminal window, connect to the server using `telnet`
```bash
$ telnet localhost 8888
Trying ::1....
Connected to localhost.
```
Type commands in the telnet CLI
```bash
/name brad
> changed name to "brad"

/join general
> welcome to general
> brad has joined the room

/rooms
> available rooms are: general

/msg hello there!
> brad: hello there!

/users
> users in this room: brad (you)

/leave
> you have left general

/quit
> goodbye, come back soon!

Connection to host lost.
```

# Commands
- `/name <name>` -- Set your username. Otherwise, you will remain anonymous.
- `/join <room name>` -- Join a chat room. If the room doesn't exist, a new one will be created.
- `/rooms` -- Show list of available rooms to join.
- `/desc <room description>` -- Set the description of the current room.
- `/info [room name]` -- Show the information of either the current room, or [room name] if it is specified.
- `/msg <message>` -- Broadcast message to everyone in current room.
- `/users` -- List the users that are in the current room.
- `/leave` -- Leave the current room.
- `/quit` -- Disconnect from the chat server.
- `/help` -- List available commands.
