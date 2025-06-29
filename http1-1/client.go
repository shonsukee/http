package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

func Client(w http.ResponseWriter, r *http.Request) {
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		http.Error(w, "Error loading client certificate", http.StatusInternalServerError)
		return
	}

	// CA証明書を読み込む
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		http.Error(w, "Error reading CA certificate", http.StatusInternalServerError)
		return
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{cert},
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// 通信作成
	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		http.Error(w, "Error dumping response", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(dump))
}
