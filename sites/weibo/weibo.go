package weibo

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/coolestowl/Alhajoth/sites"
)

type Weibo interface {
	FriendsFeed() ([]sites.Item, error)

	UserFeed(uid, containerId string) ([]sites.Item, error)

	ContainerID(uid, title string) (string, error)
}

type impl struct {
	cookie string
}

func New(cookie string) Weibo {
	return &impl{
		cookie: cookie,
	}
}

func (i *impl) newReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("Cookie", i.cookie)
	return req, nil
}

func (i *impl) fetchJSON(req *http.Request, ret interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, ret)
}
