package form

func BoolToString(b bool) string {
	if b {
		return "on"
	}
	return "off"
}