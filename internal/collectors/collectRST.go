package collectors

import (
	"checker/internal/parsers/rst"
	"checker/internal/sources"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	iowrap "github.com/spf13/afero"

	log "github.com/sirupsen/logrus"
)

var (
	FS                  iowrap.Fs
	FSUtil              *iowrap.Afero
	basepath            string
	sharedConstantRegex = regexp.MustCompile(`\{\+([[:alnum:]\p{P}\p{S}]+)\+\}`)
)

func init() {
	FS = iowrap.NewOsFs()
	FSUtil = &iowrap.Afero{Fs: FS}
}

func exists(path string) bool {

	if _, err := FS.Stat(path); os.IsNotExist(err) {
		log.Errorf("%s does not exist", path)
		return false
	}
	return true
}

func snootyTomlExists(path string) bool {
	return exists(filepath.Join(path, "snooty.toml"))
}

func sourceDirectoryExists(path string) bool {

	return exists(filepath.Join(path, "source"))
}

func GatherFiles(path string) []string {
	basepath = path
	if !snootyTomlExists(path) || !sourceDirectoryExists(path) {
		log.Panic("snooty.toml or source directory does not exist")
	}

	files := make([]string, 0)

	exts := []string{".rst", ".txt", ".yml", ".yaml"}
	validExt := func(s string) bool {
		for _, ext := range exts {
			if strings.Contains(s, ext) {
				return true
			}
		}
		return false
	}

	err := FSUtil.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "draft" {
			return filepath.SkipDir
		}
		if validExt(filepath.Ext(path)) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return files
}

func gather(files []string, fn func(filename string, data []byte)) {
	for _, file := range files {
		dat, err := FSUtil.ReadFile(file)
		if err != nil {
			log.Panic(err)
		}

		fileName := strings.Replace(file, basepath, "", 1)
		fn(fileName, dat)
	}
}

type RstRoleMap map[rst.RstRole]string

func GatherRoles(files []string) RstRoleMap {
	roles := make(map[rst.RstRole]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, role := range rst.ParseForRoles(data) {
			roles[role] = filename
		}
	})
	return roles
}

func (r *RstRoleMap) Get(key string) (*rst.RstRole, bool) {
	for k := range *r {
		if k.Name == key {
			return &k, true
		}
	}
	return nil, false
}

func (r *RstRoleMap) Union(other RstRoleMap) {
	for k, v := range other {
		(*r)[k] = v
	}
}

func GatherConstants(files []string) map[rst.RstConstant]string {
	consts := make(map[rst.RstConstant]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, con := range rst.ParseForConstants(data) {
			consts[con] = filename
		}
	})
	return consts
}

func GatherHTTPLinks(files []string) map[rst.RstHTTPLink]string {
	links := make(map[rst.RstHTTPLink]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, link := range rst.ParseForHTTPLinks(data) {
			links[link] = filename
		}
	})
	return links
}

type RefTargetMap map[rst.RefTarget]string

func GatherLocalRefs(files []string) RefTargetMap {
	refs := make(map[rst.RefTarget]string, len(files))
	gather(files, func(filename string, data []byte) {
		for _, ref := range rst.ParseForLocalRefs(data) {
			refs[ref] = filename
		}
	})
	return refs
}

func (r *RefTargetMap) Get(ref *rst.RstRole) (*rst.RefTarget, bool) {
	for k := range *r {
		if k.Name == ref.Target {
			return &k, true
		}
	}
	return nil, false
}

func (r *RefTargetMap) Union(other RefTargetMap) {
	for k, v := range other {
		(*r)[k] = v
	}
}
func (r RefTargetMap) SSLtoTLS() RefTargetMap {
	for k, v := range r {
		if strings.Contains(k.Name, "ssl") {
			tlsK := rst.RefTarget{Name: strings.Replace(k.Name, "ssl", "tls", 1)}
			r[tlsK] = v
		}
	}
	return r
}

func GatherSharedIncludes(files []string) []rst.SharedInclude {
	includes := make([]rst.SharedInclude, 0)
	gather(files, func(filename string, data []byte) {
		includes = append(includes, rst.ParseForSharedIncludes(data)...)
	})
	return includes
}

func GatherSharedRefs(input []byte, defs sources.TomlConfig) RstRoleMap {
	roles := make(RstRoleMap, len(input))
	for _, role := range rst.ParseForRoles(input) {
		allFound := sharedConstantRegex.FindAllString(role.Target, -1)
		for _, match := range allFound {
			for _, inner := range sharedConstantRegex.FindAllStringSubmatch(match, -1) {
				role.Target = strings.Replace(role.Target, inner[0], defs.Constants[inner[1]], 1)
			}
		}
		roles[role] = "shared"
	}
	return roles
}

func GatherSharedLocalRefs(input []byte, defs sources.TomlConfig) RefTargetMap {
	refs := make(map[rst.RefTarget]string, len(input))
	for _, ref := range rst.ParseForLocalRefs(input) {
		allFound := sharedConstantRegex.FindAllString(ref.Name, -1)
		for _, match := range allFound {
			for _, inner := range sharedConstantRegex.FindAllStringSubmatch(match, -1) {
				ref.Name = strings.Replace(ref.Name, inner[0], defs.Constants[inner[1]], 1)
			}
		}
		refs[ref] = "shared"
	}
	return refs
}

func (r RstRoleMap) ConvertConstantRefs(defs sources.TomlConfig) RstRoleMap {
	for k, v := range r {
		allFound := sharedConstantRegex.FindAllString(k.Target, -1)
		for _, match := range allFound {
			for _, inner := range sharedConstantRegex.FindAllStringSubmatch(match, -1) {
				delete(r, k)
				k.Target = strings.Replace(k.Target, inner[0], defs.Constants[inner[1]], 1)
				k.Name = strings.Replace(k.Name, inner[0], defs.Constants[inner[1]], 1)
				r[k] = v
			}
		}
	}
	return r
}
