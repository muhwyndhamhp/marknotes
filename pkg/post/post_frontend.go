package post

import (
	"fmt"
	"github.com/muhwyndhamhp/marknotes/db"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/middlewares"
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	pub_postlist "github.com/muhwyndhamhp/marknotes/pub/components/postlist"
	pub_post_index "github.com/muhwyndhamhp/marknotes/pub/pages/post_index"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	templateRenderer "github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/cloudbucket"
	"github.com/muhwyndhamhp/marknotes/utils/jwt"
	"github.com/muhwyndhamhp/marknotes/utils/params"
	"github.com/muhwyndhamhp/marknotes/utils/resp"
	"github.com/muhwyndhamhp/marknotes/utils/storage"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

type PostFrontend struct {
	repo   models.PostRepository
	bucket *cloudbucket.S3Client
}

func NewPostFrontend(g *echo.Group,
	repo models.PostRepository,
	bucket *cloudbucket.S3Client,
	htmxMid echo.MiddlewareFunc,
	authMid echo.MiddlewareFunc,
	authDescribeMid echo.MiddlewareFunc,
	byIDMiddleware echo.MiddlewareFunc,
	cacheControlMid echo.MiddlewareFunc,
) {
	fe := &PostFrontend{
		repo:   repo,
		bucket: bucket,
	}

	g.GET("/posts", fe.PostsGet, authDescribeMid)
	g.GET("/posts_index", fe.PostsIndex, authDescribeMid)

	g.GET("/posts/:id", fe.GetPostByID, authDescribeMid, byIDMiddleware)

	// Upload and Download Media
	g.POST("/posts/:id/media/upload", fe.PostMediaUpload, authDescribeMid)
	g.GET("/posts/media/:filename", fe.PostMediaGet)

	// Alias `articles` for `posts`
	g.GET("/articles", fe.PostsIndex, authDescribeMid, cacheControlMid)
}

func (fe *PostFrontend) PostMediaGet(c echo.Context) error {
	filename := c.Param("filename")
	path, err := storage.ServeFile(filename)
	if err != nil {
		return c.String(http.StatusNotFound, "File not found")
	}

	return c.File(path)
}

func (fe *PostFrontend) PostMediaUpload(c echo.Context) error {
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

	url, err := fe.bucket.UploadMedia(ctx, f, fmt.Sprintf("%d", id), ct, size)
	if err != nil {
		return err
	}

	return resp.HTTPOk(c, struct {
		URL string `json:"url"`
	}{
		URL: url,
	})
}

func (fe *PostFrontend) PostsIndex(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := fe.repo.Get(ctx,
		db.Paginate(1, 10),
		db.OrderBy("published_at", db.Descending),
		db.WithStatus(values.Published),
		db.PostIndexedOnly(),
	)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		posts[len(posts)-1].AppendFormMeta(2, values.Published, "published_at", "")
	}

	search := pub_variables.SearchBar{
		SearchPlaceholder: "Search Articles...",
		SearchPath:        "/posts?page=1&pageSize=10&sortBy=published_at&status=published",
	}

	userID := jwt.AppendAndReturnUserID(c, map[string]interface{}{})

	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		Component:     nil,
	}

	postIndex := pub_post_index.PostIndex(bodyOpts, posts, search)

	return templateRenderer.AssertRender(c, http.StatusOK, postIndex)
}

func (fe *PostFrontend) GetPostByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := c.Get(middlewares.ByIDKey).(int)

	post, err := fe.repo.GetByID(ctx, uint(id))
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%s.html", post.Slug))
}

func (fe *PostFrontend) PostsGet(c echo.Context) error {
	ctx := c.Request().Context()
	page, pageSize, sortBy, statusStr, keyword, loadNext := params.GetCommonParams(c)
	status := values.PostStatus(statusStr)

	scp := []db.QueryScope{
		db.OrderBy(tern.String(sortBy, "created_at"), db.Descending),
		db.Paginate(page, pageSize),
	}

	if keyword != "" {
		scp = append(scp, db.PostDeepMatch(keyword, status))
	} else {
		scp = append(scp, db.WithStatus(status))
	}

	posts, err := fe.repo.Get(ctx, scp...)
	if err != nil {
		return err
	}

	if len(posts) > 0 && loadNext {
		posts[len(posts)-1].AppendFormMeta(page+1, status, sortBy, keyword)
	}

	postList := pub_postlist.PostList(posts)

	return templateRenderer.AssertRender(c, http.StatusOK, postList)
}
