// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks_importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"io"
	"sync"
)

// Ensure, that UploadServiceBackendMock does implement importer.UploadServiceBackend.
// If this is not the case, regenerate this file with moq.
var _ importer.UploadServiceBackend = &UploadServiceBackendMock{}

// UploadServiceBackendMock is a mock implementation of importer.UploadServiceBackend.
//
// 	func TestSomethingThatUsesUploadServiceBackend(t *testing.T) {
//
// 		// make and configure a mocked importer.UploadServiceBackend
// 		mockedUploadServiceBackend := &UploadServiceBackendMock{
// 			UploadFunc: func(ctx context.Context, fileContent io.ReadCloser, metadata upload.Metadata) error {
// 				panic("mock out the Upload method")
// 			},
// 		}
//
// 		// use mockedUploadServiceBackend in code that requires importer.UploadServiceBackend
// 		// and then make assertions.
//
// 	}
type UploadServiceBackendMock struct {
	// UploadFunc mocks the Upload method.
	UploadFunc func(ctx context.Context, fileContent io.ReadCloser, metadata upload.Metadata) error

	// calls tracks calls to the methods.
	calls struct {
		// Upload holds details about calls to the Upload method.
		Upload []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FileContent is the fileContent argument value.
			FileContent io.ReadCloser
			// Metadata is the metadata argument value.
			Metadata upload.Metadata
		}
	}
	lockUpload sync.RWMutex
}

// Upload calls UploadFunc.
func (mock *UploadServiceBackendMock) Upload(ctx context.Context, fileContent io.ReadCloser, metadata upload.Metadata) error {
	if mock.UploadFunc == nil {
		panic("UploadServiceBackendMock.UploadFunc: method is nil but UploadServiceBackend.Upload was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		FileContent io.ReadCloser
		Metadata    upload.Metadata
	}{
		Ctx:         ctx,
		FileContent: fileContent,
		Metadata:    metadata,
	}
	mock.lockUpload.Lock()
	mock.calls.Upload = append(mock.calls.Upload, callInfo)
	mock.lockUpload.Unlock()
	return mock.UploadFunc(ctx, fileContent, metadata)
}

// UploadCalls gets all the calls that were made to Upload.
// Check the length with:
//     len(mockedUploadServiceBackend.UploadCalls())
func (mock *UploadServiceBackendMock) UploadCalls() []struct {
	Ctx         context.Context
	FileContent io.ReadCloser
	Metadata    upload.Metadata
} {
	var calls []struct {
		Ctx         context.Context
		FileContent io.ReadCloser
		Metadata    upload.Metadata
	}
	mock.lockUpload.RLock()
	calls = mock.calls.Upload
	mock.lockUpload.RUnlock()
	return calls
}
