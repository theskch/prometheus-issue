package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"go.uber.org/ratelimit"
)

const url = "http://localhost:8080/v1/info"

const (
	maxIdleConnsPerHost = 100
)

func main() {

	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
		},
	}

	rl := ratelimit.New(2)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		rl.Take()
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := pickRandom([]int{1, 2, 3, 4, 5})
			if err := makeRequest(httpClient, id); err != nil {
				fmt.Printf("error: %s\n", err)
			}
		}()
	}
	wg.Wait()

	fmt.Println("DONE!")
}

func pickRandom(numbers []int) int {
	randomIndex := rand.Intn(len(numbers))

	return numbers[randomIndex]
}

func makeRequest(client http.Client, id int) error {
	resp, err := http.Get(fmt.Sprintf("%s/%d", url, id))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not 200: %d", resp.StatusCode)
	}

	return nil
}
