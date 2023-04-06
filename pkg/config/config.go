package config

import (
	"coffeeshop/pkg/util"
	"flag"
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
)

var k = koanf.New(".")

type Config struct {
	Shop     ShopCfg                `koanf:"shop"`
	Beans    map[string]BeanCfg     `koanf:"beans"`
	Brewers  map[string]*BrewerCfg  `koanf:"brewers"`
	Grinders map[string]*GrinderCfg `koanf:"grinders"`
}

type ShopCfg struct {
	CashRegisterTimeMS int `koanf:"cashRegisterTimeMS"`
	BaristaCount       int `koanf:"baristaCount"`
	CustomerCount      int `koanf:"customerCount"`
	OrderPipeDepth     int `koanf:"orderPipeDepth"`
}

type BeanCfg struct {
	BeanType string `koanf:"type"`
}

type GrinderCfg struct {
	BeanId              string `koanf:"beanId"`
	BeanCfg             *BeanCfg
	GrindGramsPerSecond int `koanf:"grindGramsPerSecond"`
	AddGramsPerSecond   int `koanf:"addGramsPerSecond"`
	HopperSize          int `koanf:"hopperSize"`
	RefillPercentage    int `koanf:"refillPercentage"`
}

type BrewerCfg struct {
	OuncesPerSecond int `koanf:"ouncesPerSecond"`
}

// BeanTypes helper
func (c *Config) BeanTypes() map[string]bool {
	beanTypes := map[string]bool{}
	for _, bt := range c.Beans {
		beanTypes[bt.BeanType] = true
	}
	return beanTypes
}

// Load uses this lib to get config: https://github.com/knadh/koanf#api
func Load(filename string) (*Config, error) {
	log := util.NewLogger("Config")

	f := flag.NewFlagSet("config", flag.ExitOnError)
	confFile := f.String("conf", filename, "path to a .yaml config file")
	custCount := f.Int("customers", -1, "number of customers. -1 means use config file")
	baristaCount := f.Int("baristas", -1, "number of baristas. -1 means use config file")
	if err := f.Parse(os.Args[1:]); err != nil {
		f.Usage()
		log.Infof("fuck")
		return nil, err
	}

	if err := k.Load(file.Provider(*confFile), yaml.Parser()); err != nil {
		log.Errorf("error loading config: %v", err)
		return nil, err
	}

	var out Config
	if err := k.Unmarshal("", &out); err != nil {
		log.Errorf("error loading config: %v", err)
		return nil, err
	}

	if *custCount >= 0 {
		out.Shop.CustomerCount = *custCount
	}
	if *baristaCount >= 0 {
		out.Shop.BaristaCount = *baristaCount
	}

	for _, g := range out.Grinders {
		b, have := out.Beans[g.BeanId]
		if have {
			g.BeanCfg = &b
		} else {
			return nil, fmt.Errorf("grinder wants unknown beanId %v", g.BeanId)
		}
	}

	return &out, nil
}
