package fileprocessor

// File represents the contents and the path to a file
type File struct {
	Name string
	Data []byte
}

// FileProcessor can process a template file and produce an output file
type FileProcessor interface {
	ProcessFile(file File) (*File, error)
}
