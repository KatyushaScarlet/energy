//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//

package install

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/energye/energy/v2/cmd/internal/command"
	"github.com/energye/energy/v2/cmd/internal/consts"
	"github.com/energye/energy/v2/cmd/internal/env"
	progressbar "github.com/energye/energy/v2/cmd/internal/progress-bar"
	"github.com/energye/energy/v2/cmd/internal/tools"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type downloadInfo struct {
	fileName      string
	frameworkPath string
	downloadPath  string
	url           string
	success       bool
	isSupport     bool
	module        string
}

func Install(c *command.Config) {
	// 初始配置和安装目录
	initInstall(c)
	// 检查环境
	willInstall := checkInstallEnv(c)
	var (
		goRoot                      string
		goSuccessCallback           func()
		cefFrameworkRoot            string
		cefFrameworkSuccessCallback func()
		nsisRoot                    string
		nsisSuccessCallback         func()
		upxRoot                     string
		upxSuccessCallback          func()
	)
	if len(willInstall) > 0 {
		println("Following will be installed")
		for _, name := range willInstall {
			println("\t", name)
		}
		println("Press Enter to start installation...")
		var s string
		fmt.Scanln(&s)
	}
	// 安装Go开发环境
	goRoot, goSuccessCallback = installGolang(c)
	// 设置 go 环境变量
	if goRoot != "" {
		env.SetGoEnv(goRoot)
	}
	// 安装nsis安装包制作工具, 仅windows - amd64
	nsisRoot, nsisSuccessCallback = installNSIS(c)
	// 设置nsis环境变量
	if nsisRoot != "" {
		env.SetNSISEnv(nsisRoot)
	}
	// 安装upx, 内置, 仅windows, linux
	upxRoot, upxSuccessCallback = installUPX(c)
	// 设置upx环境变量
	if upxRoot != "" {
		env.SetUPXEnv(upxRoot)
	}
	// 安装CEF二进制框架
	cefFrameworkRoot, cefFrameworkSuccessCallback = installCEFFramework(c)
	// 设置 energy cef 环境变量
	if cefFrameworkRoot != "" && c.Install.IsSame {
		env.SetEnergyHomeEnv(cefFrameworkRoot)
	}
	// success 输出
	if nsisSuccessCallback != nil || goSuccessCallback != nil || upxSuccessCallback != nil || cefFrameworkSuccessCallback != nil {
		println("-----------------------------------------------------\n-----------------------------------------------------")
		println("Hint: Reopen the cmd window for command to take effect.")
		println("-----------------------------------------------------\n-----------------------------------------------------")
	}
	if nsisSuccessCallback != nil {
		nsisSuccessCallback()
	}
	if goSuccessCallback != nil {
		goSuccessCallback()
	}
	if upxSuccessCallback != nil {
		upxSuccessCallback()
	}
	if cefFrameworkSuccessCallback != nil {
		cefFrameworkSuccessCallback()
	}
	// end
	//if len(willInstall) > 0 {
	//	print("was installed: ")
	//	for i, name := range willInstall {
	//		if i > 0 {
	//			print("|")
	//		}
	//		print(name)
	//	}
	//	println()
	//}
}

func cefInstallPathName(c *command.Config) string {
	if c.Install.IsSame {
		return filepath.Join(c.Install.Path, consts.ENERGY, c.Install.Name)
	} else {
		return filepath.Join(c.Install.Path, consts.ENERGY, fmt.Sprintf("%s_%s%s", c.Install.Name, c.Install.OS, c.Install.Arch))
	}
}

func goInstallPathName(c *command.Config) string {
	return filepath.Join(c.Install.Path, consts.ENERGY, "go")
}

func nsisInstallPathName(c *command.Config) string {
	return filepath.Join(c.Install.Path, consts.ENERGY, "nsis")
}

func upxInstallPathName(c *command.Config) string {
	return filepath.Join(c.Install.Path, consts.ENERGY, "upx")
}

func nsisIsInstall() bool {
	return consts.IsWindows && !consts.IsARM64
}

func upxIsInstall() bool {
	return (consts.IsWindows && !consts.IsARM64) || (consts.IsLinux)
}

// 检查当前环境
//  golang, nsis, cef, upx
//  golang: all os
//  nsis: windows
//  cef: all os
//  upx: windows amd64, 386, linux amd64, arm64
func checkInstallEnv(c *command.Config) (result []string) {
	skip := c.Install.All
	var check = func(chkInstall func() bool, name string, yes func()) {
		if chkInstall() {
			println(" ", name, "installed")
		} else {
			fmt.Printf("  %s: Not installed, install %s ? (Y/n): ", name, name)
			var s string
			if !skip { // 跳过输入Y,
				fmt.Scanln(&s)
			} else {
				s = "y"
			}
			if strings.ToLower(s) == "y" {
				result = append(result, name)
				yes()
			} else {
				println(" ", name, "install skip")
			}
		}
	}

	// go
	check(func() bool {
		return tools.CommandExists("go")
	}, "Golang", func() {
		c.Install.IGolang = true
	})

	// nsis
	check(func() bool {
		if nsisIsInstall() {
			return tools.CommandExists("makensis")
		} else {
			print("  Non Windows amd64 skipping NSIS.")
			return true
		}
	}, "NSIS", func() {
		c.Install.INSIS = true
	})

	// cef
	var cefName = fmt.Sprintf("CEF Framework %s%s", c.Install.OS, c.Install.Arch)
	check(func() bool {
		if c.Install.IsSame {
			// 检查环境变量是否配置
			return tools.CheckCEFDir()
		}
		// 非当月系统架构时检查一下目标安装路径是否已经存在
		var lib = func() string {
			if c.Install.OS.IsWindows() {
				return "libcef.dll"
			} else if c.Install.OS.IsLinux() {
				return "libcef.so"
			} else if c.Install.OS.IsDarwin() {
				return "cef_sandbox.a"
			}
			return ""
		}()
		if lib != "" {
			s := filepath.Join(cefInstallPathName(c), lib)
			return tools.IsExist(s)
		} else {
			print("Unsupported system architecture")
			return true
		}
	}, cefName, func() {
		c.Install.ICEF = true
	})

	// upx
	check(func() bool {
		if upxIsInstall() {
			return tools.CommandExists("upx")
		} else {
			print("  Non Windows skipping UPX.")
			return true
		}
	}, "UPX", func() {
		c.Install.IUPX = true
	})
	return
}

func initInstall(c *command.Config) {
	if c.Install.Path == "" {
		// current dir
		c.Install.Path = c.Wd
	}
	if c.Install.Version == "" {
		// latest
		c.Install.Version = "latest"
	}
	if c.Install.OS == "" {
		c.Install.OS = command.OS(runtime.GOOS)
	}
	if c.Install.Arch == "" {
		c.Install.Arch = command.Arch(runtime.GOARCH)
	}
	if string(c.Install.OS) == runtime.GOOS && string(c.Install.Arch) == runtime.GOARCH {
		c.Install.IsSame = true
	}
	// 创建安装目录
	os.MkdirAll(c.Install.Path, fs.ModePerm)
	os.MkdirAll(cefInstallPathName(c), fs.ModePerm)
	os.MkdirAll(goInstallPathName(c), fs.ModePerm)
	if nsisIsInstall() {
		os.MkdirAll(nsisInstallPathName(c), fs.ModePerm)
	}
	if upxIsInstall() {
		os.MkdirAll(upxInstallPathName(c), fs.ModePerm)
	}
	os.MkdirAll(filepath.Join(c.Install.Path, consts.FrameworkCache), fs.ModePerm)
}

func filePathInclude(compressPath string, files ...any) (string, bool) {
	if len(files) == 0 {
		return compressPath, true
	} else {
		for _, file := range files {
			f := file.(string)
			tIdx := strings.LastIndex(f, "/*")
			if tIdx != -1 {
				f = f[:tIdx]
				if f[0] == '/' {
					if strings.Index(compressPath, f[1:]) == 0 {
						return compressPath, true
					}
				}
				if strings.Index(compressPath, f) == 0 {
					return strings.Replace(compressPath, f, "", 1), true
				}
			} else {
				if f[0] == '/' {
					if compressPath == f[1:] {
						return f, true
					}
				}
				if compressPath == f {
					f = f[strings.Index(f, "/")+1:]
					return f, true
				}
			}
		}
	}
	return "", false
}

func dir(path string) string {
	path = strings.ReplaceAll(path, "\\", string(filepath.Separator))
	lastSep := strings.LastIndex(path, string(filepath.Separator))
	return path[:lastSep]
}

func tarFileReader(filePath string) (*tar.Reader, func(), error) {
	reader, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: cannot open file, error=[%v]\n", err)
		return nil, nil, err
	}
	if filepath.Ext(filePath) == ".gz" {
		gr, err := gzip.NewReader(reader)
		if err != nil {
			fmt.Printf("error: cannot open gzip file, error=[%v]\n", err)
			return nil, nil, err
		}
		return tar.NewReader(gr), func() {
			gr.Close()
			reader.Close()
		}, nil
	} else {
		return tar.NewReader(reader), func() {
			reader.Close()
		}, nil
	}
}

func ExtractUnTar(filePath, targetPath string, files ...any) error {
	tarReader, close, err := tarFileReader(filePath)
	if err != nil {
		return err
	}
	defer close()
	bar := progressbar.NewBar(100)
	bar.SetNotice("\t")
	bar.HideRatio()
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error: cannot read tar file, error=[%v]\n", err)
			return err
		}
		// 去除压缩包内的一级目录
		compressPath := filepath.Clean(header.Name[strings.Index(header.Name, "/")+1:])
		includePath, isInclude := filePathInclude(compressPath, files...)
		if !isInclude {
			continue
		}
		fmt.Println("compressPath:", compressPath)
		info := header.FileInfo()
		targetFile := filepath.Join(targetPath, includePath)
		if info.IsDir() {
			if err = os.MkdirAll(targetFile, info.Mode()); err != nil {
				fmt.Printf("error: cannot mkdir file, error=[%v]\n", err)
				return err
			}
		} else {
			fDir := dir(targetFile)
			_, err = os.Stat(fDir)
			if os.IsNotExist(err) {
				if err = os.MkdirAll(fDir, info.Mode()); err != nil {
					fmt.Printf("error: cannot file mkdir file, error=[%v]\n", err)
					return err
				}
			}
			file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				fmt.Printf("error: cannot open file, error=[%v]\n", err)
				return err
			}
			bar.SetCurrentValue(0)
			writeFile(tarReader, file, header.Size, func(totalLength, processLength int64) {
				bar.PrintBar(int((float64(processLength) / float64(totalLength)) * 100))
			})
			file.Sync()
			file.Close()
			bar.PrintBar(100)
			bar.PrintEnd()
			if err != nil {
				fmt.Printf("error: cannot write file, error=[%v]\n", err)
				return err
			}
		}
	}
	return nil
}

func ExtractUnZip(filePath, targetPath string, rmRootDir bool, files ...any) error {
	if rc, err := zip.OpenReader(filePath); err == nil {
		defer rc.Close()
		bar := progressbar.NewBar(100)
		bar.SetNotice("\t")
		bar.HideRatio()
		var createWriteFile = func(info fs.FileInfo, path string, file io.Reader) error {
			targetFileName := filepath.Join(targetPath, path)
			if info.IsDir() {
				os.MkdirAll(targetFileName, info.Mode())
				return nil
			}
			fDir := filepath.Dir(targetFileName)
			if !tools.IsExist(fDir) {
				os.MkdirAll(fDir, 0755)
			}
			if targetFile, err := os.Create(targetFileName); err == nil {
				defer targetFile.Close()
				fmt.Println("extract file: ", path)
				bar.SetCurrentValue(0)
				writeFile(file, targetFile, info.Size(), func(totalLength, processLength int64) {
					bar.PrintBar(int((float64(processLength) / float64(totalLength)) * 100))
				})
				bar.PrintBar(100)
				bar.PrintEnd()
				return nil
			} else {
				println("createWriteFile", err.Error())
				return err
			}
		}
		// 所有文件
		if len(files) == 0 {
			zipFiles := rc.File
			for _, f := range zipFiles {
				r, _ := f.Open()
				var name string
				if rmRootDir {
					// 移除压缩包内的根文件夹名
					name = filepath.Clean(f.Name[strings.Index(f.Name, "/")+1:])
				} else {
					name = filepath.Clean(f.Name)
				}
				if err := createWriteFile(f.FileInfo(), name, r.(io.Reader)); err != nil {
					return err
				}
				_ = r.Close()
			}
			return nil
		} else {
			// 指定名字的文件
			for i := 0; i < len(files); i++ {
				if f, err := rc.Open(files[i].(string)); err == nil {
					info, _ := f.Stat()
					if err := createWriteFile(info, files[i].(string), f); err != nil {
						return err
					}
					_ = f.Close()
				} else {
					fmt.Printf("error: cannot open file, error=[%v]\n", err)
					return err
				}
			}
			return nil
		}
	} else {
		fmt.Printf("error: cannot read zip file, error=[%v]\n", err)
		return err
	}
}

// 释放bz2文件到tar
func UnBz2ToTar(name string, callback func(totalLength, processLength int64)) (string, error) {
	fileBz2, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer fileBz2.Close()
	dirName := fileBz2.Name()
	dirName = dirName[:strings.LastIndex(dirName, ".")]
	if !tools.IsExist(dirName) {
		r := bzip2.NewReader(fileBz2)
		w, err := os.Create(dirName)
		if err != nil {
			return "", err
		}
		defer w.Close()
		writeFile(r, w, 0, callback)
	} else {
		println("File already exists")
	}
	return dirName, nil
}

func writeFile(r io.Reader, w *os.File, totalLength int64, callback func(totalLength, processLength int64)) {
	buf := make([]byte, 1024*10)
	var written int64
	for {
		nr, err := r.Read(buf)
		if nr > 0 {
			nw, err := w.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			callback(totalLength, written)
			if err != nil {
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if err != nil {
			break
		}
	}
}

func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		return true
	}
	os.Remove(filename)
	return false
}

// 下载文件
func downloadFile(url string, localPath string, callback func(totalLength, processLength int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 1024*10)
		written int64
	)
	tmpFilePath := localPath + ".download"
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("download-error=[%v]\n", err)
		return err
	}
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Printf("download-error=[%v]\n", err)
		return err
	}
	if isFileExist(localPath, fsize) {
		println("File already exists")
		return nil
	}
	println("Save path: [", localPath, "] file size:", fsize)
	file, err := os.Create(tmpFilePath)
	if err != nil {
		fmt.Printf("download-error=[%v]\n", err)
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		fmt.Printf("Download-error=[body is null]\n")
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, err := file.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			callback(fsize, written)
			if err != nil {
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err == nil {
		file.Sync()
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		if err != nil {
			return err
		}
	}
	return err
}
