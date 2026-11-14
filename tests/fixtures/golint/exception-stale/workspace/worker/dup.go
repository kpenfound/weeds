package worker

func First(values []int) (int, int, int) {
	low := values[0]
	high := values[0]
	sum := 0
	for _, v := range values {
		if v < low {
			low = v
		}
		if v > high {
			high = v
		}
		sum += v
	}
	span := high - low
	if span < 0 {
		span = -span
	}
	mean := sum / len(values)
	if mean < low {
		mean = low
	}
	if mean > high {
		mean = high
	}
	return low, span, mean
}

func Second(values []int) (int, int, int) {
	low := values[0]
	high := values[0]
	sum := 0
	for _, v := range values {
		if v < low {
			low = v
		}
		if v > high {
			high = v
		}
		sum += v
	}
	span := high - low
	if span < 0 {
		span = -span
	}
	mean := sum / len(values)
	if mean < low {
		mean = low
	}
	if mean > high {
		mean = high
	}
	return low, span, mean
}
