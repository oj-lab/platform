package judge_service

import (
	"context"

	"github.com/google/uuid"
	judge_model "github.com/oj-lab/oj-lab-platform/models/judge"
	gorm_agent "github.com/oj-lab/oj-lab-platform/modules/agent/gorm"
)

func UpsertJudgeScoreCache(ctx context.Context, uid uuid.UUID, verdict judge_model.JudgeVerdict) (*judge_model.JudgeScoreCache, error) {
	db := gorm_agent.GetDefaultDB()
	judge, err := judge_model.GetJudge(db, uid)
	if err != nil {
		return nil, err
	}
	// log_module.AppLogger().WithField("judge", judge).Debug("getjudge")
	var scoreCache *judge_model.JudgeScoreCache
	for {
		scoreCache, err = judge_model.GetJudgeScoreCache(db, judge.UserAccount, judge.ProblemSlug)
		if err != nil {
			// previous empty
			// log_module.AppLogger().Debug("previous empty")

			scoreCache := judge_model.NewJudgeScoreCache(judge.UserAccount, judge.ProblemSlug)
			if verdict == judge_model.JudgeVerdictAccepted {
				scoreCache.IsAccepted = true
				scoreCache.SolveTime = judge.CreateAt
			}
			newScoreCache, err := judge_model.CreateJudgeScoreCache(db, scoreCache)
			if err != nil { // create fail, get data again and continue the update logic.
				continue
			}
			return newScoreCache, nil
		} else {
			break
		}
	}

	// log_module.AppLogger().WithField("scoreCache", scoreCache).Debug("get scoreCache")

	// previous no ac || current more early
	// need to update
	if !scoreCache.IsAccepted || judge.CreateAt.Before(*scoreCache.SolveTime) {
		if verdict == judge_model.JudgeVerdictAccepted {
			scoreCache.SubmissionCount, err = judge_model.GetBeforeSubmission(db, *judge) // rescan to count previous finished
			if err != nil {
				return nil, err
			}
			scoreCache.IsAccepted = true
			scoreCache.SolveTime = judge.CreateAt
		} else {
			scoreCache.SubmissionCount += 1
		}

		// log_module.AppLogger().WithField("scoreCache", scoreCache).Debug("update scoreCache")

		err = judge_model.UpdateJudgeScoreCache(db, *scoreCache)
		if err != nil {
			return nil, err
		}
	}
	// if no early, no need update, just a query
	return scoreCache, nil
}
