package client

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bbruun/galt/common"
	"github.com/ryanuber/columnize"
)

type Client struct {
	fs *flag.FlagSet

	name string

	log    *log.Logger
	server string
	port   int
	debug  bool
}

func (c *Client) Name() string {
	return c.fs.Name()
}

func (c *Client) Init(args []string) error {
	return c.fs.Parse(args)
}

func NewClient() *Client {
	gc := &Client{
		fs: flag.NewFlagSet("client", flag.ExitOnError),
	}
	log := common.ReturnLogger()
	log.Debugf("parsing parameters")
	gc.fs.StringVar(&gc.server, "server", "127.0.0.1", "Server IP - default 127.0.0.1")
	gc.fs.StringVar(&gc.server, "s", "0.0.0.0", "Server IP - default 127.0.0.1")
	gc.fs.IntVar(&gc.port, "port", 8080, "Server port to connect to - default 8080")
	gc.fs.IntVar(&gc.port, "p", 8080, "Server port to connect to - default 8080")
	gc.fs.BoolVar(&gc.debug, "debug", false, "Enable debug messages")
	gc.fs.BoolVar(&gc.debug, "d", false, "Enable debug messages")

	return gc

}

func (c *Client) Usage() {
	var out []string
	out = append(out, "Galt Clent syntax:")
	out = append(out, " --server <ipAddress> | (optional) | The server IP to connect to - default 0.0.0.0")
	out = append(out, " -s <ipAddress> | |")
	out = append(out, "")
	out = append(out, " --port <port | (optional) | The port the server runs on - default 8080")
	out = append(out, " -p <port> | |")
	out = append(out, "")
	out = append(out, " --debug | (optional) | Enable debug")
	out = append(out, " -d | |")
	out = append(out, "")

	fmt.Print(columnize.SimpleFormat(out))
	os.Exit(0)
}

func (c *Client) Run() error {
	log := common.ReturnLogger()
	log.Debug("client started")
	return nil
}
