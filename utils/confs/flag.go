package confs

import "flag"

// pg db flags
var FlagPGName string
var FlagPGUser string
var FlagPGPassword string
var FlagPGHost string
var FlagPGPort string

// service flag
var FlagSericePort string

func init() {
	flag.StringVar(&FlagPGName, "pg-dbname", ConfigMap["pg-dbname"], "please specify pg db name")
	flag.StringVar(&FlagPGUser, "pg-user", ConfigMap["pg-user"], "please specify pg user")
	flag.StringVar(&FlagPGPassword, "pg-pwd", ConfigMap["pg-pwd"], "please specify pg password")
	flag.StringVar(&FlagPGHost, "pg-host", ConfigMap["pg-host"], "please specify pg host")
	flag.StringVar(&FlagPGPort, "pg-port", ConfigMap["pg-port"], "please specify pg port")
	flag.StringVar(&FlagSericePort, "service-port", ConfigMap["service-port"], "please specify gin service port")

	flag.Parse()
}
