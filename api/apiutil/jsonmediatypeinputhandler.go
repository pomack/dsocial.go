package apiutil

import (
    "encoding/json"
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "io"
    "net/http"
    //"log"
)

type JSONObjectInputHandler interface {
    HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, inputObj jsonhelper.JSONObject) (int, http.Header, io.WriterTo)
}

type JSONMediaTypeInputHandler struct {
    charset             string
    language            string
    handler             JSONObjectInputHandler
    reader              io.Reader
    writtenStatusHeader bool
}

func NewJSONMediaTypeInputHandler(charset, language string, handler JSONObjectInputHandler, reader io.Reader) *JSONMediaTypeInputHandler {
    return &JSONMediaTypeInputHandler{
        charset:  charset,
        language: language,
        handler:  handler,
        reader:   reader,
    }
}

func (p *JSONMediaTypeInputHandler) MediaTypeInput() string {
    return wm.MIME_TYPE_JSON
}

func (p *JSONMediaTypeInputHandler) MediaTypeHandleInputFrom(req wm.Request, cxt wm.Context) (int, http.Header, io.WriterTo) {
    defer func() {
        if p.reader != nil {
            if closer, ok := p.reader.(io.Closer); ok {
                closer.Close()
            }
        }
    }()
    //log.Printf("[JSONMTIH]: Calling OutputTo with reader %v\n", p.reader)
    if p.reader == nil {
        return p.handler.HandleJSONObjectInputHandler(req, cxt, nil)
    }
    obj := jsonhelper.NewJSONObject()
    dec := json.NewDecoder(p.reader)
    err := dec.Decode(&obj)
    if err != nil {
        headers := make(http.Header)
        return OutputErrorMessage(err.Error(), nil, 500, headers)
    }
    return p.handler.HandleJSONObjectInputHandler(req, cxt, obj)
}

func (p *JSONMediaTypeInputHandler) MediaTypeInputHandler() wm.MediaTypeInputHandler {
    return p
}
