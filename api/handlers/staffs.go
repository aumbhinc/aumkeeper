package handlers

import (
	"net/http"
	"time"
)

// StaffsHandler renders the Staffs Manager page
func StaffsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "staffs", map[string]interface{}{
		"Title":          "Staff Manager",
		"Year":           time.Now().Year(),
		"IncludeStaffJS": true, // ensure staffs CSS & JS are loaded
	})
}
