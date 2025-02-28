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
	// Definição dos flags:
	// -c: ativa o modo clusterbomb (substitui todos os parâmetros)
	// -a: ativa o modo append (concatena a fuzz_word ao valor existente)
	// -ignore-path: ignora o path na verificação de duplicatas (no modo clusterbomb)
	clusterbomb := flag.Bool("c", false, "Ativa o modo clusterbomb: substitui todos os parâmetros pela fuzz_word")
	appendMode := flag.Bool("a", false, "Anexa a fuzz_word ao valor existente do parâmetro")
	ignorePath := flag.Bool("ignore-path", false, "Ignora o path ao considerar duplicatas")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Uso: %s [opções] <fuzz_word>\n", os.Args[0])
		os.Exit(1)
	}
	fuzzWord := flag.Arg(0)

	// Para o modo clusterbomb, usamos um mapa para evitar duplicatas
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		u, err := url.Parse(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao interpretar URL '%s': %v\n", line, err)
			continue
		}

		origQuery := u.Query()

		if *clusterbomb {
			// Cria uma chave baseada em hostname, (opcionalmente) path e parâmetros ordenados
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

			// Substitui TODOS os parâmetros pela fuzz_word (ou os concatena se -a estiver ativo)
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
			// Modo normal: gera uma URL para cada combinação, substituindo individualmente cada parâmetro
			params := make([]string, 0, len(origQuery))
			for p := range origQuery {
				params = append(params, p)
			}
			sort.Strings(params)
			for _, param := range params {
				// Cria uma cópia dos valores da query
				newQuery := url.Values{}
				for k, v := range origQuery {
					newQuery[k] = append([]string(nil), v...)
				}
				// Substitui apenas o parâmetro atual
				if *appendMode {
					newQuery.Set(param, newQuery.Get(param)+fuzzWord)
				} else {
					newQuery.Set(param, fuzzWord)
				}
				// Prepara a nova URL com a query modificada
				newURL := *u
				newURL.RawQuery = newQuery.Encode()
				fmt.Println(newURL.String())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro lendo a entrada: %v\n", err)
	}
}
