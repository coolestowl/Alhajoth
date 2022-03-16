package weibo_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/coolestowl/Alhajoth/sites"
	"github.com/coolestowl/Alhajoth/sites/weibo"
)

var _ = Describe("Weibo", func() {
	var wb weibo.Weibo

	BeforeEach(func() {
		wb = weibo.New("SUB=_2A25PKe7SDeRhGeNO71IY9ifMzD6IHXVs1fKarDV6PUJbkdANLUfwkW1NTuu3_Q-o3HVf-ErMRxF09AX_pAZLn1VV;")
	})

	Describe("Collecting friends feed", func() {
		var (
			items []sites.Item
			err   error
		)

		BeforeEach(func() {
			items, err = wb.FriendsFeed()
		})

		Context("err check", func() {
			It("should be nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("items length", func() {
			It("should have data", func() {
				Expect(len(items)).NotTo(BeZero())
			})
		})
	})

	Describe("Testing ContainerID", func() {
		var (
			id  string
			err error
		)

		BeforeEach(func() {
			id, err = wb.ContainerID("5582985423", "weibo")
		})

		Context("err check", func() {
			It("should be nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("id check", func() {
			It("should be correct id", func() {
				Expect(id).To(Equal("1076035582985423"))
			})
		})
	})

	Describe("Collecting friends feed", func() {
		var (
			items []sites.Item
			err   error
		)

		BeforeEach(func() {
			items, err = wb.UserFeed("5582985423", "1076035582985423")
		})

		Context("err check", func() {
			It("should be nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("items length", func() {
			It("should have data", func() {
				Expect(len(items)).NotTo(BeZero())
			})
		})
	})
})
