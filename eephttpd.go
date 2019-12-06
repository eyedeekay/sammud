package sammud

import (
	"log"

	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
    "github.com/zrma/mud/logging"
    "github.com/eyedeekay/mud/server"
    "strconv"
)

//SAMMud is a structure which automatically configured the forwarding of
//a local service to i2p over the SAM API.
type SAMMud struct {
	*samforwarder.SAMForwarder
    *server.Server
    logging.Logger
	up       bool
}

var err error

func (f *SAMMud) GetType() string {
	return "sammud"
}

func (f *SAMMud) ServeParent() {
	log.Println("Starting eepsite server", f.Base32())
	if err = f.SAMForwarder.Serve(); err != nil {
		f.Cleanup()
	}
}

//Serve starts the SAM connection and and forwards the local host:port to i2p
func (f *SAMMud) Serve() error {
	go f.ServeParent()
	if f.Up() {
		log.Println("Starting web server", f.Target())
        f.Server.Run()
	}
	return nil
}

func (f *SAMMud) Up() bool {
	return f.up
}

//Close shuts the whole thing down.
func (f *SAMMud) Close() error {
	return f.SAMForwarder.Close()
}

func (s *SAMMud) Load() (samtunnel.SAMTunnel, error) {
	if !s.up {
		log.Println("Started putting tunnel up")
	}
	f, e := s.SAMForwarder.Load()
	if e != nil {
		return nil, e
	}
	s.SAMForwarder = f.(*samforwarder.SAMForwarder)
    port, _ := strconv.Atoi(s.SAMForwarder.Config().TargetPort)
    s.Server = server.New(s.Logger, s.SAMForwarder.Config().TargetHost, port)
	s.up = true
	log.Println("Finished putting tunnel up")
	return s, nil
}

//NewSAMMud makes a new SAM forwarder with default options, accepts host:port arguments
func NewSAMMud(host, port string) (*SAMMud, error) {
	return NewSAMMudFromOptions(SetHost(host), SetPort(port))
}

//NewSAMMudFromOptions makes a new SAM forwarder with default options, accepts host:port arguments
func NewSAMMudFromOptions(opts ...func(*SAMMud) error) (*SAMMud, error) {
	var s SAMMud
	s.SAMForwarder = &samforwarder.SAMForwarder{}
	log.Println("Initializing eephttpd")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	l, e := s.Load()
	//log.Println("Options loaded", s.Print())
	if e != nil {
		return nil, e
	}
	return l.(*SAMMud), nil
}
