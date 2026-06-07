package render

import (
"bytes"
"fmt"
"html/template"
"log"
"net/http"
)

# /*

# TYPE-SAFE VIEW MODELS

*/

type FormData struct {
FirstName  string
MiddleName string
LastName   string

```
Role string

Email       string
PhoneNumber string

Street string
City   string
State  string
ZipCode string

SSN string

DependentClaims string

Wage             string
PaymentFrequency string

Comments string
```

}

# /*

# RENDER DATA

*/

type RenderData struct {
Title       string
Description string
Page        string

```
FormData FormData
Errors   map[string]string
```

}

# /*

# RENDERER

*/

type Renderer struct {
Tmpl  *template.Template
Debug bool
}

# /*

# CONSTRUCTOR

*/

func NewRenderer(
t *template.Template,
debug bool,
) *Renderer {
return &Renderer{
Tmpl:  t,
Debug: debug,
}
}

# /*

# CORE RENDER

*/

func (r *Renderer) Render(
w http.ResponseWriter,
page string,
data *RenderData,
) {
if r == nil || r.Tmpl == nil {
http.Error(
w,
"renderer not initialized",
http.StatusInternalServerError,
)
return
}

```
if data == nil {
	data = &RenderData{}
}

if data.Errors == nil {
	data.Errors = map[string]string{}
}

if r.Debug {
	log.Println("🧠 [RENDER START]")
	log.Println("➡️ Page:", page)
	log.Println("➡️ Title:", data.Title)
}

t := r.Tmpl.Lookup(page)
if t == nil {
	msg := fmt.Sprintf(
		"template not found: %s",
		page,
	)

	log.Println("❌", msg)

	http.Error(
		w,
		msg,
		http.StatusInternalServerError,
	)

	return
}

var buf bytes.Buffer

if err := t.Execute(&buf, data); err != nil {
	log.Println(
		"❌ TEMPLATE EXEC ERROR:",
		err,
	)

	http.Error(
		w,
		"template execution error",
		http.StatusInternalServerError,
	)

	return
}

w.Header().Set(
	"Content-Type",
	"text/html; charset=utf-8",
)

w.WriteHeader(http.StatusOK)

if _, err := buf.WriteTo(w); err != nil {
	log.Println(
		"❌ RESPONSE WRITE ERROR:",
		err,
	)

	return
}

if r.Debug {
	log.Println("✅ [RENDER SUCCESS]")
}
```

}

# /*

# HELPERS

*/

func (r *Renderer) OK(
w http.ResponseWriter,
page string,
data *RenderData,
) {
r.Render(w, page, data)
}
