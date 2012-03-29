package contacts

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    "github.com/pomack/dsocial.go/backend/authentication"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "io"
    "net/http"
    "net/url"
    "time"
)

type ViewContactRequestHandler struct {
    wm.DefaultRequestHandler
    ds         acct.DataStore
    authDS     authentication.DataStore
    contactsDS bc.DataStoreService
}

type ViewContactContext interface {
    AuthUser() *dm.User
    SetAuthUser(user *dm.User)
    User() *dm.User
    SetUser(user *dm.User)
    ContactId() string
    SetContactId(contactId string)
    Contact() *dm.Contact
    SetContact(contact *dm.Contact)
    Result() jsonhelper.JSONObject
    SetResult(result jsonhelper.JSONObject)
    LastModified() time.Time
    SetLastModified(lastModified time.Time)
    SetETag(etag string)
    ETag() string
}

type viewContactContext struct {
    authUser     *dm.User
    user         *dm.User
    contact      *dm.Contact
    contactId    string
    result       jsonhelper.JSONObject
    lastModified time.Time
    etag         string
}

func NewViewContactContext() ViewContactContext {
    return new(viewContactContext)
}

func (p *viewContactContext) AuthUser() *dm.User {
    return p.authUser
}

func (p *viewContactContext) SetAuthUser(user *dm.User) {
    p.authUser = user
}

func (p *viewContactContext) User() *dm.User {
    return p.user
}

func (p *viewContactContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *viewContactContext) Contact() *dm.Contact {
    return p.contact
}

func (p *viewContactContext) SetContact(contact *dm.Contact) {
    p.contact = contact
}

func (p *viewContactContext) ContactId() string {
    return p.contactId
}

func (p *viewContactContext) SetContactId(contactId string) {
    p.contactId = contactId
}

func (p *viewContactContext) Result() jsonhelper.JSONObject {
    return p.result
}

func (p *viewContactContext) SetResult(result jsonhelper.JSONObject) {
    p.result = result
}

func (p *viewContactContext) LastModified() time.Time {
    return p.lastModified
}

func (p *viewContactContext) SetLastModified(lastModified time.Time) {
    p.lastModified = lastModified
}

func (p *viewContactContext) ETag() string {
    return p.etag
}

func (p *viewContactContext) SetETag(etag string) {
    p.etag = etag
}

func NewViewContactRequestHandler(ds acct.DataStore, authDS authentication.DataStore, contactsDS bc.DataStoreService) *ViewContactRequestHandler {
    return &ViewContactRequestHandler{ds: ds, authDS: authDS, contactsDS: contactsDS}
}

func (p *ViewContactRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) ViewContactContext {
    if vcc, ok := cxt.(ViewContactContext); ok {
        return vcc
    }
    return NewViewContactContext()
}

func (p *ViewContactRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/u/<uid>/contacts/list
    // /u/<uid>/contacts/list
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 9 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "u" && path[6] == "contacts" && path[7] == "view" {
            return p
        }
    }
    if pathLen == 6 {
        if path[0] == "" && path[1] == "u" && path[3] == "contacts" && path[4] == "view" {
            return p
        }
    }
    return nil
}

func (p *ViewContactRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    spac := p.GenerateContext(req, cxt)
    return req, spac
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *ViewContactRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    vcc := cxt.(ViewContactContext)
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    contactId := ""
    if pathLen == 9 {
        contactId = path[8]
    } else if pathLen == 6 {
        contactId = path[5]
    }
    vcc.SetContactId(contactId)
    if contactId == "" || vcc.User() == nil || vcc.User().Id == "" {
        return false, req, cxt, 0, nil
    }
    contact, _, err := p.contactsDS.RetrieveDsocialContact(vcc.User().Id, contactId)
    vcc.SetContact(contact)
    if contact != nil {
        vcc.SetETag(contact.Etag)
        if contact.ModifiedAt > 0 {
            vcc.SetLastModified(time.Unix(contact.ModifiedAt, 0).UTC())
        }
    } else {
        vcc.SetETag("")
        vcc.SetLastModified(time.Time{})
    }
    httpStatus := 0
    if err != nil {
        httpStatus = http.StatusInternalServerError
    }
    return contact != nil, req, cxt, httpStatus, err
}

func (p *ViewContactRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, error) {
    return []string{wm.GET, wm.HEAD}, req, cxt, 0, nil
}

func (p *ViewContactRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, error) {
    vcc := cxt.(ViewContactContext)
    hasSignature, authUserId, _, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if authUserId != "" {
        authUser, _ := p.ds.RetrieveUserAccountById(authUserId)
        vcc.SetAuthUser(authUser)
    }
    userId := apiutil.UserIdFromRequestUrl(req)
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        vcc.SetUser(user)
    }
    return true, "", req, cxt, 0, nil
}

func (p *ViewContactRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    vcc := cxt.(ViewContactContext)
    if vcc.AuthUser() != nil && vcc.AuthUser().Accessible() && vcc.User() != nil && vcc.User().Accessible() && vcc.AuthUser().Id == vcc.User().Id {
        return false, req, cxt, 0, nil
    }
    // Cannot find user with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *ViewContactRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *ViewContactRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}
*/

func (p *ViewContactRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, error) {
    genFunc := func() (jsonhelper.JSONObject, time.Time, string, int, http.Header) {
        vcc := cxt.(ViewContactContext)
        jsonObj := vcc.Result()
        headers := apiutil.AddNoCacheHeaders(nil)
        return jsonObj, vcc.LastModified(), vcc.ETag(), http.StatusOK, headers
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandlerWithGenerator(genFunc, time.Time{}, "")}, req, cxt, 0, nil
}

func (p *ViewContactRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, error) {
    arr := []wm.MediaTypeInputHandler{
        apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body()),
        apiutil.NewUrlEncodedMediaTypeInputHandler("", "", p),
    }
    return arr, req, cxt, 0, nil
}

/*
func (p *ViewContactRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ViewContactRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ViewContactRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ViewContactRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewContactRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ViewContactRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *ViewContactRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *ViewContactRequestHandler) LastModified(req wm.Request, cxt wm.Context) (time.Time, wm.Request, wm.Context, int, error) {
    vcc := cxt.(ViewContactContext)
    return vcc.LastModified(), req, cxt, 0, nil
}

/*
func (p *ViewContactRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *ViewContactRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, error) {
    vcc := cxt.(ViewContactContext)
    return vcc.ETag(), req, cxt, 0, nil
}

/*
func (p *ViewContactRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ViewContactRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *ViewContactRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *ViewContactRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    vcc := cxt.(ViewContactContext)
    return p.HandleInputHandlerAfterSetup(vcc)
}

func (p *ViewContactRequestHandler) HandleUrlEncodedInputHandler(req wm.Request, cxt wm.Context, inputObj url.Values) (int, http.Header, io.WriterTo) {
    vcc := cxt.(ViewContactContext)
    return p.HandleInputHandlerAfterSetup(vcc)
}

func (p *ViewContactRequestHandler) HandleInputHandlerAfterSetup(cxt ViewContactContext) (int, http.Header, io.WriterTo) {
    obj := jsonhelper.NewJSONObject()
    contactObj, _ := jsonhelper.Marshal(cxt.Contact())
    obj.Set("contact", contactObj)
    obj.Set("type", "contact")
    cxt.SetResult(obj)
    return 0, nil, nil
}
