package middleware

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go-php-config/pkg/config"
	"log"
	"time"
)

var (
	EctdClient *clientv3.Client
)

func InitEtcd()error  {
	var err error
	EctdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{config.EtcdCfg.Address},
		DialTimeout: time.Second * 5,
	})
	if err != nil{
		log.Printf("InitEtcd err:%s",err)
	}

	fmt.Println("etcd connect success!")
	return nil
}