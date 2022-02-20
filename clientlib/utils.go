package clientlib

func CheckLimits(data []byte, ttl int, hits int) string {
	if len(data) < 1 {
		return "Cannot share nothing"
	}
	if len(data) > (5 * 1024 * 1024) {
		return "Cannot send more than 5MB of data"
	}
	if ttl < 1 || ttl > (2*24*60*60) {
		return "The link must expire between 1 second and 2 days"
	}
	if hits < 1 || hits > 1000000 {
		return "The link can have between 1 and 1.000.000 views"
	}
	return ""
}
