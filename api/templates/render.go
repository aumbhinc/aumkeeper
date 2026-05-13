package templates

import "aumkeeper/api/viewdata"

func ResolvePage(l *viewdata.Layout) {
	if tpl, ok := PageRegistry[l.Page]; ok {
		l.ResolvedTemplate = tpl
		return
	}

	panic("Template Error: Unknown or missing page → " + l.Page)
}