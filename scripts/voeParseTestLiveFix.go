//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"

	"mfg-dl/internal/sites/modules/voe"
)

var urls = []string{
	"https://bryantenunder.com/e/o6pq69uxu7wl",
	"https://bryantenunder.com/e/6q1ngfkyvebf",
	"https://bryantenunder.com/e/xqo6271z5tev",

	"https://bryantenunder.com/e/lutc9vijlalf",
	"https://bryantenunder.com/e/wark5ec0rq4w",
	"https://bryantenunder.com/e/kukevkqehqyg",

	"https://bryantenunder.com/e/dfwfjaxudbea",
}

func main() {
	client := &http.Client{}
	for _, url := range urls {
		resp, _ := client.Get(url)
		resp.Body.Close()

		resolvedUrl := resp.Request.URL.String()

		resp2, _ := client.Get(resolvedUrl)
		body, _ := io.ReadAll(resp2.Body)
		resp2.Body.Close()

		redirect, _ := voe.VoeRedirect(string(body))
		fmt.Printf("\"%s\",\n", redirect)
	}
}
