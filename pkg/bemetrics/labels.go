package bemetrics

// Содержит набор дополнительных статичных меток которые добавляются к метрикам
type Labels []Label
type Label struct {
	Key   string
	Value string
}

// Создает пустой экземпляр Labels
func NewLabels(customLabels Labels) *Labels {
	return &Labels{}
}

// Возвращает все ключи меток
func (l Labels) Keys() []string {
	keys := make([]string, len(l))
	for i, label := range l {
		keys[i] = label.Key
	}
	return keys
}

// Возвращает все значения меток
func (l Labels) Values() []string {
	values := make([]string, len(l))
	for i, label := range l {
		values[i] = label.Value
	}
	return values
}
