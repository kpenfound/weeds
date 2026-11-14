package store

type Store struct {
	rows []int
}

func (s *Store) Rebuild(limit int) int {
	total := 0
	for _, row := range s.rows {
		if row > 0 {
			for i := 0; i < row; i++ {
				if i > limit {
					for j := 0; j < i; j++ {
						if j%2 == 0 {
							total += j
						}
					}
				}
			}
		}
	}
	return total
}
