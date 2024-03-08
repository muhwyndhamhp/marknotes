package renderfile

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/muhwyndhamhp/marknotes/pkg/admin"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	pub_post_detail "github.com/muhwyndhamhp/marknotes/pub/pages/post_detail/post_detail"
	pub_variables "github.com/muhwyndhamhp/marknotes/pub/variables"
	"github.com/muhwyndhamhp/marknotes/template"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
)

func RenderPost(ctx context.Context, post *models.Post) {
	userID := uint(0)

	post.FormMeta = map[string]interface{}{
		"UserID": userID,
	}
	bodyOpts := pub_variables.BodyOpts{
		HeaderButtons: admin.AppendHeaderButtons(userID),
		FooterButtons: admin.AppendFooterButtons(userID),
		Component:     nil,
	}

	js, _ := json.MarshalIndent(post, "", "   ")
	fmt.Println(string(js))

	postDetail := pub_post_detail.PostDetail(bodyOpts, *post)

	// check if public/articles path exists
	if _, err := os.Stat("public/articles"); os.IsNotExist(err) {
		err := os.Mkdir("public/articles", 0o755)
		if err != nil {
			fmt.Println(err)
		}
	}

	err := template.RenderPost(postDetail, "public/articles", post.Slug, post.ID)
	if err != nil {
		fmt.Println(err)
	}
}

func RenderPosts(ctx context.Context, repo models.PostRepository) {
	err := DeletAllFiles("public/articles")
	if err != nil {
		fmt.Println(err)
	}

	posts, err := repo.Get(ctx, scopes.Where("status = ?", values.Published))
	if err != nil {
		fmt.Println(err)
	}

	for _, post := range posts {
		RenderPost(ctx, &post)
	}
}

func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DeletAllFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := os.Remove(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return nil
		}
	}

	return nil
}
