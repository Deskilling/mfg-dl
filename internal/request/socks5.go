package request

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/net/proxy"
)

func TestSocks5Latency(target string, proxyAddr string) (time.Duration, error) {
	client, _ := getSock5Client(proxyAddr)

	req, err := http.NewRequest("HEAD", target, nil)
	if err != nil {
		return 0, err
	}

	log.Debug("start measure", "proxyAddr", proxyAddr, "target", target)
	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	log.Debug("finished")
	return time.Since(start), nil
}

func getSock5Client(proxyAddr string) (client *http.Client, err error) {
	mu.Lock()
	defer mu.Unlock()

	log.Debug("check for socks5 client", "proxyAddr", proxyAddr)
	client, ok := clients[proxyAddr]
	if ok {
		log.Debug("already found client", "prodxyAddr", proxyAddr)
		return client, nil
	}

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	client = &http.Client{
		Transport: &http.Transport{
			Dial:                dialer.Dial,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 15 * time.Second,
	}

	clients[proxyAddr] = client
	log.Debug("created socks5 client", "proxyAddr", proxyAddr)
	return client, nil
}
