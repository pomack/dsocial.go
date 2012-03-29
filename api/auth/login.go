package auth

import (
    "errors"
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "io"
    "net/http"
    "net/url"
    "time"
)

var (
    ERR_INVALID_USERNAME_PASSWORD_COMBO = errors.New("Invalid combination of username/password")
    ERR_MUST_SPECIFY_USERNAME           = errors.New("Must specify username")
    ERR_MUST_SPECIFY_PASSWORD           = errors.New("Must specify password")
    ERR_VALUE_ERRORS                    = errors.New("Value errors. See result")
)

type LoginAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds     acct.DataStore
    authDS authentication.DataStore
}

type LoginAccountContext interface {
    SetFromJSON(obj jsonhelper.JSONObject)
    SetFromUrlEncoded(values url.Values)
    ValidateLogin(acctDS acct.DataStore, authDS authentication.DataStore, errors map[string][]error) (*dm.User, error)
    User() *dm.User
    Username() string
    Password() string
    InputValidated() bool
    SetResult(obj jsonhelper.JSONObject)
    Result() jsonhelper.JSONObject
}

type loginAccountContext struct {
    username       string
    password       string
    inputValidated bool
    user           *dm.User
    result         jsonhelper.JSONObject
}

func NewLoginAccountContext() LoginAccountContext {
    return new(loginAccountContext)
}

func (p *loginAccountContext) SetFromJSON(obj jsonhelper.JSONObject) {
    p.username = obj.GetAsString("username")
    p.password = obj.GetAsString("password")
    p.inputValidated = false
    p.user = nil
    p.result = nil
}

func (p *loginAccountContext) SetFromUrlEncoded(values url.Values) {
    p.username = values.Get("username")
    p.password = values.Get("password")
    p.inputValidated = false
    p.user = nil
    p.result = nil
}

func (p *loginAccountContext) ValidateLogin(acctDS acct.DataStore, authDS authentication.DataStore, errors map[string][]error) (*dm.User, error) {
    if errors == nil {
        errors = make(map[string][]error)
    }
    if p.username == "" {
        errors["username"] = []error{ERR_MUST_SPECIFY_USERNAME}
    }
    if p.password == "" {
        errors["password"] = []error{ERR_MUST_SPECIFY_PASSWORD}
    }
    p.inputValidated = true
    if len(errors) != 0 {
        return nil, nil
    }
    user, err := acctDS.FindUserAccountByUsername(p.username)
    if user == nil || err != nil || user.Id == "" {
        return nil, ERR_INVALID_USERNAME_PASSWORD_COMBO
    }
    pwd, err := authDS.RetrieveUserPassword(user.Id)
    if pwd == nil || err != nil {
        return nil, ERR_INVALID_USERNAME_PASSWORD_COMBO
    }
    if !user.Accessible() || !pwd.CheckPassword(p.password) {
        return nil, ERR_INVALID_USERNAME_PASSWORD_COMBO
    }
    p.user = user
    return user, nil
}

func (p *loginAccountContext) User() *dm.User {
    return p.user
}

func (p *loginAccountContext) Username() string {
    return p.username
}

func (p *loginAccountContext) Password() string {
    return p.password
}

func (p *loginAccountContext) InputValidated() bool {
    return p.inputValidated
}

func (p *loginAccountContext) Result() jsonhelper.JSONObject {
    return p.result
}

func (p *loginAccountContext) SetResult(result jsonhelper.JSONObject) {
    p.result = result
}

func NewLoginAccountRequestHandler(ds acct.DataStore, authDS authentication.DataStore) *LoginAccountRequestHandler {
    return &LoginAccountRequestHandler{ds: ds, authDS: authDS}
}

func (p *LoginAccountRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) LoginAccountContext {
    if lac, ok := cxt.(LoginAccountContext); ok {
        return lac
    }
    return NewLoginAccountContext()
}

func (p *LoginAccountRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/auth/login
    // /auth/login
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 6 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "auth" && path[5] == "login" {
            return p
        }
    }
    if pathLen == 3 {
        if path[0] == "" && path[1] == "auth" && path[2] == "login" {
            return p
        }
    }
    return nil
}

func (p *LoginAccountRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    lac := p.GenerateContext(req, cxt)
    return req, lac
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *LoginAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, error) {
    return []string{wm.POST}, req, cxt, 0, nil
}

/*
func (p *LoginAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    return true, "", req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *LoginAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *LoginAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}

func (p *LoginAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, error) {
    genFunc := func() (jsonhelper.JSONObject, time.Time, string, int, http.Header) {
        lac := cxt.(LoginAccountContext)
        jsonObj := lac.Result()
        headers := apiutil.AddNoCacheHeaders(nil)
        return jsonObj, time.Time{}, "", http.StatusOK, headers
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandlerWithGenerator(genFunc, time.Time{}, "")}, req, cxt, 0, nil
}

func (p *LoginAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, error) {
    arr := []wm.MediaTypeInputHandler{
        apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body()),
        apiutil.NewUrlEncodedMediaTypeInputHandler("", "", p),
    }
    return arr, req, cxt, 0, nil
}

/*
func (p *LoginAccountRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LoginAccountRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LoginAccountRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LoginAccountRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *LoginAccountRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LoginAccountRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LoginAccountRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *LoginAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LoginAccountRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *LoginAccountRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *LoginAccountRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *LoginAccountRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *LoginAccountRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    lac := cxt.(LoginAccountContext)
    lac.SetFromJSON(inputObj)
    return p.HandleInputHandlerAfterSetup(lac)
}

func (p *LoginAccountRequestHandler) HandleUrlEncodedInputHandler(req wm.Request, cxt wm.Context, inputObj url.Values) (int, http.Header, io.WriterTo) {
    lac := cxt.(LoginAccountContext)
    lac.SetFromUrlEncoded(inputObj)
    return p.HandleInputHandlerAfterSetup(lac)
}

func (p *LoginAccountRequestHandler) HandleInputHandlerAfterSetup(lac LoginAccountContext) (int, http.Header, io.WriterTo) {
    errors := make(map[string][]error)
    user, err := lac.ValidateLogin(p.ds, p.authDS, errors)
    if len(errors) > 0 {
        if err != nil {
            return apiutil.OutputErrorMessage(err.Error(), errors, http.StatusBadRequest, nil)
        }
        return apiutil.OutputErrorMessage(ERR_VALUE_ERRORS.Error(), errors, http.StatusUnauthorized, nil)
    }
    if err == ERR_INVALID_USERNAME_PASSWORD_COMBO {
        return apiutil.OutputErrorMessage(err.Error(), nil, http.StatusUnauthorized, nil)
    }
    if err != nil {
        return apiutil.OutputErrorMessage("Unable to process login request: ", nil, http.StatusInternalServerError, nil)
    }
    if user == nil {
        return apiutil.OutputErrorMessage("Unable to process login request: no such username", nil, http.StatusUnauthorized, nil)
    }
    accessKey, err := p.authDS.StoreAccessKey(dm.NewAccessKey(user.Id, ""))
    if err != nil {
        return apiutil.OutputErrorMessage("Unable to process login request: "+err.Error(), nil, http.StatusInternalServerError, nil)
    }
    obj := jsonhelper.NewJSONObject()
    obj.Set("user_id", user.Id)
    obj.Set("username", user.Username)
    obj.Set("name", user.Name)
    obj.Set("access_key_id", accessKey.Id)
    obj.Set("private_key", accessKey.PrivateKey)
    lac.SetResult(obj)
    return 0, nil, nil
}
