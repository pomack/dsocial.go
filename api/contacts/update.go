package contacts

import (
    "container/list"
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    auth "github.com/pomack/dsocial.go/backend/authentication"
    bc "github.com/pomack/dsocial.go/backend/contacts"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    "json"
    //"log"
    "os"
    "time"
)

type UpdateContactRequestHandler struct {
    wm.DefaultRequestHandler
    ds         acct.DataStore
    authDS     auth.DataStore
    contactsDS bc.DataStoreService
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
    contactId       string
    contact         *dm.Contact
    originalContact *dm.Contact
    result          jsonhelper.JSONObject
    inputValidated  bool
}

func NewUpdateContactContext() UpdateContactContext {
    return new(updateContactContext)
}

func (p *updateContactContext) SetFromJSON(obj jsonhelper.JSONObject) {
    contact := new(dm.Contact)
    if obj == nil {
        contact = nil
    } else {
        json.Unmarshal([]byte(obj.String()), contact)
    }
    p.contact = contact
}

func (p *updateContactContext) CleanInput(createdByUser *dm.User, originalUser *dm.Contact) {
    if p.user != nil && p.authUser != nil && p.user.Id == p.authUser.Id {
        p.contact.Id = p.contactId
        p.contact.UserId = p.user.Id
    }
    p.inputValidated = true
}

func (p *updateContactContext) AuthUser() *dm.User {
    return p.authUser
}

func (p *updateContactContext) SetAuthUser(authUser *dm.User) {
    p.authUser = authUser
}

func (p *updateContactContext) User() *dm.User {
    return p.user
}

func (p *updateContactContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *updateContactContext) LastModified() *time.Time {
    var lastModified *time.Time
    if p.contact != nil && p.contact.ModifiedAt > 0 {
        lastModified = time.SecondsToUTC(p.contact.ModifiedAt)
    } else if p.originalContact != nil && p.originalContact.ModifiedAt > 0 {
        lastModified = time.SecondsToUTC(p.originalContact.ModifiedAt)
    }
    return lastModified
}

func (p *updateContactContext) ETag() string {
    if p.contact != nil && p.contact.Etag != "" {
        return p.contact.Etag
    }
    if p.originalContact != nil && p.originalContact.Etag != "" {
        return p.originalContact.Etag
    }
    return ""
}

func (p *updateContactContext) ContactId() string {
    return p.contactId
}

func (p *updateContactContext) SetContactId(contactId string) {
    p.contactId = contactId
}

func (p *updateContactContext) Contact() *dm.Contact {
    return p.contact
}

func (p *updateContactContext) SetContact(contact *dm.Contact) {
    p.contact = contact
}

func (p *updateContactContext) OriginalContact() *dm.Contact {
    return p.originalContact
}

func (p *updateContactContext) SetOriginalContact(originalContact *dm.Contact) {
    p.originalContact = originalContact
}

func (p *updateContactContext) Result() jsonhelper.JSONObject {
    return p.result
}

func (p *updateContactContext) SetResult(result jsonhelper.JSONObject) {
    p.result = result
}

func (p *updateContactContext) InputValidated() bool {
    return p.inputValidated
}

func NewUpdateContactRequestHandler(ds acct.DataStore, authDS auth.DataStore, contactsDS bc.DataStoreService) *UpdateContactRequestHandler {
    return &UpdateContactRequestHandler{ds: ds, authDS: authDS, contactsDS: contactsDS}
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
    if pathLen == 9 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "u" && path[5] != "" && path[6] == "contacts" && path[7] == "update" && path[8] != "" {
            return p
        }
    } else if pathLen == 6 {
        if path[0] == "" && path[1] == "u" && path[2] != "" && path[3] == "contacts" && path[4] == "update" && path[5] != "" {
            return p
        }
    }
    return nil
}

func (p *UpdateContactRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    ucc := p.GenerateContext(req, cxt)
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    var userId string
    var contactId string
    switch pathLen {
    case 9:
        userId = path[5]
        contactId = path[8]
    case 6:
        userId = path[2]
        contactId = path[5]
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        ucc.SetUser(user)
        if contactId != "" {
            contact, _, _ := p.contactsDS.RetrieveDsocialContact(userId, contactId)
            ucc.SetOriginalContact(contact)
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
    return ucc.OriginalContact() != nil, req, cxt, 0, nil
}

func (p *UpdateContactRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST, wm.PUT}, req, cxt, 0, nil
}

func (p *UpdateContactRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    hasSignature, userId, _, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        ucc.SetAuthUser(user)
    }
    return userId != "", "", req, cxt, 0, nil
}

func (p *UpdateContactRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    ucc := cxt.(UpdateContactContext)
    if ucc.AuthUser() != nil && ucc.AuthUser().Accessible() && (ucc.AuthUser().Role == dm.ROLE_ADMIN || (ucc.User() != nil && ucc.AuthUser().Id == ucc.User().Id)) {
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
    obj := ucc.Contact()
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
    ucc := cxt.(UpdateContactContext)
    return ucc.ETag(), req, cxt, 0, nil
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
    ucc.CleanInput(ucc.AuthUser(), ucc.OriginalContact())
    //log.Print("[UARH]: HandleJSONObjectInputHandler()")
    errors := make(map[string][]os.Error)
    var obj interface{}
    var err os.Error
    contactsDS := p.contactsDS
    contact := ucc.Contact()
    origContact := ucc.OriginalContact()
    if contact != nil && origContact != nil && ucc.User() != nil {
        dsocialUserId := ucc.User().Id
        contact.Id = ucc.ContactId()
        contact.UserId = dsocialUserId
        contact.Validate(false, errors)
        if len(errors) == 0 {
            l := new(list.List)
            origContact.GenerateChanges(origContact, contact, nil, l)
            allowAdd, allowDelete, allowUpdate := true, true, true
            pipeline := bc.NewPipeline()
            l = pipeline.RemoveUnacceptedChanges(l, allowAdd, allowDelete, allowUpdate)
            changes := make([]*dm.Change, l.Len())
            for i, iter := 0, l.Front(); iter != nil; i, iter = i+1, iter.Next() {
                changes[i] = iter.Value.(*dm.Change)
            }
            changeset := &dm.ChangeSet{
                CreatedAt:      time.UTC().Format(dm.UTC_DATETIME_FORMAT),
                ChangedBy:      ucc.AuthUser().Id,
                ChangeImportId: ucc.ContactId(),
                RecordId:       ucc.ContactId(),
                Changes:        changes,
            }
            _, err = contactsDS.StoreContactChangeSet(dsocialUserId, changeset)
            if err == nil {
                contact, err = contactsDS.StoreDsocialContact(contact.UserId, contact.Id, contact)
            }
        }
    }
    if len(errors) > 0 {
        return apiutil.OutputErrorMessage("Value errors. See result", errors, http.StatusBadRequest, nil)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(err.String(), nil, http.StatusInternalServerError, nil)
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    return apiutil.OutputJSONObject(jsonObj, ucc.LastModified(), ucc.ETag(), http.StatusOK, nil)
}
