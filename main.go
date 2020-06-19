package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var (
	HTTP_PORT = ":8000"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		if r.Method != "POST" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("must POST"))
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		type params struct {
			Location string `json:"location"`
		}
		p := params{}
		err = json.Unmarshal(b, &p)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		denoPath, err := exec.LookPath("deno")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		args := []string{"run", p.Location}

		// TODO enforce a timeout? Make this configurable?
		ctx, _ := context.WithTimeout(r.Context(), 30*time.Second)
		cmd := exec.CommandContext(ctx, denoPath, args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err := cmd.Start(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		scanner := bufio.NewScanner(stdout)
		flusher, isFlusher := w.(http.Flusher)
		for scanner.Scan() {
			w.Write(append(scanner.Bytes(), []byte("\n")...))
			if isFlusher {
				flusher.Flush()
			}
		}

		// TODO not doing anything with the output on stderr right now, maybe report back somehow?
		_, err = ioutil.ReadAll(stderr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err := cmd.Wait(); err != nil {
			switch err.(type) {
			case *exec.ExitError:
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}

		log.Printf("ran: %s for %s\n", p.Location, time.Now().Sub(startTime))

	})

	log.Fatal(http.ListenAndServe(HTTP_PORT, nil))
}
