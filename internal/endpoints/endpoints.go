package endpoints

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Repo struct {
	ID              int    `json:"id"`
	NodeId          string `json:"node_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	FullName        string `json:"full_name"`
	HTMLUrl         string `json:"html_url"`
	Url             string `json:"url"`
	StarGazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	Visibility      string `json:"visibility"`
}

type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	HTMLUrl     string `json:"html_url"`
	DownloadUrl string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
}

func GetRepos() []Repo {
	url := "https://api.github.com/user/repos?type=owner"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	r.Header = map[string][]string{
		"Accept":        {"application/vnd.github+json"},
		"Authorization": {fmt.Sprintf("Bearer %s", os.Getenv("PAT"))},
	}

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	defer res.Body.Close()

	var repo []Repo

	json.NewDecoder(res.Body).Decode(&repo)

	return repo
}

func GetReadme(url string) string {
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/readme", url), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	r.Header = map[string][]string{
		"Accept":        {"application/vnd.github+json"},
		"Authorization": {fmt.Sprintf("Bearer %s", os.Getenv("PAT"))},
	}

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	defer res.Body.Close()

	var readme struct {
		Content string `json:"content"`
	}

	json.NewDecoder(res.Body).Decode(&readme)

	readmeString := strings.Join(strings.Split(readme.Content, "\n"), "")
	dec, err := base64.StdEncoding.DecodeString(readmeString)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	return string(dec)
}

func ImageFromMD(md string) string {
	exp := regexp.MustCompile(`\[[\w]+\]\(([.\w\/\\]+)\)`)

	if matches := exp.FindStringSubmatch(md); len(matches) > 0 {
		return strings.Replace(matches[1], "./", "", 1)
	}

	return ""
}

func GetFile(url, path string) File {
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/contents/%s", url, path), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	r.Header = map[string][]string{
		"Accept":        {"application/vnd.github+json"},
		"Authorization": {fmt.Sprintf("Bearer %s", os.Getenv("PAT"))},
	}

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	defer res.Body.Close()

	var file File

	json.NewDecoder(res.Body).Decode(&file)

	return file
}

func GetLanguages(url string) []string {
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/languages", url), nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	r.Header = map[string][]string{
		"Accept":        {"application/vnd.github+json"},
		"Authorization": {fmt.Sprintf("Bearer %s", os.Getenv("PAT"))},
	}

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	defer res.Body.Close()

	var langs map[string]int

	json.NewDecoder(res.Body).Decode(&langs)

	var languages []string

	for k := range langs {
		languages = append(languages, k)
	}

	return languages
}
