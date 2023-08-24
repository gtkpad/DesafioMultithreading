package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCepResponse struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func requestViaCep(ch chan<- ViaCepResponse, cep string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		close(ch)
		return
	}

	var response ViaCepResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		close(ch)
		return
	}

	ch <- response
}

func requestApiCep(ch chan<- ApiCepResponse, cep string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		close(ch)
		return
	}

	var response ApiCepResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		close(ch)
		return
	}

	ch <- response
}

func main() {

	chViaCep := make(chan ViaCepResponse)
	chApiCep := make(chan ApiCepResponse)

	go requestViaCep(chViaCep, "96170-000")
	go requestApiCep(chApiCep, "96170-000")

	select {
	case resp := <-chViaCep:
		fmt.Printf("[ViaCep]: %+v", resp)
	case resp := <-chApiCep:
		fmt.Printf("[ApiCep]: %+v", resp)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}
