package viewdata

// Layout is the main page container
type Layout struct {
	Title       string
	Description string
	Year        int

	Page            string
	ResolvedTemplate string

	PageData any

	Stats              []DashboardStat
	RightPanelSections map[string][]RightPanelItem

	IncludeDashboardJS bool
	IncludeStaffJS     bool
}

// ✅ REQUIRED: DashboardStat must exist
type DashboardStat struct {
	Name   string
	Icon   string
	Value  string
	Change string
}

// ✅ REQUIRED: Right panel items must exist
type RightPanelItem struct {
	Icon  string
	Label string
	Alert string
	Value string
}