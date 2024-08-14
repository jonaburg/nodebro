package requests

import (
    "encoding/json"
    "fmt"
    "time"
    "sync"
    "net/http"
)

type Release struct {
	Owner        string `json:"owner"`
	Repo         string `json:"repo"`
	TagName      string `json:"tag_name"`
	URL          string `json:"html_url"`
    Body         string `json:"body"`
    Prerelease bool      `json:"prerelease"`
    CreatedAt  time.Time `json:"created_at"`
    TimeDifference  string `json:"hours_diff"`
}

type RepoInfo struct {
	Owner string
	Repo  string
}

func calculateTimeDiff(release *Release) {
    release.CreatedAt = release.CreatedAt.UTC() 
    currentTime := time.Now().UTC()            
    difference := currentTime.Sub(release.CreatedAt)
    if int(difference.Hours()) > 24 {
      days := int(difference.Hours()) / 24
      release.TimeDifference = fmt.Sprintf("%dd ago", days)
    } else {
      release.TimeDifference = fmt.Sprintf("%dh ago", int(difference.Hours()))
    }
}


func FetchLatestRelease(repo RepoInfo, ch chan Release, wg *sync.WaitGroup, authToken string) {
    defer wg.Done()
    apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repo.Owner, repo.Repo)
    req, err := http.NewRequest("GET", apiUrl, nil)
    if err != nil {
        fmt.Printf("Error creating HTTP request for %s/%s: %s\n", repo.Owner, repo.Repo, err)
        return
    }
    req.Header.Set("Authorization", "Bearer "+authToken)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("Error making HTTP request for %s/%s: %s\n", repo.Owner, repo.Repo, err)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Failed to fetch release notes for %s/%s: %s\n", repo.Owner, repo.Repo, resp.Status)
        return
    }
    var release Release
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        fmt.Printf("Error parsing JSON for %s/%s: %s\n", repo.Owner, repo.Repo, err)
        return
    }
    release.Owner = repo.Owner
    release.Repo = repo.Repo
    calculateTimeDiff(&release)
    ch <- release
}
