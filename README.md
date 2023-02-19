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
goshare server --port 2202 --ip localhost
```

### TLS enabled

#### Generate Certificate
Skip this step if you already have certificate.  
You can use `cert` command to generate self-signed key.pem and cert.pem files.  
```
goshare cert --host localhost --dst ./cert/
```
The above command will create two files in `./cert/` directory.  

#### Enable TLS on server:  
```
goshare server --port 2202 --ip localhost --with-tls --cert-path ./cert/
```

#### Enable TLS on clients (sender and receiver):  
``` 
goshare receive -k KEY_CODE_FROM_SENDER --with-tls 
goshare share -f <path/to/file> -s server1:2022 -s server2:2030 --with-tls
```

By default the above command will use system's root CA set. In other cases use `ca-path` as below:       
Use `ca-path` to specify the CA's certificate that has been used to sign the server's self-signed certificate.  
``` 
goshare receive -k KEY_CODE_FROM_SENDER --with-tls --ca-path ./cert/cert.pem 
goshare share -f <path/to/file> -s server1:2022 -s server2:2030 --with-tls --ca-path ./cert/cert.pem
```

To skip the certificate verification in clients, you can use `--skip-verify` option (NOT recommended)



## Todo
- STUN

## License
This project is licensed under the MIT License - see the LICENSE file for details.
