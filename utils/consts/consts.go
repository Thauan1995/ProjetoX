package consts

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var (
	IDProjeto      = os.Getenv("GOOGLE_CLOUD_PROJECT")
	IDServico      = os.Getenv("GAE_SERVICE")
	IDVersao       = os.Getenv("GAE_VERSION")
	BaseURL        = fmt.Sprintf("%s.appspot.com", IDProjeto)
	Location       = "us-central1"
	Region         = "us"
	GAEApplication = os.Getenv("GAE_APPLICATION")
)

func CloudProject(c context.Context) (*cloudresourcemanager.Project, error) {
	client, err := google.DefaultClient(c, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cloudService, err := cloudresourcemanager.New(client)
	if err != nil {
		return nil, err
	}
	return cloudService.Projects.Get(IDProjeto).Context(c).Do()
}

func ContaServico(c context.Context) (string, error) {
	project, err := CloudProject(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d-compute@developer.gserviceaccount.com", project.ProjectNumber), nil
}
