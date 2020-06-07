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
	CIDRS   string `short:"c" long:"cidrs" required:"true" description:"CIDRS to match with, line separated"`
	Threads int    `short:"t" long:"threads" default:"20" description:"Number of concurrent threads"`
	Verbose bool   `short:"v" long:"verbose" description:"Turns on verbose logging (shows the scope and resolved IP(s))"`
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
		file, err := os.Open(opts.CIDRS)
		scanner := bufio.NewScanner(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERR: could not load CIDR file please try again", err)
			os.Exit(1)
		}
		for scanner.Scan() {
			r := scanner.Text()
			_, subnet, err := net.ParseCIDR(r)
			if err != nil {
				scopeSubnets = append(scopeSubnets, subnet)
			} else {
				fmt.Fprintln(os.Stderr, "ERR: Error parsing", r, "with fault", err)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "ERR: reading cidrs failed:", err)
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
	defer wg.Done() // done after we exit this func
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
				if opts.Verbose == true {
					fmt.Printf("%s, is within scope of %s\n", input, scope)
				} else {
					fmt.Println(input)
				}
				return true
			}
		}
	} else {
		var ips []net.IP
		var err error
		if valid.IsDNSName(input) {
			ips, err = net.LookupIP(input)
			if err != nil {
				if opts.Verbose {
					fmt.Fprintln(os.Stderr, "ERR: identified DNSName but couldn't parse IP's without error", err)
				}
				return false
			}
		} else if valid.IsURL(input) {
			if u, _ := url.Parse(input); u != nil {
				ips, err = net.LookupIP(u.Hostname())
				if err != nil {
					if opts.Verbose {
						fmt.Fprintln(os.Stderr, "ERR: identified URL but couldn't parse IP's without error", err)
					}
					return false
				}
			}
		}
		if ips != nil {
			for _, ip := range ips {
				for _, scope := range scopeSubnets {
					if scope.Contains(ip) {
						if opts.Verbose {
							fmt.Printf("%s, resolved to %s (%v) which is within scope of %s\n", input, ip, ips, scope)
						} else {
							fmt.Println(input)
						}
						return true
					}
				}
			}
		}
	}
	return false
}
