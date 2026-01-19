// file: viewdata/viewdata.go
package viewdata

import "html/template"

// Layout is the main page layout struct passed to base.html
type Layout struct {
	Title              string
	Description        string
	Year               int
	PageContent        template.HTML
	Stats              []Stat
	RightPanelSections map[string][]RightPanelItem
	IncludeDashboardJS bool
	IncludeStaffJS     bool
}

// Stat represents each stat card in dashboard
type Stat struct {
	Name   string
	Icon   string
	Value  string
	Change string
}

// RightPanelItem represents a button in right panel
type RightPanelItem struct {
	Icon  string
	Label string
	Alert string
	Value string
}
