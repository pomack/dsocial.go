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

type DeleteAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds  acct.DataStore
    authDS auth.DataStore
}

type DeleteAccountContext interface {
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
    MarkAsDeleted()
    Deleted() bool
}

type deleteAccountContext struct {
    theType      string
    user         *dm.User
    consumer     *dm.Consumer
    externalUser *dm.ExternalUser
    requestingUser     *dm.User
    requestingConsumer *dm.Consumer
    wasDeleted  bool
}

func NewDeleteAccountContext() DeleteAccountContext {
    return new(deleteAccountContext)
}

func (p *deleteAccountContext) Type() string {
    return p.theType
}

func (p *deleteAccountContext) SetType(theType string) {
    p.theType = theType
}

func (p *deleteAccountContext) User() *dm.User {
    return p.user
}

func (p *deleteAccountContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *deleteAccountContext) Consumer() *dm.Consumer {
    return p.consumer
}

func (p *deleteAccountContext) SetConsumer(consumer *dm.Consumer) {
    p.consumer = consumer
}

func (p *deleteAccountContext) ExternalUser() *dm.ExternalUser {
    return p.externalUser
}

func (p *deleteAccountContext) SetExternalUser(externalUser *dm.ExternalUser) {
    p.externalUser = externalUser
}

func (p *deleteAccountContext) LastModified() *time.Time {
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

func (p *deleteAccountContext) ToObject() interface{} {
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

func (p *deleteAccountContext) ETag() string {
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

func (p *deleteAccountContext) RequestingUser() *dm.User {
    return p.requestingUser
}

func (p *deleteAccountContext) SetRequestingUser(user *dm.User) {
    p.requestingUser = user
}

func (p *deleteAccountContext) RequestingConsumer() *dm.Consumer {
    return p.requestingConsumer
}

func (p *deleteAccountContext) SetRequestingConsumer(consumer *dm.Consumer) {
    p.requestingConsumer = consumer
}

func (p *deleteAccountContext) MarkAsDeleted() {
    p.wasDeleted = true
}
func (p *deleteAccountContext) Deleted() bool {
    return p.wasDeleted
}


func NewDeleteAccountRequestHandler(ds acct.DataStore, authDS auth.DataStore) *DeleteAccountRequestHandler {
    return &DeleteAccountRequestHandler{ds: ds, authDS: authDS}
}

func (p *DeleteAccountRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) DeleteAccountContext {
    if dac, ok := cxt.(DeleteAccountContext); ok {
        return dac
    }
    return NewDeleteAccountContext()
}

func (p *DeleteAccountRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/account/(user|consumer|external_user)/delete/
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen >= 8 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "account" && path[6] == "delete" {
            switch path[5] {
            case "user", "consumer", "external_user":
                return p
            }
        }
    }
    return nil
}

func (p *DeleteAccountRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    dac := p.GenerateContext(req, cxt)
    path := req.URLParts()
    pathLen := len(path)
    if pathLen >= 8 {
        dac.SetType(path[5])
        var id string
        if path[pathLen-1] == "" {
            id = strings.Join(path[7:pathLen-1], "/")
        } else {
            id = strings.Join(path[7:], "/")
        }
        switch dac.Type() {
        case "user":
            user, _ := p.ds.RetrieveUserAccountById(id)
            dac.SetUser(user)
        case "consumer":
            consumer, _ := p.ds.RetrieveConsumerAccountById(id)
            dac.SetConsumer(consumer)
        case "external_user":
            externalUser, _ := p.ds.RetrieveExternalUserAccountById(id)
            dac.SetExternalUser(externalUser)
        }
    }
    return req, dac
}

/*
func (p *CreateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *DeleteAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    dac := cxt.(DeleteAccountContext)
    return dac.ToObject() != nil, req, cxt, 0, nil
}

func (p *DeleteAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST, wm.PUT}, req, cxt, 0, nil
}

func (p *DeleteAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    dac := cxt.(DeleteAccountContext)
    hasSignature, userId, consumerId, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, 401, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        dac.SetRequestingUser(user)
    }
    if consumerId != "" {
        consumer, _ := p.ds.RetrieveConsumerAccountById(consumerId)
        dac.SetRequestingConsumer(consumer)
    }
    return true, "", req, cxt, 0, nil
}


func (p *DeleteAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    dac := cxt.(DeleteAccountContext)
    if dac.RequestingUser() == nil || dac.RequestingUser().Role != dm.ROLE_ADMIN {
        // Cannot find user or consumer with specified id
        return true, req, cxt, 0, nil
    }
    return false, req, cxt, 0, nil
}

/*
func (p *DeleteAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteAccountRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteAccountRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *DeleteAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    dac := cxt.(DeleteAccountContext)
    var err os.Error
    if dac.User() != nil {
        _, err = p.ds.DeleteUserAccount(dac.User())
    } else if dac.Consumer() != nil {
        _, err = p.ds.DeleteConsumerAccount(dac.Consumer())
    } else if dac.ExternalUser() != nil {
        _, err = p.ds.DeleteExternalUserAccount(dac.ExternalUser())
    }
    dac.MarkAsDeleted()
    if err != nil {
        return false, req, cxt, 500, err
    }
    return true, req, cxt, 0, nil
}

/*
func (p *DeleteAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *DeleteAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}

/*
func (p *DeleteAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/


func (p *DeleteAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return p.DeleteResource(req, cxt)
}


func (p *DeleteAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    cac := cxt.(DeleteAccountContext)
    obj := cac.ToObject()
    lastModified := cac.LastModified()
    etag := cac.ETag()
    var jsonObj jsonhelper.JSONObject
    if obj != nil {
        theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
        jsonObj, _ = theobj.(jsonhelper.JSONObject)
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, lastModified, etag)}, req, cac, 0, nil
}

func (p *DeleteAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    arr := []wm.MediaTypeInputHandler{apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body())}
    return arr, req, cxt, 0, nil
}

/*
func (p *DeleteAccountRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *DeleteAccountRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *DeleteAccountRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *DeleteAccountRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteAccountRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *DeleteAccountRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
  return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *DeleteAccountRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *DeleteAccountRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *DeleteAccountRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *DeleteAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}

/*
func (p *DeleteAccountRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *DeleteAccountRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *DeleteAccountRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *DeleteAccountRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/


func (p *DeleteAccountRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}


func (p *DeleteAccountRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, writer io.Writer, inputObj jsonhelper.JSONObject) (int, http.Header, os.Error) {
    dac := cxt.(DeleteAccountContext)
    
    var obj interface{}
    var err os.Error
    if !dac.Deleted() {
        _, req, cxt, _, err = p.DeleteResource(req, cxt)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(writer, err.String(), nil, 500, nil)
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    return apiutil.OutputJSONObject(writer, jsonObj, dac.LastModified(), dac.ETag(), 0, nil)
}
