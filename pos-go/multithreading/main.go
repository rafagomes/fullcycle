package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func getAddress(url *string) string {
	resp, err := http.Get(*url)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func getFromBrasilAPI(cep string, data chan Address) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s/", cep)
	body := getAddress(&url)

	// time.Sleep(time.Second * 2)

	data <- Address{
		url:  url,
		body: string(body),
	}
}

func getFromViaCep(cep string, data chan Address) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)
	body := getAddress(&url)

	// time.Sleep(time.Second * 2)

	data <- Address{
		url:  url,
		body: string(body),
	}
}

type Address struct {
	url  string
	body string
}

func main() {
	var cep string = "24325330"
	channel1 := make(chan Address)
	channel2 := make(chan Address)

	go getFromBrasilAPI(cep, channel1)
	go getFromViaCep(cep, channel2)

	select {
	case resp1 := <-channel1:
		fmt.Printf("URL: %s and data: %s \n", resp1.url, resp1.body)
	case resp2 := <-channel2:
		fmt.Printf("URL: %s and data: %s \n", resp2.url, resp2.body)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}
