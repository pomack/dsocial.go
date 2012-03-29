package account_test

import (
    "encoding/json"
    "github.com/pomack/dsocial.go/api/account"
    "github.com/pomack/dsocial.go/api/apiutil"
    "github.com/pomack/dsocial.go/backend/authentication"
    "github.com/pomack/dsocial.go/backend/datastore/inmemory"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/webmachine.go/webmachine"
    "net/http"
    "testing"
)

func initializeViewUserAccountDS() (ds *inmemory.InMemoryDataStore, wm webmachine.WebMachine) {
    ds = inmemory.NewInMemoryDataStore()
    gw, _ := ds.CreateUserAccount(&dm.User{
        Role:        dm.ROLE_ADMIN,
        Name:        "George Washington",
        Username:    "firstpresident",
        Email:       "george@washington.com",
        PhoneNumber: "+1-405-555-5555",
        Address:     "Valley Forge",
        AllowLogin:  true,
    })
    ds.CreateUserAccount(&dm.User{
        Role:        dm.ROLE_STANDARD,
        Name:        "Thomas Jefferson",
        Username:    "secondpresident",
        Email:       "thomas@jefferson.com",
        PhoneNumber: "+1-401-555-5555",
        Address:     "Virginia",
        AllowLogin:  true,
    })
    ja, _ := ds.CreateUserAccount(&dm.User{
        Role:        dm.ROLE_TECHNICAL_SUPPORT,
        Name:        "John Adams",
        Username:    "thirdpresident",
        Email:       "john@adams.com",
        PhoneNumber: "+1-402-555-5555",
        Address:     "Boston, MA",
        AllowLogin:  true,
    })
    authentication.GenerateNewAccessKey(ds, gw.Id, "")
    authentication.GenerateNewAccessKey(ds, ja.Id, "")
    wm = webmachine.NewWebMachine()
    wm.AddRouteHandler(account.NewViewAccountRequestHandler(ds, ds))
    return
}

func TestViewUserAccount(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser := gw
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/"+otherUser.Id, nil)
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
        t.Error("Error while unmarshaling JSON: ", err.Error())
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
            t.Error("Error trying to find user account by id: ", err.Error())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(otherUser.Username); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by username ", otherUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.Error())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(otherUser.Email, nil, 1000); err != nil || len(theusers) != 1 {
        if len(theusers) != 1 {
            t.Error("Found ", len(theusers), " User accounts by email for ", otherUser.Email, " rather than 1: ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.Error())
        }
    }
}

func TestViewUserAccountAsAdmin(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("thirdpresident")
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/"+otherUser.Id, nil)
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
        t.Error("Error while unmarshaling JSON: ", err.Error())
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
            t.Error("Error trying to find user account by id: ", err.Error())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(otherUser.Username); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by username ", otherUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.Error())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(otherUser.Email, nil, 1000); err != nil || len(theusers) != 1 {
        if len(theusers) != 1 {
            t.Error("Found ", len(theusers), " User accounts by email for ", otherUser.Email, " rather than 1: ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.Error())
        }
    }
}

func TestViewUserAccountAsNonAdminSelf(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    ja, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(ja.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("thirdpresident")
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/"+otherUser.Id, nil)
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
        t.Error("Error while unmarshaling JSON: ", err.Error())
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
            t.Error("Error trying to find user account by id: ", err.Error())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(otherUser.Username); err != nil || theuser == nil {
        if theuser == nil {
            t.Error("Unable to find User account by username ", otherUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.Error())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(otherUser.Email, nil, 1000); err != nil || len(theusers) != 1 {
        if len(theusers) != 1 {
            t.Error("Found ", len(theusers), " User accounts by email for ", otherUser.Email, " rather than 1: ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.Error())
        }
    }
}

func TestViewUserAccountAsNonAdminForOtherUser(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    ja, _ := ds.FindUserAccountByUsername("thirdpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(ja.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    otherUser, _ := ds.FindUserAccountByUsername("secondpresident")
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/"+otherUser.Id, nil)
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

func TestViewUserAccountMissingId(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1)
    accessKey := accessKeys[0]
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/", nil)
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

func TestViewUserAccountInvalidUserId(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1)
    accessKey := accessKeys[0]
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/sdflsjflsjfslf", nil)
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

func TestViewUserAccountMissingSignature(t *testing.T) {
    ds, wm := initializeViewUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    req, _ := http.NewRequest(webmachine.GET, "http://localhost/api/v1/json/account/user/view/"+gw.Id, nil)
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
