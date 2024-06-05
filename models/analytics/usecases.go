package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/muhwyndhamhp/marknotes/analytics"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/typeext"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func CacheAnalytics(ctx context.Context, db *gorm.DB, c *analytics.Client) error {
	var slugs []string
	err := db.WithContext(ctx).
		Model(&models.Post{}).
		Where("status = ?", values.Published).
		Pluck("slug", &slugs).
		Error
	if err != nil {
		return errs.Wrap(err)
	}

	eg, egctx := errgroup.WithContext(ctx)

	for i := range slugs {
		eg.Go(func() error {
			path := fmt.Sprintf("/articles/%s.html", slugs[i])
			a, err := c.GetAnalytics(ctx, slugs[i])
			if err != nil {
				return errs.Wrap(err)
			}

			bin, err := typeext.StructToJSONB(a)
			if err != nil {
				return errs.Wrap(err)
			}
			err = db.WithContext(egctx).Save(&models.Analytics{
				CaptureDate: time.Now(),
				Path:        path,
				Data:        bin,
			}).Error
			if err != nil {
				return errs.Wrap(err)
			}

			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}
