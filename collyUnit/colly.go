package collyUnit

import (
	"cloud.google.com/go/storage"
	"genosha/utils/myLogger"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
)

func CollyRun() {
	start := time.Now()
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{DomainGlob:  "*.douban.*", Parallelism: 5})

	//// create the redis storage
	//storage := &redisstorage.Storage{
	//	Address:  "127.0.0.1:6379",
	//	Password: "",
	//	DB:       0,
	//	Prefix:   "douban250",
	//}
	//
	//// add storage to the collector
	//err := c.SetStorage(storage)
	//if err != nil {
	//	panic(err)
	//}
	//// delete previous data from storage
	//if err := storage.Clear(); err != nil {
	//	myLogger.Log.Info("error", zap.Any("error", err))
	//}

	// close redis client
	defer storage.Client.Close()

	c.OnError(func(_ *colly.Response, err error) {
		myLogger.Log.Info("error", zap.Any("error", err))
	})
	c.OnRequest(func(r *colly.Request) {
		myLogger.Log.Info("Visiting", zap.Any("url", r.URL))
		//r.Headers.Set("User-Agent", RandomString())
	})

	// On every a element which has href attribute call callback
	c.OnHTML(".hd", func(e *colly.HTMLElement) {
		movieUrl:=e.ChildAttr("a", "href")
		id:= strings.Split(movieUrl, "/")[4]
		title:=strings.TrimSpace(e.DOM.Find("span.title").Eq(0).Text())
		myLogger.Log.Info("id",zap.String("doubanID",id))
		myLogger.Log.Info("title",zap.String("title",title))
	})
	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnScraped(func(r *colly.Response) {
		myLogger.Log.Info("finished")
	})
	c.Visit("https://movie.douban.com/top250?start=0&filter=")

	c.Wait()
	took := time.Since(start)
	myLogger.Log.Info("Took",zap.Any("Took",took))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
