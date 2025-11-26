package restapi_test

import (
	"bytes"
	_ "embed"
	"net/http"
	"os"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rlibaert/service-example-go/restapi"
)

// nopAdapter is a no-op implementation of [huma.Adapter].
type nopAdapter struct{}

func (nopAdapter) Handle(*huma.Operation, func(huma.Context))   {}
func (nopAdapter) ServeHTTP(http.ResponseWriter, *http.Request) {}

//go:embed openapi.yaml
var openapi []byte

func TestOpenAPI(t *testing.T) {
	api := huma.NewAPI(huma.DefaultConfig("test", "dev"), nopAdapter{})
	huma.AutoRegister(api, restapi.ServiceRegisterer{})

	b, err := api.OpenAPI().YAML()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, openapi) {
		f, ferr := os.CreateTemp("", "openapi_*.yaml") //nolint: usetesting // would remove the file after the test
		if ferr != nil {
			t.Error(ferr)
		}
		defer f.Close()

		f.Write(b)
		t.Error(t.Name(), "failed, written expected results to", f.Name())
	}
}
