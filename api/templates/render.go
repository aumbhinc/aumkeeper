package templates

import (
	"aumkeeper/api/viewdata"
)

// ResolvePage is deprecated.
// Renderer + base.html now control routing.
// This is kept only for backward compatibility.
func ResolvePage(data *viewdata.Layout) error {
	if data == nil {
		return nil
	}

	// NO routing transformation anymore
	return nil
}