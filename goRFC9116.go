package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var inputFile string = "security.txt"
var timeout int = 10 // Seconds to timeout

func checkURL(URL string) bool {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	resp, err := client.Get(URL)
	if err != nil {
		// fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func checkDomain(domain string, wg *sync.WaitGroup) {
	var loc [2]string
	var URL string
	var found bool
	defer wg.Done()

	// golang http will follow redirects, so no need to add [www.]
	loc[0] = "/.well-known/security.txt"
	loc[1] = "/security.txt"

	for _, path := range loc {
		found = false
		URL = "https://" + domain + path
		found = checkURL(URL)
		if found == true {
			break
		}
	}
	if found {
		fmt.Printf("[")
		fmt.Printf(("\033[32m"))
		fmt.Printf("OK")
		fmt.Printf(string("\033[0m]"))
		fmt.Printf("  %s %s\n", domain, URL)
	} else {
		fmt.Printf("[")
		fmt.Printf(("\033[31m"))
		fmt.Printf("NOK")
		fmt.Printf(string("\033[0m]"))
		fmt.Printf(" %s %s\n", domain, URL)
	}
}

func main() {
	var domain string
	var wg sync.WaitGroup

	file, err := os.Open(inputFile)
	if err != nil {
		log.Panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain = scanner.Text()
		wg.Add(1)
		go checkDomain(domain, &wg)
		time.Sleep(250 * time.Millisecond)
	}
	wg.Wait()
}
