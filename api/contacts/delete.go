package contacts

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    auth "github.com/pomack/dsocial.go/backend/authentication"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    "os"
    "time"
)

type DeleteContactRequestHandler struct {
    wm.DefaultRequestHandler
    ds  acct.DataStore
    authDS auth.DataStore
    contactsDS bc.DataStoreService
}

type DeleteContactContext interface {
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
    ETag() string
    SetETag(etag string)
    LastModified() *time.Time
    SetLastModified(lastModified *time.Time)
    MarkAsDeleted()
    Deleted() bool
}

type deleteContactContext struct {
    authUser        *dm.User
    user            *dm.User
    contact         *dm.Contact
    lastModified    *time.Time
    etag            string
    contactId       string
    deleted         bool
    result          jsonhelper.JSONObject
}

func NewDeleteContactContext() DeleteContactContext {
    return new(deleteContactContext)
}

func (p *deleteContactContext) AuthUser() *dm.User {
    return p.authUser
}

func (p *deleteContactContext) SetAuthUser(authUser *dm.User) {
    p.authUser = authUser
}

func (p *deleteContactContext) User() *dm.User {
    return p.user
}

func (p *deleteContactContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *deleteContactContext) ContactId() string {
    return p.contactId
}

func (p *deleteContactContext) SetContactId(contactId string) {
    p.contactId = contactId
}

func (p *deleteContactContext) Contact() *dm.Contact {
    return p.contact
}

func (p *deleteContactContext) SetContact(contact *dm.Contact) {
    p.contact = contact
}

func (p *deleteContactContext) Result() jsonhelper.JSONObject {
    return p.result
}

func (p *deleteContactContext) SetResult(result jsonhelper.JSONObject) {
    p.result = result
}

func (p *deleteContactContext) ETag() string {
    return p.etag
}

func (p *deleteContactContext) SetETag(etag string) {
    p.etag = etag
}

func (p *deleteContactContext) LastModified() *time.Time {
    return p.lastModified
}

func (p *deleteContactContext) SetLastModified(lastModified *time.Time) {
    p.lastModified = lastModified
}

func (p *deleteContactContext) MarkAsDeleted() {
    p.deleted = true
}

func (p *deleteContactContext) Deleted() bool {
    return p.deleted
}


func NewDeleteContactRequestHandler(ds acct.DataStore, authDS auth.DataStore, contactsDS bc.DataStoreService) *DeleteContactRequestHandler {
    return &DeleteContactRequestHandler{ds: ds, authDS: authDS, contactsDS: contactsDS}
}

func (p *DeleteContactRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) DeleteContactContext {
    if dcc, ok := cxt.(DeleteContactContext); ok {
        return dcc
    }
    return NewDeleteContactContext()
}

func (p *DeleteContactRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/account/(user|consumer|external_user)/delete/(id)
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 9 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "u" && path[6] == "contacts" && path[7] == "delete" {
            return p
        }
    }
    if pathLen == 6 {
        if path[0] == "" && path[1] == "u" && path[3] == "contacts" && path[4] == "delete" {
            return p
        }
    }
    return nil
}

func (p *DeleteContactRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    dcc := p.GenerateContext(req, cxt)
    return req, dcc
}

/*
func (p *CreateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *DeleteContactRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    dcc := cxt.(DeleteContactContext)
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
    dcc.SetContactId(contactId)
    if contactId == "" || dcc.User() == nil || dcc.User().Id == "" {
        return false, req, cxt, 0, nil
    }
    contact, _, err := p.contactsDS.RetrieveDsocialContact(dcc.User().Id, contactId)
    dcc.SetContact(contact)
    if contact != nil {
        dcc.SetETag(contact.Etag)
        if contact.ModifiedAt > 0 {
            dcc.SetLastModified(time.SecondsToUTC(contact.ModifiedAt))
        }
    } else {
        dcc.SetETag("")
        dcc.SetLastModified(nil)
    }
    httpStatus := 0
    if err != nil {
        httpStatus = http.StatusInternalServerError
    }
    return contact != nil, req, cxt, httpStatus, err
}

func (p *DeleteContactRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST, wm.DELETE}, req, cxt, 0, nil
}

func (p *DeleteContactRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    dcc := cxt.(DeleteContactContext)
    hasSignature, authUserId, _, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if authUserId != "" {
        authUser, _ := p.ds.RetrieveUserAccountById(authUserId)
        dcc.SetAuthUser(authUser)
    }
    userId := apiutil.UserIdFromRequestUrl(req)
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        dcc.SetUser(user)
    }
    return true, "", req, cxt, 0, nil
}


func (p *DeleteContactRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    dcc := cxt.(DeleteContactContext)
    if dcc.AuthUser() != nil && dcc.AuthUser().Accessible() && dcc.User() != nil && dcc.User().Accessible() && (dcc.AuthUser().Id == dcc.User().Id || dcc.AuthUser().Role == dm.ROLE_ADMIN) {
        return false, req, cxt, 0, nil
    }
    // Cannot find user or consumer with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *DeleteContactRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *DeleteContactRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    dcc := cxt.(DeleteContactContext)
    _, err := p.contactsDS.DeleteDsocialContact(dcc.User().Id, dcc.ContactId())
    if err != nil {
        return false, req, cxt, http.StatusInternalServerError, err
    }
    return true, req, cxt, 0, nil
}

/*
func (p *DeleteContactRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/


func (p *DeleteContactRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    _, req, cxt, httpCode, httpError := p.DeleteResource(req, cxt)
    return req, cxt, httpCode, nil, nil, httpError
}


func (p *DeleteContactRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    dcc := cxt.(DeleteContactContext)
    obj := dcc.Result()
    lastModified := dcc.LastModified()
    etag := dcc.ETag()
    var jsonObj jsonhelper.JSONObject
    if obj != nil {
        theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
        jsonObj, _ = theobj.(jsonhelper.JSONObject)
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, lastModified, etag)}, req, dcc, 0, nil
}

func (p *DeleteContactRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    arr := []wm.MediaTypeInputHandler{apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body())}
    return arr, req, cxt, 0, nil
}

/*
func (p *DeleteContactRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
  return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteContactRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *DeleteContactRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/


func (p *DeleteContactRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}


func (p *DeleteContactRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    dcc := cxt.(DeleteContactContext)
    var err os.Error
    if !dcc.Deleted() {
        _, req, cxt, _, err = p.DeleteResource(req, cxt)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(err.String(), nil, http.StatusInternalServerError, nil)
    }
    obj := jsonhelper.NewJSONObject()
    obj.Set("type", "contact")
    obj.Set("contact", dcc.Contact())
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    return apiutil.OutputJSONObject(jsonObj, nil, "", 0, nil)
}
