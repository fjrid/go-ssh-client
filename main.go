package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

//Client for saving state client ssh
type Client struct {
	client *ssh.Client
}

func main() {

	key, err := ioutil.ReadFile("path_to_private_key")

	if err != nil {
		panic(fmt.Errorf("Error reading key: %v", err))
	}

	publicKey, err := ssh.ParsePrivateKey(key)

	if err != nil {
		panic(fmt.Errorf("Error parsing key: %v", err))
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

	client, err := Dial("127.0.0.1:22", configuration)

	if err != nil {
		panic(fmt.Errorf("Erro dialing server: %v", err))
	}

	defer client.client.Close()

	session, err := client.client.NewSession()

	if err != nil {
		panic(fmt.Errorf("Erro creating session: %v", err))
	}

	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run("ls")

	if err != nil {
		panic(fmt.Errorf("Erro executing command: %v", err))
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
