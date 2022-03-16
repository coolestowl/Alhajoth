package weibo_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/coolestowl/Alhajoth/sites/weibo"
)

const (
	UserCookie = "SUB=_2A25PKe7S;"
)

func demo() weibo.Weibo {
	return weibo.New(UserCookie)
}

func TestFriendsFeed(t *testing.T) {
	items, err := demo().FriendsFeed()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 20 {
		t.Error(len(items))
	}

	for _, item := range items {
		fmt.Println(item.ID(), item.Name(), item.CreatedAt())
	}
}

func TestContainerID(t *testing.T) {
	got, err := demo().ContainerID("5582985423", "weibo")
	if err != nil {
		t.Error(err)
	}

	if got != "1076035582985423" {
		t.Error(got)
	}
}

func TestUserFeed(t *testing.T) {
	items, err := demo().UserFeed("5582985423", "1076035582985423")
	if err != nil {
		t.Error(err)
	}

	for _, item := range items {
		log.Println(item.ID(), item.Name(), item.CreatedAt(), item.Text())
	}
}
