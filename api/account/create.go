package account

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    auth "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "io"
    "net/http"
    "time"
)

type CreateAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds     acct.DataStore
    authDS auth.DataStore
}

type CreateAccountContext interface {
    SetFromJSON(obj jsonhelper.JSONObject)
    CleanInput(createdByUser *dm.User)
    Type() string
    SetType(theType string)
    User() *dm.User
    Consumer() *dm.Consumer
    ExternalUser() *dm.ExternalUser
    LastModified() time.Time
    ETag() string
    ToObject() interface{}
    RequestingUser() *dm.User
    SetRequestingUser(user *dm.User)
    RequestingConsumer() *dm.Consumer
    SetRequestingConsumer(consumer *dm.Consumer)
    Password() string
}

type createAccountContext struct {
    theType            string
    user               *dm.User
    consumer           *dm.Consumer
    externalUser       *dm.ExternalUser
    requestingUser     *dm.User
    requestingConsumer *dm.Consumer
    password           string
}

func NewCreateAccountContext() CreateAccountContext {
    return new(createAccountContext)
}

func (p *createAccountContext) SetFromJSON(obj jsonhelper.JSONObject) {
    p.user = nil
    p.consumer = nil
    p.externalUser = nil
    p.password = ""
    theType := p.theType
    if theType == "" {
        theType = obj.GetAsString("type")
    }
    switch theType {
    case "user":
        p.user = new(dm.User)
        p.user.InitFromJSONObject(obj)
        p.password = obj.GetAsString("password")
    case "consumer":
        p.consumer = new(dm.Consumer)
        p.consumer.InitFromJSONObject(obj)
    case "external_user":
        p.externalUser = new(dm.ExternalUser)
        p.externalUser.InitFromJSONObject(obj)
    }
}

func (p *createAccountContext) CleanInput(createdByUser *dm.User) {
    if p.user != nil {
        p.user.Id = ""
        p.user.CleanFromUser(createdByUser, nil)
    } else if p.consumer != nil {
        p.consumer.Id = ""
        p.consumer.CleanFromUser(createdByUser, nil)
    } else if p.externalUser != nil {
        p.externalUser.Id = ""
        p.externalUser.CleanFromUser(createdByUser, nil)
    }
}

func (p *createAccountContext) Type() string {
    return p.theType
}

func (p *createAccountContext) SetType(theType string) {
    p.theType = theType
}

func (p *createAccountContext) User() *dm.User {
    return p.user
}

func (p *createAccountContext) Consumer() *dm.Consumer {
    return p.consumer
}

func (p *createAccountContext) ExternalUser() *dm.ExternalUser {
    return p.externalUser
}

func (p *createAccountContext) LastModified() time.Time {
    var lastModified time.Time
    if p.user != nil && p.user.ModifiedAt != 0 {
        lastModified = time.Unix(p.user.ModifiedAt, 0).UTC()
    } else if p.consumer != nil && p.consumer.ModifiedAt != 0 {
        lastModified = time.Unix(p.consumer.ModifiedAt, 0).UTC()
    } else if p.externalUser != nil && p.externalUser.ModifiedAt != 0 {
        lastModified = time.Unix(p.externalUser.ModifiedAt, 0).UTC()
    }
    return lastModified
}

func (p *createAccountContext) ToObject() interface{} {
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

func (p *createAccountContext) ETag() string {
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

func (p *createAccountContext) RequestingUser() *dm.User {
    return p.requestingUser
}

func (p *createAccountContext) SetRequestingUser(user *dm.User) {
    p.requestingUser = user
}

func (p *createAccountContext) RequestingConsumer() *dm.Consumer {
    return p.requestingConsumer
}

func (p *createAccountContext) SetRequestingConsumer(consumer *dm.Consumer) {
    p.requestingConsumer = consumer
}

func (p *createAccountContext) Password() string {
    return p.password
}

func NewCreateAccountRequestHandler(ds acct.DataStore, authDS auth.DataStore) *CreateAccountRequestHandler {
    return &CreateAccountRequestHandler{ds: ds, authDS: authDS}
}

func (p *CreateAccountRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) CreateAccountContext {
    if cac, ok := cxt.(CreateAccountContext); ok {
        return cac
    }
    return NewCreateAccountContext()
}

func (p *CreateAccountRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/account/(user|consumer|external_user)/create
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 7 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "account" && path[6] == "create" {
            switch path[5] {
            case "user", "consumer", "external_user":
                return p
            }
        }
    }
    return nil
}

func (p *CreateAccountRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    cac := p.GenerateContext(req, cxt)
    path := req.URLParts()
    if len(path) >= 6 {
        cac.SetType(path[5])
    }
    return req, cac
}

/*
func (p *CreateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    return false, req, cxt, 0, nil
}

func (p *CreateAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, error) {
    return []string{wm.POST, wm.PUT}, req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    return true, "", req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    cac := cxt.(CreateAccountContext)
    hasSignature, userId, consumerId, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if err != nil {
        return true, req, cxt, 403, err
    }
    if hasSignature {
        if userId != "" {
            user, _ := p.ds.RetrieveUserAccountById(userId)
            cac.SetRequestingUser(user)
        }
        if consumerId != "" {
            consumer, _ := p.ds.RetrieveConsumerAccountById(consumerId)
            cac.SetRequestingConsumer(consumer)
        }
        if (userId != "" && (cac.RequestingUser() == nil || !cac.RequestingUser().Accessible())) && (consumerId != "" && (cac.RequestingConsumer() == nil || !cac.RequestingConsumer().Accessible())) {
            // Cannot find user or consumer with specified id
            return true, req, cxt, 0, nil
        }
    }
    return false, req, cxt, 0, nil
}

func (p *CreateAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    return true, req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *CreateAccountRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}

/*
func (p *CreateAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    return true, req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}

func (p *CreateAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, error) {
    cac := cxt.(CreateAccountContext)
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

func (p *CreateAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, error) {
    arr := []wm.MediaTypeInputHandler{apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body())}
    return arr, req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *CreateAccountRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *CreateAccountRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *CreateAccountRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *CreateAccountRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *CreateAccountRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
  return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *CreateAccountRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *CreateAccountRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *CreateAccountRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *CreateAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (time.Time, wm.Request, wm.Context, int, error) {
    cac := cxt.(CreateAccountContext)
    return cac.LastModified(), req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *CreateAccountRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *CreateAccountRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *CreateAccountRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *CreateAccountRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    cac := cxt.(CreateAccountContext)
    cac.SetFromJSON(inputObj)
    cac.CleanInput(cac.RequestingUser())

    errors := make(map[string][]error)
    var obj map[string]interface{}
    var accessKey *dm.AccessKey
    var err error
    ds := p.ds
    authDS := p.authDS
    if user := cac.User(); user != nil {
        var userPassword *dm.UserPassword
        user.Validate(true, errors)
        if len(errors) == 0 {
            user, err = ds.CreateUserAccount(user)
            if err == nil && user != nil {
                accessKey, err = authDS.StoreAccessKey(dm.NewAccessKey(user.Id, ""))
            }
        }
        if cac.Password() != "" && user != nil && user.Id != "" {
            userPassword = dm.NewUserPassword(user.Id, cac.Password())
            userPassword.Validate(true, errors)
            if len(errors) == 0 && err == nil {
                userPassword, err = authDS.StoreUserPassword(userPassword)
            }
        }
        obj = make(map[string]interface{})
        obj["user"] = user
        obj["type"] = "user"
        obj["key"] = accessKey
    } else if user := cac.Consumer(); user != nil {
        user.Validate(true, errors)
        if len(errors) == 0 {
            user, err = ds.CreateConsumerAccount(user)
            if err == nil && user != nil {
                accessKey, err = authDS.StoreAccessKey(dm.NewAccessKey("", user.Id))
            }
        }
        obj = make(map[string]interface{})
        obj["consumer"] = user
        obj["type"] = "consumer"
        obj["key"] = accessKey
    } else if user := cac.ExternalUser(); user != nil {
        user.Validate(true, errors)
        if len(errors) == 0 {
            user, err = ds.CreateExternalUserAccount(user)
            if err == nil && user != nil {
                accessKey, err = authDS.StoreAccessKey(dm.NewAccessKey(user.Id, user.ConsumerId))
            }
        }
        obj = make(map[string]interface{})
        obj["external_user"] = user
        obj["type"] = "external_user"
        obj["key"] = accessKey
    } else {
        return apiutil.OutputErrorMessage("\"type\" must be \"user\", \"consumer\", or \"external_user\"", nil, 400, nil)
    }
    if len(errors) > 0 {
        return apiutil.OutputErrorMessage("Value errors. See result", errors, http.StatusBadRequest, nil)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(err.Error(), nil, http.StatusInternalServerError, nil)
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    return apiutil.OutputJSONObject(jsonObj, cac.LastModified(), cac.ETag(), 0, nil)
}
