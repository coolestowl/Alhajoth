package cq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coolestowl/Alhajoth/sender"
	"github.com/coolestowl/Alhajoth/sites"
)

type cq struct {
	groupID int
	api     string
}

func New(api string, group int) sender.Sender {
	return &cq{
		api:     api,
		groupID: group,
	}
}

type cqNode struct {
	Type string `json:"type"`
	Data struct {
		Name    string `json:"name"`
		UIN     string `json:"uin"`
		Content string `json:"content"`
	} `json:"data"`
}

func infoNode(name, text string, tm time.Time) cqNode {
	node := cqNode{}
	node.Type = "node"
	node.Data.Name = "bot"
	node.Data.UIN = "1234567"
	node.Data.Content = fmt.Sprintf("%s\n\n%s\n\n%s", name, tm.Format("2006-01-02 15:04:05"), text)
	return node
}

func imgNode(url string) cqNode {
	node := cqNode{}
	node.Type = "node"
	node.Data.Name = "bot"
	node.Data.UIN = "1234567"
	node.Data.Content = fmt.Sprintf("[CQ:image,file=c5930b6d2299e650ca0d28,url=%s,cache=0]", url)
	return node
}

type cqGroupForwardMsg struct {
	GroupID  int64    `json:"group_id"`
	Messages []cqNode `json:"messages"`
}

func (c *cq) payload(name, text string, tm time.Time, urls []string) []byte {
	msg := cqGroupForwardMsg{
		GroupID:  int64(c.groupID),
		Messages: make([]cqNode, 0, len(urls)),
	}

	msg.Messages = append(msg.Messages, infoNode(name, text, tm))

	for _, url := range urls {
		msg.Messages = append(msg.Messages, imgNode(url))
	}

	data, _ := json.Marshal(msg)
	return data
}

func (c *cq) Send(item sites.Item) error {
	payload := c.payload(item.Name(), item.Text(), item.CreatedAt(), item.Images())

	_, err := http.Post(c.api, "application/json", bytes.NewBuffer(payload))

	return err
}
