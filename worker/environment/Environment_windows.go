package environment

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ddkwork/golibrary/std/mylog"
	"golang.org/x/sys/windows/registry"
)

// 功能就是刷新windows所有版本的系统path环境变量
// [HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment]
type (
	Interface interface {
		WalkDirs(roots ...string) (ok bool) // 遍历需要加入到path的文件夹集生成目录切片,map+goruntine
		Orig() (ok bool)                    // 读取系统Environment并用map去重+delete bad path
		Update() (ok bool)                  // sort by strings
	}
	object struct {
		paths    []string
		pathsMap map[string]string
		key      registry.Key
	}
)

// ComSpec=%SystemRoot%\system32\cmd.exe
var (
	name    = "Path"
	EnvPath = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	skip    = []string{
		"genx",
		"steam",
		"todo",
		"vmware",
		"ndk",
		"媒体",
		"apk",
		"vs2022",
		"clone",
	}
	defalutRoots = []string{
		"C:\\Users\\Admin\\go\\bin",
		// "D:\\bin",
		// "D:\\codespace",
	}
)

func (o *object) Orig() (ok bool) {
	key := mylog.Check2(registry.OpenKey(registry.LOCAL_MACHINE, EnvPath, registry.ALL_ACCESS))

	o.key = key
	value, _ := mylog.Check3(key.GetStringValue(name))

	split := strings.Split(value, ";")
	o.paths = append(o.paths, split...)
	return true
}

func (o *object) WalkDirs(roots ...string) (ok bool) {
	if len(roots) == 0 {
		roots = defalutRoots
	}
	for _, root := range roots {
		filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
			for _, s := range skip {
				if strings.Contains(path, s) {
					return nil
				}
			}
			ext := filepath.Ext(path)
			switch ext {
			case ".exe", ".cmd", ".bat":
				dir := filepath.Dir(path)
				o.pathsMap[dir] = dir
			}
			return err
		})
	}
	for _, s2 := range o.pathsMap {
		o.paths = append(o.paths, s2)
	}
	return true
}

func (o *object) Update() {
	sort.Strings(o.paths)
	for _, path := range o.paths {
		mylog.Info("update", path)
	}
	mylog.Check(o.key.SetStringValue(name, strings.Join(o.paths, ";")))
}

func New_() *object {
	return &object{
		paths:    make([]string, 0),
		pathsMap: make(map[string]string),
	}
}
