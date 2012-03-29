package auth

import (
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

type SetPasswordRequestHandler struct {
    wm.DefaultRequestHandler
    ds     acct.DataStore
    authDS authentication.DataStore
}

type SetPasswordContext interface {
    User() *dm.User
    SetUser(user *dm.User)
    Password() string
    SetPassword(password string)
    SetFromJSON(obj jsonhelper.JSONObject)
    SetFromUrlEncoded(obj url.Values)
    Result() jsonhelper.JSONObject
    SetResult(result jsonhelper.JSONObject)
}

type setPasswordContext struct {
    user     *dm.User
    password string
    result   jsonhelper.JSONObject
}

func NewSetPasswordContext() SetPasswordContext {
    return new(setPasswordContext)
}

func (p *setPasswordContext) User() *dm.User {
    return p.user
}

func (p *setPasswordContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *setPasswordContext) Password() string {
    return p.password
}

func (p *setPasswordContext) SetPassword(password string) {
    p.password = password
}

func (p *setPasswordContext) Result() jsonhelper.JSONObject {
    return p.result
}

func (p *setPasswordContext) SetResult(result jsonhelper.JSONObject) {
    p.result = result
}

func (p *setPasswordContext) SetFromJSON(obj jsonhelper.JSONObject) {
    p.password = obj.GetAsString("password")
    p.result = nil
}

func (p *setPasswordContext) SetFromUrlEncoded(values url.Values) {
    p.password = values.Get("password")
    p.result = nil
}

func NewSetPasswordRequestHandler(ds acct.DataStore, authDS authentication.DataStore) *SetPasswordRequestHandler {
    return &SetPasswordRequestHandler{ds: ds, authDS: authDS}
}

func (p *SetPasswordRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) SetPasswordContext {
    if spac, ok := cxt.(SetPasswordContext); ok {
        return spac
    }
    return NewSetPasswordContext()
}

func (p *SetPasswordRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/auth/set_password
    // /auth/set_password
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 6 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "auth" && path[5] == "set_password" {
            return p
        }
    }
    if pathLen == 3 {
        if path[0] == "" && path[1] == "auth" && path[2] == "set_password" {
            return p
        }
    }
    return nil
}

func (p *SetPasswordRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    spac := p.GenerateContext(req, cxt)
    return req, spac
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *SetPasswordRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, error) {
    return []string{wm.POST}, req, cxt, 0, nil
}

func (p *SetPasswordRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, error) {
    spac := cxt.(SetPasswordContext)
    hasSignature, userId, _, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        spac.SetUser(user)
    }
    return true, "", req, cxt, 0, nil
}

func (p *SetPasswordRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, error) {
    spac := cxt.(SetPasswordContext)
    if spac.User() != nil && spac.User().Accessible() {
        return false, req, cxt, 0, nil
    }
    // Cannot find user with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *SetPasswordRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *SetPasswordRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *SetPasswordRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, error) {
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}

func (p *SetPasswordRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, error) {
    genFunc := func() (jsonhelper.JSONObject, time.Time, string, int, http.Header) {
        spac := cxt.(SetPasswordContext)
        jsonObj := spac.Result()
        headers := apiutil.AddNoCacheHeaders(nil)
        return jsonObj, time.Time{}, "", http.StatusOK, headers
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandlerWithGenerator(genFunc, time.Time{}, "")}, req, cxt, 0, nil
}

func (p *SetPasswordRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, error) {
    arr := []wm.MediaTypeInputHandler{
        apiutil.NewJSONMediaTypeInputHandler("", "", p, req.Body()),
        apiutil.NewUrlEncodedMediaTypeInputHandler("", "", p),
    }
    return arr, req, cxt, 0, nil
}

/*
func (p *SetPasswordRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *SetPasswordRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *SetPasswordRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *SetPasswordRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *SetPasswordRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *SetPasswordRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *SetPasswordRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *SetPasswordRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *SetPasswordRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *SetPasswordRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *SetPasswordRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *SetPasswordRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}

func (p *SetPasswordRequestHandler) HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo) {
    lac := cxt.(SetPasswordContext)
    lac.SetFromJSON(inputObj)
    return p.HandleInputHandlerAfterSetup(lac)
}

func (p *SetPasswordRequestHandler) HandleUrlEncodedInputHandler(req wm.Request, cxt wm.Context, inputObj url.Values) (int, http.Header, io.WriterTo) {
    lac := cxt.(SetPasswordContext)
    lac.SetFromUrlEncoded(inputObj)
    return p.HandleInputHandlerAfterSetup(lac)
}

func (p *SetPasswordRequestHandler) HandleInputHandlerAfterSetup(cxt SetPasswordContext) (int, http.Header, io.WriterTo) {
    errors := make(map[string][]error)
    var obj jsonhelper.JSONObject
    var err error
    authDS := p.authDS
    if user := cxt.User(); user != nil {
        var userPassword *dm.UserPassword
        if user != nil {
            userPassword = dm.NewUserPassword(user.Id, cxt.Password())
        } else {
            userPassword = dm.NewUserPassword("", cxt.Password())
        }
        userPassword.Validate(true, errors)
        if len(errors) == 0 {
            userPassword, err = authDS.StoreUserPassword(userPassword)
        }
        obj = jsonhelper.NewJSONObject()
        userObj, _ := jsonhelper.Marshal(user)
        obj.Set("user", userObj)
        obj.Set("type", "user")
        obj.Set("message", "password changed")
    } else {
        return apiutil.OutputErrorMessage(ERR_MUST_SPECIFY_USERNAME.Error(), time.Time{}, http.StatusBadRequest, nil)
    }
    if len(errors) > 0 {
        return apiutil.OutputErrorMessage("Value errors. See result", errors, http.StatusBadRequest, nil)
    }
    if err != nil {
        return apiutil.OutputErrorMessage(err.Error(), time.Time{}, http.StatusInternalServerError, nil)
    }
    cxt.SetResult(obj)
    return 0, nil, nil
}
