package weibo

import (
	"fmt"

	"github.com/coolestowl/Alhajoth/sites"
)

const (
	APIContainerIndex = "https://m.weibo.cn/api/container/getIndex"
)

type UserInfoResp struct {
	Ok   int `json:"ok"`
	Data struct {
		TabInfo struct {
			Tabs []struct {
				ContainerId string `json:"containerid"`
				TabType     string `json:"tab_type"`
			} `json:"tabs"`
		} `json:"tabsInfo"`
	} `json:"data"`
}

func (i *impl) ContainerID(uid, tabType string) (string, error) {
	url := fmt.Sprintf("%s?type=%s&value=%s", APIContainerIndex, "uid", uid)

	req, err := i.newReq(url)
	if err != nil {
		return "", err
	}

	var rsp UserInfoResp
	if err = i.fetchJSON(req, &rsp); err != nil {
		return "", err
	}

	if rsp.Ok != 1 {
		return "", fmt.Errorf("got err, rsp: %v", rsp)
	}

	for _, tab := range rsp.Data.TabInfo.Tabs {
		if tab.TabType == tabType {
			return tab.ContainerId, nil
		}
	}
	return "", nil
}

type UserStatusesResp struct {
	Ok   int `json:"ok"`
	Data struct {
		Scheme string `json:"scheme"`
		Cards  []struct {
			CardType int         `json:"card_type"`
			Mblog    weiboStatus `json:"mblog"`
		} `json:"cards"`
	} `json:"data"`
}

func (i *impl) UserFeed(uid, containerId string) ([]sites.Item, error) {
	url := fmt.Sprintf("%s?type=%s&value=%s&containerid=%s", APIContainerIndex, "uid", uid, containerId)

	req, err := i.newReq(url)
	if err != nil {
		return nil, err
	}

	var rsp UserStatusesResp
	if err = i.fetchJSON(req, &rsp); err != nil {
		return nil, err
	}

	if rsp.Ok != 1 {
		return nil, fmt.Errorf("got err, rsp: %v", rsp)
	}

	items := make([]sites.Item, 0, len(rsp.Data.Cards))
	for idx := range rsp.Data.Cards {
		if rsp.Data.Cards[idx].CardType != 9 {
			continue
		}
		items = append(items, &rsp.Data.Cards[idx].Mblog)
	}
	return items, nil
}
