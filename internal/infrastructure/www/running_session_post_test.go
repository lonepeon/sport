package www_test

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/testutils/gomockutils"
	"github.com/lonepeon/golib/web/webtest"
	"github.com/lonepeon/sport/internal/infrastructure/job"
	"github.com/lonepeon/sport/internal/infrastructure/job/jobtest"
	"github.com/lonepeon/sport/internal/infrastructure/www"
)

func TestRunningSessionPostInvalidDate(t *testing.T) {
	tcs := map[string]string{
		"randomString":           "not a date",
		"invalidFormat":          "1st of February 2022 21:21:12",
		"goodFormatInvalidValue": "2022-42-20T21:12",
		"fomatIncludesSeconds":   "2022-11-20T21:12:45",
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := webtest.NewMockContext(ctrl)
			w := httptest.NewRecorder()

			var body bytes.Buffer
			bodyWriter := multipart.NewWriter(&body)
			testutils.AssertNoError(t, bodyWriter.WriteField("date", tc), "can't write date to form")
			bodyWriter.Close()

			r := httptest.NewRequest("POST", "/running-session/", &body)
			r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

			ctx.EXPECT().AddFlash(gomock.All(
				webtest.MatchFlashErrorContains("date format"),
				webtest.MatchFlashErrorContains("2006-01-02T15:04"),
			))

			expectedResponse := webtest.MockedResponse("redirection")
			ctx.EXPECT().Redirect(w, 303, "/running-session/new").Return(expectedResponse)

			response := www.RunningSessionPost(nil, "")(ctx, w, r)

			webtest.AssertResponse(t, expectedResponse, response, "unexpected response")
			testutils.AssertContainsString(t, "date format", response.LogMessage, "unexpected log message")
			testutils.AssertContainsString(t, tc, response.LogMessage, "unexpected log message")
		})
	}
}

func TestRunningSessionPostMissingGPXFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()

	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)
	testutils.AssertNoError(t, bodyWriter.WriteField("date", "2022-02-20T21:27"), "can't write date to form")
	bodyWriter.Close()

	r := httptest.NewRequest("POST", "/running-session/", &body)
	r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	ctx.EXPECT().AddFlash(webtest.MatchFlashErrorContains("gpx file must be sent"))

	expectedResponse := webtest.MockedResponse("redirection")
	ctx.EXPECT().Redirect(w, 303, "/running-session/new").Return(expectedResponse)

	response := www.RunningSessionPost(nil, "")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected response")
	testutils.AssertContainsString(t, "can't get gpx", response.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "no such file", response.LogMessage, "unexpected log message")
}

func TestRunningSessionPostHugeGPXFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()

	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)
	testutils.AssertNoError(t, bodyWriter.WriteField("date", "2022-02-20T21:27"), "can't write date to form")
	gpxFile, err := bodyWriter.CreateFormFile("gpx", "my-huge-file.gpx")
	testutils.AssertNoError(t, err, "can't create form file")
	for i := 0; i < 1024*1024; i++ {
		fmt.Fprintf(gpxFile, "6 char")
	}
	bodyWriter.Close()

	r := httptest.NewRequest("POST", "/running-session/", &body)
	r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	ctx.EXPECT().AddFlash(gomock.All(
		webtest.MatchFlashErrorContains("gpx file"),
		webtest.MatchFlashErrorContains("too big"),
	))

	expectedResponse := webtest.MockedResponse("redirection")
	ctx.EXPECT().Redirect(w, 303, "/running-session/new").Return(expectedResponse)

	response := www.RunningSessionPost(nil, "")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected response")
	testutils.AssertContainsString(t, "too big", response.LogMessage, "unexpected log message")
}

func TestRunningSessionPostInvalidUploadFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()

	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)
	testutils.AssertNoError(t, bodyWriter.WriteField("date", "2022-02-20T21:27"), "can't write date to  form")
	gpxFile, err := bodyWriter.CreateFormFile("gpx", "my-huge-file.gpx")
	testutils.AssertNoError(t, err, "can't create form file")
	for i := 0; i < 1024; i++ {
		fmt.Fprintf(gpxFile, "6 char")
	}
	bodyWriter.Close()

	r := httptest.NewRequest("POST", "/running-session/", &body)
	r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	expectedResponse := webtest.MockedResponse("server error")
	ctx.EXPECT().InternalServerErrorResponse(gomock.All(
		gomockutils.ContainsString("temporary gpx"),
		gomockutils.ContainsString("/an/invalid/path/on/the/system"),
	)).Return(expectedResponse)

	response := www.RunningSessionPost(nil, "/an/invalid/path/on/the/system")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected response")
}

func TestRunningSessionPostCantEnqueueJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	enqueuer := jobtest.NewMockEnqueuer(ctrl)
	uploadFolder, err := os.MkdirTemp("", "test-cant-enqueue")
	testutils.AssertNoError(t, err, "can't create temp folder")
	defer os.RemoveAll(uploadFolder)

	when := "2022-02-20T21:27"
	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)
	testutils.AssertNoError(t, bodyWriter.WriteField("date", when), "can't write date to form")
	gpxFile, err := bodyWriter.CreateFormFile("gpx", "my-huge-file.gpx")
	testutils.AssertNoError(t, err, "can't create form file")
	for i := 0; i < 1024; i++ {
		fmt.Fprintf(gpxFile, "6 char")
	}
	bodyWriter.Close()

	r := httptest.NewRequest("POST", "/running-session/", &body)
	r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	enqueuer.EXPECT().Enqueue(gomock.Any()).Return(errors.New("boom"))

	expectedResponse := webtest.MockedResponse("server error")
	ctx.EXPECT().InternalServerErrorResponse(
		gomockutils.ContainsString("can't enqueue"),
		gomock.Any(),
	).Return(expectedResponse)

	response := www.RunningSessionPost(enqueuer, uploadFolder)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected response")
}

func TestRunningSessionPostSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	enqueuer := jobtest.NewMockEnqueuer(ctrl)
	uploadFolder, err := os.MkdirTemp("", "test-cant-enqueue")
	testutils.AssertNoError(t, err, "can't create temp folder")
	defer os.RemoveAll(uploadFolder)

	when := "2022-02-20T21:27"
	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)
	testutils.AssertNoError(t, bodyWriter.WriteField("date", when), "can't write date to form")
	gpxFile, err := bodyWriter.CreateFormFile("gpx", "my-huge-file.gpx")
	testutils.AssertNoError(t, err, "can't create form file")
	for i := 0; i < 1024; i++ {
		fmt.Fprintf(gpxFile, "6 char")
	}
	bodyWriter.Close()

	r := httptest.NewRequest("POST", "/running-session/", &body)
	r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	enqueuer.EXPECT().Enqueue(jobtest.NewJobMatcher(
		"track-running-session-job",
		&job.TrackRunningSessionJobInput{},
		func(arg interface{}) bool {
			input := arg.(*job.TrackRunningSessionJobInput)

			return strings.HasPrefix(input.GPXFilepath, uploadFolder) &&
				input.When.Format("2006-01-02T15:04") == when
		},
	)).Return(nil)

	ctx.EXPECT().AddFlash(webtest.MatchFlashSuccessContains("being processed"))

	expectedResponse := webtest.MockedResponse("server error")
	ctx.EXPECT().Redirect(w, 303, "/").Return(expectedResponse)

	response := www.RunningSessionPost(enqueuer, uploadFolder)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected response")
}
