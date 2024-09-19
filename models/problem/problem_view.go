package problem_model

import "github.com/oj-lab/platform/models"

type ProblemTagView struct {
	models.MetaFields
	Name string `json:"name"`
}

func (pt ProblemTag) ToProblemTagView() ProblemTagView {
	return ProblemTagView{
		MetaFields: pt.MetaFields,
		Name:       pt.Name,
	}
}

func GetProblemTagViewList(tags []*ProblemTag) []ProblemTagView {
	tagView := []ProblemTagView{}
	for _, tag := range tags {
		tagView = append(tagView, tag.ToProblemTagView())
	}
	return tagView
}

var ProblemInfoSelection = append([]string{"slug", "title", "difficulty"}, models.MetaFieldsSelection...)

type ProblemInfoView struct {
	Slug       string            `json:"slug"`
	Title      string            `json:"title"`
	Difficulty ProblemDifficulty `json:"difficulty"`
	Tags       []ProblemTagView  `json:"tags"`
	Solved     *bool             `json:"solved,omitempty"`
}

func (p Problem) ToProblemInfo() ProblemInfoView {
	return ProblemInfoView{
		Slug:       p.Slug,
		Title:      p.Title,
		Difficulty: p.Difficulty,
		Tags:       GetProblemTagViewList(p.Tags),
		Solved:     p.Solved,
	}
}

func GetProblemInfoViewList(problems []Problem) []ProblemInfoView {
	problemInfoList := make([]ProblemInfoView, len(problems))
	for i, problem := range problems {
		problemInfoList[i] = problem.ToProblemInfo()
	}
	return problemInfoList
}
