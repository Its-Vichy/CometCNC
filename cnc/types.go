package main

import (
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Premium   bool               `bson:"premium"`
	ThemeName string             `bson:"tfx_theme"`

	Concurent int  `bson:"concurent"`
	Cooldown  int  `bson:"cooldown"`
	MaxTime   int  `bson:"max_time"`
	Admin     bool `bson:"admin"`
	Banned    bool `bson:"banned"`

	Logged        bool
	Retry         int
	UsedConcurent int
	CooldownTime  int
	Theme         theme

	LastedAttackTime     int
	LastedUsedMethod     string
	LastedTargetIp       string
	LastedTargetPort     string
	LastedStringToRender string
}

type bot struct {
	connection net.Conn
	connected  bool
	ip         string
	build      string
	arch       string
}

type theme_content struct {
	Banner             string
	Line               string
	WelcomeBanner      string
	AdminCommandBanner string
	AttackSendBanner   string
	HelpBanner         string
	InvalidArgsBanner  string
	MethodsBanner      string
	NotAnAdminBanner   string
	PortBanner         string
	ProfilBanner       string
	Text_banner        string
	ToolsBanner        string
	BotTypeBanner      string
}

type theme struct {
	Name        string
	Description string
	Content     theme_content
}

type bots_type struct {
	x86_64        int
	x86_32        int
	ARM2          int
	ARM3          int
	ARM4T         int
	ARM5          int
	ARM6T2        int
	ARM6          int
	ARM7          int
	ARM64         int
	MIPS          int
	SUPERH        int
	POWERPC       int
	SPARC         int
	M68K          int
	UNKNOWN_Build int

	BIG_ENDIAN      int
	LITTLE_ENDIAN   int
	BIG_ENDIAN_W    int
	LITTLE_ENDIAN_W int
	UNKNOWN_Endian  int
}

// SOON
type attacks struct {
	ip     string
	port   int
	time   int
	method string
	sender string
}
