package common

import (
    config "github.com/goconf"
    log "code.google.com/p/log4go"
)

var Conf configFile2

const (
	COMMON    = "common"
	TAOKE = "taoke"
)

type configFile2 struct {
	conf *config.ConfigFile
}

func init() {
	if err := Conf.LoadConfigFile("conf/taoke.conf"); err != nil {
		panic(err)
	}
}

func (cf *configFile2) LoadConfigFile(file string) (err error) {
	cf.conf, err = config.ReadConfigFile(file)
	return
}

func (cf *configFile2) Int(section, option string, def int) (int, error) {
	value, err := cf.conf.GetInt(section, option)
	if err != nil {
		if e, ok := err.(config.GetError); !ok || e.Reason != config.OptionNotFound {
			return 0, err
		}
		// option not found, find common.
		value, err = cf.conf.GetInt("common", option)
		if err != nil {
			if e, ok := err.(config.GetError); !ok || e.Reason != config.OptionNotFound {
				return 0, err
			}
			value = def
		}
	}
	log.Info("CONF INFO, SECTION: %s, %s = %d", section, option, value)
	return value, nil
}

func (cf *configFile2) String(section, option string, def string) (string, error) {
	value, err := cf.conf.GetString(section, option)
	if err != nil {
		if e, ok := err.(config.GetError); !ok || e.Reason != config.OptionNotFound {
			return "", err
		}
		// option not found, find common.
		value, err = cf.conf.GetString("common", option)
		if err != nil {
			if e, ok := err.(config.GetError); !ok || e.Reason != config.OptionNotFound {
				return "", err
			}
			value = def
		}
	}
	log.Info("CONF INFO, SECTION: %s, %s = %s", section, option, value)
	return value, nil
}
