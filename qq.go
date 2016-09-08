package m3u8

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/tiantour/gotour/lib/requests"
)

// M3u8 M3u8
func (q *qq) M3u8(videoURL string) (string, error) {
	vid := q.id(videoURL)
	vnum := strconv.Itoa(rand.Intn(1000000))
	vtime := strconv.FormatInt(int64(time.Now().Unix()), 10)
	return q.query(vnum, vtime, vid)
}

// id id
func (q *qq) id(url string) string {
	reg := regexp.MustCompile(`vid=(\w+)`)
	return reg.FindStringSubmatch(url)[1]
}

// query query
func (q *qq) query(vnum, vtime, vid string) (string, error) {
	requestURL, requestData, requestHeader := requests.Options()
	requestURL = "http://h5vv.video.qq.com/getinfo?callback=tvp_request_getinfo_callback" + vnum + "&platform=11001&charge=0&otype=json&sb=1&nocache=0&_rnd=" + vtime + "&vids=" + vid + "&defaultfmt=auto&sdtfrom=v3010"
	body, err := requests.Get(requestURL, requestData, requestHeader)
	if err != nil {
		return "", err
	}
	temp := body[35 : len(body)-1]
	result := map[string]interface{}{}
	err = json.Unmarshal(temp, &result)
	if err != nil {
		return "", err
	}
	videoVl := result["vl"].(map[string]interface{})
	videoVi := videoVl["vi"].([]interface{})
	videoTemp := videoVi[0].(map[string]interface{})
	videoBr := videoTemp["br"].(float64)
	videoName := videoTemp["fn"].(string)
	videoKey := videoTemp["fvkey"].(string)
	videoLevel := videoTemp["level"].(float64)
	videoHost := videoTemp["ul"].(map[string]interface{})["ui"].([]interface{})[0].(map[string]interface{})["url"].(string)
	videoM3u8 := videoHost + videoName + "?vkey=" + videoKey + "&br=" + strconv.FormatFloat(videoBr, 'f', -1, 64) + "&platform=2&fmt=auto&level=" + strconv.FormatFloat(videoLevel, 'f', -1, 64) + "&sdtfrom=v3010"
	return videoM3u8, nil
}
