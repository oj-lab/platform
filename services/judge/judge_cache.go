package judge_service

import (
	"context"

	"github.com/google/uuid"
	judge_model "github.com/oj-lab/platform/models/judge"
	problem_model "github.com/oj-lab/platform/models/problem"
	gorm_agent "github.com/oj-lab/platform/modules/agent/gorm"
)

func UpsertJudgeCache(ctx context.Context, uid uuid.UUID, verdict judge_model.JudgeVerdict) error {
	db := gorm_agent.GetDefaultDB()
	judge, err := judge_model.GetJudge(db, uid)
	if err != nil {
		return err
	}
	// log_module.AppLogger().WithField("judge", judge).Debug("getjudge")
	var problem *problem_model.Problem
	problem, err = problem_model.GetProblem(db, judge.ProblemSlug)
	if err != nil {
		return err
	}

	var scoreCache *judge_model.JudgeScoreCache
	var rankCache *judge_model.JudgeRankCache
	for {
		rankCache, err = judge_model.GetJudgeRankCache(db, judge.UserAccount)
		if err != nil { // previous empty
			_, err = judge_model.CreateJudgeRankCache(db, judge_model.NewJudgeRankCache(judge.UserAccount))
			if err != nil { // create fail, exists -> get data again and continue the update logic.
				continue
			}
		} else {
			break
		}
	}

	for {
		scoreCache, err = judge_model.GetJudgeScoreCache(db, judge.UserAccount, judge.ProblemSlug)
		if err != nil { // previous empty
			_, err := judge_model.CreateJudgeScoreCache(db, judge_model.NewJudgeScoreCache(judge.UserAccount, judge.ProblemSlug))
			if err != nil { // create fail, exists -> get data again and continue the update logic.
				continue
			}

		} else {
			break
		}
	}

	// log_module.AppLogger().WithField("scoreCache", scoreCache).Debug("get scoreCache")

	// previous no ac || current more early
	// need to update
	if !scoreCache.IsAccepted || judge.CreateAt.Before(*scoreCache.SolveTime) {
		extraPoint := 0
		if !scoreCache.IsAccepted {
			extraPoint = 1
		}
		preSubmissionCount := scoreCache.SubmissionCount
		if verdict == judge_model.JudgeVerdictAccepted {
			scoreCache.SubmissionCount, err = judge_model.GetBeforeSubmission(db, *judge) // rescan to count previous finished judge
			if err != nil {
				return err
			}
			scoreCache.IsAccepted = true
			rankCache.Points = rankCache.Points + extraPoint
			problem.AcceptCount += extraPoint
			scoreCache.SolveTime = judge.CreateAt
		} else {
			scoreCache.SubmissionCount += 1
		}
		rankCache.TotalSubmissions += scoreCache.SubmissionCount - preSubmissionCount
		problem.SubmitCount += scoreCache.SubmissionCount - preSubmissionCount
		// log_module.AppLogger().WithField("scoreCache", scoreCache).Debug("update scoreCache")

		err = problem_model.UpdateProblem(db, *problem)
		if err != nil {
			return err
		}
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
