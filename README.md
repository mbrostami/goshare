# GoShare
GoShare is a terminal tool for sharing files through one or multiple servers.  
Servers are acting as a relay to stream data from sender to receiver. 

![screenshot](https://github.com/mbrostami/goshare/blob/master/screenshot.png?raw=true)

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
The above command gives you a key code. Receiver needs this code to download the file. 

### Receiving files
To check if there are any files available for you to receive, use the following command:

```
goshare receive -k KEY_CODE_FROM_SENDER
```

### Server Configuration
To configure the server, you can provide the following options:

```
goshare server
--port : listening port 
--ip   : listening ip
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.
