package viewdata

import "html/template"

type StatCard struct {
	Name   string
	Icon   string
	Value  string
	Change string
}

type FuncButton struct {
	Icon  string
	Label string
	Alert string
	Value string
}

type Layout struct {
	Title          string
	Description    string
	Year           int
	PageContent    template.HTML
	Stats          []StatCard
	RightPanelSections map[string][]FuncButton
}
