package handlers

import (
	"html/template"
	"net/http"
	"time"
)

// Staff represents one staff member
type Staff struct {
	Name       string
	Access     string
	TaskScore  string
	Pay        string
	Task       string
}

// StaffsPageData is passed to template
type StaffsPageData struct {
	Title string
	Year  int
	Staff []Staff
}

// StaffsHandler renders staff page
func StaffsHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		staffs := []Staff{
			{Name: "John Doe", Access: "General Manager", TaskScore: "90/91", Pay: "$9000/month", Task: "Design the UI for new project"},
			{Name: "Jane Smith", Access: "Staff", TaskScore: "90/91", Pay: "$17/hour", Task: "Design the UI for new project"},
			{Name: "Rachel Ying", Access: "Manager", TaskScore: "100/91", Pay: "$5000/month", Task: "Design the UI for new project"},
		}

		data := StaffsPageData{
			Title: "Staff Manager",
			Year:  time.Now().Year(),
			Staff: staffs,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "staffs", data); err != nil {
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
