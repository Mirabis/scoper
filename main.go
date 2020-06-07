package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"sync"

	valid "github.com/asaskevich/govalidator"
	flags "github.com/jessevdk/go-flags"
)

var scopeSubnets []*net.IPNet

var opts struct {
	CIDRS   string `short:"c" long:"cidrs" description:"CIDRS to match with, line separated"`
	Threads int    `short:"t" long:"threads" default:"40" description:"Number of concurrent threads"`
}

func main() {
	// get Arguments
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		flags.WroteHelp(err)
		panic(err)
	} else if opts.CIDRS == "" {
		fmt.Println("Scope is required in form of --cidrs (-c)")
		os.Exit(1)
	}
	// Get Scope File's once STDIN
	if opts.CIDRS != "" {
		file, _ := os.Open(opts.CIDRS)
		scanner := bufio.NewScanner(bufio.NewReader(file))
		for scanner.Scan() {
			_, subnet, _ := net.ParseCIDR(scanner.Text())
			scopeSubnets = append(scopeSubnets, subnet)
		}
		file.Close()
	}
	//get STDIN
	numWorkers := opts.Threads
	work := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		// for each line on STDIN
		for scanner.Scan() {
			work <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "ERR: reading input failed:", err)
		}
		close(work)
	}()
	// Create a waiting group
	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go doWork(work, wg) //Schedule the work
	}
	wg.Wait() //Wait for it all to complete
}

func doWork(work chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	//for each item to check scope
	for toScope := range work {
		procesInput(toScope, nil)
	}
}

func procesInput(input string, override []*net.IPNet) bool {
	if override != nil {
		scopeSubnets = override
	}

	if ip := net.ParseIP(input); ip != nil {
		for _, scope := range scopeSubnets {
			if scope.Contains(ip) {
				println(input)
				return true
			}
		}
	} else {
		var ips []net.IP
		var err error
		if valid.IsDNSName(input) {
			ips, err = net.LookupIP(input)

		} else if valid.IsURL(input) {
			if u, _ := url.Parse(input); u != nil {
				ips, err = net.LookupIP(u.Hostname())
			}
		}
		if err == nil {
			for _, ip := range ips {
				for _, scope := range scopeSubnets {
					if scope.Contains(ip) {
						println(input)
						return true
					}
				}
			}
		}
	}
	return false
}
