# GoShare
GoShare is a terminal tool for securely sharing files over SSH. It allows users to authenticate using their ssh public keys and share files with other users through a server that acts as a relay. The server does not store the data and can be configured to encrypt the data stream.

## Features
Authentication using ssh public keys
Secure file sharing through a relay server
Data stream encryption (optional)
Data storage on the server (optional)


## Installation
Make sure you have Go installed on your system. You can download it from here
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
To register in a server, you will need to provide your username and the path to your ssh public key.
```
goshare register --username <username> --server <address>:<port> --key <path/to/ssh/public/key>
```

### Authentication
To login to the server, you will need to provide your username and the path to your ssh public key.

```
goshare login --username <username> --key <path/to/ssh/public/key>
```
### Register a server
```
goshare register --server <address>:<port> --name <servername>

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
-port : listening port 
-persist: to configure the server to persist the incoming data
-encrypt : to configure the server to encrypt the data stream
-public : to configure the server to be used without registering users
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
This tool is built with crypto/ssh package.
