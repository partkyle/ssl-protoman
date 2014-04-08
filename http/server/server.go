package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/partkyle/gossl"
	"github.com/sendgrid/go-sglog"
)

var logger = sglog.NewLevelLogger("httpssl", sglog.INFO)

var (
	addr = flag.String("addr", "localhost:8443", "host and port to listen on")
	cert = flag.String("cert", "", "cert file")
	key  = flag.String("key", "", "key file")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Fatal("error listening", sglog.Values{"err": err})
	}

	logger.Info("starting listener", sglog.Values{"addr": listener.Addr()})

	context := gossl.NewContext(gossl.SSLv23ServerMethod())
	context.UseCertificateFile(*cert, gossl.FILETYPE_PEM)
	context.UsePrivateKeyFile(*key, gossl.FILETYPE_PEM)
	context.SetCipherList("ALL:RC4+RSA:+HIGH:+MEDIUM:+LOW:+SSLv2:+EXP")

	sslListener, err := gossl.NewListener(listener, context)
	if err != nil {
		logger.Fatal("error starting ssl listener", sglog.Values{"err": err})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "giddy up", 409)
	})

	if err := http.Serve(sslListener, nil); err != nil {
		logger.Fatal("error serving http", sglog.Values{"err": err})
	}
}
