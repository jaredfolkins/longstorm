package app

import (
	"net/http"

	"github.com/unrolled/render"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ren := render.New(render.Options{
		Extensions:      []string{".tmpl", ".html"},
		Directory:       "templates",    // Specify what path to load the templates from.
		Layout:          "layouts/main", // Specify a layout template. Layouts can call {{ yield }} to render the current template or {{ partial "css" }} to render a partial from the current template.
		RequirePartials: true,           // Return an error if a template is missing a partial used in a layout.
	})

	var tp TmpPost
	db.One("ID", 1, &tp)

	view := NewView(w, r)
	view.Active.Tweet = "active"
	view.Active.CreateTweet = "active"
	view.TmpPost = &tp

	ren.HTML(w, http.StatusOK, "index/index", view)
}
