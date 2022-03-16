package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/coolestowl/Alhajoth/sender"
	"github.com/coolestowl/Alhajoth/sender/cq"
	"github.com/coolestowl/Alhajoth/sites"
	"github.com/coolestowl/Alhajoth/sites/weibo"

	"github.com/go-redis/redis/v8"
)

const (
	RedisKeyWeiboIDRecord = "__weibo_id_record__"
	RedisKeyMonitorUID    = "__weibo_uids__"
	RedisKeyUID2Container = "__weibo_uid_%s__"
)

type scheduler struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	rds *redis.Client

	errChan  chan error
	itemChan chan sites.Item

	sender sender.Sender
}

func (s *scheduler) loggerRoutine() {
	s.wg.Add(1)
	defer s.wg.Done()

loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		case err := <-s.errChan:
			log.Println("[err]", err)
		}
	}
}

func (s *scheduler) senderRoutine() {
	s.wg.Add(1)
	defer s.wg.Done()

loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		case item := <-s.itemChan:
			addNum := s.rds.SAdd(context.Background(), RedisKeyWeiboIDRecord, item.ID()).Val()
			if addNum > 0 && len(item.Images()) > 0 {
				if err := s.sender.Send(item); err != nil {
					s.errChan <- err
				}
			}
		}
	}
}

func (s *scheduler) spiderRoutine(w weibo.Weibo, tk *time.Ticker) {
	s.wg.Add(1)
	defer s.wg.Done()

	uid2ContainerID := func(uid string) (string, error) {
		id := s.rds.Get(context.Background(), fmt.Sprintf(RedisKeyUID2Container, uid)).Val()
		if len(id) > 0 {
			return id, nil
		}

		return w.ContainerID(uid, "weibo")
	}

	task := func() {
		uids := s.rds.SMembers(context.Background(), RedisKeyMonitorUID).Val()
		for _, uid := range uids {
			containerId, err := uid2ContainerID(uid)
			if err != nil {
				s.errChan <- err
				continue
			}

			items, err := w.UserFeed(uid, containerId)
			if err != nil {
				s.errChan <- err
				continue
			}

			go func() {
				for _, item := range items {
					s.itemChan <- item
				}
			}()
		}

		items, err := w.FriendsFeed()
		if err != nil {
			s.errChan <- err
		} else {
			for _, item := range items {
				s.itemChan <- item
			}
		}
	}

loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		case <-tk.C:
			task()
		}
	}
}

func main() {
	var (
		redisAddr   = os.Getenv("REDIS_ADDR")             // 0.0.0.0:6379
		cqGroup, _  = strconv.Atoi(os.Getenv("GROUP_ID")) // 1234567
		cqAPI       = os.Getenv("CQ_API")                 // https://xxx.xxx.xxx/xxx/xxx
		weiboCookie = os.Getenv("COOKIE")                 // SUB=_Ke7SDeRhGeNO7.....;
		tickSecs, _ = strconv.Atoi(os.Getenv("TICK"))     // 300

		wb          = weibo.New(weiboCookie)
		tk          = time.NewTicker(time.Second * time.Duration(tickSecs))
		ctx, cancel = context.WithCancel(context.Background())
	)

	s := &scheduler{
		ctx:    ctx,
		cancel: cancel,
		wg:     sync.WaitGroup{},

		rds: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),

		errChan:  make(chan error, 1),
		itemChan: make(chan sites.Item, 1),

		sender: cq.New(cqAPI, cqGroup),
	}

	go s.loggerRoutine()
	go s.senderRoutine()
	go s.spiderRoutine(wb, tk)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	cancel()
	s.wg.Wait()
}
