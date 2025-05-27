package service

import (
	"Effective-Mobile/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Enricher interface {
	Enrich(ctx context.Context, p *domain.Person) error
}

type enrichSvc struct {
	http *http.Client
}

func NewEnricher(timeout time.Duration) Enricher {
	return &enrichSvc{
		http: &http.Client{Timeout: timeout},
	}
}

func (e *enrichSvc) Enrich(ctx context.Context, p *domain.Person) error {
	if p.Name == "" {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		url := fmt.Sprintf("https://api.agify.io/?name=%s", p.Name)
		var resp struct {
			Age *int `json:"age"`
		}
		return fetch(ctx, e.http, url, &resp, func() { p.Age = resp.Age })
	})

	g.Go(func() error {
		url := fmt.Sprintf("https://api.genderize.io/?name=%s", p.Name)
		var resp struct {
			Gender *string `json:"gender"`
		}
		return fetch(ctx, e.http, url, &resp, func() { p.Gender = resp.Gender })
	})

	g.Go(func() error {
		url := fmt.Sprintf("https://api.nationalize.io/?name=%s", p.Name)
		var resp struct {
			Country []struct {
				CountryID   string  `json:"country_id"`
				Probability float32 `json:"probability"`
			} `json:"country"`
		}
		return fetch(ctx, e.http, url, &resp, func() {
			if len(resp.Country) > 0 {
				p.CountryID = &resp.Country[0].CountryID
				p.Probability = &resp.Country[0].Probability
			}
		})
	})

	return g.Wait()
}

func fetch[T any](ctx context.Context, cli *http.Client, url string, out *T, apply func()) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status %s", res.Status)
	}
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return err
	}
	apply()
	return nil
}
