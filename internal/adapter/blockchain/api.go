package bc

import (
	"BlockchainCurrency/config"
	"BlockchainCurrency/internal/domain"
	"BlockchainCurrency/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Blockchain struct {
	cfg    *config.Config
	client *http.Client
	logger logger.Logger
}

func NewBclockchain(cfg *config.Config,logger logger.Logger) *Blockchain {
	return &Blockchain{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.HTTP.ReadTimeout},
		logger: logger,
	}
}

type DataBlockchainResponse struct {
	Symbol         string  `json:"symbol"`
	Price24h       float64 `json:"price_24h"`
	Volume24h      float64 `json:"volume_24h"`
	LastTradePrice float64 `json:"last_trade_price"`
}

func (s *Blockchain) FetchQuotesData(ctx context.Context, symbol string) (domain.Quotes, error) {
	url := fmt.Sprintf("%s/%s", s.cfg.BlockchainUrl.URL, symbol)
	s.logger.Infof("Requesting quote: %s", url)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Quotes{}, fmt.Errorf("failed to create request for %s: %v", symbol, err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.Quotes{}, fmt.Errorf("request failed for %s: %v", symbol, err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Quotes{}, fmt.Errorf("bad status for %s: %d", symbol, resp.StatusCode)
	}

	var newQuotation DataBlockchainResponse
	if err := json.NewDecoder(resp.Body).Decode(&newQuotation); err != nil {
		return domain.Quotes{}, fmt.Errorf("failed to decode response for %s: %v", symbol, err)
	}

	return MapToQuotation(newQuotation), nil
}

func MapToQuotation(b DataBlockchainResponse) domain.Quotes {
	return domain.Quotes{
		Symbol:         b.Symbol,
		Price24h:       b.Price24h,
		Volume24h:      b.Volume24h,
		LastTradePrice: b.LastTradePrice,
	}
}
