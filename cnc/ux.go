package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

func load_themes() {
	theme_name, err := ioutil.ReadDir("./themes")

	if err != nil {
		panic(err)
	}

	for _, theme_n := range theme_name {
		theme_file, err := ioutil.ReadDir(fmt.Sprintf("./themes/%s", theme_n.Name()))
		t := theme{}

		if err != nil {
			panic(err)
		}

		t.Name = theme_n.Name()

		for _, tfx_file := range theme_file {
			tfx_raw, _ := ioutil.ReadFile(fmt.Sprintf("./themes/%s/%s", theme_n.Name(), tfx_file.Name()))
			tfx_content := string(tfx_raw)

			if tfx_file.Name() == "description.txt" {
				t.Description = tfx_content
			}

			if tfx_file.Name() == "banner.tfx" {
				t.Content.Banner = tfx_content
			}

			if tfx_file.Name() == "welcome_banner.tfx" {
				t.Content.WelcomeBanner = tfx_content
			}

			if tfx_file.Name() == "line.tfx" {
				t.Content.Line = tfx_content
			}

			if tfx_file.Name() == "admin_command_banner.tfx" {
				t.Content.AdminCommandBanner = tfx_content
			}

			if tfx_file.Name() == "attack_send_banner.tfx" {
				t.Content.AttackSendBanner = tfx_content
			}

			if tfx_file.Name() == "help_banner.tfx" {
				t.Content.HelpBanner = tfx_content
			}

			if tfx_file.Name() == "invalid_args_banner.tfx" {
				t.Content.InvalidArgsBanner = tfx_content
			}

			if tfx_file.Name() == "methods_banner.tfx" {
				t.Content.MethodsBanner = tfx_content
			}

			if tfx_file.Name() == "not_an_admin_banner.tfx" {
				t.Content.NotAnAdminBanner = tfx_content
			}

			if tfx_file.Name() == "port_banner.tfx" {
				t.Content.PortBanner = tfx_content
			}

			if tfx_file.Name() == "profil_banner.tfx" {
				t.Content.ProfilBanner = tfx_content
			}

			if tfx_file.Name() == "text_banner.tfx" {
				t.Content.Text_banner = tfx_content
			}
			if tfx_file.Name() == "tools_banner.tfx" {
				t.Content.ToolsBanner = tfx_content
			}
			if tfx_file.Name() == "bot_types.tfx" {
				t.Content.BotTypeBanner = tfx_content
			}
		}

		themes = append(themes, t)
	}
}

func get_rgb(R int, G int, B int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", R, G, B)
}

func get_title(content string) string {
	return fmt.Sprintf("\033]0; %s\007", content)
}

func fade(content string, gradian int) string {
	final := ""
	count := 0
	pass := true

	for _, char := range content {
		if char == '*' {
			if pass {
				if reset_fade {
					count = 0
				}
				
				pass = false
				continue
			} else {
				pass = true
				continue
			}
		}

		if !pass {
			final = fmt.Sprintf("%s%s%c", final, white, char)
			continue
		}

		count += gradian

		switch char {
		case '♥':
			final = fmt.Sprintf("%s%s%c", final, red, char)

		default:
			final = fmt.Sprintf("%s%s%c", final, get_rgb(255, 52, count), char)
		}
	}

	return final
}

func render_theme(content string, client *user) string {
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(content, "<n>", "\n"), "<r>", "\r"), "<user>", client.Username)
	content = strings.ReplaceAll(content, "<MaxTime>", fmt.Sprintf("%d", client.MaxTime))
	content = strings.ReplaceAll(content, "<UsedConcurent>", fmt.Sprintf("%d", client.UsedConcurent))
	content = strings.ReplaceAll(content, "<Concurent>", fmt.Sprintf("%d", client.Concurent))
	content = strings.ReplaceAll(content, "<CooldownTime>", fmt.Sprintf("%d", client.CooldownTime))
	content = strings.ReplaceAll(content, "<Premium>", fmt.Sprintf("%t", client.Premium))
	content = strings.ReplaceAll(content, "<Cooldown>", fmt.Sprintf("%d", client.Cooldown))
	content = strings.ReplaceAll(content, "<LastedUsedMethod>", client.LastedUsedMethod)
	content = strings.ReplaceAll(content, "<LastedAttackTime>", fmt.Sprintf("%d", client.LastedAttackTime))
	content = strings.ReplaceAll(content, "<LastedTargetIp>", client.LastedTargetIp)
	content = strings.ReplaceAll(content, "<LastedTargetPort>", client.LastedTargetPort)
	content = strings.ReplaceAll(content, "<LastedStringToRender>", client.LastedStringToRender)

	content = strings.ReplaceAll(content, "<x86_64>", fmt.Sprintf("%d", bots_types.x86_64))
	content = strings.ReplaceAll(content, "<x86_32>", fmt.Sprintf("%d", bots_types.x86_32))
	content = strings.ReplaceAll(content, "<ARM2>", fmt.Sprintf("%d", bots_types.ARM2))
	content = strings.ReplaceAll(content, "<ARM3>", fmt.Sprintf("%d", bots_types.ARM3))
	content = strings.ReplaceAll(content, "<ARM4T>", fmt.Sprintf("%d", bots_types.ARM4T))
	content = strings.ReplaceAll(content, "<ARM5>", fmt.Sprintf("%d", bots_types.ARM5))
	content = strings.ReplaceAll(content, "<ARM6T2>", fmt.Sprintf("%d", bots_types.ARM6T2))
	content = strings.ReplaceAll(content, "<ARM6>", fmt.Sprintf("%d", bots_types.ARM6))
	content = strings.ReplaceAll(content, "<ARM7>", fmt.Sprintf("%d", bots_types.ARM7))
	content = strings.ReplaceAll(content, "<ARM64>", fmt.Sprintf("%d", bots_types.ARM64))
	content = strings.ReplaceAll(content, "<MIPS>", fmt.Sprintf("%d", bots_types.MIPS))
	content = strings.ReplaceAll(content, "<SUPERH>", fmt.Sprintf("%d", bots_types.SUPERH))
	content = strings.ReplaceAll(content, "<POWERPC>", fmt.Sprintf("%d", bots_types.POWERPC))
	content = strings.ReplaceAll(content, "<SPARC>", fmt.Sprintf("%d", bots_types.SPARC))
	content = strings.ReplaceAll(content, "<M68K>", fmt.Sprintf("%d", bots_types.M68K))
	content = strings.ReplaceAll(content, "<UNKNOWN_Build>", fmt.Sprintf("%d", bots_types.UNKNOWN_Build))
	content = strings.ReplaceAll(content, "<BIG_ENDIAN>", fmt.Sprintf("%d", bots_types.BIG_ENDIAN))
	content = strings.ReplaceAll(content, "<LITTLE_ENDIAN>", fmt.Sprintf("%d", bots_types.LITTLE_ENDIAN))
	content = strings.ReplaceAll(content, "<BIG_ENDIAN_W>", fmt.Sprintf("%d", bots_types.BIG_ENDIAN_W))
	content = strings.ReplaceAll(content, "<LITTLE_ENDIAN_W>", fmt.Sprintf("%d", bots_types.LITTLE_ENDIAN_W))
	content = strings.ReplaceAll(content, "<UNKNOWN_Endian>", fmt.Sprintf("%d", bots_types.UNKNOWN_Endian))


	if content == "<void>" {
		return ""
	}

	tmp := ""
	for _, line := range strings.Split(content, "\n") {
		tmp += "\r\n" + fade(line, 5)
	}
	content = tmp

	return fmt.Sprintf("%s%s", content, white)
}

func get_line(client *user) string {
	return render_theme(client.Theme.Content.Line, client)
}

func get_banner(client *user) string {
	return render_theme(client.Theme.Content.Banner, client)
}

func get_welcome_banner(client *user) string {
	return render_theme(client.Theme.Content.WelcomeBanner, client)
}

func clear_terminal(client *user, connection net.Conn) {
	connection.Write([]byte(clear))
	connection.Write([]byte(get_banner(client)))
	connection.Write([]byte(get_welcome_banner(client)))
}

func get_help_banner(client *user) string {
	return render_theme(client.Theme.Content.HelpBanner, client)
}

func get_port_banner(client *user) string {
	return render_theme(client.Theme.Content.PortBanner, client)
}

func get_methods_banner(client *user) string {
	/*
		Amplification: SNMP WSD DVR NTP ARD
		User Datagram: UDPBYPASS
		Transmission Control: TCPSYN TCPACK TCPTFO TCPTLS TCPAMP
		Layer3: IGMP GRE ESP IP RAND
		Special: VALVE FIVEM OVHAMP OVHTCP OVHUDP
		Layer7: SOCKET TLSv1 TLSv2
	*/
	return render_theme(client.Theme.Content.MethodsBanner, client)
}

func get_tools_banner(client *user) string {
	return render_theme(client.Theme.Content.ToolsBanner, client)
}

func get_profil_banner(client *user) string {
	return render_theme(client.Theme.Content.ProfilBanner, client)
}

func get_admin_command_banner(client *user) string {
	return render_theme(client.Theme.Content.AdminCommandBanner, client)
}

func get_not_an_admin_banner(client *user) string {
	return render_theme(client.Theme.Content.NotAnAdminBanner, client)
}

func get_invalid_args_banner(usage string, client *user) string {
	client.LastedStringToRender = usage

	return render_theme(client.Theme.Content.InvalidArgsBanner, client)
}

func get_text_banner(color string, text string, client *user) string {
	client.LastedStringToRender = text

	return render_theme(client.Theme.Content.Text_banner, client)
}

func get_attack_send_banner(client *user, method string, time int, target string, port string) string {
	client.LastedAttackTime = time
	client.LastedUsedMethod = method
	client.LastedTargetIp = target
	client.LastedTargetPort = port

	return render_theme(client.Theme.Content.AttackSendBanner, client)
}

func get_bots_type_banner(client *user) string {
	return render_theme(client.Theme.Content.BotTypeBanner, client)
}