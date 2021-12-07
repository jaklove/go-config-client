package filehelper

import (
	"bufio"
	"io/ioutil"
	"os"
)

//创建文件
func CrateFileByPath(path string)error  {
	var err error
	var file *os.File
	file,err = os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}


//读取文件内容
func ReadContentByFilePath(path string)([]byte,error)  {
	file, err := os.Open(path)
	if err != nil {
		return nil,err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil{
		return nil,err
	}
	return bytes,nil
}

//写入内容到文件
func WriteBytesToFile(path string,bytes []byte)error  {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil{
		return err
	}

	defer file.Close()

	// 写入时 使用带缓存的 *Writer
	writer := bufio.NewWriter(file)
	writer.WriteString(string(bytes))
	writer.Flush()
	return nil
}
