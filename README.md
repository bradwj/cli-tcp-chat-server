# CLI TCP Chat Server
A command-line interface where you can send and receive messages through a TCP chat server built with Go.

# Commands
- `/name <name>` -- Set your username. Otherwise, you will remain anonymous.
- `/join <room name>` -- Join a chat room. If the room doesn't exist, a new one will be created.
- `/rooms` -- Show list of available rooms to join.
- `/msg <message>` -- Broadcast message to everyone in current room.
- `/quit` -- Disconnect from the chat server.
