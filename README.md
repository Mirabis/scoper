# scoper 

Small, fast, simple tool for cleaning other tool output (IPs, URLs, hostnames) to only in-scope IPs

You feed **SCOPER** IPs, URLs or Hostnames and it returns them only if they are inside the defined scope

This can be a useful way of finding usernames belonging to a company using Azure Directories from their e-mail addresses.

## Installation

```sh
go get github.com/mirabis/scoper
```

## Usage
The most basic usage is to simply pipe a list of hosts, ips or urls into the tool, for example:

```sh
mirabis~$ cat osint-domains-and-ips.txt | scoper -c scope.txt
STDOUT will only contain ip\'s in-scope or hosts/urls that resolve to ip\'s in scope
...
```

### Parameters

```sh
mirabis~$ scoper -h

Usage:
  scoper [OPTIONS]

███████╗ ██████╗ ██████╗ ██████╗ ███████╗██████╗ 
██╔════╝██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔══██╗
███████╗██║     ██║   ██║██████╔╝█████╗  ██████╔╝d
╚════██║██║     ██║   ██║██╔═══╝ ██╔══╝  ██╔══██╗
███████║╚██████╗╚██████╔╝██║     ███████╗██║  ██║

Application Options: (/* windows, -* Unix)
  /c, /cidrs:    CIDRS to match with, line separated
  /t, /threads:  Number of concurrent threads (default: 20)

Help Options:
  /?             Show this help message
  /h, /help      Show this help message
```


## Credits
- [tomnomnom](https://github.com/tomnomnom/unfurl/blob/master/main.go) my inspiration for the tool and the switch to golang

### Contribution & License
You can contribute in following ways:

- Report bugs
- Give suggestions to make it better (I'm new to golang)
- Fix issues & submit a pull request

Do you want to have a conversation in private? Hit me up on my [twitter](https://twitter.com/iMirabis/), inbox is open :)

**scoper** is licensed under [GPL v3.0 license](https://www.gnu.org/licenses/gpl-3.0.en.html)