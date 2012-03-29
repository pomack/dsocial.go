package contacts

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    "github.com/pomack/dsocial.go/backend/authentication"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    "os"
    "strconv"
    "time"
    "url"
)

type ListContactsRequestHandler struct {
    wm.DefaultRequestHandler
    ds         acct.DataStore
    authDS     authentication.DataStore
    contactsDS bc.DataStoreService
}

type ListContactsContext interface {
    AuthUser() *dm.User
    SetAuthUser(user *dm.User)
    User() *dm.User
    SetUser(user *dm.User)
    ListFrom() string
    SetListFrom(listFrom string)
    ListNext() string
    SetListNext(listNext string)
    Count() int
    SetCount(count int)
    SetFromJSON(obj jsonhelper.JSONObject)
    SetFromUrlEncoded(obj url.Values)
    Result() jsonhelper.JSONObject
    SetResult(result jsonhelper.JSONObject)
}

type listContactsContext struct {
    authUser *dm.User
    user     *dm.User
    listFrom string
    listNext string
    result   jsonhelper.JSONObject
    count    int
}

func NewListContactsContext() ListContactsContext {
    return new(listContactsContext)
}

func (p *listContactsContext) AuthUser() *dm.User {
    return p.authUser
}

func (p *listContactsContext) SetAuthUser(user *dm.User) {
    p.authUser = user
}

func (p *listContactsContext) User() *dm.User {
    return p.user
}

func (p *listContactsContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *listContactsContext) ListFrom() string {
    return p.listFrom
}

func (p *listContactsContext) SetListFrom(listFrom string) {
    p.listFrom = listFrom
}

func (p *listContactsContext) ListNext() string {
    return p.listNext
}

func (p *listContactsContext) SetListNext(listNext string) {
    p.listNext = listNext
}

func (p *listContactsContext) Count() int {
    return p.count
}

func (p *listContactsContext) SetCount(count int) {
    p.count = count
}

func (p *listContactsContext) Result() jsonhelper.JSONObject {
    return p.result
}

func (p *listContactsContext) SetResult(result jsonhelper.JSONObject) {
    p.result = result
}

func (p *listContactsContext) SetFromJSON(obj jsonhelper.JSONObject) {
    p.listFrom = obj.GetAsString("next")
    p.count = obj.GetAsInt("count")
    p.result = nil
}

func (p *listContactsContext) SetFromUrlEncoded(values url.Values) {
    p.listFrom = values.Get("next")
    p.count, _ = strconv.Atoi(values.Get("count"))
    p.result = nil
}

func NewListContactsRequestHandler(ds acct.DataStore, authDS authentication.DataStore, contactsDS bc.DataStoreService) *ListContactsRequestHandler {
    return &ListContactsRequestHandler{ds: ds, authDS: authDS, contactsDS: contactsDS}
}

func (p *ListContactsRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) ListContactsContext {
    if lcc, ok := cxt.(ListContactsContext); ok {
        return lcc
    }
    return NewListContactsContext()
}

func (p *ListContactsRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/u/<uid>/contacts/list
    // /u/<uid>/contacts/list
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 8 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "u" && path[6] == "contacts" && path[7] == "list" {
            return p
        }
    }
    if pathLen == 5 {
        if path[0] == "" && path[1] == "u" && path[3] == "contacts" && path[4] == "list" {
            return p
        }
    }
    return nil
}

func (p *ListContactsRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    spac := p.GenerateContext(req, cxt)
    return req, spac
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *ListContactsRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.GET, wm.HEAD}, req, cxt, 0, nil
}

func (p *ListContactsRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    lcc := cxt.(ListContactsContext)
    hasSignature, authUserId, _, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if authUserId != "" {
        authUser, _ := p.ds.RetrieveUserAccountById(authUserId)
        lcc.SetAuthUser(authUser)
    }
    userId := apiutil.UserIdFromRequestUrl(req)
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        lcc.SetUser(user)
    }
    return true, "", req, cxt, 0, nil
}

func (p *ListContactsRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    lcc := cxt.(ListContactsContext)
    if lcc.AuthUser() != nil && lcc.AuthUser().Accessible() && lcc.User() != nil && lcc.User().Accessible() && lcc.AuthUser().Id == lcc.User().Id {
        return false, req, cxt, 0, nil
    }
    // Cannot find user with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *ListContactsRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *ListContactsRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}
*/

func (p *ListContactsRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    genFunc := func() (jsonhelper.JSONObject, *time.Time, string, int, http.Header) {
        lcc := cxt.(ListContactsContext)
        jsonObj := lcc.Result()
        headers := apiutil.AddNoCacheHeaders(nil)
        return jsonObj, nil, "", http.StatusOK, headers
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandlerWithGenerator(genFunc, nil, "")}, req, cxt, 0, nil
}

func (p *ListContactsRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    arr := []wm.MediaTypeInputHandler{
        apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body()),
        apiutil.NewUrlEncodedMediaTypeInputHandler("", "", p),
    }
    return arr, req, cxt, 0, nil
}

/*
func (p *ListContactsRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ListContactsRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ListContactsRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ListContactsRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ListContactsRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ListContactsRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ListContactsRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ListContactsRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ListContactsRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ListContactsRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ListContactsRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *ListContactsRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *ListContactsRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    lcc := cxt.(ListContactsContext)
    lcc.SetFromJSON(inputObj)
    return p.HandleInputHandlerAfterSetup(lcc)
}

func (p *ListContactsRequestHandler) HandleUrlEncodedInputHandler(req wm.Request, cxt wm.Context, inputObj url.Values) (int, http.Header, io.WriterTo) {
    lcc := cxt.(ListContactsContext)
    lcc.SetFromUrlEncoded(inputObj)
    return p.HandleInputHandlerAfterSetup(lcc)
}

func (p *ListContactsRequestHandler) HandleInputHandlerAfterSetup(cxt ListContactsContext) (int, http.Header, io.WriterTo) {
    contactsDS := p.contactsDS
    user := cxt.User()
    from := cxt.ListFrom()
    dsocialContacts, next, err := contactsDS.ListDsocialContacts(user.Id, from, cxt.Count())
    if err != nil {
        return apiutil.OutputErrorMessage(err.String(), nil, http.StatusInternalServerError, nil)
    }
    obj := jsonhelper.NewJSONObject()
    contactsArr, _ := jsonhelper.Marshal(dsocialContacts)
    nextObj, _ := jsonhelper.Marshal(next)
    obj.Set("contacts", contactsArr)
    obj.Set("type", "contacts")
    obj.Set("next", nextObj)
    cxt.SetResult(obj)
    return 0, nil, nil
}
