package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const CACertificatePath = "/home/habanero/sandbox/teleport/mtls-test/cert.pem"
const CACertificateKey = "/home/habanero/sandbox/teleport/mtls-test/key.pem"

func main() {
	/* Read in client certificate for presenting to the server */
	cert, err := tls.LoadX509KeyPair(CACertificatePath, CACertificateKey)
	if err != nil {
		log.Fatal(err)
	}

	/* Create a CA certificate pool to trust */
	caCert, err := os.ReadFile(CACertificatePath)
	if err != nil {
		log.Fatal(err)
	}

	// This will add the public key to the existing CAs to trust
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	/* Create the client */
	// Now, we can create a proper HTTPs client using the above
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			MinVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{cert},

			// CipherSuites are not configurable under TLS 1.3 in Go
			// See https://github.com/golang/go/issues/29349 and
			// https://pkg.go.dev/crypto/tls#Config
		},
	}
	client := &http.Client{Transport: transport}

	/* Make the request */
	// Request /hello over port 8080 via the GET method
	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}
