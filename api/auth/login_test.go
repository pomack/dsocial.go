package auth_test

import (
    "github.com/pomack/dsocial.go/api/auth"
    //"github.com/pomack/dsocial.go/api/apiutil"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/dsocial.go/backend/authentication"
    "github.com/pomack/dsocial.go/backend/datastore/inmemory"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/webmachine.go/webmachine"
    "bytes"
    "http"
    "json"
    "testing"
)

func initializeAuthUserAccountDS() (ds *inmemory.InMemoryDataStore, wm webmachine.WebMachine) {
    ds = inmemory.NewInMemoryDataStore()
    gw, _ := ds.CreateUserAccount(&dm.User{
        Role: dm.ROLE_ADMIN,
        Name: "George Washington",
        Username: "firstpresident",
        Email: "george@washington.com",
        PhoneNumber: "+1-405-555-5555",
        Address: "Valley Forge",
        AllowLogin: true,
    })
    tj, _ := ds.CreateUserAccount(&dm.User{
        Role: dm.ROLE_STANDARD,
        Name: "Thomas Jefferson",
        Username: "secondpresident",
        Email: "thomas@jefferson.com",
        PhoneNumber: "+1-401-555-5555",
        Address: "Virginia",
        AllowLogin: false,
    })
    ja, _ := ds.CreateUserAccount(&dm.User{
        Role: dm.ROLE_TECHNICAL_SUPPORT,
        Name: "John Adams",
        Username: "thirdpresident",
        Email: "john@adams.com",
        PhoneNumber: "+1-402-555-5555",
        Address: "Boston, MA",
        AllowLogin: true,
    })
    authentication.GenerateNewAccessKey(ds, gw.Id, "")
    authentication.SetUserPassword(ds, gw.Id, "number one")
    authentication.GenerateNewAccessKey(ds, tj.Id, "")
    authentication.SetUserPassword(ds, tj.Id, "number two")
    authentication.GenerateNewAccessKey(ds, ja.Id, "")
    authentication.SetUserPassword(ds, ja.Id, "number three")
    wm = webmachine.NewWebMachine()
    wm.AddRouteHandler(auth.NewLoginAccountRequestHandler(ds, ds))
    wm.AddRouteHandler(auth.NewLogoutAccountRequestHandler(ds, ds))
    return
}

func TestAuthLoginAdmin(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("username", user.Username)
    jsonobj.Set("password", "number one")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.StatusCode != http.StatusOK {
        t.Error("Expected ", http.StatusOK, " status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "success" {
        t.Error("Expected successful operation, but had status:", status)
    }
    result := obj.GetAsObject("result")
    if result == nil {
        t.Error("Expected an object for result, but was nil")
    } else {
        accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
        if len(accessKeys2) != 2 {
            t.Error("Expected 2 access keys after logging in, but found", len(accessKeys2))
        } else {
            var checkAccessKey *dm.AccessKey
            if accessKeys2[0].Id == accessKey.Id {
                checkAccessKey = accessKeys2[1]
            } else {
                checkAccessKey = accessKeys2[0]
            }
            if access_key_id := result.GetAsString("access_key_id"); access_key_id != checkAccessKey.Id {
                t.Error("Expected access_key_id with value", checkAccessKey.Id, "but was", access_key_id)
            }
            if private_key := result.GetAsString("private_key"); private_key != checkAccessKey.PrivateKey {
                t.Error("Expected private_key with value", checkAccessKey.PrivateKey, "but was", private_key)
            }
        }
        if username := result.GetAsString("username"); username != user.Username {
            t.Error("Expected username", user.Username, "but was", username)
        }
        if name := result.GetAsString("name"); name != user.Name {
            t.Error("Expected name", user.Name, "but was", name)
        }
        if user_id := result.GetAsString("user_id"); user_id != user.Id {
            t.Error("Expected user_id", user.Id, "but was", user_id)
        }
    }
}


func TestAuthLoginUser(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("username", user.Username)
    jsonobj.Set("password", "number three")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.StatusCode != http.StatusOK {
        t.Error("Expected ", http.StatusOK, " status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "success" {
        t.Error("Expected successful operation, but had status:", status)
    }
    result := obj.GetAsObject("result")
    if result == nil {
        t.Error("Expected an object for result, but was nil")
    } else {
        accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
        if len(accessKeys2) != 2 {
            t.Error("Expected 2 access keys after logging in, but found", len(accessKeys2))
        } else {
            var checkAccessKey *dm.AccessKey
            if accessKeys2[0].Id == accessKey.Id {
                checkAccessKey = accessKeys2[1]
            } else {
                checkAccessKey = accessKeys2[0]
            }
            if access_key_id := result.GetAsString("access_key_id"); access_key_id != checkAccessKey.Id {
                t.Error("Expected access_key_id with value", checkAccessKey.Id, "but was", access_key_id)
            }
            if private_key := result.GetAsString("private_key"); private_key != checkAccessKey.PrivateKey {
                t.Error("Expected private_key with value", checkAccessKey.PrivateKey, "but was", private_key)
            }
        }
        if username := result.GetAsString("username"); username != user.Username {
            t.Error("Expected username", user.Username, "but was", username)
        }
        if name := result.GetAsString("name"); name != user.Name {
            t.Error("Expected name", user.Name, "but was", name)
        }
        if user_id := result.GetAsString("user_id"); user_id != user.Id {
            t.Error("Expected user_id", user.Id, "but was", user_id)
        }
    }
}


func TestAuthLoginDisabledUser(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("secondpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("username", user.Username)
    jsonobj.Set("password", "number two")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "error" {
        t.Error("Expected error operation, but had status:", status)
    }
    if result := obj.Get("result"); result != nil {
        t.Error("Expected result to be nil, but was", result)
    }
    if message := obj.GetAsString("message"); message != auth.ERR_INVALID_USERNAME_PASSWORD_COMBO.String() {
        t.Error("Expected ERR_INVALID_USERNAME_PASSWORD_COMBO for message, but was", message)
    }
    if accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000); len(accessKeys2) != 1 {
        t.Error("Expected 1 access key after logging in, but found", len(accessKeys2))
    }
}


func TestAuthLoginNoUsername(t *testing.T) {
    _, wm := initializeAuthUserAccountDS()
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("password", "number two")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "error" {
        t.Error("Expected error operation, but had status:", status)
    }
    result := obj.GetAsObject("result")
    if result.Len() != 1 {
        t.Error("Expected a result object with 1 entry, but has", result.Len(), "entries as:", result)
    }
    if username := result.GetAsArray("username"); len(username) != 1 || username[0] != auth.ERR_MUST_SPECIFY_USERNAME.String() {
        t.Error("Expected one error for missing username, but was", result)
    }
    if message := obj.GetAsString("message"); message != auth.ERR_VALUE_ERRORS.String() {
        t.Error("Expected ERR_VALUE_ERRORS for message, but was", message)
    }
}


func TestAuthLoginNoPassword(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("username", user.Username)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "error" {
        t.Error("Expected error operation, but had status:", status)
    }
    result := obj.GetAsObject("result")
    if result.Len() != 1 {
        t.Error("Expected a result object with 1 entry, but has", result.Len(), "entries as:", result)
    }
    if password := result.GetAsArray("password"); len(password) != 1 || password[0] != auth.ERR_MUST_SPECIFY_PASSWORD.String() {
        t.Error("Expected one error for missing password, but was", result)
    }
    if message := obj.GetAsString("message"); message != auth.ERR_VALUE_ERRORS.String() {
        t.Error("Expected ERR_VALUE_ERRORS for message, but was", message)
    }
    if accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000); len(accessKeys2) != 1 {
        t.Error("Expected 1 access key after logging in, but found", len(accessKeys2))
    }
}

func TestAuthLoginNoUsernameNorPassword(t *testing.T) {
    _, wm := initializeAuthUserAccountDS()
    jsonobj := jsonhelper.NewJSONObject()
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "error" {
        t.Error("Expected error operation, but had status:", status)
    }
    result := obj.GetAsObject("result")
    if result.Len() != 2 {
        t.Error("Expected a result object with 2 entries, but has", result.Len(), "entries as:", result)
    }
    if username := result.GetAsArray("username"); len(username) != 1 || username[0] != auth.ERR_MUST_SPECIFY_USERNAME.String() {
        t.Error("Expected one error for missing username, but was", result)
    }
    if password := result.GetAsArray("password"); len(password) != 1 || password[0] != auth.ERR_MUST_SPECIFY_PASSWORD.String() {
        t.Error("Expected one error for missing password, but was", result)
    }
    if message := obj.GetAsString("message"); message != auth.ERR_VALUE_ERRORS.String() {
        t.Error("Expected ERR_VALUE_ERRORS for message, but was", message)
    }
}


func TestAuthLoginBadPassword(t *testing.T) {
    ds, wm := initializeAuthUserAccountDS()
    user, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("username", user.Username)
    jsonobj.Set("password", "blah blah")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "error" {
        t.Error("Expected error operation, but had status:", status)
    }
    if result := obj.Get("result"); result != nil {
        t.Error("Expected result to be nil, but was", result)
    }
    if message := obj.GetAsString("message"); message != auth.ERR_INVALID_USERNAME_PASSWORD_COMBO.String() {
        t.Error("Expected ERR_INVALID_USERNAME_PASSWORD_COMBO for message, but was", message)
    }
    if accessKeys2, _, _ := ds.RetrieveUserKeys(user.Id, nil, 1000); len(accessKeys2) != 1 {
        t.Error("Expected 1 access key after logging in, but found", len(accessKeys2))
    }
}



func TestAuthLoginAccountDoesNotExist(t *testing.T) {
    _, wm := initializeAuthUserAccountDS()
    jsonobj := jsonhelper.NewJSONObject()
    jsonobj.Set("username", "dudewhatever")
    jsonobj.Set("password", "blah blah")
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/auth/login/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    if err != nil {
        t.Error("Unable to unmarshal login response due to error:", err.String())
    }
    if status := obj.GetAsString("status"); status != "error" {
        t.Error("Expected error operation, but had status:", status)
    }
    if result := obj.Get("result"); result != nil {
        t.Error("Expected result to be nil, but was", result)
    }
    if message := obj.GetAsString("message"); message != auth.ERR_INVALID_USERNAME_PASSWORD_COMBO.String() {
        t.Error("Expected ERR_INVALID_USERNAME_PASSWORD_COMBO for message, but was", message)
    }
}


