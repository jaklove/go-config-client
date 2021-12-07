package php

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go-php-config/pkg/config"
	"go-php-config/pkg/filehelper"
	"go-php-config/pkg/middleware"
)

var oldFileInfoMap map[string]interface{}
var eventMagMap  map[string]interface{}
var writeBytes []byte


var PhpDeleteEventChan chan string

var PhpPutEventChan chan string

//监听文件内容
func InitPhpConfigInfo() {

	PhpDeleteEventChan = make(chan string,5)
	PhpPutEventChan   = make(chan string,5)

	go WatchChan()

	var watchParamsChan clientv3.WatchChan
	watchParamsChan = middleware.EctdClient.Watch(context.TODO(), config.PhpCfg.PhpParamsKey, clientv3.WithPrefix())

	//监听etcd变动
	for wresp := range watchParamsChan {
		for _, ev := range wresp.Events {
			fmt.Println("ev.Kv.Value:", ev.Kv.Value)
			switch ev.Type {
			case mvccpb.DELETE:
				PhpDeleteEventChan <- string(ev.Kv.Value)
			case mvccpb.PUT:
				fmt.Println("put操作")
				fmt.Println(string(ev.Kv.Value))
				PhpPutEventChan <- string(ev.Kv.Value)
			}
		}
	}
}

func WatchChan() {
	for {
		select {
		case msg := <- PhpDeleteEventChan:
		    //读取文件内容
			fileBytes, _:= filehelper.ReadContentByFilePath(config.PhpCfgPath.Path)
			json.Unmarshal(fileBytes,&oldFileInfoMap)
			json.Unmarshal([]byte(msg),&eventMagMap)

			fmt.Println("oldFileInfoMap:",oldFileInfoMap)
			fmt.Println("eventMagMap:",eventMagMap)

			//删除文件读取的内容
			for key,_ := range oldFileInfoMap{
				if _,ok := eventMagMap[key];ok{
					delete(oldFileInfoMap,key)
				}
			}

			//写入文件内容
			writeBytes,_ = json.Marshal(oldFileInfoMap)
			filehelper.WriteBytesToFile(config.PhpCfgPath.Path,writeBytes)

		case msg := <- PhpPutEventChan:
			//读取文件内容
			fileBytes, _:= filehelper.ReadContentByFilePath(config.PhpCfgPath.Path)
			json.Unmarshal(fileBytes,&oldFileInfoMap)
			json.Unmarshal([]byte(msg),&eventMagMap)
			//fmt.Println("oldFileInfoMap:",oldFileInfoMap)
			//fmt.Println("eventMagMap:",eventMagMap)

			//不喜欢
			for fileKv,_ := range oldFileInfoMap{
				_,ok := eventMagMap[fileKv]
				if ok{
					oldFileInfoMap[fileKv] = eventMagMap[fileKv]
				}else {
					//赋值到fileMap
					for minfoKey,minfoValue := range eventMagMap{
						oldFileInfoMap[minfoKey] = minfoValue
					}
				}
			}

			//写入文件内容
			writeBytes,_ = json.Marshal(oldFileInfoMap)
			filehelper.WriteBytesToFile(config.PhpCfgPath.Path,writeBytes)
		}
	}
}
