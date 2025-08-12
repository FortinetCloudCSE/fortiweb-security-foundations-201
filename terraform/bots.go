// BotLab - menu-driven bot traffic simulator for lab use
// Modes: crawlers, scanners, crawl (link-follow), slow timeouts, scrape burst, mixed
// Features: profiles (JSON), CSV logging, random/sticky XFF IPs, UA rotation, Juice Shop optimization
// NOTE: Only test against systems you own or have permission to test.

package main

import (
	"bufio"
	"context"
	cryptoRand "crypto/rand"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/html"
)

/* =================== Modes & UAs =================== */

type Mode string

const (
	ModeCrawlers Mode = "crawlers"
	ModeScanners Mode = "scanners"
	ModeCrawl    Mode = "crawl"
	ModeSlow     Mode = "slow"
	ModeScrape   Mode = "scrape"
	ModeMixed    Mode = "mixed"
)

var crawlerUAs = []string{
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"bingbot/2.0 (+http://www.bing.com/bingbot.htm)",
	"CCBot/2.0 (https://commoncrawl.org/faq/)",
	"ClaudeBot/1.0 (+https://www.anthropic.com/)", // simulated
	"Diffbot/1.0 (+https://www.diffbot.com/our-robots/)",
	"SemrushBot/7~bl",
	"AhrefsBot/7.0 (+http://ahrefs.com/robot/)",
	"YandexBot/3.0 (+http://yandex.com/bots)",
}
var scannerUAs = []string{
	"sqlmap/1.7.2#stable (http://sqlmap.org)",
	"Acunetix Web Vulnerability Scanner",
	"nikto/2.1.6 (Evasion: 1)",
	"Arachni/v1.6",
	"w3af.org",
	"OWASP ZAP",
	"Burp Suite Professional",
	"commix/3.5-dev",
	"BilboScanner/0.3", // simulated
}
var legitUAs = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
}

/* =================== Config / Stats =================== */

type Config struct {
	Target         string        `json:"target,omitempty"`
	Mode           Mode          `json:"mode"`
	Duration       time.Duration `json:"duration"`
	RPS            int           `json:"rps"`
	Concurrency    int           `json:"concurrency"`
	MaxDepth       int           `json:"max_depth"`
	SameIP         bool          `json:"same_ip"`
	HTTP10         bool          `json:"http10"`
	FollowSameHost bool          `json:"same_host_only"`
	ShowStatsEvery time.Duration `json:"progress"`
	SlowCount      int           `json:"slow_count"`
	SlowWindow     time.Duration `json:"slow_window"`
	SlowBytes      int           `json:"slow_bytes"`
	SlowInterval   time.Duration `json:"slow_interval"`
	BurstRequests  int           `json:"burst_requests"`
	BurstWindow    time.Duration `json:"burst_window"`
	Timeout        time.Duration `json:"timeout"`
	LogCSV         string        `json:"log_csv,omitempty"`
	JuiceOptimized bool          `json:"juice_optimized"`
}

type stats struct {
	sent     uint64
	ok       uint64
	err      uint64
	timeouts uint64
}

/* =================== CSV Logger =================== */

type CSVLogger struct {
	mu     sync.Mutex
	w      *csv.Writer
	file   *os.File
	closed bool
}

func NewCSVLogger(path string) (*CSVLogger, error) {
	if path == "" {
		return nil, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	l := &CSVLogger{file: f, w: csv.NewWriter(f)}
	if fi, _ := f.Stat(); fi.Size() == 0 {
		_ = l.w.Write([]string{"ts", "mode", "method", "url", "status", "bytes", "ms", "xff", "ua", "error"})
		l.w.Flush()
	}
	return l, nil
}
func (l *CSVLogger) Log(mode, method, urlStr string, status, nbytes int, ms int64, xff, ua, errStr string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return
	}
	_ = l.w.Write([]string{
		time.Now().Format(time.RFC3339Nano), mode, method, urlStr,
		strconv.Itoa(status), strconv.Itoa(nbytes), strconv.FormatInt(ms, 10), xff, ua, errStr,
	})
	l.w.Flush()
}
func (l *CSVLogger) Close() error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.closed = true
	l.w.Flush()
	return l.file.Close()
}

/* =================== Profiles (JSON) =================== */

type Profiles map[string]Config

func loadProfiles(path string) (Profiles, error) {
	if path == "" {
		return Profiles{}, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Profiles
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func saveProfile(path, name string, cfg Config) error {
	if path == "" || name == "" {
		return errors.New("profiles path and name required")
	}
	m := Profiles{}
	if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
		_ = json.Unmarshal(b, &m)
	}
	cfg.Target = "" // keep reusable
	m[name] = cfg
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

/* =================== Juice Shop tuning =================== */

var juiceSearchTerms = []string{
	"apple", "banana", "strawberry", "orange", "avocado", "lemon", "lime", "cherry",
	"coffee", "tea", "milk", "almond", "chocolate", "cookie", "cake", "juice", "smoothie",
	"honey", "pepper", "cheese", "yogurt", "ginger", "mango", "papaya", "pomegranate",
}

type weightedPath struct {
	Path   string
	Weight int
}

var juiceWeighted = []weightedPath{
	{"/rest/products/search", 40}, // with ?q=
	{"/rest/products", 15},
	{"/api/Products", 10},
	{"/assets/i18n/en.json", 8},
	{"/assets/public/images/banner-1.png", 4},
	{"/", 5},
	{"/#/search", 8},
	{"/#/contact", 3},
	{"/#/login", 3},
	{"/#/about", 2},
}

func pickJuicePath(base string) string {
	total := 0
	for _, w := range juiceWeighted {
		total += w.Weight
	}
	r := rand.Intn(total)
	acc := 0
	var p string
	for _, w := range juiceWeighted {
		acc += w.Weight
		if r < acc {
			p = w.Path
			break
		}
	}
	u := strings.TrimRight(base, "/") + p
	if strings.Contains(p, "/rest/products/search") {
		q := url.Values{}
		q.Set("q", pick(juiceSearchTerms))
		u += "?" + q.Encode()
	}
	return u
}

/* =================== Main =================== */

func main() {
	nonInteractive := flag.Bool("non-interactive", false, "Use flags instead of interactive menu")
	profilesPath := flag.String("profiles", "profiles.json", "Path to profiles.json")
	flag.Parse()

	seedRand()

	cfg := defaultConfig()
	reader := bufio.NewReader(os.Stdin)

	if *nonInteractive {
		fmt.Println("Non-interactive mode expects full flags; interactive is recommended.")
		return
	}

	fmt.Println("=== BotLab (menu-driven with profiles, CSV, Juice Shop optimization) ===")
	fmt.Println("⚠️  Only test against systems you own or have permission to test.")

	profiles, _ := loadProfiles(*profilesPath)

	useProfile := askBool(reader, "Load a saved profile?", len(profiles) > 0)
	if useProfile && len(profiles) == 0 {
		useProfile = false
	}
	if useProfile {
		names := make([]string, 0, len(profiles))
		for k := range profiles {
			names = append(names, k)
		}
		choice := choose(reader, "Select profile", names, "")
		cfg = mergeIntoDefault(profiles[choice])
		fmt.Printf("Loaded profile '%s'. You can still tweak values.\n", choice)
	}

	for {
		cfg.Target = ask(reader, "Target URL (e.g., https://<host>)", cfg.Target)
		if _, err := url.ParseRequestURI(cfg.Target); err != nil {
			fmt.Println("Invalid URL. Try again.")
			continue
		}
		break
	}

	cfg.JuiceOptimized = askBool(reader, "Optimize for OWASP Juice Shop?", cfg.JuiceOptimized)

	mode := choose(reader, "Mode", []string{
		string(ModeCrawlers), string(ModeScanners), string(ModeCrawl),
		string(ModeSlow), string(ModeScrape), string(ModeMixed),
	}, string(cfg.Mode))
	cfg.Mode = Mode(mode)

	cfg.LogCSV = ask(reader, "CSV log file (blank to skip)", cfg.LogCSV)

	cfg.Duration = askDuration(reader, "Total run duration (e.g., 60s, 2m)", cfg.Duration)
	cfg.Concurrency = askInt(reader, "Concurrency (workers)", cfg.Concurrency)
	cfg.Timeout = askDuration(reader, "HTTP timeout per request", cfg.Timeout)
	cfg.ShowStatsEvery = askDuration(reader, "Progress interval (0=off)", cfg.ShowStatsEvery)
	cfg.SameIP = askBool(reader, "Use one sticky IP in X-Forwarded-For?", cfg.SameIP)
	cfg.HTTP10 = askBool(reader, "Force HTTP/1.0 (mimic some bots)?", cfg.HTTP10)

	switch cfg.Mode {
	case ModeCrawlers, ModeScanners, ModeMixed:
		cfg.RPS = askInt(reader, "Requests per second", cfg.RPS)
	case ModeCrawl:
		cfg.MaxDepth = askInt(reader, "Max link depth", cfg.MaxDepth)
		cfg.FollowSameHost = askBool(reader, "Follow same host only?", cfg.FollowSameHost)
	case ModeSlow:
		cfg.SlowCount = askInt(reader, "Slow timeouts to trigger in window", cfg.SlowCount)
		cfg.SlowWindow = askDuration(reader, "Slow window (e.g., 100s)", cfg.SlowWindow)
		cfg.SlowBytes = askInt(reader, "Bytes to trickle (total)", cfg.SlowBytes)
		cfg.SlowInterval = askDuration(reader, "Delay between trickle chunks", cfg.SlowInterval)
	case ModeScrape:
		cfg.BurstRequests = askInt(reader, "Requests in burst (e.g., 170–200)", cfg.BurstRequests)
		cfg.BurstWindow = askDuration(reader, "Burst window (e.g., 30s)", cfg.BurstWindow)
	}

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Target: %s\nMode: %s\nJuice Optimized: %v\nDuration: %s\nConcurrency: %d\nTimeout: %s\nHTTP/1.0: %v\nSticky XFF IP: %v\n",
		cfg.Target, cfg.Mode, cfg.JuiceOptimized, cfg.Duration, cfg.Concurrency, cfg.Timeout, cfg.HTTP10, cfg.SameIP)
	switch cfg.Mode {
	case ModeCrawlers, ModeScanners, ModeMixed:
		fmt.Printf("RPS: %d\n", cfg.RPS)
	case ModeCrawl:
		fmt.Printf("MaxDepth: %d  SameHostOnly: %v\n", cfg.MaxDepth, cfg.FollowSameHost)
	case ModeSlow:
		fmt.Printf("SlowCount: %d  SlowWindow: %s  SlowBytes: %d  SlowInterval: %s\n", cfg.SlowCount, cfg.SlowWindow, cfg.SlowBytes, cfg.SlowInterval)
	case ModeScrape:
		fmt.Printf("BurstRequests: %d  BurstWindow: %s\n", cfg.BurstRequests, cfg.BurstWindow)
	}
	if !askBool(reader, "Start now?", true) {
		fmt.Println("Aborted.")
		return
	}

	if askBool(reader, "Save these settings as a profile for next time?", false) {
		name := ask(reader, "Profile name", "")
		if name != "" {
			if err := saveProfile(*profilesPath, name, cfg); err != nil {
				fmt.Println("Save failed:", err)
			} else {
				fmt.Printf("Saved profile '%s' to %s\n", name, *profilesPath)
			}
		}
	}

	run(cfg)
}

/* =================== Defaults / Menu helpers =================== */

func defaultConfig() Config {
	return Config{
		Mode:           ModeCrawlers,
		Duration:       60 * time.Second,
		RPS:            20,
		Concurrency:    20,
		MaxDepth:       2,
		SameIP:         false,
		HTTP10:         false,
		FollowSameHost: true,
		ShowStatsEvery: 10 * time.Second,
		SlowCount:      7,
		SlowWindow:     100 * time.Second,
		SlowBytes:      4096,
		SlowInterval:   1500 * time.Millisecond,
		BurstRequests:  170,
		BurstWindow:    30 * time.Second,
		Timeout:        12 * time.Second,
		JuiceOptimized: true,
	}
}

func mergeIntoDefault(in Config) Config {
	def := defaultConfig()
	if in.Mode != "" {
		def.Mode = in.Mode
	}
	if in.Duration != 0 {
		def.Duration = in.Duration
	}
	if in.RPS != 0 {
		def.RPS = in.RPS
	}
	if in.Concurrency != 0 {
		def.Concurrency = in.Concurrency
	}
	if in.MaxDepth != 0 {
		def.MaxDepth = in.MaxDepth
	}
	if in.SameIP {
		def.SameIP = true
	}
	if in.HTTP10 {
		def.HTTP10 = true
	}
	if in.FollowSameHost == false {
		def.FollowSameHost = false
	}
	if in.ShowStatsEvery != 0 {
		def.ShowStatsEvery = in.ShowStatsEvery
	}
	if in.SlowCount != 0 {
		def.SlowCount = in.SlowCount
	}
	if in.SlowWindow != 0 {
		def.SlowWindow = in.SlowWindow
	}
	if in.SlowBytes != 0 {
		def.SlowBytes = in.SlowBytes
	}
	if in.SlowInterval != 0 {
		def.SlowInterval = in.SlowInterval
	}
	if in.BurstRequests != 0 {
		def.BurstRequests = in.BurstRequests
	}
	if in.BurstWindow != 0 {
		def.BurstWindow = in.BurstWindow
	}
	if in.Timeout != 0 {
		def.Timeout = in.Timeout
	}
	if in.LogCSV != "" {
		def.LogCSV = in.LogCSV
	}
	def.JuiceOptimized = in.JuiceOptimized || def.JuiceOptimized
	return def
}

func ask(r *bufio.Reader, prompt, def string) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", prompt, def)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	text, _ := r.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return def
	}
	return text
}
func askInt(r *bufio.Reader, prompt string, def int) int {
	for {
		s := ask(r, prompt, strconv.Itoa(def))
		v, err := strconv.Atoi(s)
		if err == nil && v >= 0 {
			return v
		}
		fmt.Println("Enter a non-negative integer.")
	}
}
func askBool(r *bufio.Reader, prompt string, def bool) bool {
	d := "n"
	if def {
		d = "y"
	}
	for {
		s := strings.ToLower(ask(r, prompt+" (y/n)", d))
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
		fmt.Println("Type y or n.")
	}
}
func askDuration(r *bufio.Reader, prompt string, def time.Duration) time.Duration {
	for {
		s := ask(r, prompt, def.String())
		d, err := time.ParseDuration(s)
		if err == nil && d >= 0 {
			return d
		}
		fmt.Println("Enter a valid duration (e.g., 10s, 2m, 1h).")
	}
}
func choose(r *bufio.Reader, prompt string, options []string, def string) string {
	fmt.Printf("%s:\n", prompt)
	defIdx := 0
	for i, o := range options {
		mk := " "
		if o == def && def != "" {
			mk = "*"
			defIdx = i
		}
		fmt.Printf("  %d) %s %s\n", i+1, mk, o)
	}
	for {
		d := "1"
		if def != "" {
			d = strconv.Itoa(defIdx + 1)
		}
		s := ask(r, "Choose number", d)
		i, err := strconv.Atoi(s)
		if err == nil && i >= 1 && i <= len(options) {
			return options[i-1]
		}
		fmt.Println("Invalid choice.")
	}
}

/* =================== Runner =================== */

func run(cfg Config) {
	client := makeHTTPClient(cfg)

	logger, err := NewCSVLogger(cfg.LogCSV)
	if err != nil {
		fmt.Println("CSV logging disabled (error opening file):", err)
	}
	if logger != nil {
		defer logger.Close()
	}

	var s stats
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration)
	defer cancel()

	start := time.Now()
	var wg sync.WaitGroup

	if cfg.ShowStatsEvery > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := time.NewTicker(cfg.ShowStatsEvery)
			defer t.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					fmt.Printf("[progress] sent=%d ok=%d err=%d timeouts=%d elapsed=%s\n",
						atomic.LoadUint64(&s.sent),
						atomic.LoadUint64(&s.ok),
						atomic.LoadUint64(&s.err),
						atomic.LoadUint64(&s.timeouts),
						time.Since(start).Truncate(time.Second),
					)
				}
			}
		}()
	}

	switch cfg.Mode {
	case ModeCrawlers:
		runFixedRPS(ctx, &wg, cfg, client, &s, logger, generateCrawlerReq)
	case ModeScanners:
		runFixedRPS(ctx, &wg, cfg, client, &s, logger, generateScannerReq)
	case ModeCrawl:
		runCrawler(ctx, &wg, cfg, client, &s, logger)
	case ModeSlow:
		runSlowTimeouts(ctx, &wg, cfg, client, &s, logger)
	case ModeScrape:
		runScrapeBurst(ctx, &wg, cfg, client, &s, logger)
	case ModeMixed:
		runMixed(ctx, &wg, cfg, client, &s, logger)
	}

	wg.Wait()
	fmt.Printf("\n=== Done (%s) ===\n", cfg.Mode)
	fmt.Printf("sent=%d ok=%d err=%d timeouts=%d total=%s\n",
		s.sent, s.ok, s.err, s.timeouts, time.Since(start).Truncate(time.Millisecond))
}

/* =================== Rand / HTTP client =================== */

func seedRand() {
	var b [8]byte
	if _, err := cryptoRand.Read(b[:]); err != nil {
		rand.Seed(time.Now().UnixNano())
		return
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func makeHTTPClient(cfg Config) *http.Client {
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		MaxIdleConns:          200,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       30 * time.Second,
		DisableCompression:    false,
		ForceAttemptHTTP2:     !cfg.HTTP10,
		ResponseHeaderTimeout: cfg.Timeout,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Transport: tr,
		Timeout:   cfg.Timeout + 2*time.Second,
	}
}

/* =================== Generators =================== */

type reqGen func(base string, sameIP bool, juice bool) (*http.Request, error)

func generateCrawlerReq(base string, sameIP bool, juice bool) (*http.Request, error) {
	target := base
	if juice {
		target = pickJuicePath(base)
	} else {
		target = pickTargetPath(base, true)
	}
	req, _ := http.NewRequest(http.MethodGet, target, nil)
	ua := pick(crawlerUAs)
	commonBotHeaders(req, ua, sameIP)

	if strings.Contains(target, "/api/") || strings.Contains(target, "/rest/") {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Referer", strings.TrimRight(base, "/")+"/#/search")
	}
	randomizeQuirks(req)
	return req, nil
}

func generateScannerReq(base string, sameIP bool, juice bool) (*http.Request, error) {
	target := base
	if juice {
		target = pickJuicePath(base)
	} else {
		target = pickTargetPath(base, false)
	}
	method := http.MethodGet
	if rand.Intn(10) < 3 {
		method = http.MethodPost
	}

	u, _ := url.Parse(target)
	if strings.Contains(u.Path, "/rest/products/search") {
		q := u.Query()
		if q.Get("q") == "" {
			q.Set("q", pick(juiceSearchTerms))
		}
		u.RawQuery = q.Encode()
	} else if !juice {
		for k, v := range reqFuzzParams() {
			q := u.Query()
			q.Set(k, v)
			u.RawQuery = q.Encode()
		}
	}

	req, _ := http.NewRequest(method, u.String(), nil)
	ua := pick(scannerUAs)
	commonBotHeaders(req, ua, sameIP)
	if strings.Contains(u.Path, "/api/") || strings.Contains(u.Path, "/rest/") {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	}
	randomizeQuirks(req)
	return req, nil
}

func commonBotHeaders(req *http.Request, ua string, sameIP bool) {
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", pick([]string{"", "en-US,en;q=0.9", "en"}))
	req.Header.Set("Connection", pick([]string{"keep-alive", "close"}))
	req.Header.Set("Referer", pick([]string{"", "/", req.URL.Scheme + "://" + req.URL.Host + "/"}))
	req.Header.Set("X-Forwarded-For", pickIP(sameIP))
	if rand.Intn(10) < 3 {
		req.Header.Set("From", "crawler@example.com")
	}
	if rand.Intn(10) < 2 {
		req.Header.Set("Range", "bytes=0-1024")
	}
}
func randomizeQuirks(req *http.Request) {
	if rand.Intn(10) < 2 {
		req.Header.Set("Cache-Control", "no-cache")
	}
	if rand.Intn(10) < 2 {
		req.Header.Set("Pragma", "no-cache")
	}
	if rand.Intn(10) < 2 {
		req.Header.Set("DNT", "1")
	}
}

/* =================== Schedulers =================== */

func runFixedRPS(ctx context.Context, wg *sync.WaitGroup, cfg Config, client *http.Client, s *stats, logger *CSVLogger, gen reqGen) {
	limiter := time.NewTicker(time.Second / time.Duration(max(1, cfg.RPS)))
	defer limiter.Stop()
	work := make(chan *http.Request, cfg.Concurrency*2)
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(ctx, wg, client, s, logger, cfg.Mode, work)
	}
	for {
		select {
		case <-ctx.Done():
			close(work)
			return
		case <-limiter.C:
			req, err := gen(cfg.Target, cfg.SameIP, cfg.JuiceOptimized)
			if err != nil {
				continue
			}
			select {
			case work <- req:
			case <-ctx.Done():
				close(work)
				return
			}
		}
	}
}

func runCrawler(ctx context.Context, wg *sync.WaitGroup, cfg Config, client *http.Client, s *stats, logger *CSVLogger) {
	type node struct {
		u     string
		depth int
	}
	queue := make(chan node, 1024)
	seen := sync.Map{}
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range queue {
				if ctx.Err() != nil {
					return
				}
				req, _ := http.NewRequest(http.MethodGet, n.u, nil)
				commonBotHeaders(req, pick(append(crawlerUAs, legitUAs...)), cfg.SameIP)
				randomizeQuirks(req)
				body, _, _ := doRequestWithLog(client, req, s, logger, cfg.Mode)
				if body != nil && n.depth < cfg.MaxDepth {
					for _, link := range extractLinks(n.u, body, cfg.FollowSameHost) {
						if _, ok := seen.LoadOrStore(link, true); !ok {
							select {
							case queue <- node{u: link, depth: n.depth + 1}:
							case <-ctx.Done():
								return
							}
						}
					}
				}
			}
		}()
	}
	root := cfg.Target
	seen.Store(root, true)
	queue <- node{u: root, depth: 0}
	timer := time.NewTimer(cfg.Duration)
	<-timer.C
	close(queue)
}

func runSlowTimeouts(ctx context.Context, wg *sync.WaitGroup, cfg Config, client *http.Client, s *stats, logger *CSVLogger) {
	interval := cfg.SlowWindow / time.Duration(max(1, cfg.SlowCount))
	t := time.NewTicker(interval)
	defer t.Stop()
	for i := 0; i < cfg.SlowCount; i++ {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := doSlowPost(ctx, client, cfg, logger)
				atomic.AddUint64(&s.sent, 1)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
						atomic.AddUint64(&s.timeouts, 1)
					} else {
						atomic.AddUint64(&s.err, 1)
					}
					return
				}
				atomic.AddUint64(&s.ok, 1)
			}()
		}
	}
	<-time.After(minDur(cfg.Duration, cfg.SlowWindow+5*time.Second))
}

func runScrapeBurst(ctx context.Context, wg *sync.WaitGroup, cfg Config, client *http.Client, s *stats, logger *CSVLogger) {
	perReqDelay := cfg.BurstWindow / time.Duration(max(1, cfg.BurstRequests))
	ip := pickIP(true)
	var wgInner sync.WaitGroup
	for i := 0; i < cfg.BurstRequests; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		wgInner.Add(1)
		go func() {
			defer wgInner.Done()
			target := pickTargetPath(cfg.Target, true)
			if cfg.JuiceOptimized {
				target = pickJuicePath(cfg.Target)
			}
			req, _ := http.NewRequest(http.MethodGet, target, nil)
			commonBotHeaders(req, pick(append(crawlerUAs, scannerUAs...)), true)
			req.Header.Set("X-Forwarded-For", ip)
			if strings.Contains(target, "/api/") || strings.Contains(target, "/rest/") {
				req.Header.Set("Accept", "application/json")
				req.Header.Set("X-Requested-With", "XMLHttpRequest")
			}
			randomizeQuirks(req)
			_, _, _ = doRequestWithLog(client, req, s, logger, cfg.Mode)
		}()
		time.Sleep(perReqDelay)
	}
	wgInner.Wait()
	<-time.After(remaining(ctx))
}

func runMixed(ctx context.Context, wg *sync.WaitGroup, cfg Config, client *http.Client, s *stats, logger *CSVLogger) {
	gens := []reqGen{generateCrawlerReq, generateScannerReq}
	limiter := time.NewTicker(time.Second / time.Duration(max(1, cfg.RPS)))
	defer limiter.Stop()
	work := make(chan *http.Request, cfg.Concurrency*2)
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(ctx, wg, client, s, logger, cfg.Mode, work)
	}
	for {
		select {
		case <-ctx.Done():
			close(work)
			return
		case <-limiter.C:
			req, err := gens[rand.Intn(len(gens))](cfg.Target, cfg.SameIP, cfg.JuiceOptimized)
			if err != nil {
				continue
			}
			select {
			case work <- req:
			case <-ctx.Done():
				close(work)
				return
			}
		}
	}
}

/* =================== HTTP helpers =================== */

func worker(ctx context.Context, wg *sync.WaitGroup, client *http.Client, s *stats, logger *CSVLogger, mode Mode, in <-chan *http.Request) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-in:
			if !ok {
				return
			}
			doRequestWithLog(client, req, s, logger, mode)
		}
	}
}

func doRequestWithLog(client *http.Client, req *http.Request, s *stats, logger *CSVLogger, mode Mode) ([]byte, int, error) {
	atomic.AddUint64(&s.sent, 1)
	start := time.Now()
	resp, err := client.Do(req)
	ms := time.Since(start).Milliseconds()
	xff := req.Header.Get("X-Forwarded-For")
	ua := req.Header.Get("User-Agent")

	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			atomic.AddUint64(&s.timeouts, 1)
		} else {
			atomic.AddUint64(&s.err, 1)
		}
		if logger != nil {
			logger.Log(string(mode), req.Method, req.URL.String(), 0, 0, ms, xff, ua, err.Error())
		}
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	n := len(b)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		atomic.AddUint64(&s.ok, 1)
	} else {
		atomic.AddUint64(&s.err, 1)
	}
	if logger != nil {
		logger.Log(string(mode), req.Method, req.URL.String(), resp.StatusCode, n, ms, xff, ua, "")
	}
	return b, resp.StatusCode, nil
}

func doSlowPost(ctx context.Context, client *http.Client, cfg Config, logger *CSVLogger) error {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		chunk := []byte("x")
		total := 0
		for total < cfg.SlowBytes {
			_, _ = pw.Write(chunk)
			total += len(chunk)
			time.Sleep(cfg.SlowInterval)
		}
	}()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, cfg.Target, pr)
	commonBotHeaders(req, pick(scannerUAs), cfg.SameIP)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Expect", "100-continue")

	var dummy stats
	_, _, err := doRequestWithLog(client, req, &dummy, logger, ModeSlow)
	return err
}

/* =================== HTML / utils =================== */

func extractLinks(base string, body []byte, sameHostOnly bool) []string {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if u := normalizeLink(base, a.Val, sameHostOnly); u != "" {
						links = append(links, u)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return dedup(links)
}

func normalizeLink(base, href string, sameHostOnly bool) string {
	if href == "" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
		return ""
	}
	if strings.HasPrefix(href, "#/") { // Angular SPA hash routes
		href = "/" + href
	}
	bu, err := url.Parse(base)
	if err != nil {
		return ""
	}
	lu, err := url.Parse(href)
	if err != nil {
		return ""
	}
	u := bu.ResolveReference(lu)
	if sameHostOnly && !strings.EqualFold(bu.Hostname(), u.Hostname()) {
		return ""
	}
	if !strings.HasPrefix(u.Scheme, "http") {
		return ""
	}
	return u.String()
}

func pickTargetPath(base string, vary bool) string {
	paths := []string{
		"/", "/index.html", "/home", "/products", "/search", "/login", "/api/status",
		"/about", "/contact", "/feed", "/blog", "/docs", "/static/logo.png",
	}
	if vary {
		return strings.TrimRight(base, "/") + pick(paths)
	}
	return base
}

func reqFuzzParams() map[string]string {
	payloads := []string{
		"' OR '1'='1", "\" onmouseover=alert(1) x=\"", "../../etc/passwd", "%00", "<script>alert(1)</script>",
	}
	keys := []string{"q", "id", "search", "page", "debug", "sort"}
	out := map[string]string{}
	n := 1 + rand.Intn(3)
	for i := 0; i < n; i++ {
		out[pick(keys)] = pick(payloads)
	}
	return out
}

func dedup(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func pick[T any](arr []T) T { return arr[rand.Intn(len(arr))] }

var onceIP string

func pickIP(sticky bool) string {
	if sticky {
		if onceIP == "" {
			onceIP = randIPv4()
		}
		return onceIP
	}
	return randIPv4()
}
func randIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(223), rand.Intn(256), rand.Intn(256), 1+rand.Intn(254))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minDur(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func remaining(ctx context.Context) time.Duration {
	d, ok := ctx.Deadline()
	if !ok {
		return 0
	}
	r := time.Until(d)
	if r < 0 {
		return 0
	}
	return r
}
