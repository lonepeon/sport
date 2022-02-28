package www

import (
	"fmt"
	"net/http"

	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/sport/internal/application"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/infrastructure/job"
)

func RunningSessionsDelete(app application.Application, enqueuer job.Enqueuer) web.HandlerFunc {
	return func(ctx web.Context, w http.ResponseWriter, r *http.Request) web.Response {
		vars := ctx.Vars(r)

		slug, err := domain.NewRunnningActivitySlugFromString(vars["slug"])
		if err != nil {
			ctx.AddFlash(web.NewFlashMessageError("no activity recorded with slug '%v'", vars["slug"]))
			redirection := ctx.Redirect(w, http.StatusSeeOther, "/")
			redirection.LogMessage = fmt.Sprintf("can't parse activity time: %v", err)
			return redirection
		}

		_, err = app.GetRunningSession(ctx.StdCtx(), slug)
		if err != nil {
			ctx.AddFlash(web.NewFlashMessageError("no activity recorded with slug '%v'", vars["slug"]))
			redirection := ctx.Redirect(w, http.StatusSeeOther, "/")
			redirection.LogMessage = fmt.Sprintf("can't find activity (slug=%v): %v", vars["slug"], err)
			return redirection
		}

		input := job.DeleteRunningSessionJobInput{Slug: slug.String()}
		if err := job.EnqueueDeleteRunningSessionJob(enqueuer, input); err != nil {
			return ctx.InternalServerErrorResponse("can't enqueue running session deletion job: %v", err)
		}

		ctx.AddFlash(web.NewFlashMessageSuccess("activity recorded with slug '%s' is being deleted", vars["slug"]))
		return ctx.Redirect(w, http.StatusSeeOther, "/")
	}
}
