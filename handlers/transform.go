package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/data"
	"github.com/mholzen/information/triples/html"
	"github.com/mholzen/information/triples/transforms"
	"github.com/mholzen/information/triples/transforms/node_link"
	"github.com/russross/blackfriday/v2"
)

type Payload struct {
	Content string      `json:"content"`
	Data    interface{} `json:"data"`
	// TODO: use a context to store the map of handlers
}

type Transform (func(Payload) (Payload, error))

func ToHtml(input Payload) (Payload, error) {
	text, ok := input.Data.(string)
	if !ok {
		return input, fmt.Errorf("cannot convert '%T' to string", input.Data)
	}
	htmlContent := blackfriday.Run([]byte(text))
	response := Payload{
		Content: "text/html",
	}
	response.Data = string(htmlContent)
	return response, nil
}

func DirEntries(fileInfo FileInfo) ([]FileInfo, error) {
	entries, err := os.ReadDir(fileInfo.Name)
	if err != nil {
		return nil, err
	}
	list := []FileInfo{}

	for _, entry := range entries {

		stat, err := os.Stat(filepath.Join(fileInfo.Name, entry.Name()))

		f := FileInfo{
			Name:    entry.Name(),
			Size:    stat.Size(),
			Mode:    stat.Mode(),
			ModTime: stat.ModTime(),
			IsDir:   entry.IsDir(),
			Error:   err,
		}
		list = append(list, f)
	}
	return list, nil
}

func TextString(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	return string(bytes), err
}

func ToText(input Payload) (string, error) {
	if strings.HasPrefix(input.Content, "text/") {
		return input.Data.(string), nil
	}
	fileInfo, ok := input.Data.(FileInfo)
	if !ok {
		return "", fmt.Errorf("cannot convert '%T' to FileInfo", input.Data)
	}

	return TextString(fileInfo.Name)
}

func ToTextPayload(input Payload) (Payload, error) {
	text, err := ToText(input)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "text/plain", // TODO: use the content type of the input
		Data:    text,
	}, nil
}

func ToContent(input Payload) (Payload, error) {
	fileInfo, ok := input.Data.(FileInfo)
	if !ok {
		return input, fmt.Errorf("cannot convert '%T' to FileInfo", input.Data)
	}

	if fileInfo.IsDir {
		entries, err := DirEntries(fileInfo)
		if err != nil {
			return input, err
		}
		res := Payload{
			Content: "application/json+[]fileinfo",
			Data:    entries,
		}
		return res, nil
	}

	file, err := os.Open(fileInfo.Name)
	if err != nil {
		return input, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/octet-stream",
		Data:    bytes,
	}, nil
}

func ToHeader(input Payload) (Payload, error) {
	fileInfo, ok := input.Data.(FileInfo)
	if !ok {
		return input, fmt.Errorf("cannot convert '%T' to FileInfo", input.Data)
	}
	header, err := fileInfo.header()
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/octet-stream",
		Data:    header,
	}, nil
}

type MimeType string

func NewMimeType(input Payload) (MimeType, error) {
	switch data := input.Data.(type) {
	case FileInfo:
		res, err := data.ContentType()
		if err != nil {
			return "", err
		}
		return MimeType(res), nil
	case string:
		return MimeType(http.DetectContentType([]byte(data))), nil
	default:
		return "", fmt.Errorf("cannot convert '%T' to MimeType", input.Data)
	}
}

func ToMimeType(input Payload) (Payload, error) {
	mimeType, err := NewMimeType(input)
	if err != nil {
		return input, err
	}

	return Payload{
		Content: "text/plain",
		Data:    mimeType,
	}, nil
}

func ToTriples(input Payload) (*triples.Triples, error) {
	if strings.HasPrefix(input.Content, "application/json+triples") { // TODO: decide whether we test on type of .Data or value of .Content
		return input.Data.(*triples.Triples), nil
	}
	text, err := ToText(input)
	if err != nil {
		return nil, err
	}
	mimeType, err := NewMimeType(input)
	if err != nil {
		return nil, err
	}
	parser, err := transforms.NewParserFromContentType(string(mimeType), strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	tr := triples.NewTriples()
	tr.Transform(parser.Transformer)
	return tr, nil
}

func ToTriplesPayload(input Payload) (Payload, error) {
	res, err := ToTriples(input)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/json+triples",
		Data:    res,
	}, nil
}

func ToNodeLink(input Payload) (node_link.NodeLink, error) {
	tpls, err := ToTriples(input)
	if err != nil {
		return node_link.NodeLink{}, err
	}

	tr := node_link.NewNodeLinkTransformer()
	err = tpls.Transform(tr.Transformer)
	if err != nil {
		return tr.Result, err
	}
	return tr.Result, nil
}

func ToNodeLinkPayload(input Payload) (Payload, error) {
	res, err := ToNodeLink(input)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/json+nodelink",
		Data:    res,
	}, nil
}

func ToGraphPayload(input Payload) (Payload, error) {
	tpls, ok := input.Data.(*triples.Triples)
	if !ok {
		return input, fmt.Errorf("cannot convert '%T' to *triples.Triples", input.Data)
	}

	tr := node_link.NewNodeLinkTransformer()
	err := tpls.Transform(tr.Transformer)
	if err != nil {
		return input, err
	}
	// Should render an HTML template
	return Payload{
		Content: "application/json+triples",
		Data:    tr.Result,
	}, nil
}

func ToTableDefinition(input Payload) (*triples.Triples, error) {
	src, err := ToTriples(input)
	if err != nil {
		return nil, err
	}
	return src.Map(transforms.PredicatesSortedByString)
}

// TODO: if TableDefinition is a transformer, we should need a specific Payload method
func ToTableDefinitionPayload(input Payload) (Payload, error) {
	tr, err := ToTableDefinition(input)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/json+triples+TableDefinition",
		Data:    tr,
	}, nil
}

func ToListPayload(input Payload) (Payload, error) {
	src, err := ToTriples(input)
	if err != nil {
		return input, err
	}
	tr := transforms.NewListGenerator()
	err = src.Transform(tr.Transformer)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "text/html",
		Data:    (*tr.Result).(triples.StringNode).Value,
	}, nil
}

func ToRowsPayload(input Payload) (Payload, error) {
	res, err := ToTriples(input)
	if err != nil {
		return input, err
	}
	tripleMatch, err := transforms.NewTripleMatchFromTriples(transforms.RowQuery())
	if err != nil {
		return input, err
	}
	res, err = res.Map(transforms.Filter(tripleMatch))
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/json+triples",
		Data:    res,
	}, nil
}

func ToTableTransformPayload(input Payload) (Payload, error) {
	res, err := ToTriples(input)
	if err != nil {
		return input, err
	}

	res, err = res.Map(transforms.Table)
	if err != nil {
		return input, err
	}
	return Payload{
		Content: "application/json+triples",
		Data:    res,
	}, nil
}

func ToHtmlTransformPayload(input Payload) (Payload, error) {
	res, err := ToTriples(input)
	if err != nil {
		return input, err
	}

	res, err = res.Map(transforms.HtmlTable)
	if err != nil {
		return input, err
	}

	h, err := html.FromTriples(res)
	if err != nil {
		return input, err
	}

	return Payload{
		Content: "text/html",
		Data:    string(h),
	}, nil
}

func ToDataPayload(input Payload) (Payload, error) {
	return Payload{
		Content: "application/json+triples",
		Data:    data.Data,
	}, nil
}
