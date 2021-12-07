package phpParams

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go-php-config/pkg/config"
	"go-php-config/pkg/filehelper"
	"go-php-config/pkg/middleware"
	"log"
	"os"
)

var err error
var getResponse *clientv3.GetResponse
var phpMapParams map[string]interface{}

//初始化获取对应的配置信息
func InitPhpParams()error  {
	err = GetParamsList()
	if err != nil{
		return err
	}
	return nil;
}

//获取当前php配置信息
func GetParamsList() error {
	fmt.Println("params key:",config.PhpCfg.PhpParamsKey)

	//获取当前目录下对应的文件信息
	getResponse, err = middleware.EctdClient.Get(context.TODO(), config.PhpCfg.PhpParamsKey, clientv3.WithPrefix())
	if err != nil{
		return err
	}

	phpMapParams = make(map[string]interface{})
	//获取当前目录下的值
	for _,kv := range getResponse.Kvs{
		//获取里面的json内容
		tempMap := make(map[string]interface{})
		fmt.Println("list:",string(kv.Value))
		err = json.Unmarshal(kv.Value, &tempMap)
		if err != nil{
			log.Printf("json.Unmarshal err:%s",err)
			os.Exit(1)
		}

		for k,v := range tempMap{
			phpMapParams[k] = v
		}
	}

	return WriteParamsDataToFile(phpMapParams)
}

//写入文件内容
func WriteParamsDataToFile(phpMapParams map[string]interface{})error  {
	_, err := os.Stat(config.PhpCfgPath.Path)
	if err == nil{
		fileJsonStrByte,err := json.Marshal(phpMapParams)
		if err != nil{
			return err
		}

		return filehelper.WriteBytesToFile(config.PhpCfgPath.Path,fileJsonStrByte)
	}

	//创建文件
	return filehelper.CrateFileByPath(config.PhpCfgPath.Path)
}
