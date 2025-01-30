package datepicker

import "time"

type Entry interface {
	Date() time.Time
	hasCursor() bool
	setCursor()
	unsetCursor()
}

type DefaultEntry struct {
	date   time.Time
	cursor bool
}

func (d *DefaultEntry) Date() time.Time {
	return d.date
}

func (d *DefaultEntry) hasCursor() bool {
	return d.cursor
}

func (d *DefaultEntry) setCursor() {
	d.cursor = true
}

func (d *DefaultEntry) unsetCursor() {
	d.cursor = false
}

func defaultNewEntry() func(time.Time, bool) Entry {
	return func(date time.Time, cursor bool) Entry {
		return &DefaultEntry{
			date:   date,
			cursor: cursor,
		}
	}
}
