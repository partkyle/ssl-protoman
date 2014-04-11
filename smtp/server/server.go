package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/partkyle/gossl"
	"github.com/sendgrid/go-smtp"
)

var (
	addr = flag.String("addr", "localhost:8443", "host and port to listen on")
	cert = flag.String("cert", "../ssl/server.crt", "cert file")
	key  = flag.String("key", "../ssl/server.long.key", "key file")
)

func main() {

	flag.Parse()

	fmt.Println("Starting")
	timeout := 100 * time.Second

	// set up the smtp server
	s := smtp.NewServer("banner", nil)
	s.IdleTimeout = 100 * time.Second
	s.ShutdownTimeout = timeout

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal("Unable to start tcp listener", err)
	}

	context := gossl.NewContext(gossl.SSLv23ServerMethod())
	context.UseCertificateFile(*cert, gossl.FILETYPE_PEM)
	context.UsePrivateKeyFile(*key, gossl.FILETYPE_PEM)

	sslListener, err := gossl.NewListener(l, context)
	if err != nil {
		log.Fatal(err)
	}

	go s.Serve(sslListener)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	fmt.Println("waiting..")
	wg.Wait()

	// never should get here
	fmt.Println("stopping")
	s.Stop()
}
