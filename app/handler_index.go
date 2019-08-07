package app

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ren := NewRender("layouts/main", boxTmpl)

	var tp TmpPost
	db.One("ID", 1, &tp)

	view := NewView(w, r)
	view.Active.Tweet = "active"
	view.Active.CreateTweet = "active"
	view.TmpPost = &tp

	ren.HTML(w, http.StatusOK, "index/index", view)
}
