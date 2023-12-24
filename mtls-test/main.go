package main

// For generating ED25519 TLS certificate
// openssl req -newkey ed25519   -new -nodes -x509   -days 3650   -out cert.pem   -keyout key.pem -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost" -addext "subjectAltName = DNS:localhost"

// For configuring TLS https://pkg.go.dev/crypto/tls#Config

// For curl'ing in mTLS world
// curl -D - -o -  --cert cert.pem --key key.pem --cacert cert.pem 'https://localhost:8443/hello'

import (
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"os"
)

const CACertificatePath = "/home/habanero/sandbox/teleport/mtls-test/cert.pem"
const CACertificateKey = "/home/habanero/sandbox/teleport/mtls-test/key.pem"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	slog.InfoContext(r.Context(), "Recieved request: ", "method", r.Method, "uri", r.RequestURI, "remote-address", r.RemoteAddr, "agent", r.UserAgent())
	io.WriteString(w, "Hello, Mundo!\n")
}

func main() {
	/***** Set up TLS configuration for mTLS *****/

	/* Create a CA certificate pool to trust */
	caCert, err := os.ReadFile(CACertificatePath)
	if err != nil {
		log.Fatal(err)
	}

	// This will add the public key to the existing CAs to trust
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("Invalid certificate in provided CA PEM")
	}

	/* Create the TLS configuration for mTLS support */
	tlsConfig := &tls.Config{
		// This will require that the client present a TLS certificate
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS13,

		// CipherSuites are not configurable under TLS 1.3 in Go
		// See https://github.com/golang/go/issues/29349 and
		// https://pkg.go.dev/crypto/tls#Config

		// Additional verifications to perform on the client
		// See https://pkg.go.dev/crypto/tls#example-Config-VerifyConnection
		VerifyConnection: func(cs tls.ConnectionState) error {
			slog.Info("Verifying client:", "tlsVersion", cs.Version, "cipherSuite", cs.CipherSuite, "serverName", cs.ServerName)

			return nil
		},
	}

	/***** Run the server *****/
	/* Create the server */
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	// Listen to port 8080 and wait
	slog.Info("Starting the server . . .")
	log.Fatalln(server.ListenAndServeTLS(CACertificatePath, CACertificateKey))
}
