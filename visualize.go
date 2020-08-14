package main

type Visualize struct {
	UserName string
	ImgURL string
	TotalScore string
	ScoreByLanguage map[string]LanguageScore
	ScoreByRepository map[string]RepositoryScore
}

type LanguageScore struct {
	TotalScore Score
	MonthlyScore map[int]Score
}

type RepositoryScore struct {
	IsOSS bool
	Score Score
}

type Score struct {
	IndividualScore string
	TeamScore string
	SocietyScore string
}
