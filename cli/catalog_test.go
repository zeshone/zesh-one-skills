package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSourcePinsCatalogAndFilesToMainCommit(t *testing.T) {
	cat, files := specimen()
	commit := strings.Repeat("b", 40)
	requests := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.Path)
		switch r.URL.Path {
		case "/commits/main":
			fmt.Fprintf(w, `{"sha":%q}`, commit)
		case "/" + commit + "/skills-index.json":
			json.NewEncoder(w).Encode(cat)
		default:
			data, ok := files[strings.TrimPrefix(r.URL.Path, "/"+commit+"/")]
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Write(data)
		}
	}))
	defer server.Close()
	source := Source{Client: server.Client(), API: server.URL, Raw: server.URL}
	got, sha, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	if sha != commit || len(got.Skills) != 2 {
		t.Fatal("catalog mismatch")
	}
	if err = (Installer{Fetch: source.Fetch}).Install(t.TempDir(), got, sha, []string{"frontend/design"}, []string{"pi"}, "mockup", nil); err != nil {
		t.Fatal(err)
	}
	for _, p := range requests[1:] {
		if !strings.HasPrefix(p, "/"+commit+"/") {
			t.Fatalf("unpinned request: %s", p)
		}
	}
}
func TestSourceRejectsInvalidAndUnavailableResponses(t *testing.T) {
	for _, mode := range []string{"status", "commit", "json", "schema"} {
		t.Run(mode, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "status" {
					http.Error(w, "private internal details", 503)
					return
				}
				if r.URL.Path == "/commits/main" {
					sha := strings.Repeat("a", 40)
					if mode == "commit" {
						sha = "../bad"
					}
					fmt.Fprintf(w, `{"sha":%q}`, sha)
					return
				}
				if mode == "json" {
					fmt.Fprint(w, "invalid")
					return
				}
				fmt.Fprint(w, `{"schemaVersion":99,"skills":[]}`)
			}))
			defer server.Close()
			source := Source{Client: server.Client(), API: server.URL, Raw: server.URL}
			_, _, err := source.Load()
			if err == nil {
				t.Fatal("accepted invalid response")
			}
			if strings.Contains(err.Error(), "private internal") {
				t.Fatal("leaked response")
			}
		})
	}
}
func TestMalformedAndModifiedIndexPreservesUserText(t *testing.T) {
	for _, bad := range []string{blockEnd + blockStart, blockStart + blockStart + blockEnd, blockEnd} {
		if _, _, err := mergeBlock([]byte(bad), []byte("new")); err == nil {
			t.Fatal("malformed markers accepted")
		}
	}
}
