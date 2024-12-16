package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

// SMTP stores all information for SMTP verification lookup
type SMTP struct {
	HostExists  bool `json:"host_exists"`
	FullInbox   bool `json:"full_inbox"`
	CatchAll    bool `json:"catch_all"`
	Deliverable bool `json:"deliverable"`
	Disabled    bool `json:"disabled"`
}

// Verifier holds configuration for the email verifier
type Verifier struct {
	smtpCheckEnabled     bool
	catchAllCheckEnabled bool
	proxyURI             string
	connectTimeout       time.Duration
	operationTimeout     time.Duration
	helloName            string
	fromEmail            string
}

// CheckSMTP performs an email verification on the passed domain via SMTP
func (v *Verifier) CheckSMTP(domain, username string) (*SMTP, error) {
	if !v.smtpCheckEnabled {
		return nil, nil
	}

	var ret SMTP
	email := fmt.Sprintf("%s@%s", username, domain)

	client, mx, err := newSMTPClient(domain, v.proxyURI, v.connectTimeout, v.operationTimeout)
	if err != nil {
		return &ret, err
	}
	defer client.Close()

	if err = client.Hello(v.helloName); err != nil {
		return &ret, err
	}

	if err = client.Mail(v.fromEmail); err != nil {
		return &ret, err
	}

	ret.HostExists = true
	ret.CatchAll = true

	if v.catchAllCheckEnabled {
		randomEmail := GenerateRandomEmail(domain)
		if err = client.Rcpt(randomEmail); err != nil {
			ret.CatchAll = false
		}
	}

	if username != "" && !ret.CatchAll {
		if err = client.Rcpt(email); err == nil {
			ret.Deliverable = true
		}
	}

	return &ret, nil
}

// newSMTPClient generates a new available SMTP client
func newSMTPClient(domain, proxyURI string, connectTimeout, operationTimeout time.Duration) (*smtp.Client, *net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return nil, nil, errors.New("No MX records found")
	}

	var done bool
	var mutex sync.Mutex
	ch := make(chan interface{}, 1)
	selectedMXCh := make(chan *net.MX, 1)

	for i, r := range mxRecords {
		addr := r.Host + ":25"
		go func(index int) {
			c, err := dialSMTP(addr, proxyURI, connectTimeout, operationTimeout)
			if err != nil {
				if !done {
					ch <- err
				}
				return
			}

			mutex.Lock()
			if !done {
				done = true
				ch <- c
				selectedMXCh <- mxRecords[index]
			} else {
				c.Close()
			}
			mutex.Unlock()
		}(i)
	}

	for {
		res := <-ch
		switch r := res.(type) {
		case *smtp.Client:
			return r, <-selectedMXCh, nil
		case error:
			if len(mxRecords) == 1 {
				return nil, nil, r
			}
		default:
			return nil, nil, errors.New("Unexpected response dialing SMTP server")
		}
	}
}

// dialSMTP is a timeout wrapper for smtp.Dial
func dialSMTP(addr, proxyURI string, connectTimeout, operationTimeout time.Duration) (*smtp.Client, error) {
	var conn net.Conn
	var err error

	if proxyURI != "" {
		conn, err = establishProxyConnection(addr, proxyURI, connectTimeout)
	} else {
		conn, err = establishConnection(addr, connectTimeout)
	}
	if err != nil {
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(operationTimeout))
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// GenerateRandomEmail generates a random email address using the domain passed
func GenerateRandomEmail(domain string) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := make([]byte, 32)
	for i := range r {
		r[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return fmt.Sprintf("%s@%s", string(r), domain)
}

// establishConnection connects to the address on the named network address
func establishConnection(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, timeout)
}

// establishProxyConnection connects to the address on the named network address via proxy protocol
func establishProxyConnection(addr, proxyURI string, timeout time.Duration) (net.Conn, error) {
	u, err := url.Parse(proxyURI)
	if err != nil {
		return nil, err
	}
	dialer, err := proxy.FromURL(u, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return dialer.(proxy.ContextDialer).DialContext(ctx, "tcp", addr)
}
