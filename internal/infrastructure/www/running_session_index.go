package www

import (
	"net/http"

	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/sport/internal/application"
)

func RunningSessionsIndex(usecase application.Application) web.HandlerFunc {
	return func(ctx web.Context, w http.ResponseWriter, r *http.Request) web.Response {
		activities, err := usecase.ListRunningSessions(ctx.StdCtx())
		if err != nil {
			return ctx.InternalServerErrorResponse("can't list activities: %v", err)
		}

		return ctx.Response(200, "templates/running-sessions/index.html.tmpl", map[string]interface{}{
			"Activities": activities,
		})
	}
}
