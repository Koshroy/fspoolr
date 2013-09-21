package spoolr

type TargetFile struct {
	File string      `json:"file"`
	MimeType string  `json:"mime_type"`
}

type Artifact interface {
	Name() string
	RootDir() string
	BuildCmd() string
	Files() []string
	Target() TargetFile

	Reload() error
	Rebuild() error

	String() string
}
