package rss

import (
	"context"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

func GenerateRSS(ctx context.Context, repo models.PostRepository) error {
	posts, err := repo.Get(ctx, scopes.Where("status = ?", values.Published))
	if err != nil {
		return err
	}

	rsstr, err := generateRSSString(posts)
	if err != nil {
		return err
	}

	err = fileman.DeleteFile(constants.RSS_PATH)
	if err != nil {
		return err
	}

	file, err := os.Create(constants.RSS_PATH)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(rsstr))
	if err != nil {
		return err
	}

	return nil
}

func generateRSSString(posts []models.Post) (string, error) {
	author := &feeds.Author{
		Name:  "M Wyndham",
		Email: "business@mwyndham.dev",
	}

	postFeeds := &feeds.Feed{
		Title:       "mwyndham.dev blog",
		Link:        &feeds.Link{Href: "https://mwyndham.dev"},
		Description: "my persona blog, contain my personal thought and learning experience",
		Author:      author,
		Created:     time.Now(),
		Items:       []*feeds.Item{},
	}

	for _, post := range posts {
		postFeeds.Items = append(postFeeds.Items, &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: post.GenerateURL()},
			Author:      author,
			Description: post.Abstract,
			Updated:     post.UpdatedAt,
			Created:     post.PublishedAt,
			Content:     string(post.EncodedContent),
		})
	}

	rss, err := postFeeds.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}
