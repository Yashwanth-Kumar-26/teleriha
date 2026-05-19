package bot

import (
	"fmt"
)

// PaginationBuilder helps create inline keyboards for pagination.
type PaginationBuilder struct {
	TotalPages   int
	CurrentPage  int
	Prefix       string // The callback data prefix, e.g., "page_"
	ItemsPerRow  int
	ShowCounters bool // e.g. "1 / 5"
}

// NewPaginationBuilder creates a new pagination builder.
func NewPaginationBuilder(totalPages, currentPage int, prefix string) *PaginationBuilder {
	return &PaginationBuilder{
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		Prefix:       prefix,
		ItemsPerRow:  3,
		ShowCounters: true,
	}
}

// Build generates the InlineKeyboardMarkup for the pagination.
func (p *PaginationBuilder) Build() *InlineKeyboardMarkup {
	if p.TotalPages <= 1 {
		return &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{}}
	}

	builder := NewInlineKeyboardBuilder()
	var row []InlineKeyboardButton

	// Prev button
	if p.CurrentPage > 1 {
		row = append(row, InlineButton("«", fmt.Sprintf("%s%d", p.Prefix, p.CurrentPage-1)))
	} else {
		row = append(row, InlineButton("·", "ignore"))
	}

	// Counter
	if p.ShowCounters {
		counterText := fmt.Sprintf("%d / %d", p.CurrentPage, p.TotalPages)
		row = append(row, InlineButton(counterText, "ignore"))
	}

	// Next button
	if p.CurrentPage < p.TotalPages {
		row = append(row, InlineButton("»", fmt.Sprintf("%s%d", p.Prefix, p.CurrentPage+1)))
	} else {
		row = append(row, InlineButton("·", "ignore"))
	}

	builder.AddRow(row...)
	return builder.Build()
}
