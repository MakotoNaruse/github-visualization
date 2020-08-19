package main

import (
	"reflect"
	"testing"
	"time"
)

/*
func TestLanguageRate(t *testing.T) {
	type args struct {
		languages Languages
	}
	tests := []struct {
		name string
		args args
		want map[string]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LanguageRate(tt.args.languages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LanguageRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonthlyContributes(t *testing.T) {
	type args struct {
		contributes []Contributes
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MonthlyContributes(tt.args.contributes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MonthlyContributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
 */

func TestVisualize(t *testing.T) {
	count := 10
	type args struct {
		data GithubData
	}
	tests := []struct {
		name string
		args args
		want Visualizer
	}{
		{
			name: "empty",
			args: args{
				data: GithubData{
					GithubUser: &GithubUser{
						Login: "Test",
						ImgURL: "TestURL",
						ContributionsCollection: &ContributionsCollection{
							CommitContributions: []CommitContributions{},
							ReviewContributions: []ReviewContributions{},
						},
					},
				},
			},
			want: Visualizer{
				UserName:          "Test",
				ImgURL:            "TestURL",
				TotalScore:        Score{
					IndividualScore: 0,
					TeamScore:       0,
					SocietyScore:    0,
				},
				ScoreByLanguage: map[string]LanguageScore{},
				ScoreByRepository: map[string]RepositoryScore{},
			},
		},
		{
			name: "one commit contribution",
			args: args{
				data: GithubData{
					GithubUser: &GithubUser{
						Login: "Test",
						ImgURL: "TestURL",
						ContributionsCollection: &ContributionsCollection{
							CommitContributions: []CommitContributions{
								{
									Repository:    Repository{
										Owner:     Owner{
											Login: "Test",
										},
										Name:      "Test",
										IsPrivate: true,
										Languages: Languages{
											TotalCount:    1,
											TotalSize:     100,
											LanguageName:  []LanguageName{
												{
													Color: nil,
													Name: "Java",
												},
											},
											LanguageCount: []LanguageCount{
												{
													Count: 100,
												},
											},
										},
									},
									Contributions: Contributions{
										TotalCount:  10,
										Contributes: []Contributes{
											{
												CommitCount: &count,
												Time:        time.Now(),
											},
										},
									},
								},
							},
							ReviewContributions: []ReviewContributions{},
						},
					},
				},
			},
			want: Visualizer{
				UserName:          "Test",
				ImgURL:            "TestURL",
				TotalScore:        Score{
					IndividualScore: 10,
					TeamScore:       0,
					SocietyScore:    0,
				},
				ScoreByLanguage: map[string]LanguageScore{
					"Java": LanguageScore{
						TotalScore:   Score{
							IndividualScore: 10,
							TeamScore:       0,
							SocietyScore:    0,
						},
						MonthlyScore: map[string]Score{
							time.Now().Format("2006/1"): {
								IndividualScore: 10,
								TeamScore:       0,
								SocietyScore:    0,
							},
						},
					},
				},
				ScoreByRepository: map[string]RepositoryScore{
					"Test/Test": {
						IsOSS:        false,
						LanguageRate: map[string]float64{
							"Java": 1.00,
						},
						Score:        Score{
							IndividualScore: 10,
							TeamScore:       0,
							SocietyScore:    0,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Visualize(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Visualize() = %v, want %v", got, tt.want)
			}
		})
	}
}
