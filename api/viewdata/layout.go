package viewdata

import "html/template"

// Layout is the main page layout struct passed to base.html
// This is the ONLY struct base.html should ever know about.
type Layout struct {
	Title              string
	Description        string
	Year               int
	PageContent        template.HTML

	// Dashboard-specific data (optional, only set on dashboard)
	Stats              []DashboardStat
	RightPanelSections map[string][]RightPanelItem

	// Conditional asset flags
	IncludeDashboardJS bool
	IncludeStaffJS     bool
}

// DashboardStat represents a top stat card in the dashboard
type DashboardStat struct {
	Name   string
	Icon   string
	Value  string
	Change string
}

// RightPanelItem represents a button/action in the right panel
type RightPanelItem struct {
	Icon  string
	Label string
	Alert string
	Value string
}
