package m3u8

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/tiantour/requests"
)

// M3u8 M3u8
func (a *acfan) M3u8(videoURL string) (string, error) {
	vid := a.id(videoURL)
	return a.query(vid)
}

// id id
func (a *acfan) id(url string) string {
	reg := regexp.MustCompile(`/ac(\d+)`)
	return reg.FindStringSubmatch(url)[1]
}

// query query
func (a *acfan) query(vid string) (string, error) {
	requestURL, requestData, requestHeader := requests.Options()
	requestURL = "http://api.aixifan.com/contents/" + vid
	requestHeader.Add("Host", "api.aixifan.com")
	requestHeader.Add("Origin", "http://m.acfun.tv")
	requestHeader.Add("Referer", "http://m.acfun.tv/v/?ac=2677500")
	requestHeader.Add("DeviceType", "2")
	body, err := requests.Get(requestURL, requestData, requestHeader)
	if err != nil {
		return "", err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	videoData := result["data"].(map[string]interface{})
	videoFiles := videoData["videos"].([]interface{})
	videoAddress := videoFiles[0].(map[string]interface{})
	if len(videoFiles) > 2 {
		videoAddress = videoFiles[2].(map[string]interface{})
	}
	videoID := videoAddress["videoId"].(float64)
	requestURL = "http://api.aixifan.com/plays/" + strconv.FormatFloat(videoID, 'f', -1, 64)
	body, err = requests.Get(requestURL, requestData, requestHeader)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	sourceData := result["data"].(map[string]interface{})
	sourceID := sourceData["sourceId"].(string)
	videoM3u8 := sourceID
	return videoM3u8, nil
}
