package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

const timeout = 1
const cep_regex = "([0-9]){5}-([0-9]){3}"

func main() {
	for _, cep := range os.Args[1:] {
		viaCepChan := make(chan string)
		apiCepChan := make(chan string)

		isCepValid, err := regexp.MatchString(cep_regex, cep)
		if err != nil {
			panic(err)
		}

		if !isCepValid {
			fmt.Printf("Input: %s - ERROR: CEP Format should be 12345-678.\n", cep)
			continue
		}

		go getAPIResponse("https://cdn.apicep.com/file/apicep/"+cep+".json", apiCepChan)
		go getAPIResponse("http://viacep.com.br/ws/"+cep+"/json/", viaCepChan)

		msg := ""

		select {
		case resp := <-viaCepChan:
			msg = fmt.Sprintf("Via Cep Response: %s", resp)

		case resp := <-apiCepChan:
			msg = fmt.Sprintf("API Cep Response: %s", resp)

		case <-time.After(time.Second * timeout):
			msg = "Timeout"
		}

		fmt.Printf("Input: %s - %s\n\n", cep, msg)

		time.Sleep(time.Second * 2)
	}
}

func getAPIResponse(url string, apiChan chan<- string) {
	req, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return
	}

	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	apiChan <- string(res)
}
