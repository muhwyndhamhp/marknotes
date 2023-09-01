package site

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/StudioSol/sitemap"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
)

const (
	BaseUrl = "https://mwyndham.dev/"
)

func PingSitemap(postRepo models.PostRepository) {
	ctx := context.Background()
	lastMod := time.Now()
	pg := sitemap.NewSitemapGroup("post", false)
	pg.Configure("public/sitemap/post", false)

	posts, err := postRepo.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := range posts {
		pg.Add(sitemap.URL{Loc: appendPath(fmt.Sprintf("posts/%d", posts[i].ID)), LastMod: &lastMod})
		pg.Add(sitemap.URL{Loc: appendPath(fmt.Sprintf("articles/%s", posts[i].Slug)), LastMod: &lastMod})
	}

	pg.Add(sitemap.URL{Loc: appendPath("posts_index"), LastMod: &lastMod})
	pg.Add(sitemap.URL{Loc: appendPath("articles"), LastMod: &lastMod})
	pg.Add(sitemap.URL{Loc: appendPath("posts_manage"), LastMod: &lastMod})

	ag := sitemap.NewSitemapGroup("admin", false)
	ag.Configure("public/sitemap/admin", false)
	ag.Add(sitemap.URL{Loc: appendPath(""), LastMod: &lastMod})
	ag.Add(sitemap.URL{Loc: appendPath("unauthorized"), LastMod: &lastMod})
	ag.Add(sitemap.URL{Loc: appendPath("resume"), LastMod: &lastMod})
	ag.Add(sitemap.URL{Loc: appendPath("contact"), LastMod: &lastMod})

	postFiles := pg.Files()
	adminFiles := ag.Files()

	for file := range postFiles {
		saveFile(file)
	}

	for file := range adminFiles {
		saveFile(file)
	}

	URLs := append(pg.URLs(), ag.URLs()...)
	index := sitemap.CreateIndexBySlice(URLs, BaseUrl)

	log.Println("creating index...")
	err = sitemap.CreateSitemapIndex("public/sitemap/index.xml.gz", index)
	if err != nil {
		log.Fatal(err)
	}

	sitemap.PingSearchEngines(appendPath("public/sitemap/index.xml.gz"))
}

func saveFile(file sitemap.File) {
	log.Println(file.Name)

	f, err := os.Create(file.Name)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Write(f)
	if err != nil {
		log.Fatal(err)
	}
}

func appendPath(path string) string {
	return fmt.Sprintf("%s%s", BaseUrl, path)
}
