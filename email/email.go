package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

// Sendmail .
func Sendmail(host, user, pass string, recver []string, tls, port uint32, body []byte) (err error) {
	auth := smtp.PlainAuth(
		"",
		user,
		pass,
		host,
	)
	if tls == 0 {
		// normal
		err = mailNormal(host, port, auth, user, recver, body)
	} else if tls == 1 {
		// starttls
		err = mailTLS(host, port, auth, user, recver, body)
	} else if tls == 2 {
		// ssl
		err = mailSSL(host, port, auth, user, recver, body)
	} else {
		err = fmt.Errorf("not support tls value: %d", tls)
	}
	return
}

func mailNormal(host string, port uint32, auth smtp.Auth, from string, to []string, body []byte) (err error) {
	var c *smtp.Client
	c, err = smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		err = fmt.Errorf("smtp.Dial(%s:%d), %s", host, port, err)
		return
	}
	defer c.Close()
	err = smtpsend(c, auth, from, to, body)
	return
}

func mailTLS(host string, port uint32, auth smtp.Auth, from string, to []string, body []byte) (err error) {
	var c *smtp.Client
	c, err = smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		err = fmt.Errorf("smtp.Dial(%s:%d), %s", host, port, err)
		return
	}
	defer c.Close()
	err = smtptlssend(c, auth, host, from, to, body)
	return
}

func mailSSL(host string, port uint32, auth smtp.Auth, from string, to []string, body []byte) (err error) {
	var c *smtp.Client
	var conn *tls.Conn
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err = tls.Dial("tcp", addr, nil)
	if err != nil {
		err = fmt.Errorf("tls.Dial(%s), %s", addr, err)
		return
	}
	defer conn.Close()

	//host, _, _ := net.SplitHostPort(addr)
	c, err = smtp.NewClient(conn, host)
	if err != nil {
		err = fmt.Errorf("smtp.NewClient(conn, %s), %s", host, err)
		return
	}

	err = smtpsend(c, auth, from, to, body)

	return
}
