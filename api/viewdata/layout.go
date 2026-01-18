package viewdata

import "html/template"

type Layout struct {
	Title              string
	Description        string
	Year               int
	IncludeDashboardJS bool
	IncludeStaffJS     bool
	PageContent        template.HTML
}
