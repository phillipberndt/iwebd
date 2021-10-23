package util

import (
	"bytes"
	"crypto/sha1"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"os"
	"time"
)

// Generate a temporary, self-signed certificate/key pair for iwebd servers.
func GenKeyPair() (tls.Certificate) {
	now := time.Now()
	hostname, _ := os.Hostname()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			CommonName:         hostname,
			Country:            []string{"USA"},
			Organization:       []string{"acme"},
			OrganizationalUnit: []string{"acme"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, 1),
		BasicConstraintsValid: true,
		IsCA:        false,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
		return tls.Certificate{}
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template,
		priv.Public(), priv)
	if err != nil {
		panic(err)
		return tls.Certificate{}
	}

	var outCert tls.Certificate
	outCert.Certificate = append(outCert.Certificate, cert)
	outCert.PrivateKey = priv

	fingerprint := sha1.Sum(cert)
	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}
	fmt.Printf("Fingerprint: %s\n", buf.String())

	return outCert
	}
