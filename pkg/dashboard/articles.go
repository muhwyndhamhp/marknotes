package dashboard

import (
	"context"
	"errors"
	"fmt"
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
	pub_dashboards_articles "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/articles"
	pub_dashboards_articles_new "github.com/muhwyndhamhp/marknotes/pub/pages/dashboards/articles/create"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/renderfile"
	"github.com/muhwyndhamhp/marknotes/utils/sanitizations"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
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
	}
	articlesNew := pub_dashboards_articles_new.New(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}

func (fe *DashboardFrontend) ArticlesNew(c echo.Context) error {
	opts := pub_variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles/new"),
	}

	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, 0)

	vm := pub_dashboards_articles_new.NewVM{
		Opts:      opts,
		UploadURL: uploadURL,
	}
	articlesNew := pub_dashboards_articles_new.New(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}

func (fe *DashboardFrontend) Articles(c echo.Context) error {
	ctx := c.Request().Context()
	page, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))
	source := c.QueryParam(constants.TARGET_SOURCE)

	hx_request, _ := strconv.ParseBool(c.Request().Header.Get("Hx-Request"))

	partial := source == constants.TARGET_SOURCE_PARTIAL && hx_request

	count := fe.PostRepo.Count(ctx)

	posts, err := fe.PostRepo.Get(ctx,
		scopes.Paginate(page, pageSize),
		scopes.OrderBy("created_at", scopes.Descending),
		scopes.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, values.None, "", "")
	}

	opts := pub_variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles"),
	}

	pageSizes := fe.SizeDropdown(page, pageSize)
	pages := fe.PageDropdown(tern.Int(page, 1), tern.Int(pageSize, 10), count)
	articleVM := pub_dashboards_articles.ArticlesVM{
		Opts:       opts,
		Posts:      posts,
		PageSizes:  pageSizes,
		Pages:      pages,
		CreatePath: "/dashboard/articles/new",
	}

	dashboard := pub_dashboards_articles.Articles(articleVM)

	if !partial {
		return templates.AssertRender(c, http.StatusOK, dashboard)
	} else {
		articles := pub_dashboards_articles.ArticleOOB(posts, pageSizes, pages)
		return templates.AssertRender(c, http.StatusOK, articles)
	}
}

func (fe *DashboardFrontend) ArticlesPush(c echo.Context) error {
	claims, _ := c.Get(jwt.AuthClaimKey).(*jwt.Claims)
	if claims == nil {
		return c.JSON(http.StatusBadRequest, errors.New("claim is empty"))
	}

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

	post := tern.Struct(xp, &models.Post{})

	post.Title = req.Title
	post.Content = req.ContentJSON
	post.EncodedContent = template.HTML(sanitizedHTML)
	post.Status = status
	post.UserID = claims.UserID
	post.Slug = slug
	if status == values.Published && post.PublishedAt.IsZero() {
		now := time.Now()
		post.PublishedAt = now
	}

	tags := strings.Split(req.Tags, ",")
	post.Tags = []*models.Tag{}

	for i := range tags {
		if tags[i] == "" {
			continue
		}
		tagName := strings.ToLower(strings.TrimSpace(tags[i]))
		tagSlug := strings.ReplaceAll(tagName, " ", "-")
		res, _ := fe.TagRepo.Get(ctx, scopes.Where("slug = ?", tagSlug))
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
	}

	err = fe.PostRepo.Upsert(ctx, post)
	if err != nil {
		fail := pub_alert.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	if status == values.Published {
		go func() {
			ctx := context.Background()
			renderfile.RenderPost(ctx, post)
		}()
	} else {
		go func() {
			err := fileman.DeleteFile(fmt.Sprintf("public/articles/%s.html", post.Slug))
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	success := pub_alert.AlertSuccess("Post has been saved successfully")

	return templates.AssertRender(c, http.StatusOK, success)
}
