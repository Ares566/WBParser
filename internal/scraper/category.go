package scraper

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"math/rand"
	"time"
)

type CategoryScraper struct {
	urls  []string
	Сards []ProductCard // пул для вставки в БД
}

func NewCategoryScraper(_urls []string) *CategoryScraper {
	return &CategoryScraper{urls: _urls}
}

func (c *CategoryScraper) process(ctx context.Context) {

	for _, url := range c.urls {

		// задержка для совести
		rand.Seed(time.Now().UnixNano())
		randomSleep := 5 + rand.Intn(20)
		time.Sleep(time.Duration(randomSleep) * time.Second)

		go c.task(ctx, "https://www.wildberries.ru"+url)
	}
}

//TODO добавить пул воркеров
func (c *CategoryScraper) task(ctx context.Context, url string) {

	// ждем загрузки страницы
	err := chromedp.Run(
		ctx,
		RunWithTimeOut(&ctx, 30, chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.Sleep(555 * time.Millisecond),
			chromedp.WaitVisible("div.catalog-page__content", chromedp.ByQuery),
		}),
	)
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		fmt.Printf("%s", err)
	}

	var productCardNodes []*cdp.Node
	var productPrice []*cdp.Node
	err = chromedp.Run(
		ctx,
		RunWithTimeOut(&ctx, 20, chromedp.Tasks{
			chromedp.Nodes(`a.product-card__main`, &productCardNodes),
			chromedp.Nodes(`ins.price__lower-price`, &productPrice),
		}),
	)
	if err != nil { //&& err != context.Canceled && err != context.DeadlineExceeded
		fmt.Printf("%s", err)
	}
	fmt.Printf("Category %s\n", url)
	for i, node := range productCardNodes {
		u := node.AttributeValue("href")
		//TODO тут может быть index out of range переделать на dom.RequestChildNodes от productCardNodes
		price := productPrice[i].Children[0].NodeValue
		fmt.Printf("href = %s Price=%s\n", u, price)
	}

}
