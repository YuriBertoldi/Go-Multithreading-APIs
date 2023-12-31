package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MessageAPI struct {
	Mensagem string
	URL      string
}

var Cep = "01153000"

func main() {
	c1 := make(chan MessageAPI)
	c2 := make(chan MessageAPI)

	go RequisitarBrasilAPI(c1)
	go RequisitarViaCEP(c2)

	select {
	case Msg := <-c1:
		ExibirLogCommand(Msg)

	case Msg := <-c2:
		ExibirLogCommand(Msg)

	case <-time.After(time.Second):
		log.Fatalln("Timeout")
	}
}

func ExibirLogCommand(Log MessageAPI) {
	var mensagemCommand = "API utilizada: %s\nResultado da request: %s\n"
	fmt.Printf(mensagemCommand, Log.URL, Log.Mensagem)
}

func RequisitarBrasilAPI(chanel chan<- MessageAPI) {
	URL := "https://brasilapi.com.br/api/cep/v1/" + Cep
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	msg := MessageAPI{string(body), URL}
	chanel <- msg
}

func RequisitarViaCEP(chanel chan<- MessageAPI) {
	URL := "http://viacep.com.br/ws/" + Cep + "/json/"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	msg := MessageAPI{string(body), URL}
	chanel <- msg
}
