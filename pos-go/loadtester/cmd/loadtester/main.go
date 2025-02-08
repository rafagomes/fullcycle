package main

import (
	"flag"
	"fmt"
	"os"

	"loadtester/internal/usecase"
)

func main() {
	urlFlag := flag.String("url", "", "URL do serviço a ser testado")
	requestsFlag := flag.Int("requests", 0, "Número total de requests")
	concurrencyFlag := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *urlFlag == "" || *requestsFlag <= 0 || *concurrencyFlag <= 0 {
		fmt.Println("Uso: --url=<URL> --requests=<número total de requests> --concurrency=<número de chamadas simultâneas>")
		os.Exit(1)
	}

	loadTester := usecase.NewLoadTester()
	report, err := loadTester.Run(*urlFlag, *requestsFlag, *concurrencyFlag)
	if err != nil {
		fmt.Printf("Erro durante o teste: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("===== Relatório de Teste de Carga =====")
	fmt.Printf("Tempo total de execução: %s\n", report.TotalTime)
	fmt.Printf("Total de requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", report.Status200Count)
	fmt.Printf("Requests com erro: %d\n", report.ErrorCount)
	fmt.Println("Distribuição dos códigos de status:")
	for code, count := range report.StatusDistribution {
		if code == 0 {
			fmt.Printf("Erro: %d\n", count)
		} else {
			fmt.Printf("%d: %d\n", code, count)
		}
	}
}
