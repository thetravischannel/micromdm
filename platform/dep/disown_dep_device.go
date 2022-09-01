package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *DEPService) DisownDevice(ctx context.Context, serials []string) (*dep.DeviceStatusResponse, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.DisownDevice(serials...)
}

type disownDeviceRequest struct {
	Serials []string `json:"serials"`
}

type deviceStatusResponse struct {
	*dep.DeviceStatusResponse
	Err error `json:"err,omitempty"`
}

func (r deviceStatusResponse) Failed() error { return r.Err }

func decodeDisownDeviceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req disownDeviceRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeDeviceStatusResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp deviceStatusResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeDisownDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(disownDeviceRequest)
		details, err := svc.DisownDevice(ctx, req.Serials)
		return deviceStatusResponse{DeviceStatusResponse: details, Err: err}, nil
	}
}

func (e Endpoints) DisownDevice(ctx context.Context, serials []string) (*dep.DeviceStatusResponse, error) {
	request := disownDeviceRequest{Serials: serials}
	response, err := e.DisownDeviceEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(deviceStatusResponse).DeviceStatusResponse, response.(deviceStatusResponse).Err
}
