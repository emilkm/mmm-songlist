package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var (
	defaultSheetID = "1f11V7FE24SdahFwFJMDPjBPJ4rxV9mxs2xm1ufMhCUw"
	newColumns     = []string{"artist", "song", "userdifficulty", "usernotes", "mmmdifficulty", "mmmnotes", "style", "learn", "tuning", "mmmtutorial", "timesignature"}
)

type Row map[string]string

func main() {
	// Parse command-line flags
	doPush := flag.Bool("push", false, "also git add/commit/push data.json")
	sheetID := flag.String("sheetid", defaultSheetID, "Google Sheet ID to fetch data from")
	flag.Parse()
	// Step 1: Download CSV
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=Glossary", *sheetID)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	csvData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	os.WriteFile("data.csv", csvData, 0644)

	// Step 2: Read CSV, skip the first line (header), keep first 11 columns, rename
	f, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := csv.NewReader(bufio.NewReader(f))

	// Always discard the first line, regardless of its content
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	var records []Row
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip bad rows
		}
		if len(rec) < len(newColumns) {
			continue // skip incomplete rows
		}
		row := Row{}
		for i, col := range newColumns {
			row[col] = rec[i]
		}
		records = append(records, row)
	}

	// Step 3: Write JSON array to data.json
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		panic(err)
	}
	os.WriteFile("data.json", jsonData, 0644)

	if *doPush {
		// Only push if data.json has changes
		gitDiff := exec.Command("git", "diff", "--quiet", "--", "data.json")
		err := gitDiff.Run()
		if err == nil {
			fmt.Println("No changes to data.json, skipping git commit/push.")
		} else {
			gitAdd := exec.Command("git", "add", "data.json")
			gitAdd.Stdout = os.Stdout
			gitAdd.Stderr = os.Stderr
			_ = gitAdd.Run()

			gitCommit := exec.Command("git", "commit", "-m", "update data")
			gitCommit.Stdout = os.Stdout
			gitCommit.Stderr = os.Stderr
			_ = gitCommit.Run()

			gitPush := exec.Command("git", "push")
			gitPush.Stdout = os.Stdout
			gitPush.Stderr = os.Stderr
			_ = gitPush.Run()
		}
	}
	fmt.Println("Done.")
}
