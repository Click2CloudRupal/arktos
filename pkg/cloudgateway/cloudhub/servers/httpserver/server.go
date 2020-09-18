package httpserver

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog"
	hubconfig "k8s.io/kubernetes/pkg/cloudgateway/cloudhub/config"
	"k8s.io/kubernetes/pkg/cloudgateway/common/constants"
)

const (
	certificateBlockType = "CERTIFICATE"
	// SiteID is for the clearer log
	SiteID = "SiteID"
)

// StartHTTPServer starts the http service
func StartHTTPServer() {
	router := mux.NewRouter()
	router.HandleFunc(constants.DefaultCertURL, edgeGatewayClientCert).Methods("GET")
	router.HandleFunc(constants.DefaultCAURL, getCA).Methods("GET")

	addr := fmt.Sprintf("%s:%d", hubconfig.Config.HTTPS.Address, hubconfig.Config.HTTPS.Port)

	cert, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: certificateBlockType, Bytes: hubconfig.Config.Cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: hubconfig.Config.Key}))

	if err != nil {
		klog.Fatal(err)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequestClientCert,
		},
	}
	klog.Fatal(server.ListenAndServeTLS("", ""))
}

// getCA returns the caCertDER
func getCA(w http.ResponseWriter, r *http.Request) {
	caCertDER := hubconfig.Config.Ca
	w.Write(caCertDER)
}

// edgeGatewayClientCert will create EdgeGatewayCert and return it
func edgeGatewayClientCert(w http.ResponseWriter, r *http.Request) {
	csrContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.Errorf("fail to read file when signing the cert for edge site:%s! error:%v", r.Header.Get(SiteID), err)
	}
	csr, err := x509.ParseCertificateRequest(csrContent)
	if err != nil {
		klog.Errorf("fail to ParseCertificateRequest of edge site: %s! error:%v", r.Header.Get(SiteID), err)
	}
	subject := csr.Subject
	clientCertDER, err := signCerts(subject, csr.PublicKey)
	if err != nil {
		klog.Errorf("fail to signCerts for edge site:%s! error:%v", r.Header.Get(SiteID), err)
	}

	w.Write(clientCertDER)
}

// signCerts will create a certificate for EdgeGateway
func signCerts(subInfo pkix.Name, pbKey crypto.PublicKey) ([]byte, error) {
	cfgs := &certutil.Config{
		CommonName:   subInfo.CommonName,
		Organization: subInfo.Organization,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientKey := pbKey

	ca := hubconfig.Config.Ca
	caCert, err := x509.ParseCertificate(ca)
	if err != nil {
		return nil, fmt.Errorf("unable to ParseCertificate: %v", err)
	}

	caKeyDER := hubconfig.Config.CaKey
	caKey, err := x509.ParseECPrivateKey(caKeyDER)
	if err != nil {
		return nil, fmt.Errorf("unable to ParseECPrivateKey: %v", err)
	}

	edgeCertSigningDuration := hubconfig.Config.CloudHub.EdgeCertSigningDuration
	certDER, err := NewCertFromCa(cfgs, caCert, clientKey, caKey, edgeCertSigningDuration) //crypto.Signer(caKey)
	if err != nil {
		return nil, fmt.Errorf("unable to NewCertFromCa: %v", err)
	}

	return certDER, err
}

// PrepareAllCerts check whether the certificates exist in the local directory, generate if they don't exist
func PrepareAllCerts() error {
	// Check whether the ca exists in the local directory
	if hubconfig.Config.Ca == nil && hubconfig.Config.CaKey == nil {
		klog.Info("Ca and CaKey don't exist in local directory, and will be created by CloudGateway")
		caDER, caKey, err := NewCertificateAuthorityDer()
		if err != nil {
			klog.Errorf("failed to create Certificate Authority, error: %v", err)
			return err
		}

		caKeyDER, err := x509.MarshalECPrivateKey(caKey.(*ecdsa.PrivateKey))
		if err != nil {
			klog.Errorf("failed to convert an EC private key to SEC 1, ASN.1 DER form, error: %v", err)
			return err
		}

		UpdateConfig(caDER, caKeyDER, nil, nil)
	}

	// Check whether the CloudGateway certificates exist in the local directory
	if hubconfig.Config.Key == nil && hubconfig.Config.Cert == nil {
		klog.Infof("CloudGatewayCert and key don't exist in local directory, and will be signed by CA")
		certDER, keyDER, err := SignCerts()
		if err != nil {
			klog.Errorf("failed to sign a certificate, error: %v", err)
			return err
		}

		UpdateConfig(nil, nil, certDER, keyDER)
	}
	return nil
}