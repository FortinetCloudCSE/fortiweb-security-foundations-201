// ml-mix.go
// Juice Shop traffic mixer (menu-driven): legit terms + attack payloads via ?q=
// used to send traffic to owasp Juiceshop sitting behind FortiAppsec
// Wondyrad T version 1.0 08/25
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	code    int
	latency time.Duration
	err     bool
}

type kv struct{ Code, Count int }

func main() {
	// --- Menu with defaults ---
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("=== Juice Shop Traffic Mixer ===")
	base := promptString(reader, "Base URL", "https://669.fwebtraincse.com")
	durationStr := promptString(reader, "Duration (e.g. 5m, 300s)", "5m")
	rpsStr := promptString(reader, "Target RPS (global)", "30")
	concurrencyStr := promptString(reader, "Workers (concurrency)", "20")
	attackPctStr := promptString(reader, "Attack mix percentage (0-100)", "30")
	useRESTStr := promptString(reader, "Use /rest/products/search?q=... ? (y/n)", "n")
	insecureStr := promptString(reader, "Skip TLS verification? (y/n)", "n")
	timeoutStr := promptString(reader, "Per-request timeout (e.g. 10s)", "10s")
	verboseStr := promptString(reader, "Verbose sample logging? (y/n)", "n")

	// Parse/validate
	u, err := url.Parse(base)
	if err != nil || u.Scheme == "" || u.Host == "" {
		fmt.Println("Invalid Base URL. Example: https://669.fwebtraincse.com")
		return
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		fmt.Println("Invalid duration; try like 5m or 300s")
		return
	}
	rps, err := strconv.Atoi(strings.TrimSpace(rpsStr))
	if err != nil || rps <= 0 {
		fmt.Println("Invalid RPS; enter a positive integer")
		return
	}
	workers, err := strconv.Atoi(strings.TrimSpace(concurrencyStr))
	if err != nil || workers <= 0 {
		fmt.Println("Invalid workers; enter a positive integer")
		return
	}
	attackPct, err := strconv.Atoi(strings.TrimSpace(attackPctStr))
	if err != nil || attackPct < 0 || attackPct > 100 {
		fmt.Println("Invalid attack percentage; must be 0-100")
		return
	}
	useREST := strings.ToLower(strings.TrimSpace(useRESTStr)) == "y"
	insecureTLS := strings.ToLower(strings.TrimSpace(insecureStr)) == "y"
	reqTimeout, err := time.ParseDuration(timeoutStr)
	if err != nil || reqTimeout <= 0 {
		fmt.Println("Invalid timeout; try like 10s")
		return
	}
	verbose := strings.ToLower(strings.TrimSpace(verboseStr)) == "y"

	// HTTP client
	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecureTLS}, // testing only
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: tr, Timeout: reqTimeout}

	// Rate limiter
	interval := time.Second / time.Duration(max(rps, 1))
	tok := time.NewTicker(interval)
	defer tok.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Stats
	var sent, ok, failed uint64
	var latMu sync.Mutex
	var latencies []time.Duration
	codeMu := &sync.Mutex{}
	codeCount := map[int]int{}
	results := make(chan result, 2048)
	wg := &sync.WaitGroup{}

	// Collector
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range results {
			atomic.AddUint64(&sent, 1)
			if r.err {
				atomic.AddUint64(&failed, 1)
				continue
			}
			atomic.AddUint64(&ok, 1)
			latMu.Lock()
			latencies = append(latencies, r.latency)
			latMu.Unlock()
			codeMu.Lock()
			codeCount[r.code]++
			codeMu.Unlock()
		}
	}()

	// Workers
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)*991))
			// Determine path based on REST preference
			path := "/products/search" // Default to same endpoint as ml2.go
			if useREST {
				path = "/rest/products/search" // Use REST if requested
			}
			for {
				select {
				case <-ctx.Done():
					return
				case <-tok.C:
					qval := chooseQ(rng, attackPct)
					reqURL := buildSearchURL(u, path, qval)

					req, _ := http.NewRequest("GET", reqURL, nil)
					randomizeHeaders(req, rng, u)

					start := time.Now()
					resp, err := client.Do(req.WithContext(ctx))
					lat := time.Since(start)
					if err != nil {
						if verbose && rng.Intn(100) == 0 {
							fmt.Printf("[w%d] error: %v\n", id, err)
						}
						results <- result{err: true}
						continue
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					if verbose && rng.Intn(200) == 0 {
						fmt.Printf("[w%d] GET %s -> %d in %v\n", id, reqURL, resp.StatusCode, lat)
					}
					results <- result{code: resp.StatusCode, latency: lat}
				}
			}
		}(i)
	}

	// Minute-by-minute progress
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		start := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s := atomic.LoadUint64(&sent)
				o := atomic.LoadUint64(&ok)
				f := atomic.LoadUint64(&failed)
				elapsed := time.Since(start).Truncate(time.Second)
				fmt.Printf("[progress] elapsed=%s sent=%d ok=%d failed=%d\n", elapsed, s, o, f)
			}
		}
	}()

	<-ctx.Done()
	time.Sleep(300 * time.Millisecond)
	close(results)
	wg.Wait()

	// Summary (with your fixes: kv + sortCodes([]kv) + latities copy)
	latMu.Lock()
	p50 := percentile(latities(latencies), 0.50)
	p95 := percentile(latities(latencies), 0.95)
	p99 := percentile(latities(latencies), 0.99)
	latMu.Unlock()

	var codes []kv
	codeMu.Lock()
	for c, n := range codeCount {
		codes = append(codes, kv{c, n})
	}
	codeMu.Unlock()
	sortCodes(codes)

	fmt.Println("\n=== Summary ===")
	fmt.Printf("Target:     %s\n", base)
	fmt.Printf("Endpoint:   %s\n", map[bool]string{true: "/rest/products/search?q=", false: "/products/search?q="}[useREST])
	fmt.Printf("Duration:   %s | RPS target: %d | Workers: %d | Attack%%: %d\n", duration, rps, workers, attackPct)
	fmt.Printf("Sent:       %d | OK: %d | Failed: %d\n", atomic.LoadUint64(&sent), atomic.LoadUint64(&ok), atomic.LoadUint64(&failed))
	fmt.Printf("Latency:    p50=%v p95=%v p99=%v\n", p50, p95, p99)
	b, _ := json.MarshalIndent(codes, "", "  ")
	fmt.Printf("Codes:      %s\n", string(b))
}

// --- Prompt helpers ---
func promptString(r *bufio.Reader, label, def string) string {
	fmt.Printf("%s [%s]: ", label, def)
	text, _ := r.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return def
	}
	return text
}

// --- Request helpers ---
func buildSearchURL(base *url.URL, path, qval string) string {
	u := *base
	u.Path = strings.TrimSuffix(base.Path, "/") + path
	qs := url.Values{}
	qs.Set("q", qval)
	u.RawQuery = qs.Encode()
	return u.String()
}

func chooseQ(rng *rand.Rand, attackPct int) string {
	if rng.Intn(100) < clamp(attackPct, 0, 100) {
		return attackPayload(rng)
	}
	return legitTerm(rng)
}

func legitTerm(rng *rand.Rand) string {
	terms := []string{
		"apple", "Banana",
	}
	return terms[rng.Intn(len(terms))]
}

func attackPayload(rng *rand.Rand) string {
	switch rng.Intn(5) {
	case 0: // XSS
		xss := []string{
			`<script>alert(1)</script>`,
			`"><img src=x onerror=alert(1)>`,
			`<svg onload=confirm(1)>`,
			`"><svg/onload=alert(1)>`,
		}
		return xss[rng.Intn(len(xss))]
	case 1: // SQLi-ish
		sqli := []string{
			`' OR '1'='1`,
			`1 OR 1=1--`,
			`%' UNION SELECT NULL,NULL,NULL--`,
			`') OR ('1'='1`,
		}
		return sqli[rng.Intn(len(sqli))]
	case 2: // Template injection-ish
		tpl := []string{`{{7*7}}`, `${7*7}`, `{{constructor.constructor('return process')()}}`}
		return tpl[rng.Intn(len(tpl))]
	case 3: // Traversal-ish
		lfi := []string{`../../etc/passwd`, `..%2f..%2f..%2fetc%2fhosts`, `..\\..\\windows\\win.ini`}
		return lfi[rng.Intn(len(lfi))]
	default: // Header split-ish text
		hs := []string{`%0d%0aSet-Cookie:spl=1`, `%0d%0aLocation://evil`}
		return hs[rng.Intn(len(hs))]
	}
}

func randomizeHeaders(r *http.Request, rng *rand.Rand, base *url.URL) {
	r.Header.Set("User-Agent", randomUA(rng))
	r.Header.Set("X-Forwarded-For", randomIP(rng))
	if rng.Intn(100) < 20 {
		r.Header.Set("Accept-Language", randomLang(rng))
	}
	if rng.Intn(100) < 10 {
		r.Header.Set("Referer", base.Scheme+"://"+base.Host+"/#/")
	}
}

func randomUA(rng *rand.Rand) string {
	list := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/123 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/605.1.15 Version/17.2 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148",
		"Mozilla/5.0 (Android 14; Mobile) AppleWebKit/537.36 Chrome/124 Mobile Safari/537.36",
		"curl/8.5.0",
		"Wget/1.21.3",
	}
	return list[rng.Intn(len(list))]
}

func randomIP(rng *rand.Rand) string {
	o := func() int { return rng.Intn(256) }
	return fmt.Sprintf("%d.%d.%d.%d", 23+rng.Intn(200), o(), o(), 1+rng.Intn(254))
}

func randomLang(rng *rand.Rand) string {
	langs := []string{"en-US,en;q=0.9", "en-GB,en;q=0.8", "fr-FR,fr;q=0.7", "de-DE,de;q=0.7", "am-ET,am;q=0.8,en;q=0.5"}
	return langs[rng.Intn(len(langs))]
}

// --- Stats helpers (with your latities copy helper) ---
func latities(d []time.Duration) []time.Duration {
	cp := make([]time.Duration, len(d))
	copy(cp, d)
	return cp
}
func percentile(a []time.Duration, p float64) time.Duration {
	if len(a) == 0 {
		return 0
	}
	n := int(math.Ceil(float64(len(a))*p)) - 1
	if n < 0 {
		n = 0
	}
	quickSelect(a, n)
	return a[n]
}
func quickSelect(a []time.Duration, k int) {
	l, r := 0, len(a)-1
	for l < r {
		p := part(a, l, r)
		if p == k {
			return
		}
		if k < p {
			r = p - 1
		} else {
			l = p + 1
		}
	}
}
func part(a []time.Duration, l, r int) int {
	pv := a[r]
	i := l
	for j := l; j < r; j++ {
		if a[j] < pv {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[r] = a[r], a[i]
	return i
}
func sortCodes(c []kv) {
	for i := 1; i < len(c); i++ {
		for j := i; j > 0 && c[j-1].Code > c[j].Code; j-- {
			c[j-1], c[j] = c[j], c[j-1]
		}
	}
}

// --- misc ---
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
