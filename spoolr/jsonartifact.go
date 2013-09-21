package spoolr

import (
	"io/ioutil"
	"path"
    "os/exec"
    "encoding/json"
    "strings"
    "log"
)

type jsonartifact struct {
	rootDir string    `json:"-"`
	name string       `json:"name"`
	buildCmd string   `json:"build_cmd"`
	files []string    `json:"files"`
	target TargetFile `json:"target"`
}

func NewJsonArtifact(root string) (*jsonartifact, error) {
	artifactJsonPath := path.Join(root, "package.json")
	pkgBytes, err := ioutil.ReadFile(artifactJsonPath)
	if err != nil {
		return nil, err
	}

    j := new(jsonartifact)
    err = json.Unmarshal(pkgBytes, j)
    if err != nil {
        return nil, err
    }
    j.rootDir = root
    log.Println("new json artifact", "[" + j.name + "]")
    return j, nil
}

func (j *jsonartifact) Reload() error {
    pkgBytes, err := ioutil.ReadFile(path.Join(j.rootDir, "package.json"))
    if err != nil {
        return err
    }

    err = json.Unmarshal(pkgBytes, j)
    return err
}

func (j *jsonartifact) Rebuild() error {
    buildStrSplit := strings.Split(j.buildCmd, " ")
    cmd := exec.Command(buildStrSplit[0], buildStrSplit[1:]...)
    err := cmd.Run()
    return err
}

func (j *jsonartifact) Name() string {
    return j.name
}

func (j *jsonartifact) RootDir() string {
    return j.rootDir
}

func (j *jsonartifact) BuildCmd() string {
    return j.buildCmd
}

func (j *jsonartifact) Files() []string {
    return j.files
}

func (j *jsonartifact) Target() TargetFile {
    return j.target
}

func (j *jsonartifact) String() string {
    return "jsonartifact at dir " + j.rootDir
}