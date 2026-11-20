package config

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// newTLSServer starts a TLS server with a fresh self-signed certificate
// and returns it together with the PEM-encoded certificate.
func newTLSServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("cannot generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{CommonName: "test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("cannot create certificate: %v", err)
	}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	// Rejected handshakes are expected in these tests.
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.TLS = &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key}},
	}
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv, string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

func doGet(cfg *Config, url string) error {
	resp, err := cfg.HTTPClient().Get(url)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func TestInitConfigCACert(t *testing.T) {
	srv, caPEM := newTLSServer(t)

	testCases := []struct {
		name          string
		caCert        string
		caCertFile    string // file contents, written to a temp file
		caCertPath    string // used as VM_INSTANCE_CA_CERT_FILE as is
		cloud         bool
		expectedError string
		expectTrusted bool
	}{
		{
			name:          "not set",
			expectTrusted: false,
		},
		{
			name:          "PEM certificate",
			caCert:        caPEM,
			expectTrusted: true,
		},
		{
			name:          "PEM certificate with escaped newlines",
			caCert:        strings.ReplaceAll(caPEM, "\n", `\n`),
			expectTrusted: true,
		},
		{
			name:          "PEM certificate file",
			caCertFile:    caPEM,
			expectTrusted: true,
		},
		{
			name:          "invalid PEM",
			caCert:        "not a certificate",
			expectedError: "failed to load VM_INSTANCE_CA_CERT: cannot initialize tls: cannot parse data from `ca` value",
		},
		{
			name:          "invalid PEM file",
			caCertFile:    "not a certificate",
			expectedError: "failed to load VM_INSTANCE_CA_CERT_FILE: cannot load root CAs: cannot parse data read from `ca_file`",
		},
		{
			name:          "http URL file",
			caCertPath:    "http://example.com/ca.pem",
			expectedError: "VM_INSTANCE_CA_CERT_FILE must be a local file path, got URL \"http://example.com/ca.pem\"",
		},
		{
			name:          "https URL file",
			caCertPath:    "https://example.com/ca.pem",
			expectedError: "VM_INSTANCE_CA_CERT_FILE must be a local file path, got URL \"https://example.com/ca.pem\"",
		},
		{
			name:          "both set",
			caCert:        caPEM,
			caCertFile:    caPEM,
			expectedError: "VM_INSTANCE_CA_CERT and VM_INSTANCE_CA_CERT_FILE cannot be set at the same time",
		},
		{
			name:          "cloud API not set",
			cloud:         true,
			expectTrusted: false,
		},
		{
			name:          "cloud API PEM certificate",
			caCert:        caPEM,
			cloud:         true,
			expectTrusted: true,
		},
		{
			name:          "cloud API PEM certificate file",
			caCertFile:    caPEM,
			cloud:         true,
			expectTrusted: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cloud {
				t.Setenv("VMC_API_KEY", "test-key")
				t.Setenv("VMC_API_BASE_URL", srv.URL)
				t.Setenv("VM_INSTANCE_ENTRYPOINT", "")
			} else {
				t.Setenv("VM_INSTANCE_ENTRYPOINT", srv.URL)
				t.Setenv("VM_INSTANCE_TYPE", "single")
			}
			t.Setenv("VM_INSTANCE_CA_CERT", tc.caCert)
			caCertFile := tc.caCertPath
			if tc.caCertFile != "" {
				caCertFile = filepath.Join(t.TempDir(), "ca.pem")
				if err := os.WriteFile(caCertFile, []byte(tc.caCertFile), 0o600); err != nil {
					t.Fatalf("cannot write CA file: %v", err)
				}
			}
			t.Setenv("VM_INSTANCE_CA_CERT_FILE", caCertFile)

			cfg, err := InitConfig()
			if tc.expectedError != "" {
				if err == nil || !strings.HasPrefix(err.Error(), tc.expectedError) {
					t.Fatalf("Expected error starting with %q, got: %v", tc.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if !tc.expectTrusted && cfg.HTTPClient() != http.DefaultClient {
				t.Fatal("Expected http.DefaultClient when no CA is configured")
			}
			if tc.cloud {
				_, err = cfg.VMC().ListCloudProviders(context.Background())
			} else {
				err = doGet(cfg, srv.URL)
			}
			if tc.expectTrusted && err != nil {
				t.Fatalf("Expected TLS request to succeed, got: %v", err)
			}
			if !tc.expectTrusted && err == nil {
				t.Fatal("Expected TLS verification error for untrusted certificate, got none")
			}
		})
	}
}

func TestInitConfigCACertFileMissing(t *testing.T) {
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "https://example.com")
	t.Setenv("VM_INSTANCE_TYPE", "single")
	t.Setenv("VM_INSTANCE_CA_CERT_FILE", filepath.Join(t.TempDir(), "missing.pem"))
	_, err := InitConfig()
	if err == nil || !strings.HasPrefix(err.Error(), "failed to load VM_INSTANCE_CA_CERT_FILE: ") {
		t.Fatalf("Expected error for missing CA file, got: %v", err)
	}
}

func TestInitConfigCACertFileReload(t *testing.T) {
	srvA, caA := newTLSServer(t)
	srvB, caB := newTLSServer(t)

	caFile := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(caFile, []byte(caA), 0o600); err != nil {
		t.Fatalf("cannot write CA file: %v", err)
	}
	t.Setenv("VM_INSTANCE_ENTRYPOINT", srvA.URL)
	t.Setenv("VM_INSTANCE_TYPE", "single")
	t.Setenv("VM_INSTANCE_CA_CERT", "")
	t.Setenv("VM_INSTANCE_CA_CERT_FILE", caFile)
	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if err := doGet(cfg, srvA.URL); err != nil {
		t.Fatalf("Expected request to server A to succeed, got: %v", err)
	}
	if err := doGet(cfg, srvB.URL); err == nil {
		t.Fatal("Expected request to server B to fail before CA rotation")
	}

	if err := os.WriteFile(caFile, []byte(caB), 0o600); err != nil {
		t.Fatalf("cannot rewrite CA file: %v", err)
	}
	// The CA file is re-read at most once per second.
	deadline := time.Now().Add(5 * time.Second)
	for {
		err = doGet(cfg, srvB.URL)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("Expected request to server B to succeed after CA rotation, got: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err := doGet(cfg, srvA.URL); err == nil {
		t.Fatal("Expected request to server A to fail after CA rotation")
	}
}
