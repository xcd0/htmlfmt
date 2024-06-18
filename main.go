package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
	"github.com/yosssi/gohtml"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する
	if len(os.Args) == 1 {
		// 標準入力から読み取り、標準出力に出力する。
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(errors.Errorf("%v", err))
		}
		t := string(b)
		x := gohtml.Format(t)
		fmt.Println(x)
		return
	}
	args := &Args{}
	var err error
	parser, err = arg.NewParser(arg.Config{Program: GetFileNameWithoutExt(os.Args[0]), IgnoreEnv: false}, args)
	if err != nil {
		ShowHelp(fmt.Sprintf("%v", errors.Errorf("%v", err)))
	}
	if err := parser.Parse(os.Args[1:]); err != nil {
		if err.Error() == "help requested by user" {
			ShowHelp("")
		} else if err.Error() == "version requested by user" {
			ShowVersion()
		} else {
			panic(errors.Errorf("%v", err))
		}
	}
	//if args.Version || args.VersionSub != nil {
	if args.Version {
		ShowVersion()
	}
	if args.Debug {
		args.Print()
	}
	for _, path := range args.Input {

		x := gohtml.Format(GetText(path))
		//out := replaceExt(path, ".json", ".xml")
		out := path
		os.RemoveAll(out)
		WriteText(out, x)

		if args.Debug {
			fmt.Println(x)
		}
	}
}

func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func GetText(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(errors.Errorf("Error: %v, file: %v", err, path))
	}
	str := string(b)
	return str
}

func WriteText(file, str string) {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		panic(errors.Errorf("%v", err))
	} else {
		if _, err := f.Write([]byte(str)); err != nil {
			panic(errors.Errorf("%v", err))
		}
	}
}

var (
	Version  string = "0.0.1"
	Revision        = func() string { // {{{
		revision := ""
		modified := false
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					//return setting.Value
					revision = setting.Value
					if len(setting.Value) > 7 {
						revision = setting.Value[:7] // 最初の7文字にする
					}
				}
				if setting.Key == "vcs.modified" {
					modified = setting.Value == "true"
				}
			}
		}
		if modified {
			revision = "develop+" + revision
		}
		return revision
	}() // }}}
)

type Args struct {
	Input   []string `arg:"positional"         help:"入力ファイル。"`
	Debug   bool     `arg:"-d,--debug"         help:"デバッグ用。ログが詳細になる。"`
	Version bool     `arg:"-v,--version"       help:"バージョン情報を出力する。"`
	//VersionSub *ArgsVersion `arg:"subcommand:version" help:"バージョン情報を出力する。"`
}
type ArgsVersion struct {
}

func (args *Args) Print() {
	//	log.Printf(`
	//
	// Csv  : %v
	// Row  : %v
	// Col  : %v
	// Grep : %v
	// `, args.Csv, args.Row, args.Col, args.Grep)
}

// ShowHelp() で使う
var parser *arg.Parser

func ShowHelp(post string) {
	buf := new(bytes.Buffer)
	parser.WriteHelp(buf)
	fmt.Printf("%v\n", strings.ReplaceAll(buf.String(), "display this help and exit", "ヘルプを出力する。"))
	if len(post) != 0 {
		fmt.Println(post)
	}
	os.Exit(1)
}
func ShowVersion() {
	if len(Revision) == 0 {
		// go installでビルドされた場合、gitの情報がなくなる。その場合v0.0.0.のように末尾に.がついてしまうのを避ける。
		fmt.Printf("%v version %v\n", GetFileNameWithoutExt(os.Args[0]), Version)
	} else {
		fmt.Printf("%v version %v.%v\n", GetFileNameWithoutExt(os.Args[0]), Version, Revision)
	}
	os.Exit(0)
}
