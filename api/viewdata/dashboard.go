package viewdata

type DashboardStat struct {
	Name   string
	Value  string
	Change string
	Icon   string
}

type DashboardPanelItem struct {
	Icon  string
	Label string
	Alert string
	Value string
}

type Dashboard struct {
	Stats              []DashboardStat
	RightPanelSections map[string][]DashboardPanelItem
}
