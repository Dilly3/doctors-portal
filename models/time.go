package models

type Time struct {
	StartTime rune
	EndTime   rune
	TimeList  map[rune]string
}

func (t *Time) UpdateTimeList(start rune, end rune) {
	for i := start; i <= end; i++ {
		if i < 12 {
			t.TimeList[i] = string(i) + ":" + "AM"
		}
		if i == 12 {
			t.TimeList[i] = string(i) + ":" + "PM"
		}
		if i > 12 {
			t.TimeList[i] = string(i-12) + ":" + "PM"
		}

	}
}
