package main

import (
	"WBParser/internal/scraper"
	"github.com/rs/zerolog/log"
)

func main() {

	categoryURLsParser := scraper.NewCategories()
	categoryURLs, err := categoryURLsParser.Parse()
	if err != nil {
		log.Error().Err(err).Msg("Failed categoryURLsParser")
		return
	}

	// парсим данные товара со сраницы категорий url
	ProductScrapper := scraper.NewCategoryScraper(categoryURLs)

	// парсим остатки на складе для товара заданного url
	//BasketScraper := scraper.NewBasketScraper(productURLs)

	scraper.NewClient(ProductScrapper)

}
