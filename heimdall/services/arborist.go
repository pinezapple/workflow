package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vfluxus/dvergr/arboristutil"
	"github.com/vfluxus/heimdall/core"
	"github.com/vfluxus/heimdall/services/dto"
)

type IArborist interface {
	AuthzRequests(ctx context.Context, token string, requests []dto.AuthRequestJSON_Request) (bool, error)
	AuthzRequest(ctx context.Context, token string, resource string, method string, service string) (bool, error)
	GetShareProjects(ctx context.Context, token string) (projectAuth []string, err error)
}

func GetArboristService() IArborist {
	return abrSrv
}

type arboristService struct {
	abr        arboristutil.IArboristUtil
	requestURL string
}

var abrSrv = new(arboristService)

func SetArborist(cfg *core.ArboristConfig) {
	url := fmt.Sprintf("%s/%s", cfg.Host, "auth/request")
	abrSrv = &arboristService{
		abr:        arboristutil.NewArborist(cfg.Host),
		requestURL: url,
	}
}

func (s arboristService) AuthzRequests(ctx context.Context, token string, requests []dto.AuthRequestJSON_Request) (result bool, err error) {
	result = false
	authzRequests := dto.AuthRequestJSON{
		User: dto.AuthRequestJSON_User{
			Token: token,
		},
		Requests: requests,
	}
	requestBody, err := json.Marshal(authzRequests)
	if err != nil {
		logger.Errorf("Convert authz requests error: %s", err.Error())
		return
	}

	request, err := http.NewRequestWithContext(ctx, "POST", s.requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Errorf("Create request error: %s", err.Error())
		return
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Errorf("Call arborist service error: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	var authzResp dto.AuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authzResp); err != nil {
		logger.Errorf("Can not convert response body: %s", err.Error())
		return
	}

	return authzResp.Auth, nil
}

func (s arboristService) AuthzRequest(ctx context.Context, token string, resource string, method string, service string) (result bool, err error) {
	result = false
	authzRequest := dto.AuthRequestJSON{
		User: dto.AuthRequestJSON_User{
			Token: token,
		},
		Request: &dto.AuthRequestJSON_Request{
			Resource: resource,
			Action: dto.Action{
				Service: service,
				Method:  method,
			},
		},
	}
	requestBody, err := json.Marshal(authzRequest)
	if err != nil {
		logger.Errorf("Convert authz request error: %s", err.Error())
		return
	}

	request, err := http.NewRequestWithContext(ctx, "POST", s.requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Errorf("Create request error: %s", err.Error())
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)

	fmt.Println(authzRequest.Request.Resource)

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Errorf("Call arborist service error: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	var authzResp dto.AuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authzResp); err != nil {
		logger.Errorf("Can not convert response body: %s", err.Error())
		return
	}

	return authzResp.Auth, nil
}

func (s arboristService) CreateShareProject(ctx context.Context, token string, projectUUID uuid.UUID, roles []string) (err error) {
	// create resource
	rsName := fmt.Sprintf("Project %s", projectUUID)
	rsPath := fmt.Sprintf("/analyse/share/projects/%s", projectUUID)
	if err := s.abr.CreateResource(ctx, token, rsName, rsPath, ""); err != nil {
		if !errors.Is(err, arboristutil.ExistedErr{}) {
			return err
		}
	}

	// create policy
	policyID := fmt.Sprintf("Full Perm %s", projectUUID)
	if err := s.abr.CreatePolicy(ctx, token, policyID, []string{rsPath}, roles, ""); err != nil {
		if !errors.Is(err, arboristutil.ExistedErr{}) {
			return err
		}
	}

	// grant policy
	loggedInGroup := "logged-in"
	if err := s.abr.GrantPolicy(ctx, token, loggedInGroup, policyID, nil); err != nil {
		if !errors.Is(err, arboristutil.ExistedErr{}) {
			return err
		}
	}

	return nil
}

func (s arboristService) GetShareProjects(ctx context.Context, token string) (projectAuth []string, err error) {
	// mapping
	prefix := "/analyse/share/projects/"
	rs, err := s.abr.Mapping(ctx, token, prefix)
	if err != nil {
		return nil, err
	}

	for i := range rs {
		projectAuth = append(projectAuth, rs[i].Path)
	}

	return projectAuth, nil
}
