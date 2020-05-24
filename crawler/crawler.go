package crawler

import (
	"genosha/utils/myLogger"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	_"strings"
	"sync"
	"time"
)

func fetch(url string) *goquery.Document {
	//myLogger.Log.Info("Fetch Url:" + url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		return nil
	}
	if resp.StatusCode != 200 {
		myLogger.Log.Info("Http status code:",zap.Int("StatusCode",resp.StatusCode))
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		myLogger.Log.Error("error", zap.Error(err))
		return nil
	}
	return doc
}

func parseUrls(url string) {
	doc := fetch(url)
	doc.Find("ol.grid_view li").Find(".hd").Each(func(index int, ele *goquery.Selection) {
		//movieUrl, _ := ele.Find("a").Attr("href")
		//id:=strings.Split(movieUrl, "/")[4]
		title:=ele.Find(".title").Eq(0).Text()
		//myLogger.Log.Info("id",zap.String("doubanID",id))
		myLogger.Log.Info("title",zap.String("title",title))
	})
}

func Douban250() {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			parseUrls("https://movie.douban.com/top250?start=" + strconv.Itoa(25*i))
		}(i)
	}
	wg.Wait()
	took := time.Since(start)
	myLogger.Log.Info("Took",zap.Any("Took",took))
}
