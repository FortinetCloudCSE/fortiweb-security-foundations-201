// simple tool to prime the anomaly detection module on Fortiweb
// the tool will go after "q" parameter found at /rest/products/search?q on the juiceshop app
// it will keep sending the same value " apple"to build a model that will expect the string as input
// //Wondyrad T version 1.0 07/25
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var searchTerms = []string{
	"apple",
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
	"Mozilla/5.0 (Linux; Android 13; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.18363",
	"Mozilla/5.0 (Linux; U; Android 8.0.0; en-us; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.131 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/117.0.5938.132 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
	"Mozilla/5.0 (Linux; Android 12; SM-A505F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 9; Redmi Note 8T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.5615.137 Mobile Safari/537.36",
}

// Generate a random IPv4 address
func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(254)+1, rand.Intn(254)+1, rand.Intn(254)+1, rand.Intn(254)+1)
}

func makeRequest(baseURL string, client *http.Client) error {
	term := searchTerms[rand.Intn(len(searchTerms))]
	q := url.QueryEscape(term)
	fullURL := fmt.Sprintf("%s/rest/products/search?q=%s", baseURL, q)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	req.Header.Set("X-Forwarded-For", randomIP())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Juice Shop base URL (e.g. https://669.fwebtraincse.com): ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)

	if baseURL == "" || !strings.HasPrefix(baseURL, "http") {
		fmt.Println("Invalid URL. Exiting.")
		return
	}
	//wt modify this to control number of requests, duration and concurent request  .;/;;;;//
	const totalRequests = 5000
	const durationMinutes = 30
	const workers = 15

	baseDelay := time.Duration(durationMinutes*60*1000/totalRequests) * time.Millisecond
	jitter := func() time.Duration {
		return baseDelay + time.Duration(rand.Intn(200))*time.Millisecond // up to 200ms jitter
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	client := &http.Client{Timeout: 10 * time.Second}

	var sentCount, successCount, failCount int
	startTime := time.Now()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Progress logger
	go func() {
		for range ticker.C {
			mu.Lock()
			fmt.Printf("⏳ I have sent %d requests so far. I have %d to go... be patient.\n", sentCount, totalRequests-sentCount)
			mu.Unlock()
		}
	}()

	requestChan := make(chan struct{}, workers)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func() {
			requestChan <- struct{}{}
			err := makeRequest(baseURL, client)

			mu.Lock()
			sentCount++
			if err != nil {
				failCount++
			} else {
				successCount++
			}
			mu.Unlock()

			<-requestChan
			wg.Done()
		}()

		time.Sleep(jitter())
	}

	wg.Wait()
	elapsed := time.Since(startTime).Seconds()
	fmt.Printf("\n✅ All %d requests completed.\n\n", totalRequests)

	// Print summary
	fmt.Println("📊 Summary Report:")
	fmt.Printf("   ✅ Successful requests : %d\n", successCount)
	fmt.Printf("   ⚠️  Failed requests     : %d\n", failCount)
	fmt.Printf("   ⏱  Total time (sec)    : %.2f\n", elapsed)
	fmt.Printf("   🔢 Avg RPS             : %.2f\n", float64(totalRequests)/elapsed)
}
