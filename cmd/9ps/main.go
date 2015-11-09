package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"github.com/docker/pinata/v1/pkg/p9p"
	"golang.org/x/net/context"
)

var (
	root string
	addr string
)

func init() {
	flag.StringVar(&root, "root", "~/", "root of filesystem to serve over 9p")
	flag.StringVar(&addr, "addr", ":5640", "bind addr for 9p server, prefix with unix: for unix socket")
}

func main() {
	ctx := context.Background()
	log.SetFlags(0)
	flag.Parse()

	proto := "tcp"
	if strings.HasPrefix(addr, "unix:") {
		proto = "unix"
		addr = addr[5:]
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatalln("error listening:", err)
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatalln("error accepting:", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			ctx := context.WithValue(ctx, "conn", conn)
			log.Println("connected", conn.RemoteAddr())
			session, err := newLocalSession(ctx, root)
			if err != nil {
				log.Println("error creating session")
				return
			}

			p9pnew.Serve(ctx, conn, p9pnew.Dispatch(session))
		}(c)
	}
}

// newLocalSession returns a session to serve the local filesystem, restricted
// to the provided root.
func newLocalSession(ctx context.Context, root string) (p9pnew.Session, error) {
	// silly, just connect to ufs for now! replace this with real code later!
	log.Println("dialing", ":5640", "for", ctx.Value("conn"))
	conn, err := net.Dial("tcp", ":5640")
	if err != nil {
		return nil, err
	}

	session, err := p9pnew.NewSession(ctx, conn)
	if err != nil {
		return nil, err
	}

	return session, nil
}
