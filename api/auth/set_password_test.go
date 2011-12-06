package auth_test

import (
    //"github.com/pomack/dsocial.go/api/auth"
    "github.com/pomack/dsocial.go/api/apiutil"
    //dm "github.com/pomack/dsocial.go/models/dsocial"
    //"github.com/pomack/dsocial.go/backend/authentication"
    //"github.com/pomack/dsocial.go/backend/datastore/inmemory"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/webmachine.go/webmachine"
    "bytes"
    "http"
    "json"
    "testing"
)

func TestAuthSetPasswordAdmin(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) != 1 {
        t.Error("Expected to find 1 access key stored, but found", len(accessKeys))
    }
    accessKey := accessKeys[0]
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("password", "hi ho hi ho")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/set_password/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    reqbytes, _ := http.DumpRequest(req, true)
    t.Log("Request is:\n", string(reqbytes), "\n\n")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    t.Log("Response is:\n", resp.String(), "\n\n")
    if resp.StatusCode != http.StatusOK {
        t.Error("Expected ", http.StatusOK, " status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal setPassword response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "success" {
        t.Error("Expected successful operation, but had status:", status)
    }
    result := obj.GetAsObject("result")
    if message := result.GetAsString("message"); message != "password changed" {
        t.Error("Expected message == \"password changed\", but was \"", message, "\"")
    }
    user2 := result.GetAsObject("user")
    uid := user2.GetAsString("id")
    if uid != user.Id {
        t.Error("Expected user id of", user.Id, ", but was", uid)
    }
    accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys2) != 1 {
        t.Error("Expected to find one access key stored, but found", len(accessKeys))
    }
}

/*
func TestAuthSetPasswordUser(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKey, _ := authentication.GenerateNewAccessKey(ds, user.Id, "")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/set_password/", nil)
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    reqbytes, _ := http.DumpRequest(req, true)
    t.Log("Request is:\n", string(reqbytes), "\n\n")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    t.Log("Response is:\n", resp.String(), "\n\n")
    if resp.StatusCode != http.StatusOK {
        t.Error("Expected ", http.StatusOK, " status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal setPassword response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "success" {
        t.Error("Expected successful operation, but had status:", status)
    }
    if result := obj.Get("result"); result != nil {
        t.Error("Expected nil for result, but was", result)
    }
    accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys2) != 1 {
        t.Error("Expected to find one access key stored, but found", len(accessKeys))
    } else if len(accessKeys2) > 0 && accessKeys2[0].Id == accessKey.Id {
        t.Error("Incorrect access key was removed in setPassword")
    }
}


func TestAuthSetPasswordDisabledUser(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("secondpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/set_password/", nil)
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    reqbytes, _ := http.DumpRequest(req, true)
    t.Log("Request is:\n", string(reqbytes), "\n\n")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    t.Log("Response is:\n", resp.String(), "\n\n")
    if resp.StatusCode != http.StatusForbidden {
        t.Error("Expected ", http.StatusForbidden, " status code but received ", resp.StatusCode)
    }
    respbytes := resp.Buffer.Bytes()
    if len(respbytes) > 0 {
        t.Error("Expected a zero byte response but received:", string(respbytes))
    }
}


func TestAuthSetPasswordNoCredentials(t *testing.T) {
    _, wm := initializeAuthUserAccountDS()
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("password", "number two")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/set_password/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    reqbytes, _ := http.DumpRequest(req, true)
    t.Log("Request is:\n", string(reqbytes), "\n\n")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    t.Log("Response is:\n", resp.String(), "\n\n")
    if resp.StatusCode != http.StatusUnauthorized {
        t.Error("Expected ", http.StatusUnauthorized, " status code but received ", resp.StatusCode)
    }
    if authenticate := resp.Header().Get("Www-Authenticate"); authenticate != "dsocial" {
        t.Error("Expected header Www-Authenticate: dsocial but found value:", authenticate)
    }
    respbytes := resp.Buffer.Bytes()
    if len(respbytes) > 0 {
        t.Error("Expected a zero byte response but received:", string(respbytes))
    }
}


func TestAuthSetPasswordBadAuthorizationKey(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("firstpresident")
    user2, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    accessKeys2, _, _ := ds.RetrieveUserKeys(user2.Id, nil, 1000)
    if len(accessKeys) == 0 || len(accessKeys2) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    accessKey2 := accessKeys2[0]
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/set_password/", nil)
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey2.PrivateKey).SignRequest(req, 0)
    reqbytes, _ := http.DumpRequest(req, true)
    t.Log("Request is:\n", string(reqbytes), "\n\n")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    t.Log("Response is:\n", resp.String(), "\n\n")
    if resp.StatusCode != http.StatusForbidden {
        t.Error("Expected ", http.StatusForbidden, " status code but received ", resp.StatusCode)
    }
    respbytes := resp.Buffer.Bytes()
    if len(respbytes) > 0 {
        t.Error("Expected a zero byte response but received:", string(respbytes))
    }
}

func TestAuthSetPasswordAccountDoesNotExist(t *testing.T) {
    _, wm := initializeAuthUserAccountDS()
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/set_password/", nil)
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner("sdfsfsflsfdfsffdsdfsf", "sfsfsfssrgsgsdgdsgdgegergdggdsgdsg").SignRequest(req, 0)
    reqbytes, _ := http.DumpRequest(req, true)
    t.Log("Request is:\n", string(reqbytes), "\n\n")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    t.Log("Response is:\n", resp.String(), "\n\n")
    if resp.StatusCode != http.StatusForbidden {
        t.Error("Expected ", http.StatusForbidden, " status code but received ", resp.StatusCode)
    }
    respbytes := resp.Buffer.Bytes()
    if len(respbytes) > 0 {
        t.Error("Expected a zero byte response but received:", string(respbytes))
    }
}
*/

