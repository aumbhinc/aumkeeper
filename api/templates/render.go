package templates

import "aumkeeper/api/viewdata"

func ResolvePage(l *viewdata.Layout) {
	if tpl, ok := PageRegistry[l.Page]; ok {
		l.ResolvedTemplate = tpl
		return
	}
	l.ResolvedTemplate = "home_content" // fallback
}