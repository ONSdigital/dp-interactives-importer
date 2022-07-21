package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/log.go/v2/log"
)

type Job struct {
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

func (j *Job) Finish(logData *log.Data, event *InteractivesUploaded, uploadRootDirectory string, zipSize *int64, err *error) {
	//todo sanity check?
	l := *logData
	e := *err

	patchReq := interactives.PatchRequest{
		Attribute: interactives.PatchArchive,
		Interactive: interactives.Interactive{
			ID: event.ID,
			Archive: &interactives.Archive{
				Name: event.Path,
			},
		},
	}
	if e != nil {
		l["error"] = e.Error()
		patchReq.Interactive.Archive.ImportMessage = e.Error()
		patchReq.Interactive.Archive.UploadRootDirectory = uploadRootDirectory
	} else {
		patchReq.Interactive.Archive.ImportSuccessful = true
		patchReq.Interactive.Archive.UploadRootDirectory = uploadRootDirectory
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
