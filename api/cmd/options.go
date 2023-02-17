package cmd

type Options struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool          `short:"v" long:"verbose" description:"Show verbose debug information"`
	Server  *serverOptions  `command:"server"`
	Share   *shareOptions   `command:"share"`
	Receive *receiveOptions `command:"receive"`
	Cert    *certOptions    `command:"cert"`
}
