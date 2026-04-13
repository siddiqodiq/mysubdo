package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var domain string

func banner() {
	fmt.Println(`
███╗   ███╗██╗   ██╗███████╗██╗   ██╗██████╗ ██████╗  ██████╗
████╗ ████║╚██╗ ██╔╝██╔════╝██║   ██║██╔══██╗██╔══██╗██╔═══██╗
██╔████╔██║ ╚████╔╝ ███████╗██║   ██║██████╔╝██║  ██║██║   ██║
██║╚██╔╝██║  ╚██╔╝  ╚════██║██║   ██║██╔══██╗██║  ██║██║   ██║
██║ ╚═╝ ██║   ██║   ███████║╚██████╔╝██████╔╝██████╔╝╚██████╔╝
`)
	fmt.Println("mysubdo - subdomain automation tool")
	fmt.Println()
}

func runCommand(cmd string, args []string, outfile string) {

	fmt.Println("[+] Running", cmd)

	out, err := os.Create(outfile)
	if err != nil {
		return
	}
	defer out.Close()

	command := exec.Command(cmd, args...)
	command.Stdout = out
	command.Run()
}

func assetfinder(wg *sync.WaitGroup) {
	defer wg.Done()
	runCommand("assetfinder", []string{"-subs-only", domain}, "assetfinder.txt")
}

func subfinder(wg *sync.WaitGroup) {
	defer wg.Done()
	runCommand("subfinder", []string{"-d", domain, "-all", "-recursive"}, "subfinder.txt")
}

func crtsh(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("[+] Fetching crt.sh")

	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data []map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}

	subs := make(map[string]bool)

	for _, entry := range data {

		names := entry["name_value"].(string)

		for _, n := range strings.Split(names, "\n") {

			n = strings.ReplaceAll(n, "*.", "")

			if strings.Contains(n, domain) {
				subs[n] = true
			}
		}
	}

	f, _ := os.Create("crtsh.txt")
	defer f.Close()

	for s := range subs {
		fmt.Fprintln(f, s)
	}
}

func anubis(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("[+] Fetching anubis")

	url := fmt.Sprintf("https://anubisdb.com/anubis/subdomains/%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var subs []string

	err = json.NewDecoder(resp.Body).Decode(&subs)
	if err != nil {
		return
	}

	f, _ := os.Create("anubis.txt")
	defer f.Close()

	for _, s := range subs {
		fmt.Fprintln(f, s)
	}
}

func hackertarget(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("[+] Fetching hackertarget")

	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	f, _ := os.Create("hackertarget.txt")
	defer f.Close()

	for scanner.Scan() {

		line := scanner.Text()

		if strings.Contains(line, "error") {
			continue
		}

		parts := strings.Split(line, ",")

		if len(parts) > 0 {
			fmt.Fprintln(f, parts[0])
		}
	}
}

func urlscan(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("[+] Fetching urlscan")

	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&data)

	results, ok := data["results"].([]interface{})
	if !ok {
		return
	}

	subs := make(map[string]bool)

	for _, r := range results {

		entry := r.(map[string]interface{})

		page := entry["page"].(map[string]interface{})

		d := page["domain"].(string)

		if strings.Contains(d, domain) {
			subs[d] = true
		}
	}

	f, _ := os.Create("urlscan.txt")
	defer f.Close()

	for s := range subs {
		fmt.Fprintln(f, s)
	}
}

func merge() {

	fmt.Println("[+] Merging results")

	files := []string{
		"assetfinder.txt",
		"subfinder.txt",
		"crtsh.txt",
		"anubis.txt",
		"hackertarget.txt",
		"urlscan.txt",
	}

	all := make(map[string]bool)

	for _, file := range files {

		f, err := os.Open(file)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {

			line := strings.TrimSpace(scanner.Text())

			if strings.Contains(line, domain) {
				all[line] = true
			}
		}

		f.Close()
	}

	out, _ := os.Create("allsubs.txt")
	defer out.Close()

	for sub := range all {
		fmt.Fprintln(out, sub)
	}
}

func httpx() {

	fmt.Println("[+] Running httpx")

	cmd := exec.Command(
		"httpx",
		"-list", "allsubs.txt",
		"-status-code",
		"-content-length",
		"-content-type",
		"-title",
		"-tech-detect",
		"-server",
		"-follow-host-redirects",
		"-random-agent",
		"-o", "httpx.txt",
	)

	cmd.Run()
}

func main() {

	if len(os.Args) < 2 {

		fmt.Println("Usage: mysubdo target.com")
		return
	}

	domain = os.Args[1]

	banner()

	var wg sync.WaitGroup

	wg.Add(6)

	go assetfinder(&wg)
	go subfinder(&wg)
	go crtsh(&wg)
	go anubis(&wg)
	go hackertarget(&wg)
	go urlscan(&wg)

	wg.Wait()

	merge()

	httpx()

	fmt.Println()
	fmt.Println("[+] Recon finished")
	fmt.Println("Subdomains saved to allsubs.txt")
	fmt.Println("Live hosts saved to httpx.txt")
}
