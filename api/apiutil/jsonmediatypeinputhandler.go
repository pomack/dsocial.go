package apiutil

import (
    "github.com/pomack/jsonhelper.go/jsonhelper"
    wm "github.com/pomack/webmachine.go/webmachine"
    "http"
    "io"
    "json"
    //"log"
    "os"
)

type JSONObjectInputHandler interface {
    HandleJSONObjectInputHandler(req wm.Request, cxt wm.Context, writer io.Writer, inputObj jsonhelper.JSONObject) (int, http.Header, os.Error)
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

func (p *JSONMediaTypeInputHandler) MediaType() string {
    return wm.MIME_TYPE_JSON
}

func (p *JSONMediaTypeInputHandler) OutputTo(req wm.Request, cxt wm.Context, writer io.Writer) (int, http.Header, os.Error) {
    defer func() {
        if p.reader != nil {
            if closer, ok := p.reader.(io.Closer); ok {
                closer.Close()
            }
        }
    }()
    //log.Printf("[JSONMTIH]: Calling OutputTo with reader %v\n", p.reader)
    if p.reader == nil {
        return p.handler.HandleJSONObjectInputHandler(req, cxt, writer, nil)
    }
    obj := jsonhelper.NewJSONObject()
    dec := json.NewDecoder(p.reader)
    err := dec.Decode(&obj)
    if err != nil {
        headers := make(http.Header)
        //headers.Set("Content-Type", wm.MIME_TYPE_JSON)
        m := jsonhelper.NewJSONObject()
        w := json.NewEncoder(writer)
        m.Set("status", "error")
        m.Set("message", err.String())
        m.Set("result", nil)
        w.Encode(m)
        return 500, headers, err
    }
    return p.handler.HandleJSONObjectInputHandler(req, cxt, writer, obj)
}

func (p *JSONMediaTypeInputHandler) MediaTypeInputHandler() wm.MediaTypeInputHandler {
    return p
}
