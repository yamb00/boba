package datepicker

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	formRows    = 5
	formColumns = 7
)

// Month represents a grid of a month
type Month [formRows][formColumns]Entry

// CalendarKey is the key for Calendar Map
type CalendarKey string

// newCalendarKey creates a CalendarKey from a date
func newCalendarKey(date time.Time) CalendarKey {
	key := fmt.Sprintf("%s-%d", date.Month().String(), date.Year())
	return CalendarKey(key)
}

// Calendar is a map of month with the custom key CalendarKey
type Calendar map[CalendarKey]*Month

// Weekdays is a enum to represent the weekday order
type Weekdays []string

var (
	DefaultWeekdays = Weekdays{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}
)

// Styles the libgloss styles of the date-picker
type Styles struct {
	Header   lipgloss.Style
	Form     lipgloss.Style
	Selected lipgloss.Style
}

// Model hols the state of date picker
type Model struct {
	weekDays        Weekdays
	currentMonthKey CalendarKey
	calendar        Calendar
	grid            []*Month
	gridConfig      *gridConfig
	cursor          *Cursor
	Styles          Styles
	Title           string
	KeyMap          KeyMap
}

// New will return a new Model
func New(config *gridConfig) Model {
	cursor, calendar := initGrid(config)
	m := Model{
		weekDays:        config.weekdays,
		currentMonthKey: newCalendarKey(config.selectedDate),
		calendar:        calendar,
		gridConfig:      config,
		cursor:          cursor,
		Title:           config.selectedDate.Format("Jan 2006"),
		Styles:          DefaultStyles(),
		KeyMap:          DefaultKeyMap(),
	}
	m.setCursor(cursor)
	return m
}

// DefaultStyles returns a set of default style definitions for this date picker.
func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(lipgloss.Color("212")),
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(lipgloss.Color("212")).Background(lipgloss.Color("63")),
		Form:     lipgloss.NewStyle().Padding(0, 1),
	}
}

// View renders the date picker and returns a string
func (m Model) View() string {
	return m.renderHeader() + "\n" + m.renderBody()
}

// renderEntry return a string of the entry rendered with the provided styles
func renderEntry(entry Entry, s Styles) string {
	value := "  "
	style := s.Form
	if entry == nil {
		return style.Render(value)
	}
	value = fmt.Sprintf("%2d", entry.Date().Day())
	if entry.hasCursor() {
		style = s.Selected
	}
	return style.Render(value)
}

// renderHeader return the Header
func (m Model) renderHeader() string {
	var tableHeader string
	for _, col := range m.weekDays {
		day := col[0:2]
		tableHeader += m.Styles.Header.Render(day)
	}

	styleTitle := m.Styles.Header
	styleTitle.Align(lipgloss.Center)
	styleTitle.Width(len(tableHeader))
	month := m.cursor.Date().Month().String()
	year := m.cursor.Date().Year()
	title := styleTitle.Render(fmt.Sprintf("%s %d", month, year))

	return title + "\n" + tableHeader
}

// renderBody returns the body of date picker
func (m Model) renderBody() string {
	body := ""
	grid := m.calendar[m.currentMonthKey]
	for row := range formRows {
		for col := range formColumns {
			entry := grid[row][col]
			body += renderEntry(entry, m.Styles)
		}
		body += "\n"
	}
	return body
}

// getCurrentMonth retursn the active Month
func (m *Model) getCurrentMonth() *Month {
	return m.calendar[m.currentMonthKey]
}

// unsetCursor removes the focus of the entry from the current cursor position
func (m *Model) unsetCursor() {
	row := m.cursor.row
	col := m.cursor.col
	m.calendar[m.currentMonthKey][row][col].unsetCursor()
}

// setCursor adds the focus of the entry from the current cursor position
func (m *Model) setCursor(cursor *Cursor) {
	row := cursor.row
	col := cursor.col
	m.cursor = cursor
	m.calendar[m.currentMonthKey][row][col].setCursor()
	m.cursor.setDate(m.calendar[m.currentMonthKey][row][col].Date())
	m.cursor.updateEntry(m.getCurrentMonth())
}

// findValidDateAndSet walks in a row by direction and finds a valid date starting
// at the cursor position.
func (m *Model) findValidDate(c *Cursor, direction string) {
	if m.validDate(c) {
		m.setCursor(c)
		return
	} else {
		switch direction {
		case "left":
			c.col--

		case "right":
			c.col++
		}
	}
	m.findValidDate(c, direction)
}

func (m *Model) validDate(cursor *Cursor) bool {
	if m.getCurrentMonth()[cursor.row][cursor.col] == nil {
		return false
	}
	return true
}

// previousMonth set the previous month as current month
func (m *Model) previousMonth() {
	currentMonth := m.cursor.getMonth()
	previousMonth := currentMonth.AddDate(0, -1, 0)
	m.currentMonthKey = newCalendarKey(previousMonth)
}

// previousMonth set the next month as current month
func (m *Model) nextMonth() {
	currentMonth := m.cursor.getMonth()
	nextMonth := currentMonth.AddDate(0, 1, 0)
	m.currentMonthKey = newCalendarKey(nextMonth)
}

// CursorUp moves cursor up.
func (m *Model) CursorUp() {
	newCursor := &Cursor{
		row: m.cursor.row - 1,
		col: m.cursor.col,
	}
	if newCursor.row >= 0 && m.validDate(newCursor) {
		m.unsetCursor()
		m.setCursor(newCursor)
	}
}

// CursorUp move cursor down.
func (m *Model) CursorDown() {
	newCursor := &Cursor{
		row: m.cursor.row + 1,
		col: m.cursor.col,
	}
	if newCursor.row < formRows && m.validDate(newCursor) {
		m.unsetCursor()
		m.setCursor(newCursor)
	}
}

// CursorLeft moves cursor left.
// When col 0 the cursor will jump into the previous month
func (m *Model) CursorLeft() {
	newCursor := &Cursor{
		row: m.cursor.row,
		col: m.cursor.col - 1,
	}
	if newCursor.col >= 0 && m.validDate(newCursor) {
		m.unsetCursor()
		m.setCursor(newCursor)
	} else {
		newCursor.col = formColumns - 1
		m.unsetCursor()
		m.previousMonth()
		// m.setCursor(newCursor)
		// foundCursor := findValidDateAndSet(newCursor, m.getCurrentMonth(), "left")
		// m.setCursor(foundCursor)
		m.findValidDate(newCursor, "left")
		m.calendar = updateGrid(m.gridConfig, m.calendar, m.cursor, "previous")
	}
}

// CursorRight will move cursor to the right.
// When col 6 the cursor will jump into the next month.
func (m *Model) CursorRight() {
	newCursor := &Cursor{
		row: m.cursor.row,
		col: m.cursor.col + 1,
	}
	if newCursor.col < formColumns && m.validDate(newCursor) {
		m.unsetCursor()
		m.setCursor(newCursor)
	} else {
		newCursor.col = 0
		m.unsetCursor()
		m.nextMonth()
		m.findValidDate(newCursor, "right")
		m.calendar = updateGrid(m.gridConfig, m.calendar, m.cursor, "next")
	}
}

// Update will update the model by the command
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.CursorUp):
			m.CursorUp()
		case key.Matches(msg, m.KeyMap.CursorDown):
			m.CursorDown()
		case key.Matches(msg, m.KeyMap.CursorLeft):
			m.CursorLeft()
		case key.Matches(msg, m.KeyMap.CursorRight):
			m.CursorRight()
		}
	}
	return m, nil
}

// GetCursorDate return date of the cursor position
func (m Model) GetCursorDate() time.Time {
	return *m.cursor.date
}

func (m Model) GetCursorEntry() Entry {
	return *m.cursor.getEntry()
}
