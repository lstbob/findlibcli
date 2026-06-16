package resolver

type Library struct {
	Name        string
	Description string
	ImportPath  string
	Language    string
}

type Result struct {
	Library
	Confidence float64
}

type Resolver interface {
	Resolve(language, description string) ([]Result, error)
}
