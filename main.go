package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// Creates a new file upload http request with optional extra params
func buildUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func uploadFile(path string) {
	extraParams := map[string]string{
		"user":     "gary",
		"password": "asdf",
	}
	hostname := "127.0.0.1:8080"
	request, err := buildUploadRequest("http://"+hostname+"/upload_game", extraParams, "file", path)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(body)
}

/*
func playMatch() {
	p1 := exec.Command("lczero")
  p1_in, _ := p1.StdinPipe()
  p1_out, _ := p1.StdoutPipe()
  p1.Start()
  p1.Write("...")
}
*/

func train() {
	pid := 1

	dir, _ := os.Getwd()
	train_dir := path.Join(dir, fmt.Sprintf("data-%v", pid))
	if _, err := os.Stat(train_dir); err == nil {
		err = os.RemoveAll(train_dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	num_games := 1
	train_cmd := fmt.Sprintf("--start=train %v %v", pid, num_games)
	cmd := exec.Command(path.Join(dir, "lczero"), "--weights=weights.txt", "--randomize", "-n", "-t1", train_cmd)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stderr)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	train_file := path.Join(train_dir, "training.0.gz")
	uploadFile(train_file)
}

func main() {
	train()
}
