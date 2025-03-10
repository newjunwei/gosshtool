package gosshtool

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
)

type SSHClient struct {
	SSHClientConfig
	remoteConn *ssh.Client
	session    *ssh.Session
}

func (c *SSHClient) Connect() (conn *ssh.Client, err error) {
	if c.remoteConn != nil {
		return
	}
	port := "22"
	host := c.Host
	hstr := strings.SplitN(host, ":", 2)
	if len(hstr) == 2 {
		host = hstr[0]
		port = hstr[1]
	}

	config := makeConfig(c.User, c.Password, c.Privatekey)
	conn, err = ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return
	}
	c.remoteConn = conn
	session, err := conn.NewSession()
	if err != nil {
		return
	}
	c.session = session
	return
}

func (c *SSHClient) Cmd(cmd string) (output, errput string, err error) {
	if c.session == nil {
		_, err = c.Connect()
		if err != nil {
			return
		}
	}
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	c.session.Stdout = &stdoutBuf
	c.session.Stderr = &stderrBuf
	err = c.session.Run(cmd)
	defer c.session.Close()
	output = stdoutBuf.String()
	errput = stderrBuf.String()
	return
}

func (c *SSHClient) Pipe(rw ReadWriteCloser, pty *PtyInfo) error {
	if c.session == nil {
		_, err := c.Connect()
		if err != nil {
			return err
		}
	}
	if err := c.session.RequestPty(pty.Term, pty.H, pty.W, pty.Modes); err != nil {
		return err
	}

	wc, err := c.session.StdinPipe()
	if err != nil {
		return err
	}
	go copyIO(wc, rw)

	r, err := c.session.StdoutPipe()
	if err != nil {
		return err
	}
	go copyIO(rw, r)
	er, err := c.session.StderrPipe()
	if err != nil {
		return err
	}
	go copyIO(rw, er)
	err = c.session.Shell()
	if err != nil {
		return err
	}
	err = c.session.Wait()
	if err != nil {
		return err
	}
	defer c.session.Close()
	return nil
}

func copyIO(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
