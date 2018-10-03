package urlshortener

import (
	"net/http"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/woraphol-j/url-shortener/pkg/mongo"

	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestEndpoints(t *testing.T) {
	urlDAO := mongo.NewDAO("", "", "")
	service := NewService(urlDAO)
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	httpLogger := log.With(logger, "component", "http")
	handler := MakeHandler(service, httpLogger)

	httpExpect := httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://localhost",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	shortURLRequest := map[string]interface{}{
		"url": "www.google.com",
	}

	httpExpect.POST("/shorturls").WithJSON(shortURLRequest).
		Expect().
		Status(http.StatusOK)

}
