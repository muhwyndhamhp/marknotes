package dashboard

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/dashboard/articles"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/config"
	templates "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/constants"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
	"github.com/muhwyndhamhp/marknotes/utils/rss"
	"github.com/muhwyndhamhp/marknotes/utils/sanitizations"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/strman"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

func (fe *handler) ArticlesEdit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	post, err := fe.App.PostRepository.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	opts := variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs(fmt.Sprintf("dashboard/articles/%d/edit", id)),
	}

	baseURL := config.Get(config.BASE_URL)
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, id)

	vm := articles.NewArticleViewModel{
		Opts:      opts,
		UploadURL: uploadURL,
		Post:      post,
		BaseURL:   baseURL,
	}
	articlesNew := articles.NewArticle(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}

func (fe *handler) ArticlesNew(c echo.Context) error {
	opts := variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles/new"),
	}

	baseURL := config.Get(config.BASE_URL)
	uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, 0)

	vm := articles.NewArticleViewModel{
		Opts:      opts,
		UploadURL: uploadURL,
		BaseURL:   baseURL,
	}
	articlesNew := articles.NewArticle(vm)

	return templates.AssertRender(c, http.StatusOK, articlesNew)
}

func (fe *handler) Articles(c echo.Context) error {
	ctx := c.Request().Context()
	page, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))
	source := c.QueryParam(constants.TARGET_SOURCE)

	hx_request, _ := strconv.ParseBool(c.Request().Header.Get("Hx-Request"))

	partial := source == constants.TARGET_SOURCE_PARTIAL && hx_request

	count := fe.App.PostRepository.Count(ctx)

	posts, err := fe.App.PostRepository.Get(ctx,
		scopes.Paginate(page, pageSize),
		scopes.OrderBy("created_at", scopes.Descending),
		scopes.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}

	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, internal.PostStatusNone, "", "")
	}
	if len(posts) <= 0 && page > 1 {
		appendRoute := ""
		if source == constants.TARGET_SOURCE_PARTIAL {
			appendRoute = "&source=source-partial"
		}
		path := fmt.Sprintf("/dashboard/articles?page=%d&pageSize=%d%s", page-1, pageSize, appendRoute)
		fmt.Println(path)
		return c.Redirect(http.StatusFound, path)
	}

	opts := variables.DashboardOpts{
		Nav:         nav(0),
		BreadCrumbs: fe.Breadcrumbs("dashboard/articles"),
	}

	pageSizes := fe.SizeDropdown(page, pageSize)
	pages := fe.PageDropdown(tern.Int(page, 1), tern.Int(pageSize, 10), count)
	articleVM := articles.ArticlesViewModel{
		Opts:       opts,
		Posts:      posts,
		PageSizes:  pageSizes,
		Pages:      pages,
		CreatePath: "/dashboard/articles/new",
	}

	dashboard := articles.Articles(articleVM)

	if !partial {
		return templates.AssertRender(c, http.StatusOK, dashboard)
	} else {
		articles := articles.ArticleOOB(posts, pageSizes, pages)
		return templates.AssertRender(c, http.StatusOK, articles)
	}
}

func (fe *handler) ArticlesPush(c echo.Context) error {
	ctx := c.Request().Context()

	status := internal.PostStatus(c.QueryParam("status"))
	if status == internal.PostStatusNone {
		status = internal.PostStatusDraft
	}

	existingID, _ := strconv.Atoi(c.QueryParam("existingID"))
	var xp *internal.Post
	if existingID != 0 {
		p, err := fe.App.PostRepository.GetByID(ctx, uint(existingID))
		if err != nil {
			return errs.Wrap(err)
		}
		xp = p
	}

	var req ArticlesCreateRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		fail := common.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	if err := c.Validate(req); err != nil {
		fmt.Println(err)
		fail := common.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	sanitizedHTML := sanitizations.SanitizeHtml(req.Content)

	slug, err := strman.GenerateSlug(req.Title)
	if err != nil {
		fmt.Println(err)
		fail := common.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	count := fe.App.PostRepository.Count(ctx, scopes.Where("slug = ?", slug))
	if count > 0 && xp == nil {
		fmt.Println(count)
		slug = fmt.Sprintf("%s-%d", slug, count)
	}

	post := tern.Struct(xp, &internal.Post{})

	usr, err := fe.App.OpenAuth.GetUserFromSession(c)
	if err != nil {
		fmt.Println(err)
		fail := common.AlertFailure("Failed to save post:", err.Error())
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
	if status == internal.PostStatusPublished && post.PublishedAt.IsZero() {
		now := time.Now()
		post.PublishedAt = now
	}

	tags := strings.Split(req.Tags, ",")
	post.Tags = []*internal.Tag{}

	tagLiteral := ""
	for i := range tags {
		if tags[i] == "" {
			continue
		}
		tagName := strings.ToLower(strings.TrimSpace(tags[i]))
		tagSlug := strings.ReplaceAll(tagName, " ", "-")
		res, _ := fe.App.TagRepository.Get(ctx, scopes.Where("slug = ?", tagSlug))
		var vTag internal.Tag
		if len(res) != 0 {
			vTag = res[0]
		} else {
			tag := internal.Tag{
				Slug:  tagSlug,
				Title: tags[i],
			}
			err := fe.App.TagRepository.Upsert(c.Request().Context(), &tag)
			if err != nil {
				return errs.Wrap(err)
			}
			vTag = tag
		}

		post.Tags = append(post.Tags, &vTag)
		tagLiteral += vTag.Title + ","
	}
	post.TagsLiteral = tagLiteral

	err = fe.App.PostRepository.Upsert(ctx, post)
	if err != nil {
		fail := common.AlertFailure("Failed to save post:", err.Error())
		return templates.AssertRender(c, http.StatusOK, fail)
	}

	if status == internal.PostStatusPublished {
		go func() {
			ctx := context.Background()
			fe.App.RenderClient.RenderPost(ctx, post)
			err := fe.App.RenderClient.RenderMarkdown(post)
			if err != nil {
				fmt.Println(err)
			}

			err = rss.GenerateRSS(ctx, fe.App.PostRepository)
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
		success := common.AlertSuccess("Post has been saved successfully")

		return templates.AssertRender(c, http.StatusOK, success)
	} else {
		opts := variables.DashboardOpts{
			Nav:         nav(0),
			BreadCrumbs: fe.Breadcrumbs(fmt.Sprintf("dashboard/articles/%d/edit", post.ID)),
		}

		baseURL := config.Get(config.BASE_URL)
		uploadURL := fmt.Sprintf("%s/posts/%d/media/upload", baseURL, post.ID)
		vm := articles.NewArticleViewModel{
			Opts:      opts,
			UploadURL: uploadURL,
			Post:      post,
			BaseURL:   baseURL,
		}
		articlesNew := articles.NewArticle(vm)

		c.Response().Header().Set("HX-Replace-Url", fmt.Sprintf("/dashboard/articles/%d/edit", post.ID))
		return templates.AssertRender(c, http.StatusOK, articlesNew)
	}
}

func (fe *handler) ArticlesDelete(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))

	err := fe.App.PostRepository.Delete(ctx, uint(id))
	if err != nil {
		return err
	}

	go func() {
		err := rss.GenerateRSS(ctx, fe.App.PostRepository)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return templates.RenderEmpty(c)
}
