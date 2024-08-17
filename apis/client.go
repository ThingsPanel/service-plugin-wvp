package apis

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"plugin_wvp/model"
	"time"
)

type WvpApi struct {
	Config model.WvpForm
}

func NewWvpApi(config model.WvpForm) *WvpApi {
	return &WvpApi{Config: config}
}

type WvpDeviceListRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total int             `json:"total"`
		List  []WvpDeviceItem `json:"list"`
	} `json:"data"`
}
type WvpDeviceItem struct {
	DeviceId string `json:"deviceId"`
	Name     string `json:"name"`
	OnLine   bool   `json:"onLine"` //ture在线 false 离线
}

// GetDeviceList
// @description 获取设备列表
func (w *WvpApi) GetDeviceList(ctx context.Context, page, pageSize string) (WvpDeviceListRes, error) {
	url := "/api/device/query/devices"
	resp, err := w.get(ctx, url, map[string]interface{}{
		"page":  page,
		"count": pageSize,
	})
	logrus.Debug(resp, err)
	var result WvpDeviceListRes
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return result, err
	}
	logrus.Debug(result, err)
	return result, nil
}

type DeviceStatusRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Result       string `json:"result"`
		Online       string `json:"online"`
		Status       string `json:"status"`
		DeviceStatus string `json:"deviceStatus"`
	}
}

func (w *WvpApi) GetDeviceStatus(ctx context.Context, deviceId string) bool {
	url := fmt.Sprintf("/api/device/query/devices/%s/status", deviceId)
	resp, err := w.get(ctx, url, nil)
	logrus.Debug(resp, err)
	var result DeviceStatusRes
	if err != nil {
		logrus.Debug(result)
		return false
	}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		logrus.Debug(result)
		return false
	}
	if result.Code == 0 && result.Data.DeviceStatus == "SN" {
		return true
	}
	return false
}

type DeviceChannelsRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total int `json:"total"`
		List  []struct {
			ChannelId string  `json:"channelId"`
			DeviceId  string  `json:"deviceId"`
			StreamId  *string `json:"streamId"`
		} `json:"list"`
	} `json:"data"`
}

func (w *WvpApi) GetDeviceChannels(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("/api/device/query/devices/%s/channels", deviceId)
	var (
		ret    DeviceChannelsRes
		result = make(map[string]interface{})
	)
	resp, err := w.get(ctx, url, map[string]interface{}{
		"page":  1,
		"count": 20,
		//"online": true,
	})
	err = json.Unmarshal([]byte(resp), &ret)
	if err != nil {
		logrus.Debug(ret)
		return result, nil
	}
	for _, v := range ret.Data.List {
		if v.StreamId != nil && *v.StreamId != "" {
			result = w.getPlayURLs(v.DeviceId, v.ChannelId)
		}
	}
	return result, nil
}

func (w *WvpApi) getPlayURLs(deviceID, channelID string) map[string]interface{} {
	stream := fmt.Sprintf("%s_%s", deviceID, channelID)
	return map[string]interface{}{
		"FLV":         fmt.Sprintf("http://%s:18001/rtp/%s.live.flv", w.Config.Server, stream),
		"FLV(ws)":     fmt.Sprintf("ws://%s:18001/rtp/%s.live.flv", w.Config.Server, stream),
		"FMP4":        fmt.Sprintf("http://%s:18001/rtp/%s.live.mp4", w.Config.Server, stream),
		"FMP4(https)": fmt.Sprintf("https://%s:18443/rtp/%s.live.mp4", w.Config.Server, stream),
		"FMP4(ws)":    fmt.Sprintf("ws://%s:18001/rtp/%s.live.mp4", w.Config.Server, stream),
		"FMP4(wss)":   fmt.Sprintf("wss://%s:18443/rtp/%s.live.mp4", w.Config.Server, stream),
		"HLS":         fmt.Sprintf("http://%s:18001/rtp/%s/hls.m3u8", w.Config.Server, stream),
		"HLS(https)":  fmt.Sprintf("https://%s:18443/rtp/%s/hls.m3u8", w.Config.Server, stream),
		"HLS(ws)":     fmt.Sprintf("ws://%s:18001/rtp/%s/hls.m3u8", w.Config.Server, stream),
		"HLS(wss)":    fmt.Sprintf("wss://%s:18443/rtp/%s/hls.m3u8", w.Config.Server, stream),
		"TS":          fmt.Sprintf("http://%s:18001/rtp/%s.live.ts", w.Config.Server, stream),
		"TS(https)":   fmt.Sprintf("https://%s:18443/rtp/%s.live.ts", w.Config.Server, stream),
		"TS(ws)":      fmt.Sprintf("ws://%s:18001/rtp/%s.live.ts", w.Config.Server, stream),
		"RTC":         fmt.Sprintf("http://%s:18001/index/api/webrtc?app=rtp&stream=%s&type=play", w.Config.Server, stream),
		"RTCS":        fmt.Sprintf("https://%s:18443/index/api/webrtc?app=rtp&stream=%s&type=play", w.Config.Server, stream),
		"RTMP":        fmt.Sprintf("rtmp://%s:1935/rtp/%s", w.Config.Server, stream),
		"RTMPS":       fmt.Sprintf("rtmps://%s:19350/rtp/%s", w.Config.Server, stream),
		"RTSP":        fmt.Sprintf("rtsp://%s:554/rtp/%s", w.Config.Server, stream),
		"RTSPS":       fmt.Sprintf("rtsps://%s:332/rtp/%s", w.Config.Server, stream),
	}
}

func (w *WvpApi) get(ctx context.Context, url string, params map[string]interface{}) (string, error) {
	url = fmt.Sprintf("http://%s:%d%s", w.Config.Server, w.Config.Port, url)
	var isSymbol bool
	for k, v := range params {
		var symbol = "?"
		if isSymbol {
			symbol = "&"
		}
		isSymbol = true
		url = url + symbol + fmt.Sprintf("%s=%v", k, v)
	}
	logrus.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("access-token", w.Config.ApiToken)
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("接口异常: status code %d, body: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
