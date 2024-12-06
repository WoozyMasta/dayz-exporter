package bemetrics

// contains a set of additional static labels that are added to metrics
type Labels []Label
type Label struct {
	Key   string
	Value string
}

// creates an empty Labels instance
func NewLabels(customLabels Labels) *Labels {
	return &Labels{}
}

// returns all labels keys
func (l Labels) Keys() []string {
	keys := make([]string, len(l))
	for i, label := range l {
		keys[i] = label.Key
	}
	return keys
}

// returns all labels values
func (l Labels) Values() []string {
	values := make([]string, len(l))
	for i, label := range l {
		values[i] = label.Value
	}
	return values
}
