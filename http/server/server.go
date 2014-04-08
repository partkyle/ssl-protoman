package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/partkyle/gossl"
)

var logger = log.New(os.Stdout, "[http ssl server] ", log.LstdFlags)

var (
	addr = flag.String("addr", "localhost:8443", "host and port to listen on")
	cert = flag.String("cert", "", "cert file")
	key  = flag.String("key", "", "key file")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Fatalf("error listening: %s", err)
	}

	logger.Printf("starting listener addr=%q", listener.Addr())

	context := gossl.NewContext(gossl.SSLv23ServerMethod())
	context.UseCertificateFile(*cert, gossl.FILETYPE_PEM)
	context.UsePrivateKeyFile(*key, gossl.FILETYPE_PEM)
	context.SetCipherList("ALL:RC4+RSA:+HIGH:+MEDIUM:+LOW:+SSLv2:+EXP")

	sslListener, err := gossl.NewListener(listener, context)
	if err != nil {
		logger.Fatalf("error starting ssl listener: ", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "giddy up", 409)
	})

	if err := http.Serve(sslListener, nil); err != nil {
		logger.Fatal("error serving http: %s", err)
	}
}
