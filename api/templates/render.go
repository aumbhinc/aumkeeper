package templates

import (
	"fmt"

	"aumkeeper/api/viewdata"
)

// ResolvePage converts logical route keys
// into actual template block names.
func ResolvePage(data *viewdata.Layout) error {

	if data == nil {
		return fmt.Errorf("layout is nil")
	}

	if data.Page == "" {
		return fmt.Errorf("missing page key")
	}

	templateName, ok := PageRegistry[data.Page]
	if !ok {
		return fmt.Errorf(
			"unknown page key: %s",
			data.Page,
		)
	}

	data.Page = templateName

	return nil
}