package error

type DecisionNotFoundErr struct{}

func NewDecisionNotFoundErr() error {
	return DecisionNotFoundErr{}
}

func (e DecisionNotFoundErr) Error() string {
	return "error, decision not found"
}
