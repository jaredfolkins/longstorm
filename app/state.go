package app

import (
	"log"

	"github.com/asdine/storm"
	"github.com/gorilla/sessions"
)

var (
	fss *sessions.FilesystemStore
	db  *storm.DB
)

func init() {
	var err error

	fss = sessions.NewFilesystemStore("", []byte("do-not-put-this-application-on-any-network"))
	fss.MaxLength(1024 * 64) // this is the total size of temporary storage for the session, it is huge, don't ever use this setting for a real website

	db, err = storm.Open("longstorm.db")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: obviously we are not closing the db at all, we can figure that out later.
	// defer db.Close()
}
