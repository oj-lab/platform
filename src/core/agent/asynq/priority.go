package asynqAgent

type AsynqTaskPriority string

const (
	AsynqTaskPriorityCritical AsynqTaskPriority = "critical"
	// Default priority ruled in asynq library
	AsynqTaskPriorityDefault AsynqTaskPriority = "default"
	AsynqTaskPriorityLow     AsynqTaskPriority = "low"
)

var priorityMap = map[string]int{
	string(AsynqTaskPriorityCritical): 6,
	string(AsynqTaskPriorityDefault):  3,
	string(AsynqTaskPriorityLow):      1,
}
