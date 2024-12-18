package config

import (
	"bytes"
	"cmp"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
	"math/big"
	"net"
	"os"
	"path"
	"time"
)

// App application configuration
type App struct {
	// Listen host:port on which web service will operate
	Listen string `env:"SERVER_ADDRESS" json:"server_address,omitempty"`
	// BaseShortURL HTTP domain append to short URL
	BaseShortURL string `env:"BASE_URL" json:"base_url,omitempty"`
	// FileStoragePath Path to file database
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path,omitempty"`
	// DBDsn Database connection string
	DBDsn string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	// HTTPS config
	HTTPS *HTTPS
}

// HTTPS configuration
type HTTPS struct {
	Enable bool `env:"ENABLE_HTTPS" json:"enable_https"`
	Key    string
	Pem    string
}

// GetBaseShortURL return HTTP domain, append to short URL
func (a *App) GetBaseShortURL() string {
	return a.BaseShortURL
}

// GetStoragePath Path to file database
func (a *App) GetStoragePath() string {
	return a.FileStoragePath
}

// GetDsn Database connection string:
// Example: postgresql://postgres:postgres@localhost:5432/short
func (a *App) GetDsn() string {
	return a.DBDsn
}

// NewConfig Factory
func NewConfig() (*App, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't retrive pwd %w", err)
	}

	cfg := &App{HTTPS: &HTTPS{
		Key: path.Join(pwd, "/server.key"),
		Pem: path.Join(pwd, "/server.pem"),
	}}

	var configPath string

	flag.StringVar(&cfg.Listen, "a", ":8080", "Http service list addr")
	flag.StringVar(&cfg.BaseShortURL, "b", "http://localhost:8080", "Base short url")
	flag.StringVar(&cfg.FileStoragePath, "f", path.Join(pwd, "/data/db"), "set db file path")
	flag.StringVar(&cfg.DBDsn, "d", "", "Set postgresql connection string (DSN)")
	flag.BoolVar(&cfg.HTTPS.Enable, "s", false, "Run server in https")
	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.Parse()

	if err = env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("can't parse environment variables: %w", err)
	}

	if os.Getenv("CONFIG") != "" {
		configPath = os.Getenv("CONFIG")
	}

	if configPath != "" {
		var rawCfg []byte
		jCfg := &App{HTTPS: &HTTPS{}}

		rawCfg, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("can't load config file: %w", err)
		}

		if err = json.Unmarshal(rawCfg, jCfg); err != nil {
			return nil, fmt.Errorf("can't unmarshal config: %w", err)
		}

		cfg.Listen = cmp.Or(cfg.Listen, jCfg.Listen)
		cfg.BaseShortURL = cmp.Or(cfg.BaseShortURL, jCfg.BaseShortURL)
		cfg.FileStoragePath = cmp.Or(cfg.FileStoragePath, jCfg.FileStoragePath)
		cfg.DBDsn = cmp.Or(cfg.DBDsn, jCfg.DBDsn)
		cfg.HTTPS.Enable = cmp.Or(cfg.HTTPS.Enable, jCfg.HTTPS.Enable)
	}

	if cfg.HTTPS.Enable {
		err = createTLS(cfg.HTTPS.Pem, cfg.HTTPS.Key)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func createTLS(pemPath string, keyPath string) error {
	cert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Shortener"},
			Country:      []string{"RU"},
			Province:     []string{"Moscow"},
			Locality:     []string{"Moscow"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	certBytes, _ := x509.CreateCertificate(rand.Reader, &cert, &cert, &privateKey.PublicKey, privateKey)
	err := saveCertToFile(pemPath, "CERTIFICATE", certBytes)
	if err != nil {
		return err
	}

	err = saveCertToFile(keyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
	if err != nil {
		return err
	}

	return nil
}

func saveCertToFile(filePath string, cypherType string, cypher []byte) error {
	var (
		buf  bytes.Buffer
		file *os.File
	)

	err := pem.Encode(&buf, &pem.Block{
		Type:  cypherType,
		Bytes: cypher,
	})
	if err != nil {
		return fmt.Errorf("can't encode pem: %w", err)
	}

	file, _ = os.Create(filePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = buf.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}
