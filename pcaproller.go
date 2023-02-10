package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

const endpoint = "http://your-remote-endpoint.com/api/upload"

func main() {
	for {
		err := captureAndUpload()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(60 * time.Second)
	}
}

func captureAndUpload() error {
	// Use tcpdump to capture the packets
	cmd := exec.Command("tcpdump", "-w", "-", "-c", "100", "not port 22")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Upload the captured packets to the remote endpoint
	resp, err := http.Post(endpoint, "application/octet-stream", &out)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response from the endpoint: %s %s", resp.Status, string(body))
	}

	return nil
}

