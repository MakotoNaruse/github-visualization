package main

import (
	"fmt"
)

type Visualizer struct {
	UserName string
	ImgURL string
	TotalScore Score
	ScoreByLanguage map[string]LanguageScore
	ScoreByRepository map[string]RepositoryScore
}

type LanguageScore struct {
	TotalScore Score
	MonthlyScore map[string]Score
}

type RepositoryScore struct {
	IsOSS bool
	LanguageRate map[string] float64
	Score Score
}

type Score struct {
	IndividualScore float64
	TeamScore float64
	SocietyScore float64
}

// 新たに対応する言語があるときはここで記述
// 色はAPIから取得するのでここでは空
var languageString = map[string]string {
	"C": "",
	"Cpp": "",
	"C Sharp": "",
	"Go": "",
	"Java": "",
	"JavaScript": "",
	"PHP": "",
	"Python": "",
	"Ruby": "",
	"TypeScript": "",
}
// レビューコントリビュートあたりの重み
const reviewWeight = 10.0

func Visualize( data GithubData ) Visualizer{
	user := data.GithubUser
	totalScore := Score{
		IndividualScore: 0.0,
		TeamScore: 0.0,
		SocietyScore: 0.0,
	}
	scoreByLanguage := map[string]LanguageScore{}
	scoreByRepository := map[string]RepositoryScore{}
	visualizer := Visualizer{
		UserName: user.Login,
		ImgURL: user.ImgURL,
		TotalScore: totalScore,
		ScoreByLanguage: scoreByLanguage,
		ScoreByRepository: scoreByRepository,
	}
	commitContributions := user.ContributionsCollection.CommitContributions
	for i := 0; i < len(commitContributions); i++{
		repo := commitContributions[i].Repository
		fullName := fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name)
		isOSS := user.Login != repo.Owner.Login && !repo.IsPrivate
		languageRate := LanguageRate(repo.Languages)
		contributes := commitContributions[i].Contributions.Contributes
		monthly := MonthlyContributes(contributes)

		// このリポジトリ単体のスコア
		ind := 0.0
		team := 0.0
		soc := 0.0
		if isOSS {
			soc = float64(commitContributions[i].Contributions.TotalCount)
		} else {
			ind = float64(commitContributions[i].Contributions.TotalCount)
		}

		//まずはリポジトリのスコアを格納
		score := Score{
			IndividualScore: ind,
			TeamScore: team,
			SocietyScore: soc,
		}
		repoScore := RepositoryScore{
			IsOSS: isOSS,
			LanguageRate: languageRate,
			Score: score,
		}
		scoreByRepository[fullName] = repoScore

		//次に言語毎のスコアを出す
		//まずは全体
		for key, _ := range languageRate {
			// 比率をかける
			langInd := ind * languageRate[key]
			langTeam := team * languageRate[key]
			langSoc := soc * languageRate[key]
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
					TotalScore: totalScore,
					MonthlyScore: map[string]Score{},
				}
			}
		}
		//次に月毎
		for month, cnt := range monthly {
			//月のスコア
			monInd := 0.0
			monTeam := 0.0
			monSoc := 0.0
			if isOSS {
				monSoc = float64(cnt)
			} else {
				monInd = float64(cnt)
			}
			for key, _ := range languageRate {
				// 比率をかける
				langMonInd := monInd * languageRate[key]
				langMonTeam := monTeam * languageRate[key]
				langMonSoc := monSoc * languageRate[key]
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
	reviewContributions := user.ContributionsCollection.ReviewContributions
	for i := 0; i < len(reviewContributions); i++{
		repo := reviewContributions[i].Repository
		fullName := fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name)
		contributes := reviewContributions[i].Contributions.Contributes
		monthly := MonthlyContributes(contributes)

		// このリポジトリ単体のスコア
		ind := 0.0
		team := float64(reviewContributions[i].Contributions.TotalCount) * reviewWeight
		soc := 0.0

		//まずはリポジトリがあるならスコアをプラス
		if repoScore, ok := scoreByRepository[fullName]; ok{
			repoScore.Score.TeamScore += team
			scoreByRepository[fullName] = repoScore
		} else {
			isOSS := user.Login != repo.Owner.Login && !repo.IsPrivate
			languageRate := LanguageRate(repo.Languages)
			score := Score{
				IndividualScore: ind,
				TeamScore: team,
				SocietyScore: soc,
			}
			repoScore := RepositoryScore{
				IsOSS: isOSS,
				LanguageRate: languageRate,
				Score: score,
			}
			scoreByRepository[fullName] = repoScore
		}
		languageRate := scoreByRepository[fullName].LanguageRate

		//次に言語毎のスコアを出す
		//まずは全体
		for key, _ := range languageRate {
			// 比率をかける
			langInd := 0.0
			langTeam := team * languageRate[key]
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
					TotalScore: totalScore,
					MonthlyScore: map[string]Score{},
				}
			}
		}
		//次に月毎
		for month, cnt := range monthly {
			//月のスコア
			monInd := 0.0
			monTeam := float64(cnt)
			monSoc := 0.0
			for key, _ := range languageRate {
				// 比率をかける
				langMonInd := monInd * languageRate[key]
				langMonTeam := monTeam * languageRate[key]
				langMonSoc := monSoc * languageRate[key]
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
	return visualizer
}

func LanguageRate( languages Languages ) map[string]float64 {
	total := 0.0
	languageRate := map[string]float64{}
	n := languages.TotalCount
	for i := 0; i < n; i++ {
		// 対応言語に含まれるなら計算する
		name := languages.LanguageName[i].Name
		if _, ok := languageString[name]; ok{
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

func MonthlyContributes( contributes []Contributes ) map[string]int {
	monthly := map[string]int{}
	for i := 0; i < len(contributes); i++{
		month := fmt.Sprintf("%d/%d", contributes[i].Time.Year(), contributes[i].Time.Month())
		// CommitCountがnilだとレビューのcontribute
		cnt := 0
		if contributes[i].CommitCount != nil {
			cnt = *contributes[i].CommitCount
		} else {
			cnt = reviewWeight
		}
		if _, ok := monthly[month]; ok{
			monthly[month] += cnt
		} else {
			monthly[month] = cnt
		}
	}
	return monthly
}
