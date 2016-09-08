package m3u8

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tiantour/requests"
)

func (y *youku) M3u8(videoURL string) (string, error) {
	vid := y.id(videoURL)
	return y.query(vid)
}

func (y *youku) id(url string) string {
	reg := regexp.MustCompile(`id_(\w+)`)
	return reg.FindStringSubmatch(url)[1]
}

func (y *youku) data(vid string) (int, string, error) {
	requestURL, requestData, requestHeader := requests.Options()
	requestURL = "http://play.youku.com/play/get.json?vid=" + vid + "&ct=12&Type=Folder&ob=1"
	requestHeader.Add("Host", "play.youku.com")
	requestHeader.Add("Referer", "http://v.youku.com/v_show/id_"+vid+".html")
	body, err := requests.Get(requestURL, requestData, requestHeader)
	if err != nil {
		return 0, "", err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, "", err
	}
	resultData := result["data"].(map[string]interface{})
	resultSecurity := resultData["security"].(map[string]interface{})
	resultIP := resultSecurity["ip"].(float64)
	ip := int(resultIP)
	ep := resultSecurity["encrypt_string"].(string)
	return ip, ep, nil
}

func (y *youku) token(ep string) (sid, token string) {
	code := "becaf9be"
	epOld, _ := base64.StdEncoding.DecodeString(ep)
	epNew := y.es(code, string(epOld))
	temp := strings.Split(string([]byte(epNew)), "_")
	sid, token = temp[0], temp[1]
	return
}

func (y *youku) ep(vid, sid, token string) (ep string) {
	code := "bf7e5f01"
	temp := y.es(code, sid+"_"+vid+"_"+token)
	ep = base64.StdEncoding.EncodeToString(temp)
	return
}

func (y *youku) query(vid string) (videoM3u8 string, err error) {
	ip, epOld, err := y.data(vid)
	if err != nil {
		return "", err
	}
	sid, token := y.token(epOld)
	epNew := y.ep(vid, sid, token)
	ts := time.Now().Unix()
	videoM3u8 = "http://pl.youku.com/playlist/m3u8?vid=" + vid + "&type=flv&ts=" + strconv.FormatInt(ts, 10) + "&keyframe=0&ep=" + url.QueryEscape(epNew) + "&sid=" + sid + "&token=" + token + "&ctype=12&ev=1&oip=" + strconv.Itoa(ip)
	return
}

func (y *youku) es(a string, c string) (result []byte) {
	b := [256]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255}
	f, h := 0, 0
	for h < 256 {
		f = (f + b[h] + int(a[h%len(a)])) % 256
		b[h], b[f] = b[f], b[h]
		h = h + 1
	}
	q, f, h := 0, 0, 0
	for q < len(c) {
		h = (h + 1) % 256
		f = (f + b[h]) % 256
		b[h], b[f] = b[f], b[h]
		r := int(c[q]) ^ b[(b[h]+b[f])%256]
		result = append(result, transformIntByte(r)[3])
		q = q + 1
	}
	return
}

//字节转数字
func transformByteInt(key []byte) (value int) {
	var temp int32
	bBuf := bytes.NewBuffer(key)
	binary.Read(bBuf, binary.BigEndian, &temp)
	value = int(temp)
	return
}

// 数字转字节
func transformIntByte(key int) (value []byte) {
	temp := rune(key)
	bBuf := bytes.NewBuffer([]byte{})
	binary.Write(bBuf, binary.BigEndian, temp)
	value = bBuf.Bytes()
	return
}
