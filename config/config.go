package config

const (
	Host		= "localhost"
	Port		= 5432
	User		= "postgres"
	Password	= "securepass"
)

var PasswordTime uint32 = 1
var PasswordMemory uint32 = 64*1024
var PasswordThreads uint8 = 4
var PasswordKeyLen uint32 = 32
