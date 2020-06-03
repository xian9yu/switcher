package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type configs struct {
	//LogLevel string           `json:"log_level"`
	Rules []*rules `json:"rules"`
}

type rules struct {
	Name   string `json:"name"`
	Listen string `json:"listen"`
	//EnableRegexp bool   `json:"enable_regexp"`
	Targets []*struct {
		Regexp  string         `json:"regexp"`
		regexp  *regexp.Regexp //`json:"-"`
		Address string         `json:"address"`
	} `json:"targets"`
	//FirstPacketTimeout uint64 `json:"first_packet_timeout"`
}

var config *configs

func init() {
	var configPath = "./config.json"
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	}
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("[ERROR] failed to load config.json: %s", err)
	}

	if err := json.Unmarshal(buf, &config); err != nil {
		log.Printf("[ERROR] failed to load config.json: %s", err)
	}

	if len(config.Rules) == 0 {
		log.Println("[ERROR] empty rule", err)
	}

	for i, v := range config.Rules {
		if err := v.verify(); err != nil {
			log.Printf("[ERROR] verity rule failed at pos %d : %s", i, err)
		}
	}
}

func (r *rules) verify() error {
	if r.Name == "" {
		return fmt.Errorf("[ERROR] empty name")
	}
	if r.Listen == "" {
		return fmt.Errorf("[ERROR] invalid listen address")
	}
	if len(r.Targets) == 0 {
		return fmt.Errorf("[ERROR] invalid targets")
	}

	for i, v := range r.Targets {
		if v.Address == "" {
			return fmt.Errorf("[ERROR] invalid address at pos %d", i)
		}
		//if r.EnableRegexp {
		r, err := regexp.Compile(v.Regexp)
		if err != nil {
			return fmt.Errorf("[ERROR] invalid regexp at pos %d : %s", i, err.Error())
		}
		v.regexp = r
		//}
	}
	return nil
}
