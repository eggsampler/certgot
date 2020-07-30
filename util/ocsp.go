package util

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"golang.org/x/crypto/ocsp"
)

// TODO: watch https://github.com/golang/go/issues/40017
func IsRevoked(cert, issuer *x509.Certificate) (bool, error) {
	if len(cert.OCSPServer) == 0 {
		return false, fmt.Errorf("no ocsp servers provided")
	}
	srv := cert.OCSPServer[rand.Intn(len(cert.OCSPServer))]
	srvUrl, err := url.Parse(srv)
	if err != nil {
		return false, fmt.Errorf("error parsing ocsp server %s: %v", srv, err)
	}
	opts := &ocsp.RequestOptions{Hash: crypto.SHA1}
	buf, err := ocsp.CreateRequest(cert, issuer, opts)
	if err != nil {
		return false, fmt.Errorf("error creating ocsp request: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, srv, bytes.NewBuffer(buf))
	if err != nil {
		return false, fmt.Errorf("error creating http request: %v", err)
	}
	req.Header.Add("Content-Type", "application/ocsp-request")
	req.Header.Add("Accept", "application/ocsp-response")
	req.Header.Add("host", srvUrl.Host)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending ocsp request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading ocsp response: %v", err)
	}
	ocspResp, err := ocsp.ParseResponse(body, issuer)
	if err != nil {
		return false, fmt.Errorf("error parsing ocsp response: %v", err)
	}
	if ocspResp.Status == ocsp.Unknown {
		return false, fmt.Errorf("ocsp server returned unknown")
	} else if ocspResp.Status == ocsp.Revoked {
		return true, nil
	}
	return false, nil
}
