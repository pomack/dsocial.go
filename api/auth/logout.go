package auth

import (
    "github.com/pomack/dsocial.go/api/apiutil"
    acct "github.com/pomack/dsocial.go/backend/accounts"
    "github.com/pomack/dsocial.go/backend/authentication"
    dm "github.com/pomack/dsocial.go/models/dsocial"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    "os"
)

type LogoutAccountRequestHandler struct {
    wm.DefaultRequestHandler
    ds     acct.DataStore
    authDS authentication.DataStore
}

type LogoutAccountContext interface {
    User() *dm.User
    SetUser(user *dm.User)
    AccessKey() *dm.AccessKey
    SetAccessKey(accessKey *dm.AccessKey)
}

type logoutAccountContext struct {
    accessKey *dm.AccessKey
    user      *dm.User
}

func NewLogoutAccountContext() LogoutAccountContext {
    return new(logoutAccountContext)
}

func (p *logoutAccountContext) User() *dm.User {
    return p.user
}

func (p *logoutAccountContext) SetUser(user *dm.User) {
    p.user = user
}

func (p *logoutAccountContext) AccessKey() *dm.AccessKey {
    return p.accessKey
}

func (p *logoutAccountContext) SetAccessKey(accessKey *dm.AccessKey) {
    p.accessKey = accessKey
}

func NewLogoutAccountRequestHandler(ds acct.DataStore, authDS authentication.DataStore) *LogoutAccountRequestHandler {
    return &LogoutAccountRequestHandler{ds: ds, authDS: authDS}
}

func (p *LogoutAccountRequestHandler) GenerateContext(req wm.Request, cxt wm.Context) LogoutAccountContext {
    if lac, ok := cxt.(LogoutAccountContext); ok {
        return lac
    }
    return NewLogoutAccountContext()
}

func (p *LogoutAccountRequestHandler) HandlerFor(req wm.Request, writer wm.ResponseWriter) wm.RequestHandler {
    // /api/v1/json/auth/logout
    // /auth/logout
    path := req.URLParts()
    pathLen := len(path)
    if path[pathLen-1] == "" {
        // ignore trailing slash
        pathLen = pathLen - 1
    }
    if pathLen == 6 {
        if path[0] == "" && path[1] == "api" && path[2] == "v1" && path[3] == "json" && path[4] == "auth" && path[5] == "logout" {
            return p
        }
    }
    if pathLen == 3 {
        if path[0] == "" && path[1] == "auth" && path[2] == "logout" {
            return p
        }
    }
    return nil
}

func (p *LogoutAccountRequestHandler) StartRequest(req wm.Request, cxt wm.Context) (wm.Request, wm.Context) {
    lac := p.GenerateContext(req, cxt)
    return req, lac
}

/*
func (p *UpdateAccountRequestHandler) ServiceAvailable(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) ResourceExists(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

func (p *LogoutAccountRequestHandler) AllowedMethods(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {
    return []string{wm.POST}, req, cxt, 0, nil
}

func (p *LogoutAccountRequestHandler) IsAuthorized(req wm.Request, cxt wm.Context) (bool, string, wm.Request, wm.Context, int, os.Error) {
    lac := cxt.(LogoutAccountContext)
    hasSignature, userId, _, err := apiutil.CheckSignature(p.authDS, req.UnderlyingRequest())
    if !hasSignature || err != nil {
        return hasSignature, "dsocial", req, cxt, http.StatusUnauthorized, err
    }
    accessKey, _ := apiutil.RetrieveAccessKeyFromRequest(p.authDS, req.UnderlyingRequest())
    lac.SetAccessKey(accessKey)
    if userId != "" {
        user, _ := p.ds.RetrieveUserAccountById(userId)
        lac.SetUser(user)
    }
    return true, "", req, cxt, 0, nil
}

func (p *LogoutAccountRequestHandler) Forbidden(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    lac := cxt.(LogoutAccountContext)
    if lac.User() != nil && lac.User().Accessible() {
        return false, req, cxt, 0, nil
    }
    // Cannot find user with specified id
    return true, req, cxt, 0, nil
}

/*
func (p *LogoutAccountRequestHandler) AllowMissingPost(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) MalformedRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) URITooLong(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) DeleteResource(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, http.StatusInternalServerError, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) DeleteCompleted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) PostIsCreate(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) CreatePath(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {
    return "", req, cxt, 0, nil
}
*/

func (p *LogoutAccountRequestHandler) ProcessPost(req wm.Request, cxt wm.Context) (wm.Request, wm.Context, int, http.Header, io.WriterTo, os.Error) {
    var err os.Error
    var code int
    var headers http.Header
    var writerTo io.WriterTo
    lac := cxt.(LogoutAccountContext)
    if lac.AccessKey() != nil {
        _, err = p.authDS.DeleteAccessKey(lac.AccessKey().Id)
    }
    httpHeaders := apiutil.AddNoCacheHeaders(nil)
    if err != nil {
        code, headers, writerTo = apiutil.OutputErrorMessage("Unable to process logout request", nil, http.StatusInternalServerError, httpHeaders)
    } else {
        code, headers, writerTo = apiutil.OutputJSONObject(nil, nil, "", http.StatusOK, httpHeaders)
    }
    return req, cxt, code, headers, writerTo, nil
}

func (p *LogoutAccountRequestHandler) ContentTypesProvided(req wm.Request, cxt wm.Context) ([]wm.MediaTypeHandler, wm.Request, wm.Context, int, os.Error) {
    lac := cxt.(LogoutAccountContext)
    var err os.Error
    if lac.AccessKey() != nil {
        _, err = p.authDS.DeleteAccessKey(lac.AccessKey().Id)
    }
    if err != nil {
        return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(nil, nil, "")}, req, lac, http.StatusInternalServerError, err
    }
    return []wm.MediaTypeHandler{apiutil.NewJSONMediaTypeHandler(nil, nil, "")}, req, lac, 0, nil
}

/*
func (p *LogoutAccountRequestHandler) ContentTypesAccepted(req wm.Request, cxt wm.Context) ([]wm.MediaTypeInputHandler, wm.Request, wm.Context, int, os.Error) {
    return []wm.MediaTypeInputHandler{}, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) IsLanguageAvailable(languages []string, req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LogoutAccountRequestHandler) CharsetsProvided(charsets []string, req wm.Request, cxt wm.Context) ([]CharsetHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LogoutAccountRequestHandler) EncodingsProvided(encodings []string, req wm.Request, cxt wm.Context) ([]EncodingHandler, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LogoutAccountRequestHandler) Variances(req wm.Request, cxt wm.Context) ([]string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *LogoutAccountRequestHandler) IsConflict(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
  return false, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) MultipleChoices(req wm.Request, cxt wm.Context) (bool, http.Header, wm.Request, wm.Context, int, os.Error) {
    return false, nil, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) PreviouslyExisted(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LogoutAccountRequestHandler) MovedPermanently(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LogoutAccountRequestHandler) MovedTemporarily(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *LogoutAccountRequestHandler) LastModified(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {
    return nil, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) Expires(req wm.Request, cxt wm.Context) (*time.Time, wm.Request, wm.Context, int, os.Error) {

}
*/
/*
func (p *LogoutAccountRequestHandler) GenerateETag(req wm.Request, cxt wm.Context) (string, wm.Request, wm.Context, int, os.Error) {

}
*/

/*
func (p *LogoutAccountRequestHandler) FinishRequest(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return true, req, cxt, 0, nil
}
*/

/*
func (p *LogoutAccountRequestHandler) ResponseIsRedirect(req wm.Request, cxt wm.Context) (bool, wm.Request, wm.Context, int, os.Error) {
    return false, req, cxt, 0, nil
}
*/

func (p *LogoutAccountRequestHandler) HasRespBody(req wm.Request, cxt wm.Context) bool {
    return true
}
