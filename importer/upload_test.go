package importer_test

import (
	"context"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/maxcnunes/httpfake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNoOpUploadService(t *testing.T) {
	fakeHttp := httpfake.New()
	fakeHttp.NewHandler().Post("/v1/upload").Reply(200)
	url := fakeHttp.ResolveURL("")
	defer fakeHttp.Close()
	
	testContent := "this is some dummy content"
	size := int64(len(testContent))

	svc, err := importer.NewApiUploadService(url, size/4)
	assert.Nil(t, err)
	assert.NotNil(t, svc)

	f := &importer.File{
		Name:        "testing.css",
		ReadCloser:  ioutil.NopCloser(strings.NewReader(testContent)),
		SizeInBytes: &size,
	}

	err = svc.Send(context.TODO(), f, "title", "collectionId", "filename", "licence", "licenceUrl")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(svc.Uploads))
}
