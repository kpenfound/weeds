package worker

func Classify(a, b, c, d int) string {
	if a > 0 {
		if b > 0 {
			if c > 0 {
				if d > 0 {
					return "all"
				}
			}
		}
	}
	return "none"
}
