package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

//Client for saving state client ssh
type Client struct {
	client *ssh.Client
}

func main() {

	key, err := ioutil.ReadFile("/home/fajar/.ssh/id_rsa")

	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err))
	}

	publicKey, err := ssh.ParsePrivateKey(key)

	if err != nil {
		panic(fmt.Errorf("error parsing key: %v", err))
	}

	configuration := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(publicKey),
		},
		HostKeyCallback: ssh.HostKeyCallback(
			func(hostname string, address net.Addr, key ssh.PublicKey) error {
				fmt.Println(key)
				return nil
			},
		),
	}

	client, err := Dial("beta.doogether.id:22", configuration)

	if err != nil {
		panic(fmt.Errorf("erro dialing server: %v", err))
	}

	defer client.client.Close()

	session, err := client.client.NewSession()

	if err != nil {
		panic(fmt.Errorf("erro creating session: %v", err))
	}

	defer session.Close()

	session.Stdin = os.Stdin
	// session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)

		err = session.Run(text)

		if err != nil {
			panic(fmt.Errorf("erro executing command: %v", err))
		}
	}

	fmt.Printf("Finish\n")
}

// Dial for dialing ssh server
func Dial(address string, configuration *ssh.ClientConfig) (*Client, error) {
	client, err := ssh.Dial("tcp", address, configuration)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}
