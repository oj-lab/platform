package problem

import (
	"gorm.io/gorm"
)

func CreateProblemInfo(tx *gorm.DB, problemInfo ProblemInfo) error {
	return tx.Create(&problemInfo).Error
}

func buildGetProblemInfoTXByOptions(tx *gorm.DB, options GetProblemOptions, isCount bool) *gorm.DB {
	tagsList := []string{}
	for _, tag := range options.Tags {
		tagsList = append(tagsList, tag.Name)
	}
	tx = tx.Model(&ProblemInfo{})
	tx = tx.Distinct()

	if len(tagsList) > 0 {
		tx.Preload("Tags", "name in ?", tagsList)
	} else {
		tx = tx.Preload("AlgorithmTags")
	}
	if options.Slug != nil {
		tx = tx.Where("problem_slug = ?", *options.Slug)
	}
	if options.Title != nil {
		tx = tx.Where("title = ?", *options.Title)
	}
	if !isCount {
		if options.Offset != nil {
			tx = tx.Offset(*options.Offset)
		}
		if options.Limit != nil {
			tx = tx.Limit(*options.Limit)
		}
	}

	return tx
}

func CountProblemInfoByOptions(tx *gorm.DB, options GetProblemOptions) (int64, error) {
	var count int64

	tx = buildGetProblemInfoTXByOptions(tx, options, true)
	err := tx.Count(&count).Error

	return count, err
}

func GetProblemInfoListByOptions(
	tx *gorm.DB, options GetProblemOptions,
) ([]ProblemInfo, int64, error) {
	total, err := CountProblemInfoByOptions(tx, options)
	if err != nil {
		return nil, 0, err
	}

	problemList := []ProblemInfo{}

	tx = buildGetProblemInfoTXByOptions(tx, options, false)
	err = tx.Find(&problemList).Error
	if err != nil {
		return nil, 0, err
	}

	return problemList, total, nil
}

func DeleteProblemInfo(tx *gorm.DB, slug string) error {
	err := tx.Model(&ProblemInfo{
		ProblemSlug: slug,
	}).Association("Tags").Clear()
	if err != nil {
		return err
	}

	return tx.Delete(&ProblemInfo{
		ProblemSlug: slug,
	}).Error
}
