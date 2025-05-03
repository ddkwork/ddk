package main

import (
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

type (
	Interface interface {
		MakeEnv() string             // 生成env.sh
		UpdateBash() string          // 俩个用户bash文件每次写入第一行刷新
		LocalEnvFineName() string    // env.sh
		SystemEnvFineName() []string //".zshrc", ".bashrc"
		OpenWineRep() (ok bool)      // 开启wine仓库
	}
	object struct{ username string }
)

func New() Interface { return &object{} }

//go:generate  go build -x .

func main() {
	e := New()
	mylog.Info("MakeEnv", e.MakeEnv())
	mylog.Info("env first line", e.UpdateBash())
	e.OpenWineRep()
	mylog.Success("finish")
	select {}
}

func (o *object) MakeEnv() string {
	//bin, ok := goos.GoBin()
	//if !ok {
	//	return ""
	//}
	//env := "export PATH=${PATH}:" + bin
	//f, err2 := os.Create(o.LocalEnvFineName())
	//if !mylog.Check(err2) {
	//	return ""
	//}
	//if !mylog.Error2(f.WriteString(env)) {
	//	return ""
	//}
	//return env
	return ""
}

func (o *object) UpdateBash() string {
	//abs, err := filepath.Abs("env.sh")
	//if ! {
	//	return ""
	//}
	//bash := "source  " + abs
	//for _, s := range o.SystemEnvFineName() {
	//	dir, ok := goos.UserHomeDir()
	//	if !ok {
	//		return ""
	//	}
	//	path := dir + "/" + s
	//	buf, err := os.ReadFile(path)
	//	if err == nil {
	//		mylog.Info("path", path)
	//		b := stream.NewBufferObject(buf)
	//		lines, ok := stream.ToLines(buf)
	//		if !ok {
	//			return ""
	//		}
	//		if strings.Has(lines[0], "source") {
	//			b.Reset()
	//			if !mylog.Error2(b.WriteString(strings.Replace(string(buf), lines[0], bash, 1))) {
	//				return ""
	//			}
	//			if !stream.WriteTruncate(path, b.Bytes()) {
	//				return ""
	//			}
	//			return lines[0]
	//		} else {
	//			NewBuffer := stream.NewBufferObject("")
	//			NewBuffer.WriteStringLn(bash)
	//			if !mylog.Error2(NewBuffer.Write(b.Bytes())) {
	//				return ""
	//			}
	//			stream.WriteTruncate(path, NewBuffer.Bytes())
	//			return lines[0]
	//		}
	//	}
	//}
	return ""
}
func (o *object) LocalEnvFineName() string    { return "env.sh" }
func (o *object) SystemEnvFineName() []string { return []string{".zshrc", ".bashrc"} }

func (o *object) OpenWineRep() (ok bool) {
	if !stream.IsLinux() {
		return true
	}
	pacmanConfName := "/etc/pacman.conf"
	body := stream.NewBuffer("")
	for line := range stream.ReadFileToLines(pacmanConfName) {
		//if strings.Contains(line, "#[multilib]") {
		//	if strings.Contains(lines[i+1], "#Include = /etc/pacman.d/mirrorlist") {
		//		lines[i] = "[multilib]"
		//		lines[i+1] = "Include = /etc/pacman.d/mirrorlist"
		//	}
		//}
		body.WriteStringLn(strings.TrimPrefix(line, "#"))
	}

	stream.WriteTruncate(pacmanConfName, body.Bytes())
	install := `
now you can install wine use theme commands
			sudo pacman -Sy
			yay -S bottles
`
	mylog.Info("install", install)
	return true
}
