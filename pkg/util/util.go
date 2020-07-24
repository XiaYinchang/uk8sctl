package util

import (
	"strings"

	"encoding/json"
	"fmt"
	"runtime"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/pkg/errors"
	"github.com/ucloud/ucloud-sdk-go/private/protocol/http"
	"github.com/ucloud/ucloud-sdk-go/private/utils"
	uerr "github.com/ucloud/ucloud-sdk-go/ucloud/error"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
	"github.com/ucloud/ucloud-sdk-go/ucloud/version"
)

func SetLogLevel(logLevel string) {
	switch logLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func buildHTTPRequest(req request.Common, credential *auth.Credential) (*http.HttpRequest, error) {
	query, err := request.ToQueryMap(req)
	if err != nil {
		return nil, errors.Errorf("convert request to map failed, %s", err)
	}

	if credential == nil {
		return nil, errors.Errorf("invalid credential information, please set it before request.")
	}

	config := ucloud.NewConfig()
	httpReq := http.NewHttpRequest()
	httpReq.SetURL(config.BaseUrl)
	httpReq.SetMethod("POST")

	// set timeout with client configuration
	httpReq.SetTimeout(config.Timeout)

	// keep query string is ordered and append credential signature as the last query param
	httpReq.SetQueryString(credential.BuildCredentialedQuery(query))

	ua := fmt.Sprintf("GO/%s GO-SDK/%s %s", runtime.Version(), version.Version, config.UserAgent)
	httpReq.SetHeader("User-Agent", strings.TrimSpace(ua))

	return httpReq, nil
}

// ResponseHandler receive response and write data into this response memory area
type ResponseHandler func(req request.Common, resp response.Common, err error) (response.Common, error)

// HttpResponseHandler receive http response and return a new http response
type HttpResponseHandler func(req *http.HttpRequest, resp *http.HttpResponse, err error) (*http.HttpResponse, error)

var defaultResponseHandlers = []ResponseHandler{errorHandler, logHandler}
var defaultHttpResponseHandlers = []HttpResponseHandler{errorHTTPHandler, logDebugHTTPHandler}

// errorHandler will normalize error to several specific error
func errorHandler(req request.Common, resp response.Common, err error) (response.Common, error) {
	if err != nil {
		if _, ok := err.(uerr.Error); ok {
			return resp, err
		}
		if uerr.IsNetworkError(err) {
			return resp, uerr.NewClientError(uerr.ErrNetwork, err)
		}
		return resp, uerr.NewClientError(uerr.ErrSendRequest, err)
	}

	if resp.GetRetCode() != 0 {
		return resp, uerr.NewServerCodeError(resp.GetRetCode(), resp.GetMessage())
	}

	return resp, err
}

func errorHTTPHandler(req *http.HttpRequest, resp *http.HttpResponse, err error) (*http.HttpResponse, error) {
	if err == nil {
		return resp, err
	}

	if statusErr, ok := err.(http.StatusError); ok {
		return resp, uerr.NewServerStatusError(statusErr.StatusCode, statusErr.Message)
	}

	return resp, err
}

func logHandler(req request.Common, resp response.Common, err error) (response.Common, error) {
	action := req.GetAction()

	if err != nil {
		log.Warnf("do %s failed, %s", action, err)
	} else {
		log.Debugf("do %s successful!", action)
	}
	return resp, err
}

func logDebugHTTPHandler(req *http.HttpRequest, resp *http.HttpResponse, err error) (*http.HttpResponse, error) {
	// logging request
	log.Debugf("%s", req)

	// logging error
	if err != nil {
		log.Errorf("%s", err)
	}

	// logging response code text
	if resp != nil && resp.GetStatusCode() >= 400 {
		log.Warnf("%s", resp.GetStatusCode())
	}

	// logging response body
	if resp != nil && resp.GetStatusCode() < 400 {
		log.Debugf("%s - %v", resp.GetBody(), resp.GetStatusCode())
	}

	return resp, err
}

func DoRequest(req request.Common, res response.Common, credential *auth.Credential) error {
	httpReq, err := buildHTTPRequest(req, credential)
	if err != nil {
		return uerr.NewClientError(uerr.ErrInvalidRequest, err)
	}
	res.SetRequest(req)
	httpClient := http.NewHttpClient()
	httpResp, err := httpClient.Send(httpReq)

	// use response middleware to handle http response
	// such as convert some http status to error
	for _, handler := range defaultHttpResponseHandlers {
		httpResp, err = handler(httpReq, httpResp, err)
	}
	if err == nil {
		body := httpResp.GetBody()
		body = utils.RetCodePatcher.Patch(body)
		err = json.Unmarshal(body, res)
		if err != nil {
			return err
		}
	}
	for _, handler := range defaultResponseHandlers {
		res, err = handler(req, res, err)
	}
	if err != nil {
		return err
	}
	return nil
}
