package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

func smtpsend(c *smtp.Client, auth smtp.Auth, from string, to []string, body []byte) (err error) {
	// Auth
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				err = fmt.Errorf("smtp.Auth(), %s", err)
				return err
			}
		}
	}

	// To && From
	if err = c.Mail(from); err != nil {
		err = fmt.Errorf("smtp.Mail(%s), %s", from, err)
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			err = fmt.Errorf("smtp.Rcpt(%s), %s", addr, err)
			return err
		}
	}

	// Data
	w, err := c.Data()
	if err != nil {
		err = fmt.Errorf("smtp.Data(), %s", err)
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		err = fmt.Errorf("smtp.Write(msg), %s", err)
		return err
	}

	err = w.Close()
	if err != nil {
		err = fmt.Errorf("w.Close(), %s", err)
		return err
	}

	return c.Quit()
}

func smtptlssend(c *smtp.Client, auth smtp.Auth, host, from string, to []string, body []byte) (err error) {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	err = c.StartTLS(tlsconfig)
	if err != nil {
		err = fmt.Errorf("smtp.StartTLS(), %s", err)
		return err
	}

	err = smtpsend(c, auth, from, to, body)
	return
}
