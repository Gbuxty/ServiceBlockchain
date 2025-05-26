package service

import (
	"BlockchainCurrency/internal/domain"
	"context"
	"fmt"
)

func (s *QuotesService) Shutdown(ctx context.Context) error {
	s.shuter.mu.Lock()
	if s.shuter.isCalled {
		s.shuter.mu.Unlock()
		return domain.ErrShutdownAlreadyCalled
	}
	s.shuter.isCalled = true
	s.shuter.mu.Unlock()

	quotes := s.cache.GetAll()

	s.logger.Info("Saving quotes before shutdown...")
	
	if err := s.repo.SaveQuotes(ctx, quotes); err != nil {
		return fmt.Errorf("failed to save quotations: %v", err)
	}

	s.logger.Info("QuotesService shutdown completed.")

	s.shuter.stopChan <- struct{}{}
	close(s.shuter.stopChan)

	return nil
}
