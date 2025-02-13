package main

import (
	"fmt"
	"sync"
	"time"
)

type hostsWithProblems struct {
	Host     Host
	Problems []string
}

var (
	conf config
)

func main() {
	data, err := conf.parseConfigFile("./config/config.yaml")
	if err != nil {
		fmt.Println("Error parsing config.yaml")
		return
	}

	var (
		apiUrl        = data.ApiUrl
		token         = data.Token
		maxConcurrent = data.MaxConcurrent
	)

	startTime := time.Now()
	fmt.Printf("Сбор данных начался в %v\n", startTime.Format("2006-01-02 15:04:05"))

	hosts, err := getHosts(apiUrl, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	hostsWithProblemsChan := make(chan hostsWithProblems, len(hosts))
	sem := make(chan struct{}, maxConcurrent)

	for _, host := range hosts {
		wg.Add(1)
		go func(host Host) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			problems, err := getProblems(apiUrl, token, host)
			if err != nil {
				fmt.Println(err)
				return
			}

			var realProblems []string
			if len(problems) > 0 {
				for _, problem := range problems {
					status, _ := getTriggerStatus(apiUrl, token, problem.ObjectId)
					if status != "0" {
						continue
					}
					realProblems = append(realProblems, problem.Name)
				}
			}

			if len(realProblems) > 0 {
				hostsWithProblemsChan <- hostsWithProblems{Host: host, Problems: realProblems}
			}
		}(host)
	}

	go func() {
		wg.Wait()
		close(hostsWithProblemsChan)
	}()

	count, err := exportToCsv(fmt.Sprintf("hosts_with_problems_%s.csv", time.Now().Format("2006-01-02_15-04-05")), hostsWithProblemsChan)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Данные успешно записаны в hosts_with_problems.csv")
	}

	fmt.Printf("Экспортировано хостов: %v\n", count)

	finishTime := time.Since(startTime)
	fmt.Printf("Время выполнения скрипта: %.2f секунд\n", finishTime.Seconds())
}
