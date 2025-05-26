package service

import "context"

func (s *QuotesService) informer(ctx context.Context) {
	for {
		select {
		case quote, ok := <-s.updatesChan:
			if !ok {
				return
			}
			s.logger.Infof("UPDATE | %s | Price: %.2f | 24h Volume: %.2f | Last: %.2f",
				quote.Symbol, quote.Price24h, quote.Volume24h, quote.LastTradePrice)
		case <-ctx.Done():
			s.logger.Info("Informer: context cancelled")
			return
		}
	}
}
