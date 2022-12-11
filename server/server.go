package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/ryanuber/columnize"
)

type Server struct {
	fs *flag.FlagSet

	name string

	log    *log.Logger
	listen string
	port   int
	debug  bool
}

func (s *Server) Name() string {
	return s.fs.Name()
}

func (s *Server) Init(args []string) error {
	return s.fs.Parse(args)
}

func NewServer() *Server {
	gc := &Server{
		fs: flag.NewFlagSet("server", flag.ExitOnError),
	}

	gc.fs.StringVar(&gc.listen, "listen", "0.0.0.0", "Listen IP address - default 0.0.0.0")
	gc.fs.StringVar(&gc.listen, "l", "0.0.0.0", "Listen IP address - default 0.0.0.0")
	gc.fs.IntVar(&gc.port, "port", 8080, "Port to listen on - default 8080")
	gc.fs.IntVar(&gc.port, "p", 8080, "Port to listen on - default 8080")

	gc.fs.BoolVar(&gc.debug, "debug", false, "Enable server debug messages")
	gc.fs.BoolVar(&gc.debug, "d", false, "Enable server debug messages")

	return gc

}

func (s *Server) Usage() {
	var out []string
	out = append(out, "Galt Server syntax:")
	out = append(out, " --listen <ipAddress> | (optional) | IP to bind service to - default 0.0.0.0")
	out = append(out, " -l <ipAddress> | |")
	out = append(out, "")
	out = append(out, " --port <port | (optional) | Run service on port - default 8080")
	out = append(out, " -p <port> | |")
	out = append(out, "")
	out = append(out, " --debug | (optional) | Enable debug")
	out = append(out, " -d | |")
	out = append(out, "")

	fmt.Print(columnize.SimpleFormat(out))
	os.Exit(0)
}

func (s *Server) Run() error {

	log.SetHandler(json.New(os.Stdout))
	log.Debug("debug message")
	return nil
}
