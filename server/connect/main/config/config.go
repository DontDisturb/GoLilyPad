package config

import "io/ioutil"
import "regexp"
import yaml "launchpad.net/goyaml"
import "github.com/LilyPad/GoLilyPad/server/connect"

type Config struct {
	Bind string `yaml:"bind"`
	Logins []ConfigLogin `yaml:"logins"`
}

func (this *Config) Authenticate(username string, password string, passwordSalt string) (ok bool, err error) {
	for _, login := range this.Logins {
		var realPassword string
		if len(login.Username) != 0 && login.Username == username {
			realPassword = login.Password
		} else if len(login.Regexp) != 0 {
			if login.regexp == nil {
				login.regexp, err = regexp.Compile(login.Regexp)
				if err != nil {
					return
				}
			}
			if login.regexp.MatchString(username) {
				realPassword = login.Password
			} else {
				continue
			}
		} else {
			continue
		}
		realPassword = connect.PasswordAndSaltHash(realPassword, passwordSalt)
		if realPassword == password {
			ok = true
			return
		}
	}
	ok = false
	return
}

type ConfigLogin struct {
	Username string `yaml:"username,omitempty"`
	Regexp string `yaml:"regexp,omitempty"`
	regexp *regexp.Regexp
	Password string `yaml:"password"`
}

func DefaultConfig() (config *Config) {
	return &Config{
		Bind: ":5091",
		Logins: []ConfigLogin{
			ConfigLogin{"example", "", nil, "example"},
			ConfigLogin{"", "^example-.*$", nil, "example"},
		},
	}
}

func LoadConfig(file string) (config *Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	var cfg Config
	config = &cfg
	err = yaml.Unmarshal(data, config)
	return
}

func SaveConfig(file string, config *Config) (err error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(file, data, 0644)
	return
}
