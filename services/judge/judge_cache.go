package judge_service

import (
	"context"

	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func UpsertJudgeCache(ctx context.Context, uid uuid.UUID, verdict judge_model.JudgeVerdict) error {
	db := gorm_agent.GetDefaultDB()
	judge, err := judge_model.GetJudge(db, uid)
	if err != nil {
		return err
	}
	// log_module.AppLogger().WithField("judge", judge).Debug("getjudge")
	var scoreCache *judge_model.JudgeScoreCache
	var rankCache *judge_model.JudgeRankCache
	rankCache, err = judge_model.GetJudgeRankCache(db, judge.UserAccount)
	if err != nil {
		return err
	}

	extraPoint := 0
	for {
		scoreCache, err = judge_model.GetJudgeScoreCache(db, judge.UserAccount, judge.ProblemSlug)
		if err != nil {
			// previous empty
			// log_module.AppLogger().Debug("previous empty")
			scoreCache := judge_model.NewJudgeScoreCache(judge.UserAccount, judge.ProblemSlug)
			if verdict == judge_model.JudgeVerdictAccepted {
				extraPoint = 1
				scoreCache.IsAccepted = true
				scoreCache.SolveTime = judge.CreateAt
			}
			_, err := judge_model.CreateJudgeScoreCache(db, scoreCache)
			if err != nil { // create fail, get data again and continue the update logic.
				continue
			}

			// create success
			rankCache.Points += int64(extraPoint)
			rankCache.TotalSubmissions += 1
			err = judge_model.UpdateJudgeRankCache(db, *rankCache)
			if err != nil {
				return err
			}
			return nil
		} else {
			break
		}
	}

	// log_module.AppLogger().WithField("scoreCache", scoreCache).Debug("get scoreCache")

	// previous no ac || current more early
	// need to update
	if !scoreCache.IsAccepted || judge.CreateAt.Before(*scoreCache.SolveTime) {
		if !scoreCache.IsAccepted {
			extraPoint = 1
		}
		preSubmissionCount := scoreCache.SubmissionCount
		if verdict == judge_model.JudgeVerdictAccepted {
			scoreCache.SubmissionCount, err = judge_model.GetBeforeSubmission(db, *judge) // rescan to count previous finished
			if err != nil {
				return err
			}
			scoreCache.IsAccepted = true
			rankCache.Points = rankCache.Points + int64(extraPoint)
			scoreCache.SolveTime = judge.CreateAt
		} else {
			scoreCache.SubmissionCount += 1
		}
		rankCache.TotalSubmissions += scoreCache.SubmissionCount - preSubmissionCount
		// log_module.AppLogger().WithField("scoreCache", scoreCache).Debug("update scoreCache")

		err = judge_model.UpdateJudgeScoreCache(db, *scoreCache)
		if err != nil {
			return err
		}

		err = judge_model.UpdateJudgeRankCache(db, *rankCache)
		if err != nil {
			return err
		}
	}
	// if no early, no need update, just a query
	return nil
}
