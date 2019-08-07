package app

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

func Asset(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	fp := filepath.Join(v["dir"], v["file"])
	b, err := boxAssets.Find(fp)
	if err != nil {
		http.Error(w, fmt.Sprintf("404 Not Found: %s", v["file"]), 404)
		return
	}

	contentType := mime.TypeByExtension(path.Ext(fp))
	if contentType == "" {
		contentType = http.DetectContentType(b)
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(b)
	return
}
