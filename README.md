# GoShare
GoShare is a terminal tool for sharing files through one or multiple servers.  
Servers are acting as a relay to stream data from sender to receiver. 

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

### Sharing files
To share a file, use the following command:
```
goshare share -f <path/to/file> -s server1:2022 -s server2:2030 
```
The above command gives you a key code that needs to be shared with receiver in order to download the file. 

### Receiving files
Use the following command to receive the file:

```
goshare receive -k KEY_CODE_FROM_SENDER
```

### Server Configuration
To configure the server, you can provide the following options:

```
goshare cert --host localhost --dst ./cert/

goshare server --port 2202 --ip localhost --cert-path ./cert/
```

## Todo
- Make TLS optional
- STUN

## License
This project is licensed under the MIT License - see the LICENSE file for details.
