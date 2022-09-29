package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/app/ifce"
)

/*
	全局配置
*/

type Global struct {
	TcpServer      ifce.Isv // 当前Zinx全局的Server对象
	Host           string   // 当前服务器主机监听的IP
	TcpPort        string   // 当前服务器主机监听的端口号
	Name           string   // 当前服务器的名称
	Version        string   // 版本
	MaxConn        int      // 最大链接数
	MaxPackageSize uint32   // 最大数据包的最大值
}

var GlobalObject *Global

/*
	提供初始化方法
*/
func init() {
	GlobalObject = &Global{
		Host:           "127.0.0.1",
		TcpPort:        "3010",
		Name:           "Demo Zinx",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}

func (g *Global) Reload() {
	file, err := ioutil.ReadFile("demo/v0.1/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
