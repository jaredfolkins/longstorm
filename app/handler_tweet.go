package app

import (
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/jaredfolkins/longstorm/helpers"
	"github.com/unrolled/render"
)

func TweetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := fss.Get(r, "longstorm-session")
	if err != nil {
		log.Fatal(err)
	}

	tweet := r.FormValue("tweet")
	tp := &TmpPost{
		ID:  1,
		Txt: tweet,
	}

	err = db.Save(tp)
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return
	}

	session.AddFlash("Please review and then click 🌩️ LongStorm 🌩️ at the bottom.", "success")
	session.Save(r, w)

	http.Redirect(w, r, "/review", 302)
}

func LongStormHandler(w http.ResponseWriter, r *http.Request) {
	session, err := fss.Get(r, "longstorm-session")
	if err != nil {
		log.Fatal(err)
	}

	var tk TwitterKeys
	db.One("ID", 1, &tk)

	var tp TmpPost
	db.One("ID", 1, &tp)

	tw := helpers.NewTweetWorker(tk.ConsumerAPIKey, tk.ConsumerSecretKey, tk.AccessToken, tk.AccessTokenSecret)
	tweets := tw.Storm(tp.Txt, tk.HonorNewlines)
	tv, err := tw.FirstTweet(tweets)
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/review", 302)
		return
	}

	param := &twitter.StatusUpdateParams{
		InReplyToStatusID: tv.ID,
	}

	for k, tweet := range tweets {
		if k != 0 {
			t, err := tw.AppendTweet(tweet, param)
			if err != nil {
				session.AddFlash(err.Error(), "errors")
				session.Save(r, w)
				http.Redirect(w, r, "/review", 302)
				return
			}
			param.InReplyToStatusID = t.ID
		}
	}

	ntp := &TmpPost{ID: 1}
	err = db.Save(ntp)
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/review", 302)
		return
	}

	session.AddFlash("Your 🌩️ LongStorm 🌩️ was successfully posted.", "success")
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func ReviewHandler(w http.ResponseWriter, r *http.Request) {
	session, err := fss.Get(r, "longstorm-session")
	if err != nil {
		log.Fatal(err)
	}

	ren := render.New(render.Options{
		Extensions:      []string{".tmpl", ".html"},
		Directory:       "templates",    // Specify what path to load the templates from.
		Layout:          "layouts/main", // Specify a layout template. Layouts can call {{ yield }} to render the current template or {{ partial "css" }} to render a partial from the current template.
		RequirePartials: true,           // Return an error if a template is missing a partial used in a layout.
	})

	view := NewView(w, r)
	view.Active.Tweet = "active"
	view.Active.ReviewTweet = "active"
	var tp TmpPost
	err = db.One("ID", 1, &tp)
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return
	}

	if len(tp.Txt) < 280 {
		session.AddFlash("Your LongStorm must be at 281 characters.", "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return
	}

	var tk TwitterKeys
	err = db.One("ID", 1, &tk)
	if err != nil {
		session.AddFlash(err.Error(), "errors")
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return
	}
	tw := &helpers.TweetWorker{}
	view.Tweets = tw.Storm(tp.Txt, tk.HonorNewlines)
	ren.HTML(w, http.StatusOK, "index/review", view)
}
