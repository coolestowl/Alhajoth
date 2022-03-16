package weibo

import (
	"fmt"
	"net/http"

	"github.com/coolestowl/Alhajoth/sites"
)

const (
	APIFeedFriends = "https://m.weibo.cn/feed/friends"
)

type FeedFriendsResp struct {
	Ok       int `json:"ok"`
	HttpCode int `json:"http_code"`
	Data     struct {
		Statuses []weiboStatus `json:"statuses"`
	} `json:"data"`
}

func (i *impl) FriendsFeed() ([]sites.Item, error) {
	req, err := i.newReq("https://m.weibo.cn/feed/friends")
	if err != nil {
		return nil, err
	}

	var rsp FeedFriendsResp
	if err = i.fetchJSON(req, &rsp); err != nil {
		return nil, err
	}

	if rsp.Ok != 1 || rsp.HttpCode != http.StatusOK {
		return nil, fmt.Errorf("got err, rsp: %v", rsp)
	}

	items := make([]sites.Item, 0, len(rsp.Data.Statuses))
	for idx := range rsp.Data.Statuses {
		items = append(items, &rsp.Data.Statuses[idx])
	}
	return items, nil
}
