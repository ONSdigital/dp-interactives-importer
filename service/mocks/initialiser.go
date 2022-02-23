// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks_service

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/health"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/service"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"net/http"
	"sync"
)

// Ensure, that InitialiserMock does implement service.Initialiser.
// If this is not the case, regenerate this file with moq.
var _ service.Initialiser = &InitialiserMock{}

// InitialiserMock is a mock implementation of service.Initialiser.
//
// 	func TestSomethingThatUsesInitialiser(t *testing.T) {
//
// 		// make and configure a mocked service.Initialiser
// 		mockedInitialiser := &InitialiserMock{
// 			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
// 				panic("mock out the DoGetHTTPServer method")
// 			},
// 			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
// 				panic("mock out the DoGetHealthCheck method")
// 			},
// 			DoGetHealthClientFunc: func(name string, url string) *health.Client {
// 				panic("mock out the DoGetHealthClient method")
// 			},
// 			DoGetKafkaConsumerFunc: func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
// 				panic("mock out the DoGetKafkaConsumer method")
// 			},
// 			DoGetS3ClientFunc: func(ctx context.Context, cfg *config.Config) (importer.S3Interface, error) {
// 				panic("mock out the DoGetS3Client method")
// 			},
// 			DoGetUploadServiceBackendFunc: func(ctx context.Context, cfg *config.Config) (importer.UploadServiceBackend, error) {
// 				panic("mock out the DoGetUploadServiceBackend method")
// 			},
// 		}
//
// 		// use mockedInitialiser in code that requires service.Initialiser
// 		// and then make assertions.
//
// 	}
type InitialiserMock struct {
	// DoGetHTTPServerFunc mocks the DoGetHTTPServer method.
	DoGetHTTPServerFunc func(bindAddr string, router http.Handler) service.HTTPServer

	// DoGetHealthCheckFunc mocks the DoGetHealthCheck method.
	DoGetHealthCheckFunc func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error)

	// DoGetHealthClientFunc mocks the DoGetHealthClient method.
	DoGetHealthClientFunc func(name string, url string) *health.Client

	// DoGetKafkaConsumerFunc mocks the DoGetKafkaConsumer method.
	DoGetKafkaConsumerFunc func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error)

	// DoGetS3ClientFunc mocks the DoGetS3Client method.
	DoGetS3ClientFunc func(ctx context.Context, cfg *config.Config) (importer.S3Interface, error)

	// DoGetUploadServiceBackendFunc mocks the DoGetUploadServiceBackend method.
	DoGetUploadServiceBackendFunc func(ctx context.Context, cfg *config.Config) (importer.UploadServiceBackend, error)

	// calls tracks calls to the methods.
	calls struct {
		// DoGetHTTPServer holds details about calls to the DoGetHTTPServer method.
		DoGetHTTPServer []struct {
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// DoGetHealthCheck holds details about calls to the DoGetHealthCheck method.
		DoGetHealthCheck []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// BuildTime is the buildTime argument value.
			BuildTime string
			// GitCommit is the gitCommit argument value.
			GitCommit string
			// Version is the version argument value.
			Version string
		}
		// DoGetHealthClient holds details about calls to the DoGetHealthClient method.
		DoGetHealthClient []struct {
			// Name is the name argument value.
			Name string
			// URL is the url argument value.
			URL string
		}
		// DoGetKafkaConsumer holds details about calls to the DoGetKafkaConsumer method.
		DoGetKafkaConsumer []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
		// DoGetS3Client holds details about calls to the DoGetS3Client method.
		DoGetS3Client []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
		// DoGetUploadServiceBackend holds details about calls to the DoGetUploadServiceBackend method.
		DoGetUploadServiceBackend []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
	}
	lockDoGetHTTPServer           sync.RWMutex
	lockDoGetHealthCheck          sync.RWMutex
	lockDoGetHealthClient         sync.RWMutex
	lockDoGetKafkaConsumer        sync.RWMutex
	lockDoGetS3Client             sync.RWMutex
	lockDoGetUploadServiceBackend sync.RWMutex
}

// DoGetHTTPServer calls DoGetHTTPServerFunc.
func (mock *InitialiserMock) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	if mock.DoGetHTTPServerFunc == nil {
		panic("InitialiserMock.DoGetHTTPServerFunc: method is nil but Initialiser.DoGetHTTPServer was just called")
	}
	callInfo := struct {
		BindAddr string
		Router   http.Handler
	}{
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockDoGetHTTPServer.Lock()
	mock.calls.DoGetHTTPServer = append(mock.calls.DoGetHTTPServer, callInfo)
	mock.lockDoGetHTTPServer.Unlock()
	return mock.DoGetHTTPServerFunc(bindAddr, router)
}

// DoGetHTTPServerCalls gets all the calls that were made to DoGetHTTPServer.
// Check the length with:
//     len(mockedInitialiser.DoGetHTTPServerCalls())
func (mock *InitialiserMock) DoGetHTTPServerCalls() []struct {
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		BindAddr string
		Router   http.Handler
	}
	mock.lockDoGetHTTPServer.RLock()
	calls = mock.calls.DoGetHTTPServer
	mock.lockDoGetHTTPServer.RUnlock()
	return calls
}

// DoGetHealthCheck calls DoGetHealthCheckFunc.
func (mock *InitialiserMock) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	if mock.DoGetHealthCheckFunc == nil {
		panic("InitialiserMock.DoGetHealthCheckFunc: method is nil but Initialiser.DoGetHealthCheck was just called")
	}
	callInfo := struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}{
		Cfg:       cfg,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		Version:   version,
	}
	mock.lockDoGetHealthCheck.Lock()
	mock.calls.DoGetHealthCheck = append(mock.calls.DoGetHealthCheck, callInfo)
	mock.lockDoGetHealthCheck.Unlock()
	return mock.DoGetHealthCheckFunc(cfg, buildTime, gitCommit, version)
}

// DoGetHealthCheckCalls gets all the calls that were made to DoGetHealthCheck.
// Check the length with:
//     len(mockedInitialiser.DoGetHealthCheckCalls())
func (mock *InitialiserMock) DoGetHealthCheckCalls() []struct {
	Cfg       *config.Config
	BuildTime string
	GitCommit string
	Version   string
} {
	var calls []struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}
	mock.lockDoGetHealthCheck.RLock()
	calls = mock.calls.DoGetHealthCheck
	mock.lockDoGetHealthCheck.RUnlock()
	return calls
}

// DoGetHealthClient calls DoGetHealthClientFunc.
func (mock *InitialiserMock) DoGetHealthClient(name string, url string) *health.Client {
	if mock.DoGetHealthClientFunc == nil {
		panic("InitialiserMock.DoGetHealthClientFunc: method is nil but Initialiser.DoGetHealthClient was just called")
	}
	callInfo := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  url,
	}
	mock.lockDoGetHealthClient.Lock()
	mock.calls.DoGetHealthClient = append(mock.calls.DoGetHealthClient, callInfo)
	mock.lockDoGetHealthClient.Unlock()
	return mock.DoGetHealthClientFunc(name, url)
}

// DoGetHealthClientCalls gets all the calls that were made to DoGetHealthClient.
// Check the length with:
//     len(mockedInitialiser.DoGetHealthClientCalls())
func (mock *InitialiserMock) DoGetHealthClientCalls() []struct {
	Name string
	URL  string
} {
	var calls []struct {
		Name string
		URL  string
	}
	mock.lockDoGetHealthClient.RLock()
	calls = mock.calls.DoGetHealthClient
	mock.lockDoGetHealthClient.RUnlock()
	return calls
}

// DoGetKafkaConsumer calls DoGetKafkaConsumerFunc.
func (mock *InitialiserMock) DoGetKafkaConsumer(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
	if mock.DoGetKafkaConsumerFunc == nil {
		panic("InitialiserMock.DoGetKafkaConsumerFunc: method is nil but Initialiser.DoGetKafkaConsumer was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.Config
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetKafkaConsumer.Lock()
	mock.calls.DoGetKafkaConsumer = append(mock.calls.DoGetKafkaConsumer, callInfo)
	mock.lockDoGetKafkaConsumer.Unlock()
	return mock.DoGetKafkaConsumerFunc(ctx, cfg)
}

// DoGetKafkaConsumerCalls gets all the calls that were made to DoGetKafkaConsumer.
// Check the length with:
//     len(mockedInitialiser.DoGetKafkaConsumerCalls())
func (mock *InitialiserMock) DoGetKafkaConsumerCalls() []struct {
	Ctx context.Context
	Cfg *config.Config
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.Config
	}
	mock.lockDoGetKafkaConsumer.RLock()
	calls = mock.calls.DoGetKafkaConsumer
	mock.lockDoGetKafkaConsumer.RUnlock()
	return calls
}

// DoGetS3Client calls DoGetS3ClientFunc.
func (mock *InitialiserMock) DoGetS3Client(ctx context.Context, cfg *config.Config) (importer.S3Interface, error) {
	if mock.DoGetS3ClientFunc == nil {
		panic("InitialiserMock.DoGetS3ClientFunc: method is nil but Initialiser.DoGetS3Client was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.Config
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetS3Client.Lock()
	mock.calls.DoGetS3Client = append(mock.calls.DoGetS3Client, callInfo)
	mock.lockDoGetS3Client.Unlock()
	return mock.DoGetS3ClientFunc(ctx, cfg)
}

// DoGetS3ClientCalls gets all the calls that were made to DoGetS3Client.
// Check the length with:
//     len(mockedInitialiser.DoGetS3ClientCalls())
func (mock *InitialiserMock) DoGetS3ClientCalls() []struct {
	Ctx context.Context
	Cfg *config.Config
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.Config
	}
	mock.lockDoGetS3Client.RLock()
	calls = mock.calls.DoGetS3Client
	mock.lockDoGetS3Client.RUnlock()
	return calls
}

// DoGetUploadServiceBackend calls DoGetUploadServiceBackendFunc.
func (mock *InitialiserMock) DoGetUploadServiceBackend(ctx context.Context, cfg *config.Config) (importer.UploadServiceBackend, error) {
	if mock.DoGetUploadServiceBackendFunc == nil {
		panic("InitialiserMock.DoGetUploadServiceBackendFunc: method is nil but Initialiser.DoGetUploadServiceBackend was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.Config
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetUploadServiceBackend.Lock()
	mock.calls.DoGetUploadServiceBackend = append(mock.calls.DoGetUploadServiceBackend, callInfo)
	mock.lockDoGetUploadServiceBackend.Unlock()
	return mock.DoGetUploadServiceBackendFunc(ctx, cfg)
}

// DoGetUploadServiceBackendCalls gets all the calls that were made to DoGetUploadServiceBackend.
// Check the length with:
//     len(mockedInitialiser.DoGetUploadServiceBackendCalls())
func (mock *InitialiserMock) DoGetUploadServiceBackendCalls() []struct {
	Ctx context.Context
	Cfg *config.Config
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.Config
	}
	mock.lockDoGetUploadServiceBackend.RLock()
	calls = mock.calls.DoGetUploadServiceBackend
	mock.lockDoGetUploadServiceBackend.RUnlock()
	return calls
}
