package renderfile

import (
	"context"
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/pub"
	pub_post_detail "github.com/muhwyndhamhp/marknotes/pub/pages/post_detail/post_detail"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/fileman"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

type RenderClient struct {
	App *internal.Application
}

func NewRenderClient(app *internal.Application) *RenderClient {
	return &RenderClient{App: app}
}

func (r *RenderClient) RenderPost(ctx context.Context, post *internal.Post) {
	userID := uint(0)

	post.FormMeta = map[string]interface{}{
		"UserID": userID,
	}

	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]
	canonURL := fmt.Sprintf("%s/articles/%s.html", baseURL, post.Slug)

	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		Component:     nil,
		ExtraHeaders: []templ.Component{
			pub.CannonicalRel(canonURL),
		},
	}

	postDetail := pub_post_detail.PostDetail(bodyOpts, *post)

	if err := fileman.CheckDir(config.Get(config.POST_RENDER_PATH) + ""); err != nil {
		fmt.Println(err)
	}

	file, err := template.RenderPost(postDetail, config.Get(config.POST_RENDER_PATH)+"", post.Slug, post.ID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(file.Name())

	// if config.Get(config.ENV) != "dev" {

	prefix := ""
	if config.Get(config.ENV) != "dev" {
		prefix = "/store/"
	}

	_, err = r.App.Bucket.UploadStatic(ctx, file.Name(), prefix, "text/html")
	if err != nil {
		fmt.Println(err)
	}
	// }
}

func (r *RenderClient) RenderPosts(ctx context.Context) {
	// check last_render.txt, read the content as time format RFC3339.
	// If more than 6 hours, then continue
	// if less than 6 hours, then return
	lastRender, _ := os.ReadFile(config.Get(config.POST_RENDER_PATH) + "/last_render.txt")

	lastRenderTime, err := time.Parse(time.RFC3339, string(lastRender))
	if err != nil {
		fmt.Println(err)
	}

	if time.Since(lastRenderTime).Hours() < 1 && config.Get(config.ENV) != "dev" {
		return
	}

	err = fileman.DeleteAllFiles(config.Get(config.POST_RENDER_PATH) + "")
	if err != nil {
		fmt.Println(err)
	}

	posts, err := r.App.PostRepository.Get(ctx, scopes.Where("status = ?", values.Published), scopes.Preload("Tags"))
	if err != nil {
		fmt.Println(err)
	}

	for _, post := range posts {
		r.RenderPost(ctx, &post)
	}

	// write current time to last_render.txt
	err = os.WriteFile(
		config.Get(config.POST_RENDER_PATH)+"/last_render.txt",
		[]byte(time.Now().Format(time.RFC3339)),
		0o755,
	)
	if err != nil {
		fmt.Println(err)
	}
}
