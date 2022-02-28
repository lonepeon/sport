package www

import (
	"errors"
	"net/http"

	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/sport/internal/application"
	"github.com/lonepeon/sport/internal/domain"
)

func RunningSessionsShow(app application.Application) web.HandlerFunc {
	return func(ctx web.Context, w http.ResponseWriter, r *http.Request) web.Response {
		vars := ctx.Vars(r)

		when, err := domain.NewRunnningActivitySlugFromString(vars["slug"])
		if err != nil {
			return ctx.NotFoundResponse("can't parse activity slug (slug=%s): %v", vars["slug"], err)
		}

		activity, err := app.GetRunningSession(ctx.StdCtx(), when)
		if err != nil {
			if errors.Is(err, domain.ErrCantGetRunningSession) {
				return ctx.NotFoundResponse("can't find activity (slug=%s): %v", vars["slug"], err)
			}
			return ctx.InternalServerErrorResponse("failed while finding activity (slug=%s): %v", vars["slug"], err)
		}

		return ctx.Response(200, "templates/running-sessions/show.html.tmpl", map[string]interface{}{
			"Activity": activity,
		})
	}
}
