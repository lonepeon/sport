package www

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/sport/internal/infrastructure/job"
)

const MaxGPXFileSize = 5 * 1024 * 1024

func RunningSessionPost(enqueuer job.Enqueuer, uploadFolder string) web.HandlerFunc {
	return func(ctx web.Context, w http.ResponseWriter, r *http.Request) web.Response {
		if err := r.ParseMultipartForm(MaxGPXFileSize); err != nil {
			ctx.AddFlash(web.NewFlashMessageError("can't parse request parameters. Please try again"))

			response := ctx.Redirect(w, http.StatusSeeOther, "/running-session/new")
			response.LogMessage = fmt.Sprintf("can't parse form: %v", err)
			return response
		}

		date := r.FormValue("date")

		datetimeLayout := "2006-01-02T15:04"
		when, err := time.Parse(datetimeLayout, date)
		if err != nil {
			ctx.AddFlash(web.NewFlashMessageError("date format is expected to follow %s", datetimeLayout))

			response := ctx.Redirect(w, http.StatusSeeOther, "/running-session/new")
			response.LogMessage = fmt.Sprintf("can't parse date format (date=%s): %v", date, err)
			return response
		}

		file, header, err := r.FormFile("gpx")
		if err != nil {
			ctx.AddFlash(web.NewFlashMessageError("gpx file must be sent"))
			response := ctx.Redirect(w, http.StatusSeeOther, "/running-session/new")
			response.LogMessage = fmt.Sprintf("can't get gpx file from http form: %v", err)
			return response
		}
		defer file.Close()

		if header.Size > MaxGPXFileSize {
			ctx.AddFlash(web.NewFlashMessageError("gpx file is too big %f > 5Mb", (float64(header.Size) / 1024)))
			response := ctx.Redirect(w, http.StatusSeeOther, "/running-session/new")
			response.LogMessage = fmt.Sprintf("gpx file is too big (size: %db)", header.Size)
			return response
		}

		filepath := path.Join(uploadFolder, date+".gpx")
		if err := createTemporaryFile(file, filepath); err != nil {
			return ctx.InternalServerErrorResponse(err.Error())
		}

		input := job.TrackRunningSessionJobInput{When: when, GPXFilepath: filepath}
		if err = job.EnqueueTrackRunningSessionJob(enqueuer, input); err != nil {
			return ctx.InternalServerErrorResponse("can't enqueue running session job: %v", err)
		}

		ctx.AddFlash(web.NewFlashMessageSuccess("running session is being processed"))
		return ctx.Redirect(w, http.StatusSeeOther, "/")
	}
}

func createTemporaryFile(f io.Reader, filepath string) error {
	dest, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("can't create temporary gpx file (path=%s): %v", filepath, err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, f); err != nil {
		return fmt.Errorf("can't copy uploaded file to upload folder (path=%s): %v", filepath, err)
	}

	return nil
}
