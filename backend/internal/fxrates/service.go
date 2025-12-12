package fxrates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL       = "https://api.fastforex.io"
	cacheDuration = 24 * time.Hour // Cache rates for 24 hours
)

// Service handles FX rate operations
type Service struct {
	apiKey      string
	cache       *RateCache
	httpClient  *http.Client
}

// RateCache stores cached exchange rates
type RateCache struct {
	mu          sync.RWMutex
	rates       map[string]float64
	baseCurrency string
	lastUpdated time.Time
}

// NewService creates a new FX rates service
func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
		cache: &RateCache{
			rates: make(map[string]float64),
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRates retrieves exchange rates for a base currency
func (s *Service) GetRates(baseCurrency string) (*FXRatesResponse, error) {
	// Check cache first
	if s.isCacheValid(baseCurrency) {
		return s.getCachedRates(), nil
	}

	// Fetch fresh rates from API
	rates, err := s.fetchRatesFromAPI(baseCurrency)
	if err != nil {
		// If API fails and we have stale cache, return it
		if s.cache.baseCurrency == baseCurrency && len(s.cache.rates) > 0 {
			return s.getCachedRates(), nil
		}
		return nil, err
	}

	// Update cache
	s.updateCache(baseCurrency, rates)

	return s.getCachedRates(), nil
}

// GetRate retrieves a specific exchange rate between two currencies
func (s *Service) GetRate(from, to string) (float64, error) {
	rates, err := s.GetRates(from)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[to]
	if !exists {
		return 0, fmt.Errorf("rate not found for %s/%s", from, to)
	}

	return rate, nil
}

// Convert converts an amount from one currency to another
func (s *Service) Convert(from, to string, amount float64) (*ConversionResponse, error) {
	rate, err := s.GetRate(from, to)
	if err != nil {
		return nil, err
	}

	result := amount * rate

	return &ConversionResponse{
		From:        from,
		To:          to,
		Amount:      amount,
		Result:      result,
		Rate:        rate,
		LastUpdated: s.cache.lastUpdated,
	}, nil
}

// fetchRatesFromAPI fetches rates from FastForex API
func (s *Service) fetchRatesFromAPI(baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/fetch-all?from=%s&api_key=%s", baseURL, baseCurrency, s.apiKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error body for more details
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp FastForexAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResp.Results) == 0 {
		return nil, fmt.Errorf("API returned empty results")
	}

	return apiResp.Results, nil
}

// isCacheValid checks if the cache is still valid
func (s *Service) isCacheValid(baseCurrency string) bool {
	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	if s.cache.baseCurrency != baseCurrency {
		return false
	}

	if len(s.cache.rates) == 0 {
		return false
	}

	return time.Since(s.cache.lastUpdated) < cacheDuration
}

// updateCache updates the cache with new rates
func (s *Service) updateCache(baseCurrency string, rates map[string]float64) {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	s.cache.baseCurrency = baseCurrency
	s.cache.rates = rates
	s.cache.lastUpdated = time.Now()
}

// getCachedRates returns the cached rates
func (s *Service) getCachedRates() *FXRatesResponse {
	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	return &FXRatesResponse{
		BaseCurrency: s.cache.baseCurrency,
		Rates:        s.cache.rates,
		LastUpdated:  s.cache.lastUpdated,
	}
}

// RefreshCache forces a cache refresh
func (s *Service) RefreshCache(baseCurrency string) error {
	rates, err := s.fetchRatesFromAPI(baseCurrency)
	if err != nil {
		return err
	}

	s.updateCache(baseCurrency, rates)
	return nil
}
