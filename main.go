package main

import (
	"go-php-config/list/phpParams"
	"go-php-config/pkg/config"
	"go-php-config/pkg/middleware"
	"go-php-config/watch/php"
	"log"
)

func init(){
	var err error

	//初始化配置信息
	err = config.InitConfig()
	if err != nil{
		log.Fatalf("config.InitConfig err:%s\n",err)
	}

	//初始化etcd
	err = middleware.InitEtcd()
	if err != nil{
		log.Fatalf("middleware.InitEtcd err:%s\n",err)
	}

	//获取对应的配置文件信息
	err = phpParams.InitPhpParams()
	if err != nil{
		log.Fatalf("phpParams.InitPhpParams err:%s\n",err)
	}
}

func main()  {
	//监听文件改动
	php.InitPhpConfigInfo()
}



