package controlplane

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/awslabs/eksdemo/pkg/application"
)

type Options struct {
	application.ApplicationOptions
	TrustAnchor string
	IssuerCert  string
	IssuerKey   string
}

func newOptions() (options *Options) {
	return &Options{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2024.7.3",
				PreviousChart: "2024.7.3",
				Latest:        "edge-24.7.3",
				Previous:      "edge-24.7.3",
			},
			Namespace: "linkerd",
		},
	}
}

func (options *Options) PreInstall() error {
	// CA Cert
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}

	serialNumberUpperBound := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberUpperBound)
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}

	rootCertTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "root.linkerd.cluster.local",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * time.Duration(365)),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		rootCertTemplate,
		rootCertTemplate,
		&rootKey.PublicKey,
		rootKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}

	certBuffer := bytes.Buffer{}
	err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}
	options.TrustAnchor = certBuffer.String()

	// Issuer Key
	issuerKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}
	b, _ := x509.MarshalECPrivateKey(issuerKey)

	keyBuffer := bytes.Buffer{}
	err = pem.Encode(&keyBuffer, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}
	options.IssuerKey = keyBuffer.String()

	// Issuer Cert
	serialNumberUpperBound = new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err = rand.Int(rand.Reader, serialNumberUpperBound)
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}

	issuerCertTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "identity.linkerd.cluster.local",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * time.Duration(365)),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err = x509.CreateCertificate(
		rand.Reader,
		issuerCertTemplate,
		rootCertTemplate,
		&issuerKey.PublicKey,
		rootKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}

	certBuffer = bytes.Buffer{}
	err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return fmt.Errorf("failed to create linkerd certificates: %w", err)
	}
	options.IssuerCert = certBuffer.String()

	return nil
}
