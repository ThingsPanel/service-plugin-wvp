package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"plugin_wvp/apis"
	"plugin_wvp/cache"
	httpclient "plugin_wvp/http_client"
	"plugin_wvp/model"
	"plugin_wvp/mqtt"
)

type WvpService struct {
}

func NewWvpService() *WvpService {
	return &WvpService{}
}

func (w *WvpService) DeviceMqttPublish() {
	ctx := context.Background()
	keys, err := cache.GetWvpConfigKey(ctx)
	logrus.Debug(keys, err)
	for _, v := range keys {
		con, err2 := cache.GetWvpConfig(ctx, v)
		if err2 != nil {
			continue
		}
		w.mqttPublish(ctx, con)
		logrus.Debug(con, err)
	}
}

func (w *WvpService) mqttPublish(ctx context.Context, config model.WvpForm) {
	api := apis.NewWvpApi(config)
	resp, err := api.GetDeviceList(ctx, "1", "100")
	if err != nil || resp.Code != 0 {
		return
	}
	for _, v := range resp.Data.List {
		deviceNumber := fmt.Sprintf(viper.GetString("wvp.device_number_key"), v.DeviceId)
		// 读取设备信息
		deviceInfo, err1 := httpclient.GetDeviceConfig(deviceNumber)
		if err1 != nil || deviceInfo.Code != 200 || deviceInfo.Data.ID == "" {
			continue
		}
		//if v.OnLine {
		err = mqtt.DeviceStatusUpdate(deviceInfo.Data.ID, 1)
		if err != nil {
			logrus.Debug(err)
		}
		payload, err2 := api.GetDeviceChannels(ctx, v.DeviceId)
		if err != nil {
			logrus.Debug(err2)
			continue
		}
		//err = mqtt.PublishTelemetry(deviceInfo.Data.ID, payload)
		err = mqtt.PublishAttributes(deviceInfo.Data.ID, payload)
		if err != nil {
			logrus.Debug(err)
		}
		//} else {
		//	err = mqtt.DeviceStatusUpdate(deviceInfo.Data.ID, 0)
		//	if err != nil {
		//		logrus.Debug(err)
		//	}
		//}
	}
}
