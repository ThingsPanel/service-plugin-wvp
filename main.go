package main

import (
	"log"
	"plugin_wvp/cache"
	httpclient "plugin_wvp/http_client"
	httpservice "plugin_wvp/http_service"
	"plugin_wvp/mqtt"
	"plugin_wvp/services"
	"strings"

	"github.com/spf13/viper"
)

func main() {

	conf()
	LogInIt()
	log.Println("Starting the application...")

	//apiToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjNlNzk2NDZjNGRiYzQwODM4M2E5ZWVkMDlmMmI4NWFlIn0.eyJqdGkiOiJOZzZPR3BPeG03SGhhbzl6NlJYMmtRIiwiaWF0IjoxNzIzNzI4MzYxLCJleHAiOjE4NjE4OTExNDUsIm5iZiI6MTcyMzcyODM2MSwic3ViIjoibG9naW4iLCJhdWQiOiJBdWRpZW5jZSIsInVzZXJOYW1lIjoiYWRtaW4iLCJhcGlLZXlJZCI6MX0.ZkYCr1ffxij8wBHMCSJIqkeYoSAwe2ko-vVsnwi60SVxPA6lTosrv3s9Rnl0y06ZefwKtivDpykL7n3jn0CiqN11O20vpz02x9LwJADKksUOgtWh86M0IiSUtVwuAHPl7rxBph7Juc1m-vVQVCXdzraoEMC5CZaBjJxNaUOLNG5we6vBzxGEVlxbbjZoBZb5qe2_0FrbMsV6vR-sv2Ib8JcKdSPNkOgQHmebKofJfF3fa9PFxv3daymtH81qLWsvqevYrRiZ9_oE5od7N6dccYe7ilN53xJND4Jweg7U47yoI7_dpKgTXNUQuMWo0aMtJl-0joVDZYNarC7uazVKSg"
	//host := "104.156.140.42"
	//client := apis.NewWvpApi(model.WvpForm{
	//	Server:   host,
	//	Port:     18080,
	//	ApiToken: apiToken,
	//})
	////resp, err := client.GetDeviceList(context.Background(), 1, 10)
	////resp := client.GetDeviceStatus(context.Background(), "44010200492000000001")
	//resp, _ := client.GetDeviceChannels(context.Background(), "44010200492000000001")
	//logrus.Debug(resp)
	//
	//return
	//初始化redis
	cache.RedisInit()
	// 启动mqtt客户端
	mqtt.InitClient()
	// 启动http客户端
	httpclient.Init()
	// 启动服务
	//go services.Start()
	go services.StartHttp(services.NewChirpStack().Init())

	// 启动http服务
	httpservice.Init()
	select {}
}
func conf() {
	log.Println("加载配置文件...")
	// 设置环境变量前缀
	viper.SetEnvPrefix("plugin_wvp")
	// 使 Viper 能够读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("加载配置文件完成...")
}
