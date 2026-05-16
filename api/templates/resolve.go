package templates

import (
	"aumkeeper/api/viewdata"
)

// ResolvePage is now deprecated.
// Renderer + base.html handle routing directly.
func ResolvePage(data *viewdata.Layout) error {
	// NO-OP to prevent routing conflicts
	if data == nil {
		return nil
	}
	return nil
}