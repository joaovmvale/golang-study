package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoringRepetitions = 3
const logsFileName = "logs.txt"
const sitesFileName = "sites.txt"

func showGreetings() {
	fmt.Println("Welcome to the site scanner")
}

func showMenu() {
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("0 - Exit program")
	fmt.Println("**************************")
	fmt.Print("Select your option: ")
}

func readInput() int {
	var option int
	fmt.Scan(&option)

	return option
}

func readFile(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		fmt.Print("An error has ocurred during the file read: ")
		fmt.Println(err)
	}

	var sites []string

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}

	}

	file.Close()

	return sites
}

func writeLogs(url string, statusCode int) {
	file, err := os.OpenFile(logsFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		fmt.Print("An error has ocurred during the logs writing: ")
		fmt.Println(err)
	}

	currentDate := time.Now().Format("02/01/2006 15:04:05")
	logString := currentDate + " - SITE: " + url + " | STATUS CODE: " + strconv.Itoa(statusCode) + "\n"

	file.WriteString(logString)

	file.Close()
}

func testWebsite(url string) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Print("An error has ocurred during the request: ")
		fmt.Println(err)
	}

	statusCode := response.StatusCode

	if statusCode >= 200 && statusCode < 300 {
		fmt.Println("O site:", url, "foi carregado com sucesso!")
		fmt.Println("STATUS CODE:", statusCode)
	} else {
		fmt.Println("O site:", url, "apresentou problemas!")
		fmt.Println("STATUS CODE:", statusCode)
	}

	writeLogs(url, statusCode)
}

func startMonitoring() {
	fmt.Println("Start monitoring...")

	sites := readFile(sitesFileName)

	for i := 0; i < monitoringRepetitions; i++ {
		fmt.Println("========================")
		fmt.Println("Monitoring counter:", i+1)
		fmt.Println("========================")
		fmt.Println()

		for idx, url := range sites {
			fmt.Println("Making request to the site", idx, url)
			testWebsite(url)
			fmt.Println()
		}

		// Awaiting 5 seconds on each monitoring
		time.Sleep(5 * time.Second)
	}
}

func showLogs() {
	fmt.Println("Printing logs...")

	file, err := ioutil.ReadFile(logsFileName)

	if err != nil {
		fmt.Print("An error has ocurred during the logs file read: ")
		fmt.Println(err)
	}

	fmt.Println(string(file))

	fmt.Println("End of the file!")
}

func main() {
	showGreetings()

	for {
		showMenu()

		option := readInput()

		switch option {
		case 1:
			startMonitoring()
		case 2:
			showLogs()
		case 0:
			fmt.Println("Exiting program...")
			os.Exit(0)
		default:
			fmt.Println("Option not available!")
			os.Exit(-1)
		}

		fmt.Println()
	}
}
