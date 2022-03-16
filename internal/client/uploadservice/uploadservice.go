package uploadservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/v2/log"
)

const service = "upload-service"

// Client is an import api client which can be used to make requests to the API
type Client struct {
	cli dphttp.Clienter
	url string
}

// New creates new instance of Client with a give import api url
func New(url string) *Client {
	hcClient := healthcheck.NewClient(service, url)

	return &Client{
		cli: hcClient.Client,
		url: url,
	}
}

// ErrInvalidAPIResponse is returned when the api does not respond with a valid status
type ErrInvalidAPIResponse struct {
	actualCode int
	uri        string
	body       string
}

// Error should be called by the user to print out the stringified version of the error
func (e ErrInvalidAPIResponse) Error() string {
	return fmt.Sprintf(
		"invalid response: %d from %s: %s, body: %s",
		e.actualCode,
		service,
		e.uri,
		e.body,
	)
}

// Code returns the status code received from the api if an error is returned
func (e ErrInvalidAPIResponse) Code() int {
	return e.actualCode
}

var _ error = ErrInvalidAPIResponse{}

type UploadJob struct {
	Path                 string
	ResumableFilename    string
	IsPublishable        bool
	CollectionId         string
	Title                string
	ResumableTotalSize   int32
	ResumableType        string
	Licence              string
	LicenceUrl           string
	ResumableChunkNumber int32
	ResumableTotalChunks int32
	File                 *os.File
}

// Checker calls api health endpoint and returns a check object to the caller.
func (c *Client) Checker(ctx context.Context, check *health.CheckState) error {
	hcClient := healthcheck.Client{
		Client: c.cli,
		URL:    c.url,
		Name:   service,
	}

	return hcClient.Checker(ctx, check)
}

// Upload POST's form with file to be uploaded and metadata
func (c *Client) Upload(ctx context.Context, serviceToken string, job UploadJob) error {
	uri := fmt.Sprintf("%s/%s", c.url, "upload-new")
	logData := log.Data{"uri": uri}

	resp, err := c.doPostForm(ctx, serviceToken, uri, job)
	if err != nil {
		return err
	}
	defer closeResponseBody(ctx, resp)

	jsonBody, err := getBody(resp)
	if err != nil {
		log.Error(ctx, "Failed to read body from API", err)
		return err
	}

	logData["filename"] = job.ResumableFilename
	logData["httpCode"] = resp.StatusCode
	logData["jsonBody"] = string(jsonBody)

	if resp.StatusCode != http.StatusOK {
		return NewErrorResponse(resp, uri)
	}

	return nil
}

func (c *Client) doPostForm(ctx context.Context, serviceToken, uri string, job UploadJob) (*http.Response, error) {
	logData := log.Data{"uri": uri}

	URL, err := url.Parse(uri)
	if err != nil {
		log.Error(ctx, "Failed to create url for API call", err, logData)
		return nil, err
	}
	uri = URL.String()
	logData["url"] = uri

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("path", job.Path)
	_ = writer.WriteField("resumableFilename", job.ResumableFilename)
	_ = writer.WriteField("isPublishable", fmt.Sprintf("%v", job.IsPublishable))
	_ = writer.WriteField("collectionId", job.CollectionId)
	_ = writer.WriteField("title", job.Title)
	_ = writer.WriteField("resumableTotalSize", fmt.Sprintf("%d", job.ResumableTotalSize))
	_ = writer.WriteField("resumableType", job.ResumableType)
	_ = writer.WriteField("licence", job.Licence)
	_ = writer.WriteField("licenceUrl", job.LicenceUrl)
	_ = writer.WriteField("resumableChunkNumber", fmt.Sprintf("%d", job.ResumableChunkNumber))
	_ = writer.WriteField("resumableTotalChunks", fmt.Sprintf("%d", job.ResumableTotalChunks))

	f, err := os.Open(job.File.Name())
	defer f.Close()
	fileWriter, err := writer.CreateFormFile("file", job.File.Name())
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, payload)
	if err != nil {
		return nil, err
	}

	// add a service token to request where one has been provided
	dprequest.AddServiceTokenHeader(req, serviceToken)

	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.cli.Do(ctx, req)
	if err != nil {
		log.Error(ctx, "Failed to action API", err, logData)
		return nil, err
	}

	return resp, nil
}

// NewErrorResponse creates an error response, optionally adding body to e when status is 404
func NewErrorResponse(resp *http.Response, uri string) (e *ErrInvalidAPIResponse) {
	e = &ErrInvalidAPIResponse{
		actualCode: resp.StatusCode,
		uri:        uri,
	}
	if resp.StatusCode == http.StatusNotFound {
		body, err := getBody(resp)
		if err != nil {
			e.body = "Client failed to read response body"
			return
		}
		e.body = string(body)
	}
	return
}

func getBody(resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(ctx context.Context, resp *http.Response) {
	if resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			log.Error(ctx, "error closing http response body", err)
		}
	}
}
