package config

import "github.com/go-ini/ini"

type EtcdConfig struct {
	Address string `ini:"address"`
}

type PhpConfig struct {
	PhpParamsKey string `ini:"phpParamsKey"`
}

type PhpConfigPath struct {
	Path string `ini:"path"`
}

var EtcdCfg    = &EtcdConfig{}
var PhpCfg     = &PhpConfig{}
var PhpCfgPath = &PhpConfigPath{}

func InitConfig()error{
	cfg, err := ini.Load("./conf/conf.ini")
	if err != nil{
		return err
	}
	cfg.Section("etcd").MapTo(EtcdCfg)
	cfg.Section("php-params").MapTo(PhpCfg)
	cfg.Section("php-params-file").MapTo(PhpCfgPath)
	return nil
}
