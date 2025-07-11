package models

import "html/template"

/*
*
This is example data Model of Report
*/
type Chart struct {
	Tag         template.HTML `json:"tag"`
	Align       string        `json:"align"` // "left", "center", or "right"
	Order       int           `json:"order"`
	Description string        `json:"description"`
	Title       string        `json:"title"`
}

type Header struct {
	Title       string        `json:"title"`
	Subtitle    string        `json:"subtitle"`
	Explanation string        `json:"explanation"`
	Logo        template.HTML `json:"logo"`      // SVG/HTML tag
	LogoAlign   string        `json:"logoAlign"` // "left", "center", or "right"
}

type Footer struct {
	Note string `json:"note"`
}

type Report struct {
	Header Header            `json:"header"`
	Footer Footer            `json:"footer"`
	Data   map[string]string `json:"data"`
	Charts []Chart           `json:"charts"`
}
