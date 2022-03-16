package weibo

import (
	"regexp"
	"time"
)

type weiboStatus struct {
	IDStr   string `json:"id"`
	TextStr string `json:"text"`
	PicNum  int    `json:"pic_num"`
	Pics    []struct {
		PID   string `json:"pid"`
		Large struct {
			URL string `json:"url"`
		} `json:"large"`
	} `json:"pics"`
	CreatedAtStr string `json:"created_at"`

	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`

	Retweeted *weiboStatus `json:"retweeted_status"`
}

func (w *weiboStatus) ID() string {
	if w.Retweeted != nil {
		return w.Retweeted.IDStr
	}
	return w.IDStr
}

func (w *weiboStatus) CreatedAt() time.Time {
	tm, _ := time.Parse("Mon Jan 2 15:04:05 +0800 2006", w.CreatedAtStr)
	return tm
}

func (w *weiboStatus) Name() string {
	if w.Retweeted != nil {
		return w.Retweeted.User.ScreenName
	}
	return w.User.ScreenName
}

func (w *weiboStatus) Text() string {
	text := w.TextStr
	if w.Retweeted != nil {
		text = w.Retweeted.TextStr
	}

	return regexp.MustCompile(`<[^>]+>`).ReplaceAllString(text, "")
}

func (w *weiboStatus) Images() []string {
	if w.Retweeted != nil && w.Retweeted.PicNum > 0 {
		return w.Retweeted.Images()
	}

	if w.PicNum == 0 {
		return nil
	}

	images := make([]string, 0, w.PicNum)
	for _, i := range w.Pics {
		images = append(images, i.Large.URL)
	}

	return images
}
