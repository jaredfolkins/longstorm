package app

import (
	"log"
	"path"
	"strings"

	"github.com/asdine/storm"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var (
	fss       *sessions.FilesystemStore
	db        *storm.DB
	boxTmpl   *packr.Box
	boxAssets *packr.Box
)

func init() {
	var err error

	boxTmpl = packr.New("templates", "../templates")
	boxAssets = packr.New("assets", "../assets")

	fss = sessions.NewFilesystemStore("", []byte("do-not-put-this-application-on-any-network"))
	fss.MaxLength(1024 * 64) // this is the total size of temporary storage for the session, it is huge, don't ever use this setting for a real website

	db, err = storm.Open("longstorm.db")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: obviously we are not closing the db at all, we can figure that out later.
	// defer db.Close()
}

func NewRender(layout string, box *packr.Box) *render.Render {
	dummyDir := "__DUM__"
	return render.New(render.Options{
		Directory: dummyDir, // Specify what path to load the templates from.
		Asset: func(name string) ([]byte, error) {
			name = strings.TrimPrefix(name, dummyDir)
			name = strings.TrimPrefix(name, "/")
			return box.Find(name)
		},
		AssetNames: func() []string {
			names := box.List()
			for k, v := range names {
				pth := path.Join(dummyDir, v)
				names[k] = pth
			}
			return names
		},
		Extensions:      []string{".tmpl", ".html"},
		Layout:          layout, // Specify a layout template. Layouts can call {{ yield }} to render the current template or {{ partial "css" }} to render a partial from the current template.
		RequirePartials: true,   // Return an error if a template is missing a partial used in a layout.
	})
}
