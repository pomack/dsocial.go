package apiutil

import (
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    //"log"
    //"os"
    "url"
)

type UrlEncodedInputHandler interface {
    HandleUrlEncodedInputHandler(req wm.Request, cxt wm.Context, inputObj url.Values) (int, http.Header, io.WriterTo)
}

type UrlEncodedMediaTypeInputHandler struct {
    charset             string
    language            string
    handler             UrlEncodedInputHandler
    reader              io.Reader
    writtenStatusHeader bool
}

func NewUrlEncodedMediaTypeInputHandler(charset, language string, handler UrlEncodedInputHandler, reader io.Reader) *UrlEncodedMediaTypeInputHandler {
    return &UrlEncodedMediaTypeInputHandler{
        charset:  charset,
        language: language,
        handler:  handler,
        reader:   reader,
    }
}

func (p *UrlEncodedMediaTypeInputHandler) MediaTypeInput() string {
    return wm.MIME_TYPE_JSON
}

func (p *UrlEncodedMediaTypeInputHandler) MediaTypeHandleInputFrom(req wm.Request, cxt wm.Context) (int, http.Header, io.WriterTo) {
    defer func() {
        if p.reader != nil {
            if closer, ok := p.reader.(io.Closer); ok {
                closer.Close()
            }
        }
    }()
    //log.Printf("[UEMTIH]: Calling OutputTo with reader %v\n", p.reader)
    if p.reader == nil {
        return p.handler.HandleUrlEncodedInputHandler(req, cxt, nil)
    }
    m := req.Form()
    if m == nil || len(m) == 0 {
        if err := req.ParseForm(); err != nil {
            return OutputErrorMessage(err.String(), nil, http.StatusBadRequest, nil)
        }
        m = req.Form()
    }
    return p.handler.HandleUrlEncodedInputHandler(req, cxt, m)
}

func (p *UrlEncodedMediaTypeInputHandler) MediaTypeInputHandler() wm.MediaTypeInputHandler {
    return p
}
