package datepicker

import "time"

// Cursor represents the cursor position in the date picker grid.
// It holds the current date too.
type Cursor struct {
	row, col int
	date     *time.Time
	entry    *Entry
}

// Date will return the cursor date
func (c Cursor) Date() *time.Time {
	return c.date
}

// setDate will set the cursor date
func (c *Cursor) setDate(d time.Time) {
	c.date = &d
}

func (c *Cursor) updateEntry(month *Month) {
	c.entry = &month[c.row][c.col]
}

// getMonth returns date of the first day in the current month
func (c *Cursor) getMonth() time.Time {
	return time.Date(c.date.Year(), c.date.Month(), 1, 0, 0, 0, 0, c.date.Location())
}

// getEntry returns entry of cursor
func (c *Cursor) getEntry() *Entry {
	return c.entry
}
