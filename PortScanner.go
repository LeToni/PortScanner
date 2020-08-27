package main

import (
	"fmt"
	"net"
	"sort"
)

func scanner(port int) int {

	address := fmt.Sprintf("scanme.nmap.org:%d", port)
	fmt.Println(address)
	connection, err := net.Dial("tcp", address)

	if err != nil {
		return 0
	}

	defer connection.Close()
	return port
}

func worker(jobs <-chan int, results chan<- int) {
	for n := range jobs {
		results <- scanner(n)
	}
}

func main() {

	jobs := make(chan int, 100)
	results := make(chan int)
	var openPorts []int

	for i := 0; i < cap(jobs); i++ {
		go worker(jobs, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			jobs <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(jobs)
	close(results)

	sort.Ints(openPorts)
	for _, openPort := range openPorts {
		fmt.Printf("%d open\n", openPort)
	}
}
