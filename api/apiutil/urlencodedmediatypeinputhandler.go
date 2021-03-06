package apiutil

import (
    wm "github.com/pomack/webmachine.go/webmachine"
    "io"
    "net/http"
    "net/url"
)

type UrlEncodedInputHandler interface {
    HandleUrlEncodedInputHandler(req wm.Request, cxt wm.Context, inputObj url.Values) (int, http.Header, io.WriterTo)
}

type UrlEncodedMediaTypeInputHandler struct {
    charset             string
    language            string
    handler             UrlEncodedInputHandler
    writtenStatusHeader bool
}

func NewUrlEncodedMediaTypeInputHandler(charset, language string, handler UrlEncodedInputHandler) *UrlEncodedMediaTypeInputHandler {
    return &UrlEncodedMediaTypeInputHandler{
        charset:  charset,
        language: language,
        handler:  handler,
    }
}

func (p *UrlEncodedMediaTypeInputHandler) MediaTypeInput() string {
    return wm.MIME_TYPE_JSON
}

func (p *UrlEncodedMediaTypeInputHandler) MediaTypeHandleInputFrom(req wm.Request, cxt wm.Context) (int, http.Header, io.WriterTo) {
    m := req.Form()
    if m == nil || len(m) == 0 {
        if err := req.ParseForm(); err != nil {
            return OutputErrorMessage(err.Error(), nil, http.StatusBadRequest, nil)
        }
        m = req.Form()
    }
    return p.handler.HandleUrlEncodedInputHandler(req, cxt, m)
}

func (p *UrlEncodedMediaTypeInputHandler) MediaTypeInputHandler() wm.MediaTypeInputHandler {
    return p
}
