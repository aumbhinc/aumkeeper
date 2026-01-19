package viewdata

type Stat struct {
	Name   string
	Value  string
	Change string
	Icon   string
}

type PanelItem struct {
	Icon  string
	Label string
	Alert string
	Value string
}

type Dashboard struct {
	Stats               []Stat
	RightPanelSections  map[string][]PanelItem
}
