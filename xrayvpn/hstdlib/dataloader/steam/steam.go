package steam

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/realglebivanov/hstd/hstdlib/httpclient"
)

const asn = "AS32590"
const url = "https://stat.ripe.net/data/announced-prefixes/data.json?resource=" + asn

type response struct {
	Data struct {
		Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"prefixes"`
	} `json:"data"`
}

func Fetch() ([]string, error) {
	slog.Info("fetching Steam CIDRs", "asn", asn)

	resp, err := httpclient.Default.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", asn, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: HTTP %d", asn, resp.StatusCode)
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode %s: %w", asn, err)
	}

	cidrs := make([]string, 0, len(r.Data.Prefixes))
	for _, p := range r.Data.Prefixes {
		if p.Prefix != "" {
			cidrs = append(cidrs, p.Prefix)
		}
	}

	slog.Info("fetched Steam CIDRs", "count", len(cidrs))
	return cidrs, nil
}
