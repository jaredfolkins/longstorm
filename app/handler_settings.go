package app

import (
	"log"
	"net/http"

	"github.com/unrolled/render"
)

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	ren := render.New(render.Options{
		Extensions:      []string{".tmpl", ".html"},
		Directory:       "templates",    // Specify what path to load the templates from.
		Layout:          "layouts/main", // Specify a layout template. Layouts can call {{ yield }} to render the current template or {{ partial "css" }} to render a partial from the current template.
		RequirePartials: true,           // Return an error if a template is missing a partial used in a layout.
	})

	view := NewView(w, r)
	view.Active.Settings = "active"

	var tk TwitterKeys
	db.One("ID", 1, &tk)
	view.TwitterKeys = &tk
	ren.HTML(w, http.StatusOK, "index/settings", view)
}

func SaveSettingsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := fss.Get(r, "longstorm-session")
	if err != nil {
		log.Fatal(err)
	}

	tk := &TwitterKeys{
		ID:                1,
		AccessToken:       r.FormValue("access-token"),
		AccessTokenSecret: r.FormValue("access-token-secret"),
		ConsumerAPIKey:    r.FormValue("consumer-api-key"),
		ConsumerSecretKey: r.FormValue("consumer-secret-key"),
		HonorNewlines:     convertHonorNewlines(r.FormValue("honor-newlines")),
	}

	err = tk.Wash()
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/settings", 302)
		return
	}

	err = db.Save(tk)
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/settings", 302)
		return
	}

	session.AddFlash("Twitter keys saved successfully.", "success")
	session.Save(r, w)
	http.Redirect(w, r, "/settings", 302)
}

func convertHonorNewlines(s string) bool {
	if s == "yes" {
		return true
	}
	return false
}
