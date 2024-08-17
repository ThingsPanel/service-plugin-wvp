package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	httpclient "plugin_wvp/http_client"
	"plugin_wvp/mqtt"
)

type ChirpStackService struct {
	mux *http.ServeMux
}

func NewChirpStack() *ChirpStackService {
	return &ChirpStackService{
		mux: http.NewServeMux(),
	}
}

func (ctw *ChirpStackService) Init() *http.ServeMux {
	ctw.mux.HandleFunc("/accept/telemetry", ctw.telemetry)
	return ctw.mux
}

func (ctw *ChirpStackService) telemetry(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	var msg ChirpStackMessage
	err := decoder.Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debug("telemetry:", msg)
	if msg.Data == "" {
		return
	}
	deviceNumber := fmt.Sprintf(viper.GetString("chirp_stack.device_number_key"), msg.DeviceInfo.DevEui)
	//deviceNumber := msg.DeviceInfo.DevEui
	// 读取设备信息
	deviceInfo, err := httpclient.GetDeviceConfig(deviceNumber)
	if err != nil || deviceInfo.Code != 200 {
		// 获取设备信息失败，请检查连接包是否正确
		logrus.Error(err)
		return
	}
	logrus.Debug("deviceInfo:", deviceInfo)
	payload, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		//base64解密失败
		logrus.Error(err)
		return
	}
	telemetry := make(map[string]interface{})
	err = json.Unmarshal(payload, &telemetry)
	if err != nil {
		// 数据客户转换失败
		logrus.Error(err)
		return
	}
	logrus.Debug("telemetry:", telemetry)
	err = mqtt.PublishTelemetry(deviceInfo.Data.ID, telemetry)
	if err != nil {
		logrus.Error(err)
	}
}

type ChirpStackMessage struct {
	DeviceInfo ChirpStackDeviceInfo `json:"deviceInfo"`
	Data       string               `json:"data"`
}

type ChirpStackDeviceInfo struct {
	DevEui        string `json:"devEui"`
	ApplicationId string `json:"applicationId"`
}
