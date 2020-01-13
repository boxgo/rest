package rest

import (
	"regexp"
	"strings"
)

type (
	// Spec OpenAPI Specification
	Spec struct {
		Swagger      string       `json:"swagger,omitempty"`
		Info         Info         `json:"info,omitempty"`
		Host         string       `json:"host,omitempty"`
		BasePath     string       `json:"basePath,omitempty"`
		Schemes      []string     `json:"schemes,omitempty"`
		Consumes     []string     `json:"consumes,omitempty"`
		Produces     []string     `json:"produces,omitempty"`
		Paths        Paths        `json:"paths,omitempty"`
		Tags         []Tag        `json:"tags,omitempty"`
		ExternalDocs ExternalDocs `json:"externalDocs,omitempty"`
	}

	// Info Provides metadata about the API. The metadata can be used by the clients if needed.
	Info struct {
		Title          string  `json:"title,omitempty"`
		Description    string  `json:"description,omitempty"`
		TermsOfService string  `json:"termsOfService,omitempty"`
		Contact        Contact `json:"contact,omitempty"`
		License        License `json:"license,omitempty"`
		Version        string  `json:"version,omitempty"`
	}

	// Contact information for the exposed API.
	Contact struct {
		Name  string `json:"name,omitempty"`
		URL   string `json:"url,omitempty"`
		Email string `json:"email,omitempty"`
	}

	// License information for the exposed API.
	License struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	}

	// Paths The available paths and operations for the API.
	Paths map[string]PathItem

	// PathItem Describes the operations available on a single path
	PathItem struct {
		Get     *Operation `json:"get,omitempty"`
		Put     *Operation `json:"put,omitempty"`
		Post    *Operation `json:"post,omitempty"`
		Delete  *Operation `json:"delete,omitempty"`
		Options *Operation `json:"options,omitempty"`
		Head    *Operation `json:"head,omitempty"`
		Patch   *Operation `json:"patch,omitempty"`
	}

	// Operation Describes a single API operation on a path.
	Operation struct {
		Tags         []string     `json:"tags,omitempty"`
		Summary      string       `json:"summary,omitempty"`
		Description  string       `json:"description,omitempty"`
		ExternalDocs ExternalDocs `json:"externalDocs,omitempty"`
		OperationID  string       `json:"operationId,omitempty"`
		Consumes     []string     `json:"consumes,omitempty"`
		Produces     []string     `json:"produces,omitempty"`
		Parameters   Parameters   `json:"parameters,omitempty"`
		Responses    Responses    `json:"responses,omitempty"`
		Schemes      []string     `json:"schemes,omitempty"`
		Deprecated   bool         `json:"deprecated,omitempty"`
		Security     interface{}  `json:"security,omitempty"`
	}

	// Parameter Describes a single operation parameter.
	Parameter struct {
		Name             string        `json:"name,omitempty"`
		In               string        `json:"in,omitempty"`
		Description      string        `json:"description,omitempty"`
		Required         bool          `json:"required,omitempty"`
		Schema           Schema        `json:"schema,omitempty"`
		Type             string        `json:"type,omitempty"`
		Format           string        `json:"format,omitempty"`
		AllowEmptyValue  bool          `json:"allowEmptyValue,omitempty"`
		Items            Items         `json:"items,omitempty"`
		CollectionFormat string        `json:"collectionFormat,omitempty"`
		Default          interface{}   `json:"default,omitempty"`
		Maximum          float64       `json:"maximum,omitempty"`
		ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
		Minimum          float64       `json:"minimum,omitempty"`
		ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
		MaxLength        uint          `json:"maxLength,omitempty"`
		MinLength        uint          `json:"minLength,omitempty"`
		Pattern          string        `json:"pattern,omitempty"`
		MaxItems         uint          `json:"maxItems,omitempty"`
		MinItems         uint          `json:"minItems,omitempty"`
		UniqueItems      bool          `json:"uniqueItems,omitempty"`
		Enum             []interface{} `json:"enum,omitempty"`
		MultipleOf       uint          `json:"multipleOf,omitempty"`
	}
	// Parameters Parameters
	Parameters []Parameter

	// Items A limited subset of JSON-Schema's items object. It is used by parameter definitions that are not located in "body".
	Items struct {
		Type             string        `json:"type,omitempty"`
		Format           string        `json:"format,omitempty"`
		Description      string        `json:"description,omitempty"`
		Items            interface{}   `json:"items,omitempty"`
		CollectionFormat string        `json:"collectionFormat,omitempty"`
		Default          interface{}   `json:"default,omitempty"`
		Maximum          float64       `json:"maximum,omitempty"`
		ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
		Minimum          float64       `json:"minimum,omitempty"`
		ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
		MaxLength        uint          `json:"maxLength,omitempty"`
		MinLength        uint          `json:"minLength,omitempty"`
		Pattern          string        `json:"pattern,omitempty"`
		MaxItems         uint          `json:"maxItems,omitempty"`
		MinItems         uint          `json:"minItems,omitempty"`
		UniqueItems      bool          `json:"uniqueItems,omitempty"`
		Enum             []interface{} `json:"enum,omitempty"`
		MultipleOf       uint          `json:"multipleOf,omitempty"`
		Ref              string        `json:"$ref,omitempty"`
	}

	// Schema The Schema Object allows the definition of input and output data types.
	Schema struct {
		Title                string           `json:"title,omitempty"`
		Format               string           `json:"format,omitempty"`
		Description          string           `json:"description,omitempty"`
		Default              interface{}      `json:"default,omitempty"`
		MultipleOf           uint             `json:"multipleOf,omitempty"`
		Maximum              float64          `json:"maximum,omitempty"`
		ExclusiveMaximum     bool             `json:"exclusiveMaximum,omitempty"`
		Minimum              float64          `json:"minimum,omitempty"`
		ExclusiveMinimum     bool             `json:"exclusiveMinimum,omitempty"`
		MaxLength            uint             `json:"maxLength,omitempty"`
		MinLength            uint             `json:"minLength,omitempty"`
		Pattern              string           `json:"pattern,omitempty"`
		MaxItems             uint             `json:"maxItems,omitempty"`
		MinItems             uint             `json:"minItems,omitempty"`
		UniqueItems          bool             `json:"uniqueItems,omitempty"`
		MaxProperties        uint             `json:"maxProperties,omitempty"`
		MinProperties        uint             `json:"minProperties,omitempty"`
		Required             bool             `json:"required,omitempty"`
		Enum                 []interface{}    `json:"enum,omitempty"`
		Type                 string           `json:"type,omitempty"`
		Items                interface{}      `json:"items,omitempty"`
		AllOf                []Items          `json:"allOf,omitempty"`
		Properties           map[string]Items `json:"properties,omitempty"`
		AdditionalProperties map[string]Items `json:"additionalProperties,omitempty"`
	}

	// Response Describes a single response from an API Operation.
	Response struct {
		Description string `json:"description,omitempty"`
		Schema      Schema `json:"schema,omitempty"`
	}

	// Responses Responses
	Responses map[string]Response

	// Tag Allows adding meta data to a single tag that is used by the Operation Object. It is not mandatory to have a Tag Object per tag used there.
	Tag struct {
		Name         string       `json:"name,omitempty"`
		Description  string       `json:"description,omitempty"`
		ExternalDocs ExternalDocs `json:"externalDocs,omitempty"`
	}

	// ExternalDocs Additional external documentation.
	ExternalDocs struct {
		Description string `json:"description,omitempty"`
		URL         string `json:"url,omitempty"`
	}
)

var (
	replaceRegexp *regexp.Regexp
)

func init() {
	var err error
	replaceRegexp, err = regexp.Compile(":(\\w+)$")
	if err != nil {
		panic(err)
	}
}

// DescribeAPI describe a api
func (spec *Spec) DescribeAPI(path, method string, op Operation) {
	if len(spec.Paths) == 0 {
		spec.Paths = make(map[string]PathItem)
	}

	path = spec.gin2oai(path)

	pathItem := spec.Paths[path]
	switch strings.ToLower(method) {
	case "get":
		pathItem.Get = &op
	case "put":
		pathItem.Put = &op
	case "post":
		pathItem.Post = &op
	case "delete":
		pathItem.Delete = &op
	case "options":
		pathItem.Options = &op
	case "head":
		pathItem.Head = &op
	case "patch":
		pathItem.Patch = &op
	}

	spec.Paths[path] = pathItem
}

func (spec *Spec) gin2oai(path string) string {
	paths := []string{}
	for _, p := range strings.Split(path, "/") {
		paths = append(paths, replaceRegexp.ReplaceAllString(p, "{$1}"))
	}

	return strings.Join(paths, "/")
}

// NewSpec new a spec
func NewSpec() *Spec {
	return &Spec{
		Swagger: "2.0",
	}
}
