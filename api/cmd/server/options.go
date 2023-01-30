package server

type Options struct {
	IP   string `long:"ip" description:"ip address to listen on"`
	Port string `short:"p" long:"port" description:"port number to listen on" default:"8080"`
}
