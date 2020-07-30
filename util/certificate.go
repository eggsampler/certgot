package util

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

func ReadCertificate(path string) (*x509.Certificate, error) {
	certPem, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading certificate file %s: %v", path, err)
	}
	cert, err := readCertificatePem(certPem)
	if err != nil {
		return nil, fmt.Errorf("error read certificate file %s: %v", path, err)
	}
	return cert, nil
}

func readCertificatePem(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("no certificate present")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("not a certificate: %s", block.Type)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
