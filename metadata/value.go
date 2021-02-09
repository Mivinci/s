package metadata

func Int(v interface{}) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case float32:
		return int(t), true
	case int:
		return t, true
	default:
		return 0, false
	}
}
