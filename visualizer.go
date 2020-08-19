package main

import (
	"fmt"
)

type Visualizer struct {
	UserName          string
	ImgURL            string
	TotalScore        Score
	ScoreByLanguage   map[string]LanguageScore
	ScoreByRepository map[string]RepositoryScore
}

type LanguageScore struct {
	TotalScore   Score
	MonthlyScore map[string]Score
}

type RepositoryScore struct {
	IsOSS        bool
	LanguageRate map[string]float64
	Score        Score
}

type Score struct {
	IndividualScore float64
	TeamScore       float64
	SocietyScore    float64
}

// 新たに対応する言語があるときはここで記述
// 色はAPIから取得するのでここでは空
var languageString = map[string]string{
	"C":          "",
	"Cpp":        "",
	"Go":         "",
	"Java":       "",
	"JavaScript": "",
	"PHP":        "",
	"Python":     "",
	"Ruby":       "",
	"Swift":      "",
	"TypeScript": "",
}

// レビューコントリビュートあたりの重み
const reviewWeight = 10.0

func Visualize(data GithubData) Visualizer {
	user := data.GithubUser
	scoreByLanguage := map[string]LanguageScore{}
	scoreByRepository := map[string]RepositoryScore{}
	commitContributions := user.ContributionsCollection.CommitContributions
	for i := 0; i < len(commitContributions); i++ {
		contributions := commitContributions[i]
		contributes := contributions.Contributions.Contributes
		monthly := MonthlyContributes(contributes)
		repoScore := aggregateTotalRepoScoreByCommit(contributions, user, scoreByRepository)
		//次に言語毎のスコアを出す
		//まずは全体
		aggregateCommitRepoScore(repoScore, scoreByLanguage)
		//次に月毎
		aggregateCommitMonthlyScore(monthly, repoScore, scoreByLanguage)
	}
	reviewContributions := user.ContributionsCollection.ReviewContributions
	for i := 0; i < len(reviewContributions); i++ {
		contributions := reviewContributions[i]
		contributes := contributions.Contributions.Contributes
		monthly := MonthlyContributes(contributes)
		repoScore := aggregateTotalRepoScoreByReview(contributions, scoreByRepository, user)
		//次に言語毎のスコアを出す
		//まずは全体
		aggregateReviewRepoScore(repoScore, scoreByLanguage)
		//次に月毎
		aggregateReviewMonthlyScore(monthly, repoScore, scoreByLanguage)
	}
	totalScore := aggregateTotalScore(scoreByRepository)
	return Visualizer{
		UserName:          user.Login,
		ImgURL:            user.ImgURL,
		TotalScore:        totalScore,
		ScoreByLanguage:   scoreByLanguage,
		ScoreByRepository: scoreByRepository,
	}
}

func aggregateTotalScore(scoreByRepository map[string]RepositoryScore) Score {
	totalIndRepo := 0.0
	totalTeamRepo := 0.0
	totalSocRepo := 0.0
	for _, val := range scoreByRepository {
		totalIndRepo += val.Score.IndividualScore
		totalTeamRepo += val.Score.TeamScore
		totalSocRepo += val.Score.SocietyScore
	}
	return Score{
		IndividualScore: totalIndRepo,
		TeamScore:       totalTeamRepo,
		SocietyScore:    totalSocRepo,
	}
}

func aggregateReviewMonthlyScore(monthly map[string]int, repoScore RepositoryScore, scoreByLanguage map[string]LanguageScore) {
	for month, cnt := range monthly {
		//月のスコア
		monInd := 0.0
		monTeam := float64(cnt)
		monSoc := 0.0
		for key, _ := range repoScore.LanguageRate {
			// 比率をかける
			langMonInd := monInd * repoScore.LanguageRate[key]
			langMonTeam := monTeam * repoScore.LanguageRate[key]
			langMonSoc := monSoc * repoScore.LanguageRate[key]
			// 既に言語のスコアがあるならプラスし、ないなら新たに作る
			if monScore, ok := scoreByLanguage[key].MonthlyScore[month]; ok {
				monScore.IndividualScore += langMonInd
				monScore.TeamScore += langMonTeam
				monScore.SocietyScore += langMonSoc
				scoreByLanguage[key].MonthlyScore[month] = monScore
			} else {
				scoreByLanguage[key].MonthlyScore[month] = Score{
					IndividualScore: langMonInd,
					TeamScore:       langMonTeam,
					SocietyScore:    langMonSoc,
				}
			}
		}
	}
}

func aggregateReviewRepoScore(repoScore RepositoryScore, scoreByLanguage map[string]LanguageScore) {
	for key, _ := range repoScore.LanguageRate {
		// 比率をかける
		langInd := 0.0
		langTeam := repoScore.Score.TeamScore * repoScore.LanguageRate[key]
		langSoc := 0.0
		// 既に言語のスコアがあるならプラスし、ないなら新たに作る
		if lnScore, ok := scoreByLanguage[key]; ok {
			lnScore.TotalScore.IndividualScore += langInd
			lnScore.TotalScore.TeamScore += langTeam
			lnScore.TotalScore.SocietyScore += langSoc
			scoreByLanguage[key] = lnScore
		} else {
			totalScore := Score{
				IndividualScore: langInd,
				TeamScore:       langTeam,
				SocietyScore:    langSoc,
			}
			scoreByLanguage[key] = LanguageScore{
				TotalScore:   totalScore,
				MonthlyScore: map[string]Score{},
			}
		}
	}
}

func aggregateTotalRepoScoreByReview(contributions ReviewContributions, scoreByRepository map[string]RepositoryScore, user *GithubUser) RepositoryScore {
	// まずはリポジトリスコア
	repo := contributions.Repository
	fullName := fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name)

	ind := 0.0
	team := float64(contributions.Contributions.TotalCount) * reviewWeight
	soc := 0.0

	//まずはリポジトリがあるならスコアをプラス、ないなら作る
	repoScore, ok := scoreByRepository[fullName]
	if ok {
		repoScore.Score.TeamScore += team
		scoreByRepository[fullName] = repoScore
	} else {
		isOSS := user.Login != repo.Owner.Login && !repo.IsPrivate
		languageRate := LanguageRate(repo.Languages)
		score := Score{
			IndividualScore: ind,
			TeamScore:       team,
			SocietyScore:    soc,
		}
		repoScore := RepositoryScore{
			IsOSS:        isOSS,
			LanguageRate: languageRate,
			Score:        score,
		}
		scoreByRepository[fullName] = repoScore
	}
	return repoScore
}

func aggregateCommitMonthlyScore(monthly map[string]int, repoScore RepositoryScore, scoreByLanguage map[string]LanguageScore) {
	for month, cnt := range monthly {
		//月のスコア
		monInd := 0.0
		monTeam := 0.0
		monSoc := 0.0
		if repoScore.IsOSS {
			monSoc = float64(cnt)
		} else {
			monInd = float64(cnt)
		}
		for key, _ := range repoScore.LanguageRate {
			// 比率をかける
			langMonInd := monInd * repoScore.LanguageRate[key]
			langMonTeam := monTeam * repoScore.LanguageRate[key]
			langMonSoc := monSoc * repoScore.LanguageRate[key]
			// 既に言語のスコアがあるならプラスし、ないなら新たに作る
			if monScore, ok := scoreByLanguage[key].MonthlyScore[month]; ok {
				monScore.IndividualScore += langMonInd
				monScore.TeamScore += langMonTeam
				monScore.SocietyScore += langMonSoc
				scoreByLanguage[key].MonthlyScore[month] = monScore
			} else {
				scoreByLanguage[key].MonthlyScore[month] = Score{
					IndividualScore: langMonInd,
					TeamScore:       langMonTeam,
					SocietyScore:    langMonSoc,
				}
			}
		}
	}
}

func aggregateCommitRepoScore(repoScore RepositoryScore, scoreByLanguage map[string]LanguageScore) {

	for key, _ := range repoScore.LanguageRate {
		// 比率をかける
		langInd := repoScore.Score.IndividualScore * repoScore.LanguageRate[key]
		langSoc := repoScore.Score.SocietyScore * repoScore.LanguageRate[key]
		// 既に言語のスコアがあるならプラスし、ないなら新たに作る
		if lnScore, ok := scoreByLanguage[key]; ok {
			lnScore.TotalScore.IndividualScore += langInd
			lnScore.TotalScore.SocietyScore += langSoc
			scoreByLanguage[key] = lnScore
		} else {
			totalScore := Score{
				IndividualScore: langInd,
				TeamScore:       0,
				SocietyScore:    langSoc,
			}
			scoreByLanguage[key] = LanguageScore{
				TotalScore:   totalScore,
				MonthlyScore: map[string]Score{},
			}
		}
	}
}

func aggregateTotalRepoScoreByCommit(contributions CommitContributions, user *GithubUser, scoreByRepository map[string]RepositoryScore) RepositoryScore {
	repo := contributions.Repository
	fullName := fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name)
	isOSS := user.Login != repo.Owner.Login && !repo.IsPrivate
	languageRate := LanguageRate(repo.Languages)
	// このリポジトリ単体のスコア
	ind := 0.0
	team := 0.0
	soc := 0.0
	if isOSS {
		soc = float64(contributions.Contributions.TotalCount)
	} else {
		ind = float64(contributions.Contributions.TotalCount)
	}

	//まずはリポジトリのスコアを格納
	score := Score{
		IndividualScore: ind,
		TeamScore:       team,
		SocietyScore:    soc,
	}
	repoScore := RepositoryScore{
		IsOSS:        isOSS,
		LanguageRate: languageRate,
		Score:        score,
	}
	scoreByRepository[fullName] = repoScore
	return repoScore
}

func LanguageRate(languages Languages) map[string]float64 {
	total := 0.0
	languageRate := map[string]float64{}
	n := languages.TotalCount
	for i := 0; i < n; i++ {
		// 対応言語に含まれるなら計算する
		name := languages.LanguageName[i].Name
		if _, ok := languageString[name]; ok {
			color := languages.LanguageName[i].Color
			if color != nil {
				languageString[languages.LanguageName[i].Name] = *color
			}
			count := float64(languages.LanguageCount[i].Count)
			languageRate[name] = count
			total += count
		}
	}
	//最後に合計で割って割合化する
	for key, _ := range languageRate {
		languageRate[key] /= total
	}
	return languageRate
}

func MonthlyContributes(contributes []Contributes) map[string]int {
	monthly := map[string]int{}
	for i := 0; i < len(contributes); i++ {
		month := fmt.Sprintf("%d/%d", contributes[i].Time.Year(), contributes[i].Time.Month())
		// CommitCountがnilだとレビューのcontribute
		cnt := 0
		if contributes[i].CommitCount != nil {
			cnt = *contributes[i].CommitCount
		} else {
			cnt = reviewWeight
		}
		if _, ok := monthly[month]; ok {
			monthly[month] += cnt
		} else {
			monthly[month] = cnt
		}
	}
	return monthly
}
