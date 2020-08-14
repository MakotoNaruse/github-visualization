package main

import (
	"time"
)

type GithubWrap struct {
	GithubData *GithubData `json:"data"`
}

type GithubData struct {
	GithubUser *GithubUser `json:"viewer"`
}

type GithubUser struct {
	Login string `json:"login"`
	ContributionsCollection *ContributionsCollection `json:"contributionsCollection"`
}

type ContributionsCollection struct {
	CommitContributions []CommitContributions `json:"commitContributionsByRepository"`
	ReviewContributions []ReviewContributions `json:"pullRequestReviewContributionsByRepository"`
}

type CommitContributions struct {
	Repository Repository `json:"repository"`
	Contributions Contributions `json:"contributions"`
}

type ReviewContributions struct {
	Repository Repository `json:"repository"`
	Contributions Contributions `json:"contributions"`
}

type Repository struct {
	Owner Owner `json:"owner"`
	Name string `json:"name"`
	IsPrivate bool `json:"isPrivate"`
	Languages Languages `json:"languages"`
}

type Owner struct {
	Login string `json:"login"`
}

type Languages struct {
	TotalCount int `json:"totalCount"`
	TotalSize int `json:"totalSize"`
	LanguageName []LanguageName `json:"nodes"`
	LanguageCount []LanguageCount `json:"edges"`
}

type LanguageName struct {
	Color *string `json:"color"`
	Name string `json:"name"`
}

type LanguageCount struct {
	Count int `json:"size"`
}

type Contributions struct {
	TotalCount int `json:"totalCount"`
	Contributes []Contributes `json:"nodes"`
}

type Contributes struct {
	CommitCount int `json:"commitCount"`
	Time time.Time `json:"occurredAt"`
}