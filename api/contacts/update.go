package contacts

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    auth "github.com/pomack/dsocial.go/backend/authentication"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    //"log"
    "os"
    "strings"
    "time"
)

type UpdateContactRequestHandler struct {
    wm.DefaultRequestHandler
    ds  acct.DataStore
    authDS auth.DataStore
}

type UpdateContactContext interface {
    SetFromJSON(obj jsonhelper.JSONObject)
    CleanInput(createdByUser *dm.User, originalUser *dm.Contact)
    AuthUser() *dm.User
    SetAuthUser(user *dm.User)
    User() *dm.User
    SetUser(user *dm.User)
    LastModified() *time.Time
    ETag() string
    ContactId() string
    SetContactId(contactId string)
    Contact() *dm.Contact
    SetContact(contact *dm.Contact)
    OriginalContact() *dm.Contact
    SetOriginalContact(originalContact *dm.Contact)
    Result() jsonhelper.JSONObject
    SetResult(result jsonhelper.JSONObject)
    InputValidated() bool
}

type updateContactContext struct {
    authUser        *dm.User
    user            *dm.User
    lastModified    *time.Time
    etag            string
    contactId       string
    originalContact *dm.Contact
    contact         *dm.Contact
    result          jsonhelper.JSONObject
    inputValidated bool
}

func NewUpdateContactContext() UpdateContactContext {
    return new(updateContactContext)
}

func (p *updateContactContext) SetFromJSON(obj jsonhelper.JSONObject) {
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

func (p *updateContactContext) CleanInput(createdByUser *dm.User, originalUser interface{}) {
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

func (p *updateContactContext) Type() string {
    return p.theType
}

func (p *updateContactContext) SetType(theType string) {
    p.theType = theType
}

func (p *updateContactContext) User() *dm.User {
    return p.user
}

func (p *updateContactContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *updateContactContext) Consumer() *dm.Consumer {
    return p.consumer
}

func (p *updateContactContext) SetConsumer(consumer *dm.Consumer) {
    p.consumer = consumer
}

func (p *updateContactContext) ExternalUser() *dm.ExternalUser {
    return p.externalUser
}

func (p *updateContactContext) SetExternalUser(externalUser *dm.ExternalUser) {
    p.externalUser = externalUser
}

func (p *updateContactContext) LastModified() *time.Time {
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

func (p *updateContactContext) ToObject() interface{} {
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

func (p *updateContactContext) ETag() string {
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

func (p *updateContactContext) RequestingUser() *dm.User {
    return p.requestingUser
}

func (p *updateContactContext) SetRequestingUser(user *dm.User) {
    p.requestingUser = user
}

func (p *updateContactContext) RequestingConsumer() *dm.Consumer {
    return p.requestingConsumer
}

func (p *updateContactContext) SetRequestingConsumer(consumer *dm.Consumer) {
    p.requestingConsumer = consumer
}

func (p *updateContactContext) OriginalValue() interface{} {
    return p.originalValue
}

func (p *updateContactContext) SetOriginalValue(value interface{}) {
    p.originalValue = value
}

func (p *updateContactContext) InputValidated() bool {
    return p.inputValidated
}

func NewUpdateContactRequestHandler(ds acct.DataStore, authDS auth.DataStore) *UpdateContactRequestHandler {
    return &UpdateContactRequestHandler{ds: ds, authDS: authDS}
}

func (p *UpdateContactRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) UpdateContactContext {
    if ucc, ok := cxt.(UpdateContactContext); ok {
        return ucc
    }
    return NewUpdateContactContext()
}

func (p *UpdateContactRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
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

func (p *UpdateContactRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    ucc := p.GenerateContext(req, cxt)
    path := req.URLParts()
    pathLen := len(path)
    if pathLen >= 8 {
        ucc.SetType(path[5])
        var id string
        if path[pathLen-1] == "" {
            id = strings.Join(path[7:pathLen-1], "/")
        } else {
            id = strings.Join(path[7:], "/")
        }
        switch ucc.Type() {
        case "user":
            user, _ := p.ds.RetrieveUserAccountById(id)
            ucc.SetUser(user)
            if user == nil {
                //log.Printf("[UARH]: Setting original value for user: %#v\n", nil)
                ucc.SetOriginalValue(nil)
            } else {
                //log.Printf("[UARH]: Setting original value for user: %#v\n", user)
                ucc.SetOriginalValue(user)
            }
        case "consumer":
            consumer, _ := p.ds.RetrieveConsumerAccountById(id)
            ucc.SetConsumer(consumer)
            if consumer == nil {
                ucc.SetOriginalValue(nil)
            } else {
                ucc.SetOriginalValue(consumer)
            }
        case "external_user":
            externalUser, _ := p.ds.RetrieveExternalUserAccountById(id)
            ucc.SetExternalUser(externalUser)
            if externalUser == nil {
                ucc.SetOriginalValue(nil)
            } else {
                ucc.SetOriginalValue(externalUser)
            }
        }
    }
    return req, ucc
}

/*
func (p *UpdateContactRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *UpdateContactRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    //log.Printf("[UARH]: Checking original value: %#v vs. %v\n", ucc.OriginalValue(), ucc.OriginalValue() != nil)
    
    return ucc.OriginalValue() != nil, req, cxt, 0, nil
}

func (p *UpdateContactRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST, wm.PUT}, req, cxt, 0, nil
}

func (p *UpdateContactRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    hasSignature, userId, consumerId, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        ucc.SetRequestingUser(user)
    }
    if consumerId != "" {
        consumer, _ := p.ds.RetrieveConsumerAccountById(consumerId)
        ucc.SetRequestingConsumer(consumer)
    }
    return true, "", req, cxt, 0, nil
}

func (p *UpdateContactRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    if ucc.RequestingUser() != nil && ucc.RequestingUser().Accessible() && (ucc.RequestingUser().Role == dm.ROLE_ADMIN || (ucc.User() != nil && ucc.RequestingUser().Id == ucc.User().Id)) {
        return false, req, cxt, 0, nil
    }
    // Cannot find user or consumer with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *UpdateContactRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *UpdateContactRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *UpdateContactRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}

func (p *UpdateContactRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    obj := ucc.ToObject()
    lastModified := ucc.LastModified()
    etag := ucc.ETag()
    var jsonObj jsonhelper.JSONObject
    if obj != nil {
        theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
        jsonObj, _ = theobj.(jsonhelper.JSONObject)
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, lastModified, etag)}, req, ucc, 0, nil
}

func (p *UpdateContactRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    arr := []wm.MediaTypeInputHandler{apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body())}
    return arr, req, cxt, 0, nil
}

/*
func (p *UpdateContactRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateContactRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateContactRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateContactRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *UpdateContactRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateContactRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *UpdateContactRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *UpdateContactRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    return ucc.LastModified(), req, cxt, 0, nil
}

/*
func (p *UpdateContactRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/

func (p *UpdateContactRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    var etag string
    ucc := cxt.(UpdateContactContext)
    switch ucc.Type() {
    case "user":
        etag = ucc.User().Etag
    case "consumer":
        etag = ucc.Consumer().Etag
    case "external_user":
        etag = ucc.ExternalUser().Etag
    }
    return etag, req, cxt, 0, nil
}


/*
func (p *UpdateContactRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *UpdateContactRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *UpdateContactRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *UpdateContactRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    ucc := cxt.(UpdateContactContext)
    ucc.SetFromJSON(inputObj)
    ucc.CleanInput(ucc.RequestingUser(), ucc.OriginalValue())
    //log.Print("[UARH]: HandleJSONObjectInputHandler()")
    errors := make(map[string][]os.Error)
    var obj interface{}
    var err os.Error
    ds := p.ds
    switch ucc.Type() {
    case "user":
        if user := ucc.User(); user != nil {
            //log.Printf("[UARH]: user is not nil1: %v\n", user)
            user.Validate(false, errors)
            if len(errors) == 0 {
                user, err = ds.UpdateUserAccount(user)
                //log.Printf("[UARH]: user after errors is %v\n", user)
            }
            obj = user
            ucc.SetUser(user)
            //log.Printf("[UARH]: setUser to %v\n", user)
        }
    case "consumer":
        if user := ucc.Consumer(); user != nil {
            user.Validate(false, errors)
            if len(errors) == 0 {
                user, err = ds.UpdateConsumerAccount(user)
            }
            obj = user
            ucc.SetConsumer(user)
        }
    case "external_user":
        if user := ucc.ExternalUser(); user != nil {
            user.Validate(false, errors)
            if len(errors) == 0 {
                user, err = ds.UpdateExternalUserAccount(user)
            }
            obj = user
            ucc.SetExternalUser(user)
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
    return apiutil.OutputJSONObject(jsonObj, ucc.LastModified(), ucc.ETag(), http.StatusOK, nil)
}
