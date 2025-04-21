package metric

// GetKey - возвращает ключ по типу и наименованию метрики.
func GetKey(mType string, mName string) string {
	return string(mType) + ":" + mName
}
