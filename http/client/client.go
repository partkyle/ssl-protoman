package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/shanemhansen/gossl"

	"flag"
	"fmt"
)

var logger = log.New(os.Stdout, "[http client] ", log.LstdFlags)

var (
	url = flag.String("url", "http://localhost:443", "url to GET")

	context *gossl.Context
)

func main() {
	flag.Parse()

	context = gossl.NewContext(gossl.SSLv23ClientMethod())
	context.SetCipherList("ALL:RC4+RSA:+HIGH:+MEDIUM:+LOW:+SSLv2:+EXP")

	transport := &OpenSSLRoundTripper{context}
	client := &http.Client{Transport: transport}

	logger.Printf("retrieving url=%q", *url)
	response, err := client.Get(*url)
	if err != nil {
		logger.Printf("error with http: %+v", err)
		return
	}

	fmt.Println("Output:")
	io.Copy(os.Stdout, response.Body)
}

func OpenSSLDial(network, address string) (net.Conn, error) {
	logger.Printf("initializing connection to network=%q address=%q", network, address)
	internalConn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	logger.Printf("creating ssl connection")
	conn, err := gossl.NewConn(context, internalConn)
	if err != nil {
		return nil, err
	}

	logger.Printf("performing handshake localaddr=%q remoteaddr=%q", conn.LocalAddr(), conn.RemoteAddr())
	if err := conn.Handshake(); err != nil {
		return nil, err
	}

	logger.Printf("successful ssl connection")
	return conn, nil
}

type OpenSSLRoundTripper struct {
	ctx *gossl.Context
}

func (o *OpenSSLRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme == "http" {
		return http.DefaultTransport.RoundTrip(r)
	}

	conn, err := net.Dial("tcp", r.URL.Host)
	if err != nil {
		return nil, err
	}

	sslConn, err := gossl.NewConn(o.ctx, conn)
	if err != nil {
		return nil, err
	}

	if err := sslConn.Handshake(); err != nil {
		return nil, err
	}

	return response, nil
}
