package http

import "fmt"

const (
	CharsetUTf8 = "utf-8"
)

const (
	BaseContentTypeApplicationJS   = "application/javascript"
	BaseContentTypeApplicationJSON = "application/json"
	BaseContentTypeCSS             = "text/css"
	BaseContentTypeHTML            = "text/html"
	BaseContentTextPlain           = "text/plain"
	BaseContentTypeXML             = "text/xml"
)

var ContentTypeApplicationJSON = getFullContentType(BaseContentTypeApplicationJSON, CharsetUTf8)
var ContentTypeHTML = getFullContentType(BaseContentTypeHTML, CharsetUTf8)
var ContentTypeTextPlain = getFullContentType(BaseContentTextPlain, CharsetUTf8)

func getFullContentType(contentType string, charset string) string {
	return fmt.Sprintf("%s; charset=%s", contentType, charset)
}
