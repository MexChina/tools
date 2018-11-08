package tools

import (
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"fmt"
	"bytes"
	"strings"
)

func main() {
	runSsh()
}

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func runSsh() {
	var stdOut, stdErr bytes.Buffer
	session, err := SSHConnect("root", "123456", "192.168.1.66", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	session.Run("pwd")

	fmt.Println(stdOut.String())
	session.Run("cd /opt/wwwroot/go;ll")
	fmt.Println(stdOut.String())

	log.Println(strings.Replace(stdOut.String(), "\n", "", -1))
	//ret, err := strconv.Atoi(strings.Replace(stdOut.String(), "\n", "", -1))
	//if err != nil {
	//	log.Println(err)
	//	//panic(err)
	//}
	//fmt.Printf("%d, %s\n", ret, stdErr.String())
}
