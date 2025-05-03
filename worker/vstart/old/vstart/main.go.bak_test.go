package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

type GitRemoteUrl struct {
	Path      string
	UrlBefore string
	Url       string
	Command   string
}

func TestGitRemoteSet(t *testing.T) {
	// t.Skip()
	// git remote set-url origin https://username:token@github.com/username/repository
	// git clone https://github.com/ddkwork/golibrary
	// https://kkgithub.com/ddkwork/golibrary
	// git clone --recursive https://github.com/HyperDbg/HyperDbg.git
	// git remote set-url origin https://github.com/HyperDbg/HyperDbg.git

	var remoteUrls []GitRemoteUrl
	filepath.Walk("D:\\workspace\\workspace", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".git") {
			mylog.Check(os.Chdir(filepath.Dir(path)))
			session := stream.RunCommand("git remote get-url origin")
			u := session.Output.String()
			u = strings.ReplaceAll(u, "gitlab.com", "github.com")
			u = strings.ReplaceAll(u, "gitee.com", "github.com")
			u = strings.ReplaceAll(u, "git.homegu.com", "github.com")
			u = strings.ReplaceAll(u, "github.com", "ddkwork:xxxxxxxx@github.com")
			u = strings.TrimSuffix(u, ".git")
			u = strings.TrimSuffix(u, "/")
			u += ".git"

			before, after, found := strings.Cut(u, "github.com")
			if found {
				before = "https://"
				u = before + "ddkwork:xxxxxx@" + after
			}
			if !strings.Contains(u, "github.com") {
				u = strings.ReplaceAll(u, "@", "@github.com/")
			}
			u = strings.ReplaceAll(u, `@github.com//`, `@github.com/`)

			remoteUrls = append(remoteUrls, GitRemoteUrl{
				Path:      path,
				UrlBefore: session.Output.String(),
				Url:       u,
				Command:   "git remote set-url origin " + u,
				// git remote set-url origin https://github.com/ddkwork/gui
				// git remote set-url origin https://github.com/HyperDbg/gui
			})
		}
		return err
	})
	mylog.Check(os.Chdir("D:\\workspace\\workspace\\rep\\demo\\vstart"))
	stream.MarshalJsonToFile(remoteUrls, "remoteUrls.json")
	for _, v := range remoteUrls {
		mylog.Check(os.Chdir(filepath.Dir(v.Path)))
		stream.RunCommand(v.Command)
	}
}
