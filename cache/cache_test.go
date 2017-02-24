package cache

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCaching_NewGenericCache(t *testing.T) {

	Convey("Given I create a new memory cache", t, func() {

		c := NewMemoryCache()

		Convey("Then the cache should be initialized", func() {
			So(c.(*memoryCache).data, ShouldResemble, map[string]*cacheItem{})
			So(len(c.(*memoryCache).data), ShouldBeZeroValue)
		})

		Convey("When I Get a non cached object", func() {
			smt := c.Get("id-something")

			Convey("Then retrieved object should be nil", func() {
				So(smt, ShouldBeNil)
			})
		})

		Convey("When I GetReset a non cached object", func() {
			smt := c.GetReset("id-something")

			Convey("Then retrieved object should be nil", func() {
				So(smt, ShouldBeNil)
			})
		})

		Convey("When I Exists a non cached object", func() {
			ex := c.Exists("id-something")

			Convey("Then the result should be false", func() {
				So(ex, ShouldBeFalse)
			})
		})

		Convey("When I Del a non cached object", func() {

			Convey("It should not panic", func() {
				So(func() { c.Del("id-something") }, ShouldNotPanic)
			})
		})

		Convey("When I cache something", func() {

			something := &struct {
				Name        string
				Description string
			}{
				Name:        "something",
				Description: "that's something",
			}

			c.Set("id-something", something)

			Convey("Then the object should be cached", func() {
				So(len(c.(*memoryCache).data), ShouldEqual, 1)
			})

			Convey("When I All the cache", func() {
				ex := c.All()
				So(len(ex), ShouldEqual, 1)
			})

			Convey("When I Get a cached object", func() {
				smt := c.Get("id-something")

				Convey("Then retrieved object should be the same as the original", func() {
					So(smt, ShouldEqual, something)
				})
			})

			Convey("When I Exists a the cached object", func() {
				ex := c.Exists("id-something")

				Convey("Then the result should be true", func() {
					So(ex, ShouldBeTrue)
				})
			})

			Convey("When I Del a non cached object", func() {

				Convey("It should not panic", func() {
					So(func() { c.Del("id-something") }, ShouldNotPanic)

					Convey("When I Exists a the deleted cached object", func() {
						ex := c.Exists("id-something")

						Convey("Then the result should be false", func() {
							So(ex, ShouldBeFalse)
						})
					})
				})
			})
		})
	})
}

func TestCaching_NewGenericCacheWithDefaultExpiration(t *testing.T) {

	Convey("Given I create a new memory cache with default expiration", t, func() {

		c := NewMemoryCache()
		c.SetDefaultExpiration(2 * time.Second)

		Convey("Then the cache should be initialized", func() {
			So(c.(*memoryCache).data, ShouldResemble, map[string]*cacheItem{})
			So(len(c.(*memoryCache).data), ShouldBeZeroValue)
		})

		Convey("When I set an item that should expire by default", func() {
			c.Set("id-default", "item-default")
			Convey("When I wait for 1.5 seconds", func() {
				<-time.After(1500 * time.Millisecond)
				So(c.Get("id-default"), ShouldEqual, "item-default")

				Convey("When I wait for another 1 second", func() {
					<-time.After(1000 * time.Millisecond)
					Convey("Then the item should be gone", func() {
						So(c.Get("id-default"), ShouldEqual, nil)
					})
				})
			})
		})

		Convey("When I set an item that expired after 1sec", func() {
			c.SetWithExpiration("id", "item", 1*time.Second)

			Convey("Then the item should be present", func() {
				So(c.Get("id"), ShouldEqual, "item")

				Convey("When I wait for 1.5 seconds", func() {
					<-time.After(1500 * time.Millisecond)
					Convey("Then the item should be gone", func() {
						So(c.Get("id"), ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCaching_NewGenericCacheWithDefaultExpirationAndMultipleSets(t *testing.T) {

	Convey("Given I create a new memory cache with default expiration", t, func() {

		c := NewMemoryCache()
		c.SetDefaultExpiration(2 * time.Second)

		Convey("Then the cache should be initialized", func() {
			So(c.(*memoryCache).data, ShouldResemble, map[string]*cacheItem{})
			So(len(c.(*memoryCache).data), ShouldBeZeroValue)
		})

		Convey("When I set an item that expired after 1sec", func() {
			c.SetWithExpiration("id", "item", 1*time.Second)

			Convey("Then the item should be present", func() {
				So(c.Get("id"), ShouldEqual, "item")

				Convey("When I wait for 0.5 seconds", func() {
					<-time.After(500 * time.Millisecond)

					Convey("When I set the item again to expire after 1sec", func() {
						c.SetWithExpiration("id", "item", 1*time.Second)

						Convey("When I wait for 0.7 seconds", func() {
							<-time.After(700 * time.Millisecond)
							Convey("Then the item should be exists", func() {
								So(c.Get("id"), ShouldEqual, "item")
							})

							Convey("When I wait for 0.5 seconds", func() {
								<-time.After(500 * time.Millisecond)
								Convey("Then the item should be gone", func() {
									So(c.Get("id"), ShouldBeNil)
								})
							})
						})
					})
				})
			})
		})
	})
}

func TestCaching_NewGenericCacheWithDefaultNotifier(t *testing.T) {

	Convey("Given I create a new memory cache with default notifier", t, func() {

		expiredCalled := int32(0)
		cachedItem := "item"

		c := NewMemoryCache()
		c.SetDefaultExpirationNotifier(func(c Cacher, id string, item interface{}) {
			atomic.AddInt32(&expiredCalled, 1)
		})

		Convey("When I set an item that expires after 1sec", func() {
			c.SetWithExpiration("id", cachedItem, 1*time.Second)

			Convey("Then the item should be present", func() {
				So(c.Get("id"), ShouldEqual, "item")

				Convey("When I wait for 1.5 seconds", func() {
					<-time.After(1500 * time.Millisecond)

					Convey("Then the item should be gone", func() {
						So(c.Get("id"), ShouldBeNil)
					})

					Convey("Then the expiration notification should have been called 1 times", func() {
						So(atomic.LoadInt32(&expiredCalled), ShouldEqual, 1)
					})
				})
			})
		})
	})
}

func TestCaching_Expiration(t *testing.T) {

	Convey("Given I create a new memory cache", t, func() {

		c := NewMemoryCache()

		Convey("When I set an item that expired after 1sec", func() {
			c.SetWithExpiration("id", "item", 1*time.Second)

			Convey("Then the item should be present", func() {

				So(c.Get("id"), ShouldEqual, "item")

				Convey("When I wait for 1.5 seconds", func() {

					<-time.After(1500 * time.Millisecond)

					Convey("Then the item should be gone", func() {
						So(c.Get("id"), ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCaching_GetReset(t *testing.T) {

	Convey("Given I create a new memory cache", t, func() {

		c := NewMemoryCache()

		Convey("When I set an item that expired after 1sec", func() {

			c.SetWithExpiration("id", "item", 1*time.Second)

			Convey("Then the item should be present", func() {
				So(c.Get("id"), ShouldEqual, "item")

				Convey("When I wait for 0.5 seconds", func() {
					<-time.After(500 * time.Millisecond)

					Convey("When I GetReset the item", func() {
						So(c.GetReset("id"), ShouldEqual, "item")

						Convey("When I wait for 0.7 more seconds", func() {
							<-time.After(700 * time.Millisecond)

							Convey("Then the item should exist", func() {
								So(c.Get("id"), ShouldEqual, "item")
							})

							Convey("When I wait for 0.5 more seconds", func() {
								<-time.After(500 * time.Millisecond)

								Convey("Then the item should be gone", func() {
									So(c.Get("id"), ShouldBeNil)
								})
							})
						})
					})
				})
			})
		})
	})
}
