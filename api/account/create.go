package account

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    auth "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "bytes"
    "http"
    "io"
    "os"
    "time"
)

type CreateAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds  acct.DataStore
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
    LastModified() *time.Time
    ETag() string
    ToObject() interface{}
    RequestingUser() *dm.User
    SetRequestingUser(user *dm.User)
    RequestingConsumer() *dm.Consumer
    SetRequestingConsumer(consumer *dm.Consumer)
}

type createAccountContext struct {
    theType      string
    user         *dm.User
    consumer     *dm.Consumer
    externalUser *dm.ExternalUser
    requestingUser     *dm.User
    requestingConsumer *dm.Consumer
}

func NewCreateAccountContext() CreateAccountContext {
    return new(createAccountContext)
}

func (p *createAccountContext) SetFromJSON(obj jsonhelper.JSONObject) {
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

func (p *createAccountContext) LastModified() *time.Time {
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

func (p *CreateAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}

func (p *CreateAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST, wm.PUT}, req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    return true, "", req, cxt, 0, nil
}
*/


func (p *CreateAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
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
        if (userId != "" && cac.RequestingUser() == nil) && (consumerId != "" && cac.RequestingConsumer() == nil) {
            // Cannot find user or consumer with specified id
            return true, req, cxt, 0, nil
        }
    }
    return false, req, cxt, 0, nil
}


func (p *CreateAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
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

func (p *CreateAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 500, nil
}

/*
func (p *CreateAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}

/*
func (p *CreateAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *CreateAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        buf := bytes.NewBufferString("")
        httpCode, _, httpError := mths[0].OutputTo(req, cxt, buf)
        if httpCode > 0 {
            if httpError == nil && buf.Len() > 0 {
                return false, req, cxt, httpCode, buf
            }
        }
        return false, req, cxt, httpCode, httpError
    }
    return false, req, cxt, code, err
}

func (p *CreateAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
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

func (p *CreateAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
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

func (p *CreateAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
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
    method := req.Method()
    if method == wm.HEAD || method == wm.PUT || method == wm.DELETE {
        return false
    }
    return true
}

func (p *CreateAccountRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, writer io.Writer, inputObj jsonhelper.JSONObject) (int, http.Header, os.Error) {
    cac := cxt.(CreateAccountContext)
    cac.SetFromJSON(inputObj)
    cac.CleanInput(cac.RequestingUser())
    
    errors := make(map[string][]os.Error)
    var obj interface{}
    var err os.Error
    ds := p.ds
    if user := cac.User(); user != nil {
        user.Validate(true, errors)
        if len(errors) == 0 {
            user, err = ds.CreateUserAccount(user)
        }
        obj = user
    } else if user := cac.Consumer(); user != nil {
        user.Validate(true, errors)
        if len(errors) == 0 {
            user, err = ds.CreateConsumerAccount(user)
        }
        obj = user
    } else if user := cac.ExternalUser(); user != nil {
        user.Validate(true, errors)
        if len(errors) == 0 {
            user, err = ds.CreateExternalUserAccount(user)
        }
        obj = user
    } else {
        return apiutil.OutputErrorMessage(writer, "\"type\" must be \"user\", \"consumer\", or \"external_user\"", nil, 400, nil)
    }
    if len(errors) > 0 {
        return apiutil.OutputErrorMessage(writer, "Value errors. See result", errors, 400, nil)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(writer, err.String(), nil, 500, nil)
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    return apiutil.OutputJSONObject(writer, jsonObj, cac.LastModified(), cac.ETag(), 0, nil)
}
