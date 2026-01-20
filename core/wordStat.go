package core

type WordStat struct {
	Success int
	Failure int
}

func (ss WordStat) Total() int {
	return ss.Success + ss.Failure
}

func (ss WordStat) IncrSuccess() WordStat {
	return WordStat{
		ss.Success + 1,
		ss.Failure,
	}
}

func (ss WordStat) IncrFailure() WordStat {
	return WordStat{
		ss.Success,
		ss.Failure + 1,
	}
}

// Returns a number between 0 and 1
func (s WordStat) SuccessRate() float32 {
	total := s.Total()
	if total == 0 {
		return 0
	} else {
		return float32(s.Success) / float32(total)
	}
}
