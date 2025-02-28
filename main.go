package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

func main() {
	clusterbomb := flag.Bool("c", false, "Enable clusterbomb mode: substitute all parameters with fuzz_word")
	appendMode := flag.Bool("a", false, "Append fuzz_word to the existing parameter value")
	ignorePath := flag.Bool("ignore-path", false, "Ignore the path when considering duplicates")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <fuzz_word>\n", os.Args[0])
		os.Exit(1)
	}
	fuzzWord := flag.Arg(0)

	seen := make(map[string]bool)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		u, err := url.Parse(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing URL '%s': %v\n", line, err)
			continue
		}

		origQuery := u.Query()

		if *clusterbomb {
			params := make([]string, 0, len(origQuery))
			for p := range origQuery {
				params = append(params, p)
			}
			sort.Strings(params)
			key := u.Hostname()
			if !*ignorePath {
				key += u.EscapedPath()
			}
			key += "?" + strings.Join(params, "&")
			if seen[key] {
				continue
			}
			seen[key] = true

			for param, values := range origQuery {
				if *appendMode {
					origQuery.Set(param, values[0]+fuzzWord)
				} else {
					origQuery.Set(param, fuzzWord)
				}
			}
			u.RawQuery = origQuery.Encode()
			fmt.Println(u.String())
		} else {
			params := make([]string, 0, len(origQuery))
			for p := range origQuery {
				params = append(params, p)
			}
			sort.Strings(params)
			for _, param := range params {
				newQuery := url.Values{}
				for k, v := range origQuery {
					newQuery[k] = append([]string(nil), v...)
				}
				if *appendMode {
					newQuery.Set(param, newQuery.Get(param)+fuzzWord)
				} else {
					newQuery.Set(param, fuzzWord)
				}
				newURL := *u
				newURL.RawQuery = newQuery.Encode()
				fmt.Println(newURL.String())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
