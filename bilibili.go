package m3u8

import (
	"encoding/json"
	"regexp"

	"github.com/tiantour/requests"
)

// M3u8 M3u8
func (b *bilibili) M3u8(videoURL string) (string, error) {
	vid := b.id(videoURL)
	return b.query(vid)
}

// id id
func (b *bilibili) id(url string) string {
	reg := regexp.MustCompile(`/av(\d+)/`)
	return reg.FindStringSubmatch(url)[1]
}

// query query
func (b *bilibili) query(vid string) (string, error) {
	requestURL, requestData, requestHeader := requests.Options()
	requestURL = "http://www.bilibili.com/m/html5?aid=" + vid + "&page=1"
	body, err := requests.Get(requestURL, requestData, requestHeader)
	if err != nil {
		return "", err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	videoM3u8 := result["src"].(string)
	return videoM3u8, nil
}
