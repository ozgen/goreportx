package models

import "html/template"

type SimpleReport struct {
	Header SimpleHeader
	Chart  SimpleChart
	Footer SimpleFooter
}

type SimpleHeader struct {
	Title       string
	Subtitle    string
	Explanation string
	Logo        template.HTML
}

type SimpleChart struct {
	Image       template.HTML
	Align       string
	Title       string
	Description string
}

type SimpleFooter struct {
	Note string
}
