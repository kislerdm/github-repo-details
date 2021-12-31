package main

import (
	"log"
	"os"

	gitClient "github.com/kislerdm/github-repo-details"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	c, err := gitClient.NewClient(gitClient.WithAPIKey(token))
	if err != nil {
		log.Fatalln(err)
	}

	r, err := c.GetGraphQLData("kislerdm", "gbqschema_converter")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(r.Data.Repository.URL)
}
