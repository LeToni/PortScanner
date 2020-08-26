package main

import (
	"fmt"
	"net"
	"sort"
)

func scanner(ports, results chan int) {
	for port := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", port)
		connection, err := net.Dial("tcp", address)

		if err != nil {
			results <- 0
			continue
		}

		connection.Close()
		results <- port
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)

	var openPorts []int

	for i := 0; i < cap(ports); i++ {
		go scanner(ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openPorts)
	for _, openPort := range openPorts {
		fmt.Printf("%d open\n", openPort)
	}
}
