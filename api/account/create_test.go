package account_test

import (
    "github.com/pomack/dsocial.go/api/account"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/dsocial.go/backend/datastore/inmemory"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/webmachine.go/webmachine"
    "bytes"
    "http"
    "json"
    "testing"
)

type Stringer interface {
    String() string
}

func TestCreateUserAccount(t *testing.T) {
    ds := inmemory.NewInMemoryDataStore()
    wm := webmachine.NewWebMachine()
    wm.AddRouteHandler(account.NewCreateAccountRequestHandler(ds, ds))
    buf := bytes.NewBufferString(`{"role": 9999999999999999, "name": "George Washington", "username": "firstpresident", "email":"george@washington.com", "phone_number": "+1-405-555-5555", "address": "Valley Forge"}`)
    oldUser := new(dm.User)
    json.Unmarshal(buf.Bytes(), &oldUser)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/create/", buf)
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    resp := webmachine.NewMockResponseWriter(req)
    reqb, _ := http.DumpRequest(req, true)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != 200 {
        t.Error("Expected 200 status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    user := new(dm.User)
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    user.InitFromJSONObject(obj.GetAsObject("result").GetAsObject("user"))
    if err != nil {
        t.Error("Error while unmarshaling JSON: ", err.String())
    }
    if obj.GetAsString("status") != "success" {
        t.Error("Expected status = \"success\", but was \"", obj.GetAsString("status"), "\"")
    }
    if user.Name != oldUser.Name {
        t.Logf("Request was\n%s\n================\n", string(reqb))
        t.Log("Response is:\n", resp.String(), "\n\n")
        t.Error("Expected name = \"", oldUser.Name, "\", but was ", user.Name)
    }
    if user.Username != oldUser.Username {
        t.Error("Expected username = \"", oldUser.Username, "\", but was ", user.Username)
    }
    if user.Email != oldUser.Email {
        t.Error("Expected email = \"", oldUser.Email, "\", but was ", user.Email)
    }
    if user.PhoneNumber != oldUser.PhoneNumber {
        t.Error("Expected phone_number = \"", oldUser.PhoneNumber, "\", but was ", user.PhoneNumber)
    }
    if user.Address != oldUser.Address {
        t.Error("Expected address = \"", oldUser.Address, "\", but was ", user.Address)
    }
    if user.Role != dm.ROLE_STANDARD {
        t.Error("Expected role = ", dm.ROLE_STANDARD, " but was ", user.Role)
    }
    if user.Id == "" {
        t.Error("Expected id to be populated, but was empty")
    }
}

func TestCreateUserAccountMissingName(t *testing.T) {
    ds := inmemory.NewInMemoryDataStore()
    wm := webmachine.NewWebMachine()
    wm.AddRouteHandler(account.NewCreateAccountRequestHandler(ds, ds))
    buf := bytes.NewBufferString(`{"role": 9999999999999999, "username": "firstpresident", "email":"george@washington.com", "phone_number": "+1-405-555-5555", "address": "Valley Forge"}`)
    oldUser := new(dm.User)
    json.Unmarshal(buf.Bytes(), &oldUser)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/create/", buf)
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != 400 {
        t.Error("Expected 400 status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Error while unmarshaling JSON: ", err.String())
    }
    if obj.GetAsString("status") != "error" {
        t.Error("Expected status = \"error\", but was \"", obj.GetAsString("status"), "\"")
    }
    result := obj.GetAsObject("result")
    if result == nil {
        t.Error("Expected result != nil, but was nil")
    } else {
        if result.GetAsArray("name").Len() == 0 {
            t.Error("Expected an error on attribute \"name\", but was not found")
        }
    }
}

func TestCreateUserAccountMissingSeveralFields(t *testing.T) {
    ds := inmemory.NewInMemoryDataStore()
    wm := webmachine.NewWebMachine()
    wm.AddRouteHandler(account.NewCreateAccountRequestHandler(ds, ds))
    buf := bytes.NewBufferString(`{"role": 9999999999999999, "name": "    ", "username": "\n\r\n", "email": "hi ho hi ho", "phone_number": "+1-405-555-5555", "address": "Valley Forge"}`)
    oldUser := new(dm.User)
    json.Unmarshal(buf.Bytes(), &oldUser)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/create/", buf)
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != 400 {
        t.Error("Expected 400 status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Error while unmarshaling JSON: ", err.String())
    }
    if obj.GetAsString("status") != "error" {
        t.Error("Expected status = \"error\", but was \"", obj.GetAsString("status"), "\"")
    }
    result := obj.GetAsObject("result")
    if result == nil {
        t.Error("Expected result != nil, but was nil")
    } else {
        if result.GetAsArray("name").Len() == 0 {
            t.Error("Expected an error on attribute \"name\", but was not found")
        }
        if result.GetAsArray("username").Len() == 0 {
            t.Error("Expected an error on attribute \"username\", but was not found")
        }
        if result.GetAsArray("email").Len() == 0 {
            t.Error("Expected an error on attribute \"email\", but was not found")
        }
    }
}
