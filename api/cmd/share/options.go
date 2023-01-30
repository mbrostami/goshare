package share

type Options struct {
	File string `short:"f" long:"file" description:"file path you want to share" required:"true"`
}
