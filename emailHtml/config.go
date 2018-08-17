package emailHtml

import (
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	Server   string
	Port     int
	Email    string
	Password string
}

//Configuration file containing email server, port ,id and password from config.toml
var Configuration = config{}

func init() {
	if _, err := toml.DecodeFile("emailHtml/config.toml", &Configuration); err != nil {
		log.Fatal(err)
	}
}
