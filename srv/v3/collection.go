/* Collection management handlers. */
package v3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	var col, numParts string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "numparts", &numParts) {
		return
	}
	numPartsInt, err := strconv.Atoi(numParts)
	if err != nil {
		http.Error(w, fmt.Sprintf("numparts should be an integer, %s given", numParts), 400)
		return
	}
	V3Sync.Lock()
	defer V3Sync.Unlock()
	if err := V3DB.Create(col, numPartsInt); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	} else {
		w.WriteHeader(201)
	}
}

func All(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	cols := make(map[string]interface{}, 0)
	V3Sync.Lock()
	defer V3Sync.Unlock()
	for k, v := range V3DB.StrCol {
		cols[k] = map[string]interface{}{"partitions": v.NumChunks}
	}
	resp, err := json.Marshal(cols)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	w.Write(resp)
}

func Rename(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	var oldName, newName string
	if !Require(w, r, "old", &oldName) {
		return
	}
	if !Require(w, r, "new", &newName) {
		return
	}
	V3Sync.Lock()
	defer V3Sync.Unlock()
	if err := V3DB.Rename(oldName, newName); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	}
}

func Drop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	var col string
	if !Require(w, r, "col", &col) {
		return
	}
	V3Sync.Lock()
	defer V3Sync.Unlock()
	if err := V3DB.Drop(col); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	}
}

func Scrub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	var col string
	if !Require(w, r, "col", &col) {
		return
	}
	V3Sync.Lock()
	defer V3Sync.Unlock()
	dbCol := V3DB.Use(col)
	if dbCol == nil {
		http.Error(w, fmt.Sprintf("Collection %s does not exist", col), 400)
	} else {
		V3DB.Scrub(col)
	}
}

func Repartition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	var col, numParts string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "numparts", &numParts) {
		return
	}
	numPartsInt, err := strconv.Atoi(numParts)
	if err != nil {
		http.Error(w, fmt.Sprintf("numparts should be an integer, %s given", numParts), 400)
		return
	}
	V3Sync.Lock()
	defer V3Sync.Unlock()
	dbCol := V3DB.Use(col)
	if dbCol == nil {
		http.Error(w, fmt.Sprintf("Collection %s does not exist", col), 400)
	} else {
		V3DB.Repartition(col, numPartsInt)
	}
}

func Flush(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	V3Sync.Lock()
	defer V3Sync.Unlock()
	V3DB.Flush()
}
