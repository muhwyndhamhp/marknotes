package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/admin"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/common/variables"
	"github.com/muhwyndhamhp/marknotes/internal/handler/http/post/articles"
	"github.com/muhwyndhamhp/marknotes/internal/middlewares"

	"github.com/labstack/echo/v4"
	templateRenderer "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/params"
	"github.com/muhwyndhamhp/marknotes/utils/resp"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"github.com/muhwyndhamhp/marknotes/utils/storage"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

type handler struct {
	app *internal.Application
}

func NewHandler(g *echo.Group, app *internal.Application) {
	fe := &handler{app: app}

	g.GET("/posts", fe.PostsGet)
	g.GET("/posts_index", fe.PostsIndex)

	g.GET("/posts/:id", fe.GetPostByID, app.GetIdParamWare)

	// Upload and Download Media
	g.POST("/posts/:id/media/upload", fe.PostMediaUpload)
	g.GET("/posts/media/:filename", fe.PostMediaGet)

	// Alias `articles` for `posts`
	g.GET("/articles", fe.PostsIndex, app.CacheControlWare)
}

func (fe *handler) PostMediaGet(c echo.Context) error {
	filename := c.Param("filename")
	path, err := storage.ServeFile(filename)
	if err != nil {
		return c.String(http.StatusNotFound, "File not found")
	}

	return c.File(path)
}

func (fe *handler) PostMediaUpload(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request().Context()

	imgSize := c.QueryParam("size")
	size, _ := strconv.Atoi(imgSize)

	f, err := c.FormFile("file")
	if err != nil {
		return err
	}

	ct, valid := storage.IsValidFileType(f)
	if !valid {
		return resp.HTTPBadRequest(c, "NOT_SUPPORTED", "file type not supported")
	}

	url, err := fe.app.Bucket.UploadMedia(ctx, f, fmt.Sprintf("%d", id), ct, size)
	if err != nil {
		return err
	}

	return resp.HTTPOk(c, struct {
		URL string `json:"url"`
	}{
		URL: url,
	})
}

func (fe *handler) PostsIndex(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := fe.app.PostRepository.Get(ctx,
		scopes.Paginate(1, 10),
		scopes.OrderBy("published_at", scopes.Descending),
		internal.WithStatus(internal.PostStatusPublished),
		scopes.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, internal.PostStatusPublished, "published_at", "")
	}

	search := variables.SearchBar{
		SearchPlaceholder: "Search Articles...",
		SearchPath:        "/posts?page=1&pageSize=10&sortBy=published_at&status=published",
	}

	user, err := fe.app.OpenAuth.GetUserFromSession(c)
	if err != nil {
		log.Error(err)
	}

	if user == nil {
		user = &internal.User{}
	}

	bodyOpts := variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(user.ID),
		Component:     nil,
	}

	postIndex := articles.PostIndex(bodyOpts, posts, search)

	return templateRenderer.AssertRender(c, http.StatusOK, postIndex)
}

func (fe *handler) GetPostByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

	post, err := fe.app.PostRepository.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%s.html", post.Slug))
}

func (fe *handler) PostsGet(c echo.Context) error {
	ctx := c.Request().Context()
	page, pageSize, sortBy, statusStr, keyword, loadNext := params.GetCommonParams(c)
	status := internal.PostStatus(statusStr)

	scp := []scopes.QueryScope{
		scopes.OrderBy(tern.String(sortBy, "created_at"), scopes.Descending),
		scopes.Paginate(page, pageSize),
	}

	if keyword != "" {
		scp = append(scp, internal.PostDeepMatch(keyword, status))
	} else {
		scp = append(scp, internal.WithStatus(status))
	}

	posts, err := fe.app.PostRepository.Get(ctx, scp...)
	if err != nil {
		return err
	}

	if len(posts) > 0 && loadNext {
		posts[len(posts)-1].AppendFormMeta(page+1, status, sortBy, keyword)
	}

	postList := articles.PostList(posts)

	return templateRenderer.AssertRender(c, http.StatusOK, postList)
}
