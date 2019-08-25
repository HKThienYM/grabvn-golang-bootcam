package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/avast/retry-go"
	"github.com/myteksi/hystrix-go/hystrix"
)

const url = "http://localhost:8080/test"

func main() {
	initCircuitBreaker()

	for i := 1; i < 101; i++ {
		fmt.Printf("\nStart request: %d\n", i)
		hystrix.Do("my_command", func() error {
			return retry.Do(
				sendRequest,
				retry.Attempts(3),
				retry.OnRetry(func(n uint, err error) {
					fmt.Println("retry", n+1, "time(s)")
				}),
			)
		}, nil)
	}
}

func initCircuitBreaker() {
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:                     1000,
		MaxConcurrentRequests:       100,
		ErrorPercentThreshold:       1,
		QueueSizeRejectionThreshold: 100,
		SleepWindow:                 50,
	})
}

func sendRequest() error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	if resp.StatusCode == 500 {
		return errors.New("Request failed")
	}

	return nil
}
