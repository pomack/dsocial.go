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
    //"log"
    "os"
    "strings"
    "time"
)

type UpdateAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds  acct.DataStore
    authDS auth.DataStore
}

type UpdateAccountContext interface {
    SetFromJSON(obj jsonhelper.JSONObject)
    CleanInput(createdByUser *dm.User, originalUser interface{})
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
    OriginalValue() interface{}
    SetOriginalValue(value interface{})
    InputValidated() bool
}

type updateAccountContext struct {
    theType      string
    user         *dm.User
    consumer     *dm.Consumer
    externalUser *dm.ExternalUser
    requestingUser     *dm.User
    requestingConsumer *dm.Consumer
    originalValue interface{}
    inputValidated bool
}

func NewUpdateAccountContext() UpdateAccountContext {
    return new(updateAccountContext)
}

func (p *updateAccountContext) SetFromJSON(obj jsonhelper.JSONObject) {
    p.user = nil
    p.consumer = nil
    p.externalUser = nil
    theType := p.theType
    if theType == "" {
        theType = obj.GetAsString("type")
    }
    switch theType {
    case "user":
        p.user = new(dm.User)
        p.user.InitFromJSONObject(obj)
    case "consumer":
        p.consumer = new(dm.Consumer)
        p.consumer.InitFromJSONObject(obj)
    case "external_user":
        p.externalUser = new(dm.ExternalUser)
        p.externalUser.InitFromJSONObject(obj)
    }
}

func (p *updateAccountContext) CleanInput(createdByUser *dm.User, originalUser interface{}) {
    if p.user != nil {
        p.user.Id = ""
        p.user.CleanFromUser(createdByUser, originalUser.(*dm.User))
    } else if p.consumer != nil {
        p.consumer.Id = ""
        p.consumer.CleanFromUser(createdByUser, originalUser.(*dm.Consumer))
    } else if p.externalUser != nil {
        p.externalUser.Id = ""
        p.externalUser.CleanFromUser(createdByUser, originalUser.(*dm.ExternalUser))
    }
    p.inputValidated = true
}

func (p *updateAccountContext) Type() string {
    return p.theType
}

func (p *updateAccountContext) SetType(theType string) {
    p.theType = theType
}

func (p *updateAccountContext) User() *dm.User {
    return p.user
}

func (p *updateAccountContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *updateAccountContext) Consumer() *dm.Consumer {
    return p.consumer
}

func (p *updateAccountContext) SetConsumer(consumer *dm.Consumer) {
    p.consumer = consumer
}

func (p *updateAccountContext) ExternalUser() *dm.ExternalUser {
    return p.externalUser
}

func (p *updateAccountContext) SetExternalUser(externalUser *dm.ExternalUser) {
    p.externalUser = externalUser
}

func (p *updateAccountContext) LastModified() *time.Time {
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

func (p *updateAccountContext) ToObject() interface{} {
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

func (p *updateAccountContext) ETag() string {
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

func (p *updateAccountContext) RequestingUser() *dm.User {
    return p.requestingUser
}

func (p *updateAccountContext) SetRequestingUser(user *dm.User) {
    p.requestingUser = user
}

func (p *updateAccountContext) RequestingConsumer() *dm.Consumer {
    return p.requestingConsumer
}

func (p *updateAccountContext) SetRequestingConsumer(consumer *dm.Consumer) {
    p.requestingConsumer = consumer
}

func (p *updateAccountContext) OriginalValue() interface{} {
    return p.originalValue
}

func (p *updateAccountContext) SetOriginalValue(value interface{}) {
    p.originalValue = value
}

func (p *updateAccountContext) InputValidated() bool {
    return p.inputValidated
}

func NewUpdateAccountRequestHandler(ds acct.DataStore, authDS auth.DataStore) *UpdateAccountRequestHandler {
    return &UpdateAccountRequestHandler{ds: ds, authDS: authDS}
}

func (p *UpdateAccountRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) UpdateAccountContext {
    if uac, ok := cxt.(UpdateAccountContext); ok {
        return uac
    }
    return NewUpdateAccountContext()
}

func (p *UpdateAccountRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/account/(user|consumer|external_user)/update/(id)
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen >= 8 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "account" && path[6] == "update" {
            switch path[5] {
            case "user", "consumer", "external_user":
                return p
            }
        }
    }
    return nil
}

func (p *UpdateAccountRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    uac := p.GenerateContext(req, cxt)
    path := req.URLParts()
    pathLen := len(path)
    if pathLen >= 8 {
        uac.SetType(path[5])
        var id string
        if path[pathLen-1] == "" {
            id = strings.Join(path[7:pathLen-1], "/")
        } else {
            id = strings.Join(path[7:], "/")
        }
        switch uac.Type() {
        case "user":
            user, _ := p.ds.RetrieveUserAccountById(id)
            uac.SetUser(user)
            if user == nil {
                //log.Printf("[UARH]: Setting original value for user: %#v\n", nil)
                uac.SetOriginalValue(nil)
            } else {
                //log.Printf("[UARH]: Setting original value for user: %#v\n", user)
                uac.SetOriginalValue(user)
            }
        case "consumer":
            consumer, _ := p.ds.RetrieveConsumerAccountById(id)
            uac.SetConsumer(consumer)
            if consumer == nil {
                uac.SetOriginalValue(nil)
            } else {
                uac.SetOriginalValue(consumer)
            }
        case "external_user":
            externalUser, _ := p.ds.RetrieveExternalUserAccountById(id)
            uac.SetExternalUser(externalUser)
            if externalUser == nil {
                uac.SetOriginalValue(nil)
            } else {
                uac.SetOriginalValue(externalUser)
            }
        }
    }
    return req, uac
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *UpdateAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    uac := cxt.(UpdateAccountContext)
    //log.Printf("[UARH]: Checking original value: %#v vs. %v\n", uac.OriginalValue(), uac.OriginalValue() != nil)
    
    return uac.OriginalValue() != nil, req, cxt, 0, nil
}

func (p *UpdateAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST, wm.PUT}, req, cxt, 0, nil
}

func (p *UpdateAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    uac := cxt.(UpdateAccountContext)
    hasSignature, userId, consumerId, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        uac.SetRequestingUser(user)
    }
    if consumerId != "" {
        consumer, _ := p.ds.RetrieveConsumerAccountById(consumerId)
        uac.SetRequestingConsumer(consumer)
    }
    return true, "", req, cxt, 0, nil
}

func (p *UpdateAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    uac := cxt.(UpdateAccountContext)
    if uac.RequestingUser() != nil && uac.RequestingUser().Accessible() && (uac.RequestingUser().Role == dm.ROLE_ADMIN || (uac.User() != nil && uac.RequestingUser().Id == uac.User().Id)) {
        return false, req, cxt, 0, nil
    }
    // Cannot find user or consumer with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *UpdateAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *UpdateAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}

func (p *UpdateAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    uac := cxt.(UpdateAccountContext)
    obj := uac.ToObject()
    lastModified := uac.LastModified()
    etag := uac.ETag()
    var jsonObj jsonhelper.JSONObject
    if obj != nil {
        theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
        jsonObj, _ = theobj.(jsonhelper.JSONObject)
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, lastModified, etag)}, req, uac, 0, nil
}

func (p *UpdateAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    arr := []wm.MediaTypeInputHandler{apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body())}
    return arr, req, cxt, 0, nil
}

/*
func (p *UpdateAccountRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateAccountRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateAccountRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateAccountRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *UpdateAccountRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateAccountRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateAccountRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *UpdateAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    uac := cxt.(UpdateAccountContext)
    return uac.LastModified(), req, cxt, 0, nil
}

/*
func (p *UpdateAccountRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *UpdateAccountRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    var etag string
    uac := cxt.(UpdateAccountContext)
    switch uac.Type() {
    case "user":
        etag = uac.User().Etag
    case "consumer":
        etag = uac.Consumer().Etag
    case "external_user":
        etag = uac.ExternalUser().Etag
    }
    return etag, req, cxt, 0, nil
}


/*
func (p *UpdateAccountRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *UpdateAccountRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *UpdateAccountRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *UpdateAccountRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    uac := cxt.(UpdateAccountContext)
    uac.SetFromJSON(inputObj)
    uac.CleanInput(uac.RequestingUser(), uac.OriginalValue())
    //log.Print("[UARH]: HandleJSONObjectInputHandler()")
    errors := make(map[string][]os.Error)
    var obj interface{}
    var err os.Error
    ds := p.ds
    switch uac.Type() {
    case "user":
        if user := uac.User(); user != nil {
            //log.Printf("[UARH]: user is not nil1: %v\n", user)
            user.Validate(false, errors)
            if len(errors) == 0 {
                user, err = ds.UpdateUserAccount(user)
                //log.Printf("[UARH]: user after errors is %v\n", user)
            }
            obj = user
            uac.SetUser(user)
            //log.Printf("[UARH]: setUser to %v\n", user)
        }
    case "consumer":
        if user := uac.Consumer(); user != nil {
            user.Validate(false, errors)
            if len(errors) == 0 {
                user, err = ds.UpdateConsumerAccount(user)
            }
            obj = user
            uac.SetConsumer(user)
        }
    case "external_user":
        if user := uac.ExternalUser(); user != nil {
            user.Validate(false, errors)
            if len(errors) == 0 {
                user, err = ds.UpdateExternalUserAccount(user)
            }
            obj = user
            uac.SetExternalUser(user)
        }
    default:
        return apiutil.OutputErrorMessage("\"type\" must be \"user\", \"consumer\", or \"external_user\"", nil, 400, nil)
    }
    if len(errors) > 0 {
        return apiutil.OutputErrorMessage("Value errors. See result", errors, http.StatusBadRequest, nil)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(err.String(), nil, http.StatusInternalServerError, nil)
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    //log.Printf("[UARH]: obj was: \n%v\n", obj)
    //log.Printf("[UARH]: Going to output:\n%s\n", jsonObj)
    return apiutil.OutputJSONObject(jsonObj, uac.LastModified(), uac.ETag(), 0, nil)
}
