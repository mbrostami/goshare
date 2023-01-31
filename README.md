# GoShare
GoShare is a terminal tool for files. It allows users to authenticate using their ssh public keys and share files with other users through a server. 
The server can be configured as a simple relay to stream data with encryption.

## Features
Authentication using ssh public keys
File sharing through a relay server

### Todo
Data stream encryption (optional)
Persist data on the server for x amount of time (optional)
Share link (optional)

## Installation
Make sure you have Go installed on your system.  
Clone the repository  
```
git clone https://github.com/mbrostami/goshare.git
```
Build the binary  
```
go build -o goshare
```
Run the binary  
```
./goshare
```

## Usage
### Registration 
To register username in a server, you will need to provide your username and the path to your ssh public key.
```
goshare register --server <address>:<port> --username <username> --key <path/to/ssh/public/key>
```

### Sharing files
To share a file, use the following command:
```
goshare share --file <path/to/file> --receiver <username>
```
### Receiving files
To check if there are any files available for you to receive, use the following command:

```
goshare receive
```
### Server Configuration
To configure the server, you can provide the following options:

```
goshare server
--port : listening port 
--persist: to configure the server to persist the incoming data
--encrypt : to configure the server to encrypt the streaming data
--no-auth : to configure the server to be used without registering users
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
This tool is built with crypto/ssh package.
