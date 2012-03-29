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

func initializeDeleteUserAccountDS() (ds *inmemory.InMemoryDataStore, wm webmachine.WebMachine) {
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
    ds.CreateUserAccount(&dm.User{
        Role:        dm.ROLE_TECHNICAL_SUPPORT,
        Name:        "John Adams",
        Username:    "thirdpresident",
        Email:       "john@adams.com",
        PhoneNumber: "+1-402-555-5555",
        Address:     "Boston, MA",
        AllowLogin:  true,
    })
    authentication.GenerateNewAccessKey(ds, gw.Id, "")
    wm = webmachine.NewWebMachine()
    wm.AddRouteHandler(account.NewDeleteAccountRequestHandler(ds, ds))
    return
}

func TestDeleteUserAccount(t *testing.T) {
    ds, wm := initializeDeleteUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1000)
    if len(accessKeys) == 0 {
        t.Error("Expected to find at least one access key stored.")
    }
    accessKey := accessKeys[0]
    oldUser, _ := ds.FindUserAccountByUsername("thirdpresident")
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/delete/"+oldUser.Id, nil)
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
        t.Error("Error while unmarshaling JSON: ", err.Error())
    }
    if obj.GetAsString("status") != "success" {
        t.Error("Expected status = \"success\", but was \"", obj.GetAsString("status"), "\"")
    }
    if user.Name != oldUser.Name {
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
    if user.Role != dm.ROLE_TECHNICAL_SUPPORT {
        t.Error("Expected role = ", dm.ROLE_TECHNICAL_SUPPORT, " but was ", user.Role)
    }
    if user.Id != oldUser.Id {
        t.Error("Expected id to be ", oldUser.Id, ", but was ", user.Id)
    }
    if theuser, err := ds.RetrieveUserAccountById(oldUser.Id); err != nil || theuser != nil {
        if theuser != nil {
            t.Error("User account by id still finds ", oldUser.Id)
        }
        if err != nil {
            t.Error("Error trying to find user account by id: ", err.Error())
        }
    }
    if theuser, err := ds.FindUserAccountByUsername(oldUser.Username); err != nil || theuser != nil {
        if theuser != nil {
            t.Error("User account by username still finds ", oldUser.Username)
        }
        if err != nil {
            t.Error("Error trying to find user account by username: ", err.Error())
        }
    }
    if theusers, _, err := ds.FindUserAccountsByEmail(oldUser.Email, nil, 1000); err != nil || len(theusers) > 0 {
        if len(theusers) > 0 {
            t.Error("User accounts by email still finds ", oldUser.Email, ": ", theusers)
        }
        if err != nil {
            t.Error("Error trying to find user accounts by email: ", err.Error())
        }
    }
}

func TestDeleteUserAccountMissingId(t *testing.T) {
    ds, wm := initializeDeleteUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1)
    accessKey := accessKeys[0]
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/delete/", nil)
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

func TestDeleteUserAccountInvalidUserId(t *testing.T) {
    ds, wm := initializeDeleteUserAccountDS()
    gw, _ := ds.FindUserAccountByUsername("firstpresident")
    accessKeys, _, _ := ds.RetrieveUserKeys(gw.Id, nil, 1)
    accessKey := accessKeys[0]
    req, _ := http.NewRequest(webmachine.POST, "http://localhost/api/v1/json/account/user/delete/sdflsjflsjfslf", nil)
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
