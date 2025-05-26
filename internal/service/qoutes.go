package service

import (
	"BlockchainCurrency/config"
	"BlockchainCurrency/internal/domain"
	"BlockchainCurrency/pkg/logger"
	"context"
	"fmt"

	"sync"
	"time"
)

const (
	workersNum = 10
)

type QuotesService struct {
	cache         CacheRepository
	updatesChan   chan domain.Quotes
	wg            sync.WaitGroup
	logger        logger.Logger
	blockchainApi BlockchainFetchData
	cfg           *config.Config
	repo          QuotesRepository
	shuter        ServiceShuter
}

type ServiceShuter struct {
	stopChan chan struct{}
	mu       sync.Mutex
	isCalled bool
}

type QuotesRepository interface {
	GetAllQuotes(ctx context.Context) ([]domain.Quotes, error)
	SaveQuotes(ctx context.Context, quotes []domain.Quotes) error
}

type BlockchainFetchData interface {
	FetchQuotesData(ctx context.Context, symbol string) (domain.Quotes, error)
}

type CacheRepository interface {
	Set(symbol string, quote domain.Quotes)
	Get(symbol string) (domain.Quotes, bool)
	GetAll() []domain.Quotes
	UpdateAll(quotes []domain.Quotes)
}

func NewQuotesService(logger logger.Logger, cfg *config.Config, repo QuotesRepository, BlockchainApi BlockchainFetchData, cache CacheRepository) *QuotesService {
	return &QuotesService{
		updatesChan:   make(chan domain.Quotes),
		logger:        logger,
		cfg:           cfg,
		repo:          repo,
		blockchainApi: BlockchainApi,
		wg:            sync.WaitGroup{},
		shuter:        ServiceShuter{stopChan: make(chan struct{}), mu: sync.Mutex{}, isCalled: false},
		cache:         cache,
	}
}

func (s *QuotesService) ListenOnClose(cancel context.CancelFunc) {
	<-s.shuter.stopChan
	cancel()
	s.wg.Wait()
}

func (s *QuotesService) Start(ctx context.Context) error {
	quotes, err := s.repo.GetAllQuotes(ctx)
	if err != nil {
		return fmt.Errorf("failed GetQuotations:%w", err)
	}

	s.cache.UpdateAll(quotes)

	go s.informer(ctx)

	for i := 0; i < workersNum; i++ {
		s.wg.Add(1)
		go s.workerGrabber(ctx)
	}

	return nil
}

func (s *QuotesService) workerGrabber(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	defer s.wg.Done()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, symbol := range s.cfg.Quotes.Symbols {
				if err := s.processSymbol(ctx, symbol); err != nil {
					s.logger.Errorf("Error processing %s: %v", symbol, err)
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *QuotesService) processSymbol(ctx context.Context, symbol string) error {
	data, err := s.blockchainApi.FetchQuotesData(ctx, symbol)
	if err != nil {
		return fmt.Errorf("API fetch failed: %w", err)
	}

	if old, exists := s.cache.Get(symbol); !exists ||
		old.Price24h != data.Price24h ||
		old.Volume24h != data.Volume24h ||
		old.LastTradePrice != data.LastTradePrice {
		s.cache.Set(symbol, data)
		select {
		case <-ctx.Done():
			close(s.updatesChan)
			return nil
		default:
			s.updatesChan <- data
		}
	}
	return nil
}

func (s *QuotesService) GetQuotes() ([]domain.Quotes, error) {
	return s.cache.GetAll(), nil
}
