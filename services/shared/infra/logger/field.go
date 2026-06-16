package logger

type Field struct {
	key   string
	value any
}

func NewField(key string, value any) Field {
	return Field{key: key, value: value}
}
