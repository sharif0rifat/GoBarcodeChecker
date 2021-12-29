package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func landingpage(rw http.ResponseWriter, req *http.Request) {
	msg := "go to this link(http://localhost:8081/GenerateBarcodes) to generate barcode" + "\r\n"
	msg += "go to this link(http://localhost:8081/SearchBarcode?barcode=some_value) to search barcode" + "\r\n"
	fmt.Fprint(rw, msg)
}
func main() {
	http.HandleFunc("/SearchBarcode", SearchBarcode)
	http.HandleFunc("/GenerateBarcodes", GenerateBarcodes)
	http.HandleFunc("/", landingpage)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

//#region======Generate Barcode====
func GenerateBarcodes(rw http.ResponseWriter, req *http.Request) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		var z = "Company_" + strconv.Itoa(i)
		go GenerateBarcode(z, &wg)
	}
	wg.Wait()

	content, err := ioutil.ReadFile(".\\Log\\Logs.txt")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(rw, string(content))
}
func GenerateBarcode(fileName string, wg *sync.WaitGroup) {
	msg := "Started Generating:      " + time.Now().String() + "\r\n"
	barcode := []string{}
	for j := 1; j <= 100000; j++ {
		var s = fmt.Sprintf("%06d", j+2)
		barcode = append(barcode, fileName+"_"+s)
	}
	Writefile("Database\\"+fileName+".txt", barcode)
	msg += "Completed Generating :   " + time.Now().String() + "\r\n" + "\r\n"
	Writefile("Log\\Logs.txt", nil, msg)
	wg.Done()
}

func Writefile(filepath string, barcode []string, message ...string) {
	fmt.Println("Printing file")

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if len(barcode) > 0 {
		for _, data := range barcode {
			_, _ = f.WriteString(data + "\n")
		}
	} else {
		if _, err := f.WriteString(message[0]); err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

//#endregion

// #region====Searching Operation==
func SearchBarcode(rw http.ResponseWriter, req *http.Request) {
	params, ok := req.URL.Query()["barcode"]
	if ok {
		result := SearchInFiles(params[0])
		if ok {
			fmt.Fprint(rw, result)
		} else {
			fmt.Fprint(rw, "Barcode not found")
		}
	} else {
		fmt.Fprint(rw, "The passed param is not valid, you must send  '?barcode=some_value from Database folder'")
	}
}

func SearchInFiles(barcode string) string {
	directory := ".\\Database"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	responseResult := ""
	for _, f := range files {
		result, ok := SearchInFile(directory+"\\"+f.Name(), barcode)
		if ok {
			responseResult = "Barcode found in file: " + f.Name() + result
			break
		} else {
			responseResult = "Barcode not found"
		}
	}
	fmt.Println("the response: " + responseResult)
	return responseResult
}
func SearchInFile(path string, barcode string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening file: " + err.Error())
		return "Error happened reading file", false
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	fmt.Println(barcode)
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), barcode) {
			fmt.Println(scanner.Text())
			return " line:" + strconv.Itoa(line), true
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		// Handle the error
	}
	return "Barcode Not Fount", false
}

// #endregion
