package scraper

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"time"
)

type CategoryScraper struct {
	urls  []string
	Cards []ProductCard // пул для вставки в БД
}

func NewCategoryScraper(_urls []string) *CategoryScraper {
	return &CategoryScraper{urls: _urls}
}

func (c *CategoryScraper) process(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	for _, url := range c.urls {

		// задержка для совести
		rand.Seed(time.Now().UnixNano())
		randomSleep := 5 + rand.Intn(20)
		time.Sleep(time.Duration(randomSleep) * time.Second)
		g.Go(func() error {
			return c.task(ctx, "https://www.wildberries.ru"+url)
		})

	}
	err := g.Wait()
	if err != nil {
		log.Error().Err(err).Msg("Error while category scraping")
	}
}

//TODO добавить пул воркеров
func (c *CategoryScraper) task(ctx context.Context, url string) (err error) {

	// ждем загрузки страницы
	err = chromedp.Run(
		ctx,
		RunWithTimeOut(&ctx, 30, chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.Sleep(555 * time.Millisecond),
			chromedp.WaitVisible("div.catalog-page__content", chromedp.ByQuery),
		}),
	)
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		return
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
		return
	}
	fmt.Printf("Category %s\n", url)
	lenPrice := len(productPrice)
	for i, node := range productCardNodes {
		u := node.AttributeValue("href")
		//TODO тут может быть index out of range переделать на dom.RequestChildNodes от productCardNodes
		var price string
		if i < lenPrice {
			price = productPrice[i].Children[0].NodeValue
		}

		fmt.Printf("href = %s Price=%s\n", u, price)
	}

	return
}
