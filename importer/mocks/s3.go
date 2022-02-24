// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks_importer

import (
	"context"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	s3client "github.com/ONSdigital/dp-s3"
	"io"
	"sync"
)

// Ensure, that S3InterfaceMock does implement importer.S3Interface.
// If this is not the case, regenerate this file with moq.
var _ importer.S3Interface = &S3InterfaceMock{}

// S3InterfaceMock is a mock implementation of importer.S3Interface.
//
// 	func TestSomethingThatUsesS3Interface(t *testing.T) {
//
// 		// make and configure a mocked importer.S3Interface
// 		mockedS3Interface := &S3InterfaceMock{
// 			CheckPartUploadedFunc: func(ctx context.Context, req *s3client.UploadPartRequest) (bool, error) {
// 				panic("mock out the CheckPartUploaded method")
// 			},
// 			CheckerFunc: func(ctx context.Context, state *health.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			GetFunc: func(key string) (io.ReadCloser, *int64, error) {
// 				panic("mock out the Get method")
// 			},
// 			UploadPartFunc: func(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error {
// 				panic("mock out the UploadPart method")
// 			},
// 		}
//
// 		// use mockedS3Interface in code that requires importer.S3Interface
// 		// and then make assertions.
//
// 	}
type S3InterfaceMock struct {
	// CheckPartUploadedFunc mocks the CheckPartUploaded method.
	CheckPartUploadedFunc func(ctx context.Context, req *s3client.UploadPartRequest) (bool, error)

	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, state *health.CheckState) error

	// GetFunc mocks the Get method.
	GetFunc func(key string) (io.ReadCloser, *int64, error)

	// UploadPartFunc mocks the UploadPart method.
	UploadPartFunc func(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error

	// calls tracks calls to the methods.
	calls struct {
		// CheckPartUploaded holds details about calls to the CheckPartUploaded method.
		CheckPartUploaded []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Req is the req argument value.
			Req *s3client.UploadPartRequest
		}
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// State is the state argument value.
			State *health.CheckState
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Key is the key argument value.
			Key string
		}
		// UploadPart holds details about calls to the UploadPart method.
		UploadPart []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Req is the req argument value.
			Req *s3client.UploadPartRequest
			// Payload is the payload argument value.
			Payload []byte
		}
	}
	lockCheckPartUploaded sync.RWMutex
	lockChecker           sync.RWMutex
	lockGet               sync.RWMutex
	lockUploadPart        sync.RWMutex
}

// CheckPartUploaded calls CheckPartUploadedFunc.
func (mock *S3InterfaceMock) CheckPartUploaded(ctx context.Context, req *s3client.UploadPartRequest) (bool, error) {
	if mock.CheckPartUploadedFunc == nil {
		panic("S3InterfaceMock.CheckPartUploadedFunc: method is nil but S3Interface.CheckPartUploaded was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Req *s3client.UploadPartRequest
	}{
		Ctx: ctx,
		Req: req,
	}
	mock.lockCheckPartUploaded.Lock()
	mock.calls.CheckPartUploaded = append(mock.calls.CheckPartUploaded, callInfo)
	mock.lockCheckPartUploaded.Unlock()
	return mock.CheckPartUploadedFunc(ctx, req)
}

// CheckPartUploadedCalls gets all the calls that were made to CheckPartUploaded.
// Check the length with:
//     len(mockedS3Interface.CheckPartUploadedCalls())
func (mock *S3InterfaceMock) CheckPartUploadedCalls() []struct {
	Ctx context.Context
	Req *s3client.UploadPartRequest
} {
	var calls []struct {
		Ctx context.Context
		Req *s3client.UploadPartRequest
	}
	mock.lockCheckPartUploaded.RLock()
	calls = mock.calls.CheckPartUploaded
	mock.lockCheckPartUploaded.RUnlock()
	return calls
}

// Checker calls CheckerFunc.
func (mock *S3InterfaceMock) Checker(ctx context.Context, state *health.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("S3InterfaceMock.CheckerFunc: method is nil but S3Interface.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		State *health.CheckState
	}{
		Ctx:   ctx,
		State: state,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(ctx, state)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedS3Interface.CheckerCalls())
func (mock *S3InterfaceMock) CheckerCalls() []struct {
	Ctx   context.Context
	State *health.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		State *health.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *S3InterfaceMock) Get(key string) (io.ReadCloser, *int64, error) {
	if mock.GetFunc == nil {
		panic("S3InterfaceMock.GetFunc: method is nil but S3Interface.Get was just called")
	}
	callInfo := struct {
		Key string
	}{
		Key: key,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(key)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedS3Interface.GetCalls())
func (mock *S3InterfaceMock) GetCalls() []struct {
	Key string
} {
	var calls []struct {
		Key string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// UploadPart calls UploadPartFunc.
func (mock *S3InterfaceMock) UploadPart(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error {
	if mock.UploadPartFunc == nil {
		panic("S3InterfaceMock.UploadPartFunc: method is nil but S3Interface.UploadPart was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Req     *s3client.UploadPartRequest
		Payload []byte
	}{
		Ctx:     ctx,
		Req:     req,
		Payload: payload,
	}
	mock.lockUploadPart.Lock()
	mock.calls.UploadPart = append(mock.calls.UploadPart, callInfo)
	mock.lockUploadPart.Unlock()
	return mock.UploadPartFunc(ctx, req, payload)
}

// UploadPartCalls gets all the calls that were made to UploadPart.
// Check the length with:
//     len(mockedS3Interface.UploadPartCalls())
func (mock *S3InterfaceMock) UploadPartCalls() []struct {
	Ctx     context.Context
	Req     *s3client.UploadPartRequest
	Payload []byte
} {
	var calls []struct {
		Ctx     context.Context
		Req     *s3client.UploadPartRequest
		Payload []byte
	}
	mock.lockUploadPart.RLock()
	calls = mock.calls.UploadPart
	mock.lockUploadPart.RUnlock()
	return calls
}