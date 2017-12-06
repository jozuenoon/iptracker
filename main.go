package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var url = flag.String("url", "http://ifconfig.me/all.json", "URL to get IP address")
var path = flag.String("path", "$HOME/ipmon", "path to save logs")

func init() {
	flag.Parse()

	// This part fixes path if it get default value.
	// In other cases bash will manage this for us.
	if *path == "$HOME/ipmon" {
		val, ok := os.LookupEnv("HOME")
		if !ok {
			log.Fatal("unable to interpolate $HOME env, check if it is set")
		} else {
			*path = filepath.Join(val, "ipmon")
		}
	}
}

type ipData struct {
	IPAddr    string    `json:"ip_addr"`
	Timestamp time.Time `json:"ts"`
	Error     *string   `json:"error"`
}

func (id ipData) String() string {
	if id.Error != nil {
		return fmt.Sprintf("%s |ERROR| %s", id.Timestamp.Format("2006-01-02 15:04:05"), *id.Error)
	}
	return fmt.Sprintf("%s |INFO| %s", id.Timestamp.Format("2006-01-02 15:04:05"), id.IPAddr)
}

func saveLog(logrow *ipData, basePath string) error {
	date := logrow.Timestamp.Format("2006-01-02")
	ds := strings.Split(date, "-")
	if len(ds) != 3 {
		return errors.New("smth very bad happen to date")
	}

	fpath := filepath.Join(basePath, ds[0], ds[1])
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		os.MkdirAll(fpath, os.ModePerm)
	}

	logfile := filepath.Join(fpath, ds[2]+".log")
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(logrow.String() + "\n")
	return nil
}

func getIPData(url string) (*ipData, error) {
	var data ipData
	data.Timestamp = time.Now()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("User-Agent", "curl/7.47.0")
	resp, err := client.Do(req)
	if err != nil {
		es := err.Error()
		return &ipData{Timestamp: time.Now(), Error: &es}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		es := err.Error()
		return &ipData{Timestamp: time.Now(), Error: &es}, err
	}
	return &data, nil
}

func main() {
	di, err := getIPData(*url)
	if err != nil {
		log.Println(di)
	}
	err = saveLog(di, *path)
	if err != nil {
		log.Println(di)
		log.Println(err)
	}
}
