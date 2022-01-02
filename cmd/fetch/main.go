package main

import (
	"log"
	"os"

	"github.com/goccy/go-json"
	gitClient "github.com/kislerdm/github-repo-details"
	"github.com/kislerdm/github-repo-details/internal/fs"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	c, err := gitClient.NewClient(gitClient.WithAPIKey(token))
	if err != nil {
		log.Fatalln(err)
	}

	r, err := c.FetchDetails("kislerdm", "gbqschema_converter")
	if err != nil {
		log.Fatalln(err)
	}

	o, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}

	if err := fs.FileWrite(o, "/tmp/fetch_out.json"); err != nil {
		log.Fatalln(err)
	}
}
