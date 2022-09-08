# tcp-file-server
It is a project of a TPC custom protocol to send files between two or more clients in specific channels.

## 1. TCP custom protocol: 
Protocol based on commands, the client would send the permitted commands to subscribe to specific channels, look at the available channels, send a file or quit. 

## 2. Server
The server handles the client's request. It listens to a TCP port and is designed thinking on the TCP custom protocol to understand the protocol commands.

## 3. Client
The client is a go terminal program, it allows you to connect to the sever, subscribe to a specific channel, and send or receive files.

### Client input commands
- `/subscribe` subscribe a channel, if channel not exist it will be created.
- `/channels` list the available channels.
- `/send <filepath>` send a file to everyone in the channel.
- `/quit` disconnects from the file server.

### Client example 
![image](https://user-images.githubusercontent.com/69649613/189131921-95eae646-feca-460a-8447-bd9a01f9600b.png)


