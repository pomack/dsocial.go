package account_test

import (
    "github.com/pomack/dsocial.go/api/account"
    "github.com/pomack/dsocial.go/api/apiutil"
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

func initializeUpdateUserAccountDS() (ds *inmemory.InMemoryDataStore, wm webmachine.WebMachine) {
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
    ds.CreateUserAccount(&dm.User{
        Role: dm.ROLE_STANDARD,
        Name: "Thomas Jefferson",
        Username: "secondpresident",
        Email: "thomas@jefferson.com",
        PhoneNumber: "+1-401-555-5555",
        Address: "Virginia",
        AllowLogin: true,
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
    authentication.GenerateNewAccessKey(ds, ja.Id, "")
    wm = webmachine.NewWebMachine()
    wm.AddRouteHandler(account.NewUpdateAccountRequestHandler(ds, ds))
    return
}

func TestUpdateUserAccount1(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser := gw
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "GW")
    jsonobj.Set("email", "gw@gwu.edu")
    jsonobj.Set("address", "Pre-White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/" + otherUser.Id, bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
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
    user := new(dm.User)
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    user.InitFromJSONObject(obj.GetAsObject("result"))
    if err != nil {
        t.Error("Error while unmarshaling JSON: ", err.String())
    }
    if obj.GetAsString("status") != "success" {
        t.Error("Expected status = \"success\", but was \"", obj.GetAsString("status"), "\"")
    }
    if user.Name != otherUser.Name {
        t.Error("Expected name = \"", otherUser.Name, "\", but was ", user.Name)
    }
    if user.Username != otherUser.Username {
        t.Error("Expected username = \"", otherUser.Username, "\", but was ", user.Username)
    }
    if user.Email != otherUser.Email {
        t.Error("Expected email = \"", otherUser.Email, "\", but was ", user.Email)
    }
    if user.PhoneNumber != otherUser.PhoneNumber {
        t.Error("Expected phone_number = \"", otherUser.PhoneNumber, "\", but was ", user.PhoneNumber)
    }
    if user.Address != otherUser.Address {
        t.Error("Expected address = \"", otherUser.Address, "\", but was ", user.Address)
    }
    if user.Role != otherUser.Role {
        t.Error("Expected role = ", otherUser.Role, " but was ", user.Role)
    }
    if user.Id != otherUser.Id {
        t.Error("Expected id to be ", otherUser.Id, ", but was ", user.Id)
    }
    if theuser, err := ds.RetrieveUserAccountById(otherUser.Id); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by id ", otherUser.Id)
        }
        if err != nil {
            t.Error("Error trying to find user account by id: ", err.String())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(otherUser.Username); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by username ", otherUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.String())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(otherUser.Email, nil, 1000); err != nil || len(theusers) != 1 {
        if len(theusers) != 1 {
            t.Error("Found ", len(theusers), " User accounts by email for ", otherUser.Email, " rather than 1: ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.String())
        }
    }
}


func TestUpdateUserAccountAsAdmin(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("thirdpresident")
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "John A")
    jsonobj.Set("email", "ja@adamsu.edu")
    jsonobj.Set("address", "White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/" + otherUser.Id, bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != http.StatusOK {
        t.Error("Expected ", http.StatusOK, " status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    user := new(dm.User)
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    user.InitFromJSONObject(obj.GetAsObject("result"))
    if err != nil {
        t.Error("Error while unmarshaling JSON: ", err.String())
    }
    if obj.GetAsString("status") != "success" {
        t.Error("Expected status = \"success\", but was \"", obj.GetAsString("status"), "\"")
    }
    if user.Name != otherUser.Name {
        t.Error("Expected name = \"", otherUser.Name, "\", but was ", user.Name)
    }
    if user.Username != otherUser.Username {
        t.Error("Expected username = \"", otherUser.Username, "\", but was ", user.Username)
    }
    if user.Email != otherUser.Email {
        t.Error("Expected email = \"", otherUser.Email, "\", but was ", user.Email)
    }
    if user.PhoneNumber != otherUser.PhoneNumber {
        t.Error("Expected phone_number = \"", otherUser.PhoneNumber, "\", but was ", user.PhoneNumber)
    }
    if user.Address != otherUser.Address {
        t.Error("Expected address = \"", otherUser.Address, "\", but was ", user.Address)
    }
    if user.Role != otherUser.Role {
        t.Error("Expected role = ", otherUser.Role, " but was ", user.Role)
    }
    if user.Id != otherUser.Id {
        t.Error("Expected id to be ", otherUser.Id, ", but was ", user.Id)
    }
    if theuser, err := ds.RetrieveUserAccountById(otherUser.Id); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by id ", otherUser.Id)
        }
        if err != nil {
            t.Error("Error trying to find user account by id: ", err.String())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(otherUser.Username); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by username ", otherUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.String())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(otherUser.Email, nil, 1000); err != nil || len(theusers) != 1 {
        if len(theusers) != 1 {
            t.Error("Found ", len(theusers), " User accounts by email for ", otherUser.Email, " rather than 1: ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.String())
        }
    }
}

func TestUpdateUserAccountAsNonAdminSelf(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    ja, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(ja.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("thirdpresident")
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "John A")
    jsonobj.Set("email", "ja@adamsu.edu")
    jsonobj.Set("address", "White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/" + otherUser.Id, bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != http.StatusOK {
        t.Error("Expected ", http.StatusOK, " status code but received ", resp.StatusCode)
    }
    if resp.Header().Get("Content-Type") != req.Header.Get("Accept") {
        t.Error("Expected Content-Type \"", req.Header.Get("Accept"), "\" but received ", resp.Header().Get("Content-Type"))
    }
    user := new(dm.User)
    obj := jsonhelper.NewJSONObject()
    err := json.Unmarshal(resp.Buffer.Bytes(), &obj)
    user.InitFromJSONObject(obj.GetAsObject("result"))
    if err != nil {
        t.Error("Error while unmarshaling JSON: ", err.String())
    }
    if obj.GetAsString("status") != "success" {
        t.Error("Expected status = \"success\", but was \"", obj.GetAsString("status"), "\"")
    }
    if user.Name != otherUser.Name {
        t.Error("Expected name = \"", otherUser.Name, "\", but was ", user.Name)
    }
    if user.Username != otherUser.Username {
        t.Error("Expected username = \"", otherUser.Username, "\", but was ", user.Username)
    }
    if user.Email != otherUser.Email {
        t.Error("Expected email = \"", otherUser.Email, "\", but was ", user.Email)
    }
    if user.PhoneNumber != otherUser.PhoneNumber {
        t.Error("Expected phone_number = \"", otherUser.PhoneNumber, "\", but was ", user.PhoneNumber)
    }
    if user.Address != otherUser.Address {
        t.Error("Expected address = \"", otherUser.Address, "\", but was ", user.Address)
    }
    if user.Role != otherUser.Role {
        t.Error("Expected role = ", otherUser.Role, " but was ", user.Role)
    }
    if user.Id != otherUser.Id {
        t.Error("Expected id to be ", otherUser.Id, ", but was ", user.Id)
    }
    if theuser, err := ds.RetrieveUserAccountById(otherUser.Id); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by id ", otherUser.Id)
        }
        if err != nil {
            t.Error("Error trying to find user account by id: ", err.String())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(otherUser.Username); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by username ", otherUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.String())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(otherUser.Email, nil, 1000); err != nil || len(theusers) != 1 {
        if len(theusers) != 1 {
            t.Error("Found ", len(theusers), " User accounts by email for ", otherUser.Email, " rather than 1: ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.String())
        }
    }
}

func TestUpdateUserAccountAsNonAdminForOtherUser(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    ja, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(ja.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("secondpresident")
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "Tom J")
    jsonobj.Set("email", "tj@jeffersonacademy.edu")
    jsonobj.Set("address", "White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/" + otherUser.Id, bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != http.StatusForbidden {
        t.Error("Expected ", http.StatusForbidden, " status code but received ", resp.StatusCode)
    }
}

func TestUpdateUserAccountMissingId(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1)
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("secondpresident")
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "Tom J")
    jsonobj.Set("email", "tj@jeffersonacademy.edu")
    jsonobj.Set("address", "White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != http.StatusBadRequest {
        t.Error("Expected ", http.StatusBadRequest, " status code but received ", resp.StatusCode)
    }
}

func TestUpdateUserAccountInvalidUserId(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1)
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("secondpresident")
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "Tom J")
    jsonobj.Set("email", "tj@jeffersonacademy.edu")
    jsonobj.Set("address", "White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/sdflsjflsjfslf", bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    apiutil.NewSigner(accessKey.Id, accessKey.PrivateKey).SignRequest(req, 0)
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != http.StatusNotFound {
        t.Error("Expected ", http.StatusNotFound, " status code but received ", resp.StatusCode)
    }
}

func TestUpdateUserAccountMissingSignature(t *testing.T) {
    ds, wm := initializeUpdateUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    otherUser, _ := ds.FindUserAccountByUsername("secondpresident")
    anobj, _ := jsonhelper.Marshal(otherUser)
    jsonobj := anobj.(jsonhelper.JSONObject)
    jsonobj.Set("name", "GW")
    jsonobj.Set("email", "gw@gwu.edu")
    jsonobj.Set("address", "Pre-White House")
    otherUser = new(dm.User)
    otherUser.InitFromJSONObject(jsonobj)
    jsonbuf, _ := json.Marshal(jsonobj)
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/update/" + gw.Id, bytes.NewBuffer(jsonbuf))
    req.Header.Set("Content-Type", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept", webmachine.MIME_TYPE_JSON+"; charset=utf-8")
    req.Header.Set("Accept-Charset", "utf-8")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Accept-Language", "en-us")
    req.Header.Set("Connection", "close")
    resp := webmachine.NewMockResponseWriter(req)
    wm.ServeHTTP(resp, req)
    if resp.StatusCode != http.StatusUnauthorized {
        t.Error("Expected ", http.StatusUnauthorized, " status code but received ", resp.StatusCode)
    }
}
