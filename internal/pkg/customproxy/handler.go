package customproxy

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"github.com/DmitriyKomarovCoder/http_proxy/common/logger"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/utils"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"syscall"
	"time"
)

const (
	caMaxAge   = 5 * 365 * 24 * time.Hour
	leafMaxAge = 24 * time.Hour
	caUsage    = x509.KeyUsageDigitalSignature |
		x509.KeyUsageContentCommitment |
		x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDataEncipherment |
		x509.KeyUsageKeyAgreement |
		x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign
	leafUsage = caUsage
)

func LoadCA(certFile, keyFile, subject string) (cert tls.Certificate, err error) {
	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if os.IsNotExist(err) {
		cert, err = genFilesCA(certFile, keyFile, subject)
	}
	if err == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	}
	return
}

func genKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
}

func genFilesCA(certFile, keyFile, subject string) (cert tls.Certificate, err error) {
	err = os.MkdirAll("certs", 0700)
	if err != nil {
		return
	}
	certPEM, keyPEM, err := GenCA(subject)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = os.WriteFile(certFile, certPEM, 0400)
	if err == nil {
		err = os.WriteFile(keyFile, keyPEM, 0400)
	}
	return cert, err
}

func GenCA(name string) (certPEM, keyPEM []byte, err error) {
	now := time.Now().UTC()
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             now,
		NotAfter:              now.Add(caMaxAge),
		KeyUsage:              caUsage,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		SignatureAlgorithm:    x509.ECDSAWithSHA512,
	}
	key, err := genKeyPair()
	if err != nil {
		return
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		return
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return
	}
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: keyDER,
	})
	return
}

func GenerateFakeCertificate(host []string, Ca *tls.Certificate) (*tls.Certificate, error) {
	caCert, err := x509.ParseCertificate(Ca.Certificate[0])
	if err != nil {
		return nil, err
	}

	caPrivKey := Ca.PrivateKey

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"yngwie customproxy CA"},
		},
		DNSNames:              host,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, &priv.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  priv,
	}

	return &cert, nil
}

type ProxyHandler struct {
	apiUsecase api.Usecase
	log        logger.Logger
	CA         *tls.Certificate
}

func NewProxy(apiUsecase api.Usecase, log logger.Logger, ca *tls.Certificate) *ProxyHandler {
	return &ProxyHandler{
		apiUsecase: apiUsecase,
		log:        log,
		CA:         ca,
	}
}

type ForwardProxyTransport struct {
	http.Transport
}

var CheckRedirectDisabler = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (t *ForwardProxyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Del("Proxy-Connection")

	return t.Transport.RoundTrip(r)
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleHTTPS(w, r)
		return
	}

	p.handleHTTP(w, r)
}

func (p *ProxyHandler) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	p.log.Infof("Connect to: %v", r.Host)

	hj, ok := w.(http.Hijacker)
	if !ok {
		p.log.Error("http server doesn't support hijacking connection")
		return
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		p.log.Error("http hijacking failed")
		return
	}
	defer clientConn.Close()

	host, _, err := net.SplitHostPort(r.Host)

	if err != nil {
		p.log.Error("error splitting host/port:", err)
		return
	}

	fakeCert, err := GenerateFakeCertificate([]string{host}, p.CA)
	if err != nil {
		p.log.Fatal("can't generate certificate for", host)
		return
	}

	if _, err := clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")); err != nil {
		p.log.Error("error writing status to client:", err)
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		Certificates:     []tls.Certificate{*fakeCert},
	}

	tlsConn := tls.Server(clientConn, tlsConfig)
	defer tlsConn.Close()

	connReader := bufio.NewReader(tlsConn)

	for {
		req, err := http.ReadRequest(connReader)
		if err == io.EOF {
			break
		} else if errors.Is(err, syscall.ECONNRESET) {
			p.log.Info("This is connection reset by peer error")
			break
		} else if err != nil {
			p.log.Fatal(req, err)
		}

		if b, err := httputil.DumpRequest(req, false); err == nil {
			log.Printf("incoming request:\n%s\n", string(b))
		}

		utils.ChangeRequestToTarget(req, r.Host)

		client := http.Client{}

		var bodyBufferRequest bytes.Buffer
		_, err = io.Copy(&bodyBufferRequest, req.Body)
		if err != nil {
			p.log.Errorf("Error reading response body: %v", err)
			return
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyBufferRequest.Bytes()))

		resp, err := client.Do(req)
		if err != nil {
			p.log.Println("error sending request to target:", err)
			break
		}
		copyReq := req

		if b, err := httputil.DumpResponse(resp, false); err == nil {
			log.Printf("target response:\n%s\n", string(b))
		}
		defer resp.Body.Close()

		copyResp := *resp
		var bodyBufferResponse bytes.Buffer
		_, err = io.Copy(&bodyBufferResponse, copyResp.Body)
		if err != nil {
			p.log.Errorf("Error reading response body: %v", err)
			return
		}
		resp.Body = io.NopCloser(bytes.NewReader(bodyBufferResponse.Bytes()))

		if err := resp.Write(tlsConn); err != nil {
			p.log.Println("error writing response back:", err)
		}

		reqId, err := p.apiUsecase.SaveRequest(copyReq, bodyBufferRequest.Bytes())
		if err != nil {
			p.log.Errorf("Error save: %v", err)
			http.Error(w, "failed http save %v", http.StatusBadRequest)
			return
		}

		err = p.apiUsecase.SaveResponse(reqId, copyResp, bodyBufferResponse.Bytes())
		if err != nil {
			p.log.Errorf("error save response: %v", err)
			return
		}
	}
}

func (p *ProxyHandler) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if b, err := httputil.DumpRequest(r, true); err == nil {
		p.log.Infof("incoming request:\n%s\n", string(b))
	}

	var bodyBufferRequest bytes.Buffer
	_, err := io.Copy(&bodyBufferRequest, r.Body)
	if err != nil {
		p.log.Errorf("Error reading response body: %v", err)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBufferRequest.Bytes()))

	reqId, err := p.apiUsecase.SaveRequest(r, bodyBufferRequest.Bytes())
	if err != nil {
		p.log.Errorf("Error save: %v", err)
		http.Error(w, "failed http save %v", http.StatusBadRequest)
		return
	}
	r.RequestURI = ""

	client := http.Client{
		Transport:     &ForwardProxyTransport{},
		CheckRedirect: CheckRedirectDisabler,
	}

	resp, err := client.Do(r)
	if err != nil {
		p.log.Errorf("Error do http client: %v", err)
		http.Error(w, "failed http connect", http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if b, err := httputil.DumpResponse(resp, false); err == nil {
		log.Printf("target response:\n%s\n", string(b))
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		p.log.Fatal("http server doesn't support hijacking connection")
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		p.log.Fatal("http hijacking failed")
	}
	defer clientConn.Close()

	var bodyBufferResponse bytes.Buffer
	_, err = io.Copy(&bodyBufferResponse, resp.Body)
	if err != nil {
		p.log.Errorf("Error reading response body: %v", err)
		return
	}
	resp.Body = io.NopCloser(bytes.NewReader(bodyBufferResponse.Bytes()))

	copyResp := *resp
	if err := resp.Write(clientConn); err != nil {
		p.log.Errorf("error writing response back: %v", err)
		return
	}

	err = p.apiUsecase.SaveResponse(reqId, copyResp, bodyBufferResponse.Bytes())
	if err != nil {
		p.log.Errorf("error save response: %v", err)
		return
	}
}
