package m3u8

// Acfan Acfan
var (
	Acfan    = &acfan{}
	Bilibili = &bilibili{}
	QQ       = &qq{}
	Youku    = &youku{}
)

// M3u8er M3u8er
type (
	M3u8er interface {
		M3u8(videoURL string) (videoM3u8 string, err error)
	}
	acfan    struct{}
	bilibili struct{}
	qq       struct{}
	youku    struct{}
)
