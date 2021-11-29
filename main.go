package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type RepoDetail struct {
	Name string `json:"name"`
	Diff int    `json:"max_minutes"`
}

type RepoList struct {
	ListRepo []RepoDetail `json:"repos"`
}

type RunningWorkflow struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartedAt string `json:"run_started_at"`
	CancelUrl string `json:"cancel_url"`
	ReRunUrl  string `json:"rerun_url"`
}

type RunningActionResponse struct {
	TotalCount   int               `json:"total_count"`
	WorkflowRuns []RunningWorkflow `json:"workflow_runs"`
}

func apiCall(method string, url string, body io.Reader) []byte {
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func restartAction(flow *RunningWorkflow) {

	// Cancel workflow
	apiCall("POST", flow.CancelUrl, nil)

	// Wait for work flow to be cancelled grassfully.
	for {
		body := apiCall("GET", flow.Url, nil)
		data := RunningWorkflow{}
		if err := json.Unmarshal(body, &data); err != nil {
			panic(err)
		}
		if data.Status != "in_progress" {
			fmt.Println("Workflow " + data.Status + " so moving towards re-run the flow")
			break
		}
		fmt.Println("Workflow still in progress. -- we are waiting.")
		time.Sleep(2 * time.Second)
	}

	// ReRun workflow
	apiCall("POST", flow.ReRunUrl, nil)
	fmt.Println("Completed.")

}

func ghActionCheck(repo *RepoDetail) {

	body := apiCall("GET", "https://api.github.com/repos/"+repo.Name+"/actions/runs?status=in_progress", nil)

	fres := RunningActionResponse{}
	if err := json.Unmarshal(body, &fres); err != nil {
		panic(err)
	}

	fmt.Println("Found " + string(fres.TotalCount) + " Currently running")

	for i, item := range fres.WorkflowRuns {
		t, e := time.Parse(time.RFC3339, item.StartedAt)
		if e != nil {
			panic(e)
		}
		fmt.Println(i, item.Name)
		// calculate diff
		diff := time.Now().Sub(t).Minutes()

		fmt.Println(t, diff)
		if diff > float64(repo.Diff) {
			fmt.Println("Restarting actions.")
			restartAction(&item)
		}

	}

}

func main() {
	godotenv.Load()
	file, _ := ioutil.ReadFile("conf.json")
	data := RepoList{}
	json.Unmarshal([]byte(file), &data)

	for {
		for _, repo := range data.ListRepo {
			fmt.Println(repo)
			ghActionCheck(&repo)
		}
		time.Sleep(1 * time.Minute)
	}
}
