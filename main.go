package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Config struct {
	InHost      string `json:"inHost"`
	InPort      int    `json:"inPort"`
	OutHost     string `json:"outHost"`
	OutPort     int    `json:"outPort"`
	OutUsername string `json:"outUsername"`
	OutPassword string `json:"outPassword"`
}

var config *Config

// The Backend implements SMTP server methods.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after successful login.
type Session struct {
	From string
	To   []string
}

// AuthPlain implements authentication using SASL PLAIN.
func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.From = from
	log.Println("Mail from:", s.From)
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.To = append(s.To, to)
	log.Println("Rcpt to:", s.To)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	auth := sasl.NewPlainClient("", config.OutUsername, config.OutPassword)
	addr := fmt.Sprintf("%s:%d", config.OutHost, config.OutPort)
	err := smtp.SendMail(addr, auth, s.From, s.To, r)
	if err != nil {
		return err
	}
	log.Println("Data forwarded to:", s.To)
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	// Get config

	data, err := os.ReadFile("smtpfw.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	if config.OutHost == "" {
		log.Fatal("outHost config is required")
	}
	if config.OutPort == 0 {
		log.Fatal("outPort config is required")
	}
	if config.OutUsername == "" {
		log.Fatal("outUsername config is required")
	}
	if config.OutPassword == "" {
		log.Fatal("outPasOutPassword config is required")
	}

	if config.InHost == "" {
		config.InHost = "localhost"
	}

	if config.InPort == 0 {
		config.InPort = 25
	}

	// Start server

	s := smtp.NewServer(new(Backend))

	s.Addr = fmt.Sprintf("%s:%d", config.InHost, config.InPort)
	s.AllowInsecureAuth = true
	s.WriteTimeout = 180 * time.Second
	s.ReadTimeout = 180 * time.Second

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
