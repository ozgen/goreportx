package models

import "html/template"

type MultipleColumnReport struct {
	Header struct {
		Logo        template.HTML
		Title       string
		Subtitle    string
		Explanation string
	}
	ColumnHeaders []string   // e.g., ["Name", "Score", "Status"]
	Rows          [][]string // e.g., [["Alice", "90", "Pass"], ["Bob", "85", "Pass"]]
	Chart         struct {
		Title       string
		Image       template.HTML
		Description string
	}
	Footer struct {
		Note string
	}
}
