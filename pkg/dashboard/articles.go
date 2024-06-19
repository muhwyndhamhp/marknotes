package dashboard

import (
	"context"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/db"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	pub_alert "github.com/muhwyndhamhp/marknotes/pub/components/alert"
	pub_dashboards_articles_new "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/articles/create"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
	"github.com/muhwyndhamhp/marknotes/utils/renderfile"
	"github.com/muhwyndhamhp/marknotes/utils/rss"
	"github.com/muhwyndhamhp/marknotes/utils/sanitizations"
	"github.com/muhwyndhamhp/marknotes/utils/strman"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

func (fe *DashboardFrontend) ArticlesEdit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := fe.PostRepo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	opts := pub_variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs(fmt.Sprintf("dashboard/articles/%d/edit", id)),
	}

	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, id)

	vm := pub_dashboards_articles_new.NewVM{
		Opts:      opts,
		UploadURL: uploadURL,
		Post:      post,
		BaseURL:   baseURL,
	}
	articlesNew := pub_dashboards_articles_new.New(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}

func (fe *DashboardFrontend) ArticlesPush(c echo.Context) error {
	ctx := c.Request().Context()

	status := values.PostStatus(c.QueryParam("status"))
	if status == values.None {
		status = values.Draft
	}

	existingID, _ := strconv.Atoi(c.QueryParam("existingID"))
	var xp *models.Post
	if existingID != 0 {
		p, err := fe.PostRepo.GetByID(ctx, uint(existingID))
		if err != nil {
			return errs.Wrap(err)
		}
		xp = p
	}

	var req ArticlesCreateRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		fail := pub_alert.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	if err := c.Validate(req); err != nil {
		fmt.Println(err)
		fail := pub_alert.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	sanitizedHTML := sanitizations.SanitizeHtml(req.Content)

	slug, err := strman.GenerateSlug(req.Title)
	if err != nil {
		fmt.Println(err)
		fail := pub_alert.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	count := fe.PostRepo.Count(ctx, db.Where("slug = ?", slug))
	if count > 0 && xp == nil {
		fmt.Println(count)
		slug = fmt.Sprintf("%s-%d", slug, count)
	}

	post := tern.Struct(xp, &models.Post{})

	usr, err := fe.ClerkClient.GetUserFromSession(c, fe.UserRepo)
	if err != nil {
		fmt.Println(err)
		fail := pub_alert.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	post.Title = req.Title
	post.Content = req.ContentJSON
	post.EncodedContent = template.HTML(sanitizedHTML)
	post.Status = status
	post.UserID = usr.ID
	post.Slug = slug
	post.HeaderImageURL = req.HeaderImageURL
	post.MarkdownContent = req.MarkdownContent
	if status == values.Published && post.PublishedAt.IsZero() {
		now := time.Now()
		post.PublishedAt = now
	}

	tags := strings.Split(req.Tags, ",")
	post.Tags = []*models.Tag{}

	tagLiteral := ""
	for i := range tags {
		if tags[i] == "" {
			continue
		}
		tagName := strings.ToLower(strings.TrimSpace(tags[i]))
		tagSlug := strings.ReplaceAll(tagName, " ", "-")
		res, _ := fe.TagRepo.Get(ctx, db.Where("slug = ?", tagSlug))
		var vTag models.Tag
		if len(res) != 0 {
			vTag = res[0]
		} else {
			tag := models.Tag{
				Slug:  tagSlug,
				Title: tags[i],
			}
			err := fe.TagRepo.Upsert(c.Request().Context(), &tag)
			if err != nil {
				return errs.Wrap(err)
			}
			vTag = tag
		}

		post.Tags = append(post.Tags, &vTag)
		tagLiteral += vTag.Title + ","
	}
	post.TagsLiteral = tagLiteral

	err = fe.PostRepo.Upsert(ctx, post)
	if err != nil {
		fail := pub_alert.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	if status == values.Published {
		go func() {
			ctx := context.Background()
			renderfile.RenderPost(ctx, post, fe.Bucket)
			err := renderfile.RenderMarkdown(post)
			if err != nil {
				fmt.Println(err)
			}

			err = rss.GenerateRSS(ctx, fe.PostRepo)
			if err != nil {
				fmt.Println(err)
			}
		}()
	} else {
		go func() {
			err := fileman.DeleteFile(fmt.Sprintf(config.Get(config.POST_RENDER_PATH)+"/%s.html", post.Slug))
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	if existingID != 0 {
		success := pub_alert.AlertSuccess("Post has been saved successfully")

		return templates.AssertRender(c, http.StatusOK, success)
	} else {
		opts := pub_variables.DashboardOpts{
			Nav:         nav(0),
			BreadCrumbs: fe.Breadcrumbs(fmt.Sprintf("dashboard/articles/%d/edit", post.ID)),
		}

		baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
		uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, post.ID)
		vm := pub_dashboards_articles_new.NewVM{
			Opts:      opts,
			UploadURL: uploadURL,
			Post:      post,
			BaseURL:   baseURL,
		}
		articlesNew := pub_dashboards_articles_new.New(vm)

		c.Response().Header().Set("HX-Replace-Url", fmt.Sprintf("/dashboard/articles/%d/edit", post.ID))
		return templates.AssertRender(c, http.StatusOK, articlesNew)
	}
}

func (fe *DashboardFrontend) ArticlesDelete(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	err := fe.PostRepo.Delete(ctx, uint(id))
	if err != nil {
		return err
	}

	go func() {
		err := rss.GenerateRSS(ctx, fe.PostRepo)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return templates.RenderEmpty(c)
}
