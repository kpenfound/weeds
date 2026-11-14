package worker

func Added() int {
	return 1
}

func AlsoAdded() int {
	return 2
}

func Sum(rows [][]int, limit int) int {
	total := 0
	for _, row := range rows {
		if len(row) > 0 {
			for _, v := range row {
				if v > limit {
					for i := 0; i < v; i++ {
						if i%2 == 0 {
							total += i
						}
					}
				}
			}
		}
	}
	return total
}
