package importer_test

import (
	"context"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/internal/client/uploadservice"
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

	client := uploadservice.New(url)

	svc := importer.NewUploadService(client, 1024)
	assert.NotNil(t, svc)

	f := &importer.File{
		Name:        "testing.css",
		ReadCloser:  ioutil.NopCloser(strings.NewReader("this is some dummy content")),
		SizeInBytes: &size,
	}

	err := svc.SendFile(context.TODO(), f, "title", "collectionId", "licence", "licenceUrl")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(svc.Uploads))
}
