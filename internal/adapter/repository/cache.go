package repository

import (
	"BlockchainCurrency/internal/domain"
	"sync"
)

type QuotesCache struct {
	data map[string]domain.Quotes
	mu   sync.RWMutex
}

func NewQuotesCache() *QuotesCache {
	return &QuotesCache{
		data: make(map[string]domain.Quotes),
		mu:   sync.RWMutex{},
	}
}

func (c *QuotesCache) Set(symbol string, quote domain.Quotes) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[symbol] = quote
}

func (c *QuotesCache) Get(symbol string) (domain.Quotes, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	quote, exists := c.data[symbol]
	return quote, exists
}

func (c *QuotesCache) GetAll() []domain.Quotes {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tickers := make([]domain.Quotes, 0, len(c.data))
	for _, v := range c.data {
		tickers = append(tickers, v)
	}
	return tickers
}

func (c *QuotesCache) UpdateAll(quotes []domain.Quotes) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, q := range quotes {
		c.data[q.Symbol] = q
	}
}
