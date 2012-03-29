package auth

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "os"
)

type GeneratePrivateKeyRequestHandler struct {
    wm.DefaultRequestHandler
    ds     acct.DataStore
    authDS authentication.DataStore
}

type GeneratePrivateKeyContext interface {
    User() *dm.User
    SetUser(user *dm.User)
    Consumer() *dm.Consumer
    SetConsumer(consumer *dm.Consumer)
    AccessKey() *dm.AccessKey
    SetAccessKey(accessKey *dm.AccessKey)
}

type generatePrivateKeyContext struct {
    accessKey *dm.AccessKey
    user      *dm.User
    consumer  *dm.Consumer
}

func NewGeneratePrivateKeyContext() GeneratePrivateKeyContext {
    return new(generatePrivateKeyContext)
}

func (p *generatePrivateKeyContext) User() *dm.User {
    return p.user
}

func (p *generatePrivateKeyContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *generatePrivateKeyContext) Consumer() *dm.Consumer {
    return p.consumer
}

func (p *generatePrivateKeyContext) SetConsumer(consumer *dm.Consumer) {
    p.consumer = consumer
}

func (p *generatePrivateKeyContext) AccessKey() *dm.AccessKey {
    return p.accessKey
}

func (p *generatePrivateKeyContext) SetAccessKey(accessKey *dm.AccessKey) {
    p.accessKey = accessKey
}

func NewGeneratePrivateKeyRequestHandler(ds acct.DataStore, authDS authentication.DataStore) *GeneratePrivateKeyRequestHandler {
    return &GeneratePrivateKeyRequestHandler{ds: ds, authDS: authDS}
}

func (p *GeneratePrivateKeyRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) GeneratePrivateKeyContext {
    if gpkc, ok := cxt.(GeneratePrivateKeyContext); ok {
        return gpkc
    }
    return NewGeneratePrivateKeyContext()
}

func (p *GeneratePrivateKeyRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/auth/login
    // /auth/login
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 6 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "auth" && path[5] == "generate_private_key" {
            return p
        }
    }
    if pathLen == 3 {
        if path[0] == "" && path[1] == "auth" && path[2] == "generate_private_key" {
            return p
        }
    }
    return nil
}

func (p *GeneratePrivateKeyRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    gpkc := p.GenerateContext(req, cxt)
    return req, gpkc
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *GeneratePrivateKeyRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST}, req, cxt, 0, nil
}

func (p *GeneratePrivateKeyRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    gpkc := cxt.(GeneratePrivateKeyContext)
    hasSignature, userId, consumerId, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        gpkc.SetUser(user)
    }
    if consumerId != "" {
        consumer, _ := p.ds.RetrieveConsumerAccountById(consumerId)
        gpkc.SetConsumer(consumer)
    }
    if (userId != "" && gpkc.User() == nil) || (consumerId != "" && gpkc.Consumer() == nil) {
        gpkc.SetUser(nil)
        gpkc.SetConsumer(nil)
    }
    return true, "", req, cxt, 0, nil
}

func (p *GeneratePrivateKeyRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    gpkc := cxt.(GeneratePrivateKeyContext)
    if gpkc.User() == nil && gpkc.Consumer() == nil {
        // cannot find user or consumer with specified ids
        return true, req, cxt, 0, nil
    }
    if gpkc.User() != nil && !gpkc.User().Accessible() {
        // user is not accessible
        return true, req, cxt, 0, nil
    }
    if gpkc.Consumer() != nil && !gpkc.Consumer().Accessible() {
        // consumer is not accessible
        return true, req, cxt, 0, nil
    }
    return false, req, cxt, 0, nil
}

/*
func (p *GeneratePrivateKeyRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    return req, cxt, 0, nil, nil, err
}
*/

func (p *GeneratePrivateKeyRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    gpkc := cxt.(GeneratePrivateKeyContext)
    user := gpkc.User()
    consumer := gpkc.Consumer()
    var userId, consumerId string
    if user != nil {
        userId = user.Id
    }
    if consumer != nil {
        consumerId = consumer.Id
    }
    accessKey, err := p.authDS.StoreAccessKey(dm.NewAccessKey(userId, consumerId))
    gpkc.SetAccessKey(accessKey)
    obj := make(map[string]interface{})
    if user != nil {
        obj["user_id"] = user.Id
        obj["username"] = user.Username
        obj["name"] = user.Name
    }
    if consumer != nil {
        obj["consumer_id"] = consumer.Id
        obj["consumer_short_name"] = consumer.ShortName
    }
    if accessKey != nil {
        obj["access_key_id"] = accessKey.Id
        obj["private_key"] = accessKey.PrivateKey
    }
    theobj, _ := jsonhelper.MarshalWithOptions(obj, dm.UTC_DATETIME_FORMAT)
    jsonObj, _ := theobj.(jsonhelper.JSONObject)
    if err != nil {
        return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, nil, "")}, req, gpkc, http.StatusInternalServerError, err
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(jsonObj, nil, "")}, req, gpkc, 0, nil
}

/*
func (p *GeneratePrivateKeyRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    return []wm.MediaTypeInputHandler{}, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *GeneratePrivateKeyRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *GeneratePrivateKeyRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *GeneratePrivateKeyRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *GeneratePrivateKeyRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *GeneratePrivateKeyRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *GeneratePrivateKeyRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *GeneratePrivateKeyRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *GeneratePrivateKeyRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}
