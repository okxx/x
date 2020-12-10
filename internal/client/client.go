package client

import (
	"bufio"
	"fmt"
	"github.com/laracro/x/internal/instance"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	DefaultCiphers = []string{
		"aes128-ctr",
		"aes192-ctr",
		"aes256-ctr",
		"aes192-cbc",
		"aes256-cbc",
		"3des-cbc",
		"arcfour",
		"arcfour256",
		"arcfour128",
		"aes128-cbc",
		"blowfish-cbc",
		"cast128-cbc",
		"aes128-gcm@openssh.com",
		"chacha20-poly1305@openssh.com",
	}
)

type Client interface {
	Login()
}

type defaultClient struct {
	instance         	*instance.Instance
	clientConfig 		*ssh.ClientConfig
}

func genSSHConfig(instance *instance.Instance) *defaultClient {
	//u, err := user.Current()
	//if err != nil {
	//	log.Fatal(err)
	//	return nil
	//}

	var err error
	var authMethods []ssh.AuthMethod

	//if instance.KeyPath == "" && instance.Password == "" {}

	var pemBytes []byte
	if instance.KeyPath != "" {
		pemBytes, err = ioutil.ReadFile(instance.KeyPath)
		if err != nil {
			log.Fatal(err)
		} else {
			var signer ssh.Signer
			if instance.Passphrase != "" {
				signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(instance.Passphrase))
			} else {
				signer, err = ssh.ParsePrivateKey(pemBytes)
			}
			if err != nil {
				log.Fatal(err)
			} else {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	password := instance.GetPassword()
	if password != nil {
		authMethods = append(authMethods, password)
	}

	authMethods = append(authMethods, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, 0, len(questions))
		for i, q := range questions {
			fmt.Print(q)
			if echos[i] {
				scan := bufio.NewScanner(os.Stdin)
				if scan.Scan() {
					answers = append(answers, scan.Text())
				}
				err := scan.Err()
				if err != nil {
					return nil, err
				}
			} else {
				b, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return nil, err
				}
				fmt.Println()
				answers = append(answers, string(b))
			}
		}
		return answers, nil
	}))

	config := &ssh.ClientConfig{
		User:            instance.GetUser(),
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	config.SetDefaults()
	config.Ciphers = append(config.Ciphers, DefaultCiphers...)

	return &defaultClient{
		clientConfig: config,
		instance:         instance,
	}
}

func NewClient(instance *instance.Instance) Client {
	return genSSHConfig(instance)
}

func (c *defaultClient) Login() {
	host := c.instance.Host
	port := c.instance.GetPort()

	var client *ssh.Client
	client1, err := ssh.Dial("tcp", net.JoinHostPort(host, port), c.clientConfig)
	client = client1
	if err != nil {
		msg := err.Error()
		// use terminal password retry
		if strings.Contains(msg, "no supported methods remain") && !strings.Contains(msg, "password") {
			fmt.Printf("%s@%s's password:", c.clientConfig.User, host)
			var b []byte
			b, err = terminal.ReadPassword(int(syscall.Stdin))
			if err == nil {
				p := string(b)
				if p != "" {
					c.clientConfig.Auth = append(c.clientConfig.Auth, ssh.Password(p))
				}
				fmt.Println()
				client, err = ssh.Dial("tcp", net.JoinHostPort(host, port), c.clientConfig)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	fmt.Printf("connect server ssh -p %d %s@%s version: %s\n", c.instance.GetPort(), c.instance.GetUser(), host, string(client.ServerVersion()))

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		log.Fatal(err)
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", h, w, modes)
	if err != nil {
		log.Fatal(err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = session.Shell()
	if err != nil {
		log.Fatal(err)
		return
	}

	// change stdin to user
	go func() {
		_, err = io.Copy(stdinPipe, os.Stdin)
		log.Fatal(err)
		session.Close()
	}()

	// interval get terminal size
	// fix resize issue
	go func() {
		var (
			ow = w
			oh = h
		)
		for {
			cw, ch, err := terminal.GetSize(fd)
			if err != nil {
				break
			}

			if cw != ow || ch != oh {
				err = session.WindowChange(ch, cw)
				if err != nil {
					break
				}
				ow = cw
				oh = ch
			}
			time.Sleep(time.Second)
		}
	}()

	// send keepalive
	go func() {
		for {
			time.Sleep(time.Second * 10)
			client.SendRequest("keepalive@openssh.com", false, nil)
		}
	}()
	session.Wait()
}