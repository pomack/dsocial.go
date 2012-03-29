package account

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    auth "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    "os"
    "strings"
    "time"
)

type ViewAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds     acct.DataStore
    authDS auth.DataStore
}

type ViewAccountContext interface {
    Type() string
    SetType(theType string)
    User() *dm.User
    SetUser(user *dm.User)
    Consumer() *dm.Consumer
    SetConsumer(consumer *dm.Consumer)
    ExternalUser() *dm.ExternalUser
    SetExternalUser(externalUser *dm.ExternalUser)
    LastModified() *time.Time
    ETag() string
    ToObject() interface{}
    RequestingUser() *dm.User
    SetRequestingUser(user *dm.User)
    RequestingConsumer() *dm.Consumer
    SetRequestingConsumer(consumer *dm.Consumer)
}

type viewAccountContext struct {
    theType            string
    user               *dm.User
    consumer           *dm.Consumer
    externalUser       *dm.ExternalUser
    requestingUser     *dm.User
    requestingConsumer *dm.Consumer
}

func NewViewAccountContext() ViewAccountContext {
    return new(viewAccountContext)
}

func (p *viewAccountContext) Type() string {
    return p.theType
}

func (p *viewAccountContext) SetType(theType string) {
    p.theType = theType
}

func (p *viewAccountContext) User() *dm.User {
    return p.user
}

func (p *viewAccountContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *viewAccountContext) Consumer() *dm.Consumer {
    return p.consumer
}

func (p *viewAccountContext) SetConsumer(consumer *dm.Consumer) {
    p.consumer = consumer
}

func (p *viewAccountContext) ExternalUser() *dm.ExternalUser {
    return p.externalUser
}

func (p *viewAccountContext) SetExternalUser(externalUser *dm.ExternalUser) {
    p.externalUser = externalUser
}

func (p *viewAccountContext) LastModified() *time.Time {
    var lastModified *time.Time
    if p.user != nil && p.user.ModifiedAt != 0 {
        lastModified = time.SecondsToUTC(p.user.ModifiedAt)
    } else if p.consumer != nil && p.consumer.ModifiedAt != 0 {
        lastModified = time.SecondsToUTC(p.consumer.ModifiedAt)
    } else if p.externalUser != nil && p.externalUser.ModifiedAt != 0 {
        lastModified = time.SecondsToUTC(p.externalUser.ModifiedAt)
    }
    return lastModified
}

func (p *viewAccountContext) ToObject() interface{} {
    var user interface{}
    if p.user != nil {
        user = p.user
    } else if p.consumer != nil {
        user = p.consumer
    } else if p.externalUser != nil {
        user = p.externalUser
    }
    return user
}

func (p *viewAccountContext) ETag() string {
    var etag string
    if p.user != nil {
        etag = p.user.Etag
    } else if p.consumer != nil {
        etag = p.consumer.Etag
    } else if p.externalUser != nil {
        etag = p.externalUser.Etag
    }
    return etag
}

func (p *viewAccountContext) RequestingUser() *dm.User {
    return p.requestingUser
}

func (p *viewAccountContext) SetRequestingUser(user *dm.User) {
    p.requestingUser = user
}

func (p *viewAccountContext) RequestingConsumer() *dm.Consumer {
    return p.requestingConsumer
}

func (p *viewAccountContext) SetRequestingConsumer(consumer *dm.Consumer) {
    p.requestingConsumer = consumer
}

func NewViewAccountRequestHandler(ds acct.DataStore, authDS auth.DataStore) *ViewAccountRequestHandler {
    return &ViewAccountRequestHandler{ds: ds, authDS: authDS}
}

func (p *ViewAccountRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) ViewAccountContext {
    if dac, ok := cxt.(ViewAccountContext); ok {
        return dac
    }
    return NewViewAccountContext()
}

func (p *ViewAccountRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/account/(user|consumer|external_user)/view/(id)
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen >= 8 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "account" && path[6] == "view" {
            switch path[5] {
            case "user", "consumer", "external_user":
                return p
            }
        }
    }
    return nil
}

func (p *ViewAccountRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    vac := p.GenerateContext(req, cxt)
    path := req.URLParts()
    pathLen := len(path)
    if pathLen >= 8 {
        vac.SetType(path[5])
        var id string
        if path[pathLen-1] == "" {
            id = strings.Join(path[7:pathLen-1], "/")
        } else {
            id = strings.Join(path[7:], "/")
        }
        switch vac.Type() {
        case "user":
            user, _ := p.ds.RetrieveUserAccountById(id)
            vac.SetUser(user)
        case "consumer":
            consumer, _ := p.ds.RetrieveConsumerAccountById(id)
            vac.SetConsumer(consumer)
        case "external_user":
            externalUser, _ := p.ds.RetrieveExternalUserAccountById(id)
            vac.SetExternalUser(externalUser)
        }
    }
    return req, vac
}

/*
func (p *CreateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *ViewAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    vac := cxt.(ViewAccountContext)
    return vac.ToObject() != nil, req, cxt, 0, nil
}

func (p *ViewAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.GET, wm.HEAD}, req, cxt, 0, nil
}

func (p *ViewAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    vac := cxt.(ViewAccountContext)
    hasSignature, userId, consumerId, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        vac.SetRequestingUser(user)
    }
    if consumerId != "" {
        consumer, _ := p.ds.RetrieveConsumerAccountById(consumerId)
        vac.SetRequestingConsumer(consumer)
    }
    return true, "", req, cxt, 0, nil
}

func (p *ViewAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    vac := cxt.(ViewAccountContext)
    if vac.RequestingUser() != nil && vac.RequestingUser().Accessible() && (vac.RequestingUser().Role == dm.ROLE_ADMIN || (vac.User() != nil && vac.RequestingUser().Id == vac.User().Id)) {
        return false, req, cxt, 0, nil
    }
    // Cannot find user or consumer with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *ViewAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *ViewAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    vac := cxt.(ViewAccountContext)
    obj := vac.ToObject()
    lastModified := vac.LastModified()
    etag := vac.ETag()
    var jsonObj jsonhelper.JSONObject
    if obj != nil {
        theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
        jsonObj, _ = theobj.(jsonhelper.JSONObject)
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, lastModified, etag)}, req, vac, 0, nil
}

/*
func (p *ViewAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    return []wm.MediaTypeInputHandler{}, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewAccountRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewAccountRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewAccountRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewAccountRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
  return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewAccountRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *ViewAccountRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *ViewAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    vac := cxt.(ViewAccountContext)
    return vac.LastModified(), req, cxt, 0, nil
}

/*
func (p *ViewAccountRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *ViewAccountRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    vac := cxt.(ViewAccountContext)
    return vac.ETag(), req, cxt, 0, nil
}

/*
func (p *ViewAccountRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *ViewAccountRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *ViewAccountRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *ViewAccountRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    vac := cxt.(ViewAccountContext)

    obj := vac.ToObject()
    var err os.Error
    if err != nil {
        return apiutil.OutputErrorMessage(err.String(), nil, http.StatusInternalServerError, nil)
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    return apiutil.OutputJSONObject(jsonObj, vac.LastModified(), vac.ETag(), 0, nil)
}
