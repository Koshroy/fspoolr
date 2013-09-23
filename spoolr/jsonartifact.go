package spoolr

import (
	"io/ioutil"
    "log"
	"path"
    "path/filepath"
    "os/exec"
    "encoding/json"
    "strings"
)

type jsonartifact struct {
	rootDir string    `json:"-"`
	NameJson string       `json:"name"`
	BuildCmdJson string   `json:"build_cmd"`
	FilesJson []string    `json:"files"`
	TargetJson TargetFile `json:"target"`
    absFilesJson []string
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
    j.absFilesJson = j.absFiles()
    return j, nil
}

func (j *jsonartifact) absFiles() []string {
    retval := make([]string, 0, len(j.FilesJson))
    for _, elem := range j.FilesJson {
        absPath, err := filepath.Abs(path.Join(j.rootDir, elem))
        if err != nil {
            log.Println("error finding absolute path for relative path", elem)
            return retval
        }
        retval = append(retval, absPath)
    }
    return retval
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
    buildStrSplit := strings.Split(j.BuildCmdJson, " ")
    cmd := exec.Command(buildStrSplit[0], buildStrSplit[1:]...)
    err := cmd.Run()
    return err
}

func (j *jsonartifact) Name() string {
    return j.NameJson
}

func (j *jsonartifact) RootDir() string {
    return j.rootDir
}

func (j *jsonartifact) BuildCmd() string {
    return j.BuildCmdJson
}

func (j *jsonartifact) Files() []string {
    return j.absFilesJson
}

func (j *jsonartifact) Target() TargetFile {
    return j.TargetJson
}

func (j *jsonartifact) String() string {
    return "jsonartifact at dir " + j.rootDir
}