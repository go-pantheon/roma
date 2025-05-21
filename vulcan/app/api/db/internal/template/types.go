package template

// TemplateData holds all the data needed to execute the proto pool template.
// It includes package name, import paths, and message definitions.

type Data struct {
	Files []*File
}

type File struct {
	Dir      string
	FileName string
	Package  string
	Imports  []*Import
	Messages []*Message
}

// Import represents a single import statement in the generated Go file.
// It includes the alias and the path of the package to be imported.
type Import struct {
	Alias string
	Path  string
}

// Message represents a Protobuf message that needs a sync.Pool.
// It includes the message name, its Go package, and a list of its fields.
type Message struct {
	Name      string
	GoPackage string
	Fields    []*Field
}

// Field represents a field within a Protobuf message.
// It contains information about whether the field is a map or repeated,
// and if its elements are messages themselves (requiring pooling).
type Field struct {
	Name           string
	Type           string // Name of the message type for message's value, if applicable
	IsMessage      bool   // True if the message's value is a message type
	IsRepeated     bool
	IsMap          bool
	KeyType        string // Name of the message type for map's key, if applicable
	ValueIsMessage bool   // True if the map/repeated's value is a message type
	ValueType      string // Name of the message type for map/repeated's value, if applicable
}
