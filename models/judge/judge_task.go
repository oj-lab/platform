package judge_model

type JudgeTask struct {
	JudgeUID      string  `json:"judgeUID"`
	ProblemSlug   string  `json:"problemSlug"`
	Code          string  `json:"code"`
	Language      string  `json:"language"`
	RedisStreamID *string `json:"redisStreamID"`
}

func (jt *JudgeTask) ToStringMap() map[string]interface{} {
	return map[string]interface{}{
		"judge_uid":    jt.JudgeUID,
		"problem_slug": jt.ProblemSlug,
		"code":         jt.Code,
		"language":     jt.Language,
	}
}

func JudgeTaskFromMap(m map[string]interface{}) *JudgeTask {
	return &JudgeTask{
		JudgeUID:    m["judge_uid"].(string),
		ProblemSlug: m["problem_slug"].(string),
		Code:        m["code"].(string),
		Language:    m["language"].(string),
	}
}
