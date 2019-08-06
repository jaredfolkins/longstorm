package app

import (
	"errors"
	"log"
	"net/http"
)

type View struct {
	FlashErrors  []string
	FlashSuccess []string
	Active       *Active
	TwitterKeys  *TwitterKeys
	TmpPost      *TmpPost
	Tweets       []string
}

func NewView(w http.ResponseWriter, r *http.Request) *View {
	v := &View{
		Active:      &Active{},
		TmpPost:     &TmpPost{},
		TwitterKeys: &TwitterKeys{},
	}

	session, err := fss.Get(r, "longstorm-session")
	if err != nil {
		log.Fatal(err)
	}

	if flashes := session.Flashes("errors"); len(flashes) > 0 {
		for _, flash := range flashes {
			v.FlashErrors = append(v.FlashErrors, flash.(string))
		}
	}

	if flashes := session.Flashes("success"); len(flashes) > 0 {
		for _, flash := range flashes {
			v.FlashSuccess = append(v.FlashSuccess, flash.(string))
		}
	}

	session.Save(r, w)

	return v

}

type Active struct {
	Tweet       string
	Review      string
	History     string
	Settings    string
	CreateTweet string
	ReviewTweet string
}

type TwitterKeys struct {
	ID                int    `storm:"id"`
	AccessToken       string `storm:"unique"`
	AccessTokenSecret string `storm:"unique"`
	ConsumerAPIKey    string `storm:"unique"`
	ConsumerSecretKey string `storm:"unique"`
	HonorNewlines     bool   `storm:"unique"`
}

func (tk *TwitterKeys) Wash() error {

	if len(tk.AccessToken) < 12 {
		return errors.New("Error: AccessToken is too small.")
	}

	if len(tk.AccessTokenSecret) < 12 {
		return errors.New("Error: AccessTokenSecret is too small.")
	}

	if len(tk.ConsumerAPIKey) < 12 {
		return errors.New("Error: ConsumerAPIKey is too small.")
	}

	if len(tk.ConsumerSecretKey) < 12 {
		return errors.New("Error: ConsumerSecretKey is too small.")
	}

	return nil
}

type TmpPost struct {
	ID  int    `storm:"id"`
	Txt string `storm:"unique"`
}
