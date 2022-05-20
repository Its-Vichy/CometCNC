package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cnc_port = 444
	bot_port = 8890

	connected_clients int
	connected_bots    int
	max_retry         int = 3
	running_attacks   int = 0

	collection *mongo.Collection
	ctx        = context.TODO()

	blacklisted_target []string
	banned_ips         []string

	running_attacks_list = []attacks{}
	users                = []user{}
	bots                 = []bot{}
	themes               = []theme{}
	bots_types           = bots_type{}

	white  string = get_rgb(255, 255, 255)
	orange string = get_rgb(255, 99, 71)
	blue   string = get_rgb(88, 5, 255)
	red    string = get_rgb(255, 0, 0)
	green  string = get_rgb(0, 255, 0)
	grey   string = get_rgb(105, 105, 105)

	invisible string = "\033[30;40;196m"
	underline string = "\033[4m"
	clear     string = "\033[2J\033[1H"
	reset     string = "\033[0m"

	db_user string = "lmfao"
	db_pass string = "lmfao"

	discord_bot_token string = ".Yl3wEA."
	command_channel   string = "965748235794788432"
	use_discord       bool   = false

	reset_fade = false
)
