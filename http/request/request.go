package request

import (
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/win30221/core/basic"
	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/http/consts"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/syserrno"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	URL    string
	Data   string
	Result interface{}
	CTX    *ctx.Context
	// header
	Header        http.Header
	DefaultHeader bool
}

func GET(r *Request) (err error) {
	req, err := http.NewRequest("GET", r.URL, nil)
	if err != nil {
		err = catch.New(syserrno.HTTP, "new request error", fmt.Sprintf("new request error: %s", err.Error()))
		return
	}

	if r.Data != "" {
		req.URL.RawQuery = r.Data
	}

	if r.Header != nil {
		req.Header = r.Header
	}

	if r.DefaultHeader {
		req.Header.Add(consts.HeaderSysToken, basic.SysToken)
		req.Header.Add(consts.HeaderXRequestID, r.CTX.TraceCode)
	}

	if r.CTX.Context != nil {
		req = req.WithContext(r.CTX.Context)
	}

	return exec(req, r)
}

func POST(r *Request) (err error) {
	req, err := http.NewRequest("POST", r.URL, strings.NewReader(r.Data))
	if err != nil {
		err = catch.New(syserrno.HTTP, "new request error", fmt.Sprintf("new request error: %s", err.Error()))
		return
	}

	if r.Header != nil {
		req.Header = r.Header
	}

	if r.DefaultHeader {
		req.Header.Add(consts.HeaderSysToken, basic.SysToken)
		req.Header.Add(consts.HeaderXRequestID, r.CTX.TraceCode)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if r.CTX.Context != nil {
		req = req.WithContext(r.CTX.Context)
	}

	return exec(req, r)
}

func PUT(r *Request) (err error) {
	req, err := http.NewRequest("PUT", r.URL, strings.NewReader(r.Data))
	if err != nil {
		err = catch.New(syserrno.HTTP, "new request error", fmt.Sprintf("new request error: %s", err.Error()))
		return
	}

	if r.Header != nil {
		req.Header = r.Header
	}

	if r.DefaultHeader {
		req.Header.Add(consts.HeaderSysToken, basic.SysToken)
		req.Header.Add(consts.HeaderXRequestID, r.CTX.TraceCode)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if r.CTX.Context != nil {
		req = req.WithContext(r.CTX.Context)
	}

	return exec(req, r)
}

func PATCH(r *Request) (err error) {
	req, err := http.NewRequest("PATCH", r.URL, strings.NewReader(r.Data))
	if err != nil {
		err = catch.New(syserrno.HTTP, "new request error", fmt.Sprintf("new request error: %s", err.Error()))
		return
	}

	if r.Header != nil {
		req.Header = r.Header
	}

	if r.DefaultHeader {
		req.Header.Add(consts.HeaderSysToken, basic.SysToken)
		req.Header.Add(consts.HeaderXRequestID, r.CTX.TraceCode)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if r.CTX.Context != nil {
		req = req.WithContext(r.CTX.Context)
	}

	return exec(req, r)
}

func exec(req *http.Request, r *Request) (err error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		err = catch.New(syserrno.HTTP, err.Error(), fmt.Sprintf("err: %s, req: %+v", err.Error(), req))
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(r.Result)
	if err != nil {
		if resp.StatusCode != http.StatusOK {
			err = catch.New(
				syserrno.HTTP,
				fmt.Sprintf("call %s error: %s", r.URL, err.Error()),
				fmt.Sprintf("call %s error, resp: %v, req: %+v, header: %+v, err: %s ", r.URL, r.Data, resp, r.Header, err.Error()),
			)

			return
		}
	}

	return
}
