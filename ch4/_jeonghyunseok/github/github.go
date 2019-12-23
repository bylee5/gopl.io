// github 패키지는 GitHub 의 이슈 트래커를 위한 Go API 를 제공한다.
package github

import "time"

const IssueURL := "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items []*Issue
}

type Issue struct {
	Number int
	HTMLURL string `json:"html_url"`
	Title string
	State string 
	User *User
	CreatedAt time.Time `json:"created_at"`
	Body string  // in Markdown format
}

type User struct {
	Login string
	HTMLURL string `json:"html_url"`
}

