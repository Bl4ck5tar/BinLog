package initialize

import (
	"BinLog/server/global"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

//ConnectEs 初始化并返回一个配置好的 Elasticsearch 客户端
func ConnectEs() *elasticsearch.TypedClient {
	esConfig := global.Config.ES
	cfg := elasticsearch.Config{
		Addresses:		[]string{esConfig.URL},
		Username: 		esConfig.Username,
		Password: 		esConfig.Password,
	}

	//如果配置中指定了需要打印日志到控制台，则启动日志打印
	if esConfig.IsConsolePrint {
		cfg.Logger = &elastictransport.ColorLogger{
			Output: 			os.Stdout,	//请求日志输出到标准输出
			EnableRequestBody: 	true,		//启用请求体打印
			EnableResponseBody: true,		//启用响应体打印
		}
	}

	//创建一个新的Elasticsearch 客户端
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		global.Log.Error("Failed to connect to Elasticsearch", zap.Error(err))
	}

	return client
}