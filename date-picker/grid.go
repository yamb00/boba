package datepicker

import (
	"slices"
	"sync"
	"time"
)

type gridConfig struct {
	weekdays     Weekdays
	selectedDate time.Time
	newEntry     func(date time.Time, cursor bool) Entry
}

type monthConfig struct {
	calendarKey CalendarKey
	month       *Month
	cursor      *Cursor
	setCursor   bool
	date        time.Time
}

func NewGridConfig(selectedDate time.Time, weekdays Weekdays, newEntry func(date time.Time, cursor bool) Entry) *gridConfig {
	config := &gridConfig{selectedDate: selectedDate}
	if weekdays == nil {
		config.weekdays = DefaultWeekdays
	}

	if newEntry == nil {
		config.newEntry = defaultNewEntry()
	}
	return config
}

func initGrid(config *gridConfig) (*Cursor, Calendar) {
	calendar := Calendar{}
	var cursor *Cursor

	var wg sync.WaitGroup
	channel := make(chan *monthConfig)

	monthConfigs := []*monthConfig{
		{
			date:      config.selectedDate.AddDate(0, -1, 0),
			setCursor: false,
		},
		{
			date:      config.selectedDate,
			setCursor: true,
		},
		{
			date:      config.selectedDate.AddDate(0, 1, 0),
			setCursor: false,
		},
	}

	for _, monthConfig := range monthConfigs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			channel <- initMonth(config, monthConfig)
		}()
	}
	go func() {
		wg.Wait()
		close(channel)
	}()
	for item := range channel {
		calendar[item.calendarKey] = item.month
		if item.cursor != nil {
			cursor = item.cursor
		}
	}
	return cursor, calendar
}

func initMonth(gridConfig *gridConfig, monthConfig *monthConfig) *monthConfig {
	month := &Month{}
	weekdays := gridConfig.weekdays
	newEntry := gridConfig.newEntry
	selectedDate := monthConfig.date

	monthConfig.calendarKey = newCalendarKey(selectedDate)

	currentYear, currentMonth, _ := selectedDate.Date()
	firstDay := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, selectedDate.Location())
	startDayIdx := slices.Index(weekdays, firstDay.Weekday().String())
	currentDay := firstDay
	for r := 0; r < formRows; r++ {
		if r > 0 {
			startDayIdx = 0
		}
		for c := startDayIdx; c < formColumns; c++ {
			month[r][c] = newEntry(currentDay, false)
			currentDay = currentDay.Add(time.Hour * 24)
			if monthConfig.setCursor && currentDay.Day() == selectedDate.Day() {
				monthConfig.cursor = &Cursor{row: r, col: c, date: &currentDay}
				month[r][c].setCursor()
			}
			if currentDay.Month() != currentMonth {
				break
			}
		}
	}
	monthConfig.month = month
	return monthConfig
}

func updateGrid(gridConfig *gridConfig, calendar Calendar, cursor *Cursor, direction string) Calendar {
	var config *monthConfig
	date := cursor.getMonth()
	switch direction {
	case "previous":
		config = &monthConfig{
			date:      date.AddDate(0, -1, 0),
			setCursor: false,
		}
	case "next":
		config = &monthConfig{
			date:      date.AddDate(0, +1, 0),
			setCursor: false,
		}
	}
	newMonthConfig := initMonth(gridConfig, config)
	newMonth := newMonthConfig.month
	calendarKey := newCalendarKey(config.date)
	calendar[calendarKey] = newMonth
	return calendar
}
