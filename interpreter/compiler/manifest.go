package compiler

import (
    "fmt"
    "errors"
    "strings"
    "path/filepath"
    "encoding/json"
)


const ManifestFilenameSuffix = ".manifest.json"

type Manifest struct {
    ProjectFiles     [] string
    DependencyFiles  [] string
}
func ReadManifest(file string, fs FileSystem) (Manifest, error) {
    var manifest Manifest
    if strings.HasSuffix(file, SourceFilenameSuffix) {
        // standalone script
        manifest.ProjectFiles = [] string {
            filepath.Base(file),
        }
    } else if strings.HasSuffix(file, ManifestFilenameSuffix) {
        var content, err = ReadFile(file, fs)
        if err != nil { return Manifest{},
            fmt.Errorf("unable to open manifest file: %w", err)
        }
        { var err = json.Unmarshal(content, &manifest)
        if err != nil { return Manifest{},
            fmt.Errorf("unable to parse manifest file: %w", err)
        }}
    } else {
        return Manifest{}, errors.New("unsupported file type: " + file)
    }
    var dir = filepath.Dir(file)
    for i := range manifest.ProjectFiles {
        f := &(manifest.ProjectFiles[i])
        *f = filepath.Join(dir, *f)
    }
    for i := range manifest.DependencyFiles {
        var f = &(manifest.DependencyFiles[i])
        *f = filepath.Join(dir, *f)
    }
    return manifest, nil
}


