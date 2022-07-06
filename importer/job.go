package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/log.go/v2/log"
	"sync"
)

type Job struct {
	mu                    sync.Mutex
	archiveFiles          []*interactives.InteractiveFile
	ctx                   context.Context
	interactivesAPIClient InteractivesAPIClient
	serviceAuthToken      string
}

func NewJob(ctx context.Context, cfg *config.Config, interactivesAPIClient InteractivesAPIClient) *Job {
	return &Job{
		ctx:                   ctx,
		serviceAuthToken:      cfg.ServiceAuthToken,
		interactivesAPIClient: interactivesAPIClient,
	}
}

func (j *Job) Files() []*interactives.InteractiveFile {
	return j.archiveFiles
}

func (j *Job) Add(file *interactives.InteractiveFile) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.archiveFiles = append(j.archiveFiles, file)
}

func (j *Job) Finish(logData *log.Data, event *InteractivesUploaded, zipSize *int64, err *error) {
	//todo sanity check?
	l := *logData
	e := *err

	patchReq := interactives.PatchRequest{
		Attribute: interactives.PatchArchive,
		Interactive: interactives.Interactive{
			ID: event.ID,
			Archive: &interactives.InteractiveArchive{
				Name: event.Path,
			},
		},
	}
	if e != nil {
		l["error"] = e.Error()
		patchReq.Interactive.Archive.ImportMessage = e.Error()
	} else {
		patchReq.Interactive.Archive.ImportSuccessful = true
		patchReq.Interactive.Archive.Files = j.archiveFiles
		if zipSize != nil {
			patchReq.Interactive.Archive.Size = *zipSize
		}
	}
	// user token not valid - we auth user on api endpoints
	_, apiErr := j.interactivesAPIClient.PatchInteractive(j.ctx, "", j.serviceAuthToken, event.ID, patchReq)
	if apiErr != nil {
		//todo what if this fails - retry?
		l["apiError"] = apiErr.Error()
		log.Warn(j.ctx, "failed to update interactive", logData)
	}
}
