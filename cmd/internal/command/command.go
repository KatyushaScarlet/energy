//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package command

type OS string
type Arch string

type Config struct {
	Index     int
	Wd        string
	EnergyCfg EnergyConfig
	Install   Install `command:"install" description:"install energy development dependency environment"`
	Package   Package `command:"package" description:"energy application production and installation package"`
	Version   Version `command:"version" description:"list all release version numbers of energy"`
	Setenv    Setenv  `command:"setenv" description:"set ENERGY_ HOME framework environment"`
	Env       Env     `command:"env" description:"display ENERGY_ HOME framework environment directory"`
	Init      Init    `command:"init" description:"initialize the energy application project"`
	Build     Build   `command:"build" description:"building an energy project"`
	Bindata   Bindata `command:"bindata" description:"if the go version is less than 1.16, you can use bindata to embed static resources"`
	Gen       Gen     `command:"gen" description:"generate icons or syso commands"`
	Upg       Upgrade `command:"upg" description:"check and upgrade the current version"`
	Help      Help    `command:"help" description:"energy [cmd] help"`
	Cli       Cli     `command:"cli" description:"energy cli"`
}

type Command struct {
	Run                    func(c *Config) error
	UsageLine, Short, Long string
}

type Upgrade struct {
}

type Cli struct {
	Version bool `short:"v" long:"version" description:"energy cli version"`
	Update  bool `short:"u" long:"update" description:"energy cli update"`
}

type Gen struct {
	Icon bool `long:"icon" description:"Used to generate application icons,  can convert .png to .ico Generate pixel size of: [256, 128, 64, 48, 32, 16]"`
	Syso bool `long:"syso" description:"Generate the application program xxx.syso, and when compiling the execution file, the execution file information can be written into it"`
	// 参数
	// icon, syso
	IconFilePath string `short:"p" long:"iconFilePath" description:"Icon file directory:"`
	OutPath      string `short:"o" long:"outPath" description:"Save directory"`
	// syso
	Name             string `short:"n" long:"name" description:"Generate the syso file name and move it to the application name"`
	ManifestFilePath string `short:"m" long:"manifestFilePath" description:"Manifest file directory, if empty, will use the default template"`
	Arch             string `short:"a" long:"arch" description:"amd64 or i386 or arm64, if empty, the current system architecture"`
	InfoFilePath     string `short:"i" long:"infoFilePath" description:"Generate directory for syso information data files in JSON format."`
}

type Icon struct {
	IconFilePath string `short:"p" long:"iconFilePath" description:"Icon Directory: Supporting ICO, PNG formats, recommended as a 1024x1024 PNG"`
	OutPath      string `short:"o" long:"outPath" description:"ICO output save directory"`
}

type SYSO struct {
	Name             string `short:"n" long:"name" description:""`
	IconFilePath     string `short:"p" long:"iconFilePath" description:"Icon Directory: Supporting .ICO formats"`
	ManifestFilePath string `short:"m" long:"manifestFilePath" description:"Manifest file directory, if empty, will use the default template"`
	OutPath          string `short:"o" long:"outPath" description:".syso output save directory"`
	Arch             string `short:"a" long:"arch" description:"amd64 or i386 or arm64, if empty, the current system architecture"`
	SysoInfoFilePath string `short:"s" long:"sysoInfoFilePath" description:"Generate directory for syso information data files in JSON format."`
}

type Install struct {
	Path     string `short:"p" long:"path" description:"Installation directory Default current directory"`
	Version  string `short:"v" long:"version" description:"Specifying a version number"`
	Name     string `short:"n" long:"name" description:"Name of the framework directory after installation" default:"EnergyFramework"`
	Download string `short:"d" long:"download" description:"Download Source, 0:gitee or 1:github, Default empty" default:""`
	OS       OS     `long:"os" description:"Specify install OS: [windows, linux, darwin], default current system: os"`
	Arch     Arch   `long:"arch" description:"Specify install ARCH: [386, amd64, arm64], Default current system: architecture"`
	CEFVer   string `long:"cef" description:"Install system supports CEF version, provide 4 options, default empty. options: 109(support windows7), 106(support linux gtk2), 87(support flash)" default:""`
	IGolang  bool   // 是否安装Golang
	ICEF     bool   // 是否安装CEF
	INSIS    bool   // 是否安装nsis
	IUPX     bool   // 是否安装upx
	I7za     bool   // 是否安装7za
	IsSame   bool   // 安装的OS和Arch是否为当前系统架构, 默认当前系统架构
}

type Package struct {
	Path     string `short:"p" long:"path" description:"Project path, default current path. Can be configured in energy.json" default:""`
	Clean    bool   `short:"c" long:"clean" description:"Clear configuration and regenerate the default configuration"`
	Pkgbuild bool   `long:"pkg" description:"Using pkgbuild to create pkg development installation packages"`
}

type Env struct {
}

type Help struct {
}

type Setenv struct {
	Path string `short:"p" long:"path" description:"Energy framework dir"`
}

type Version struct {
	All bool `short:"a" long:"all" description:"show all"`
}

type Init struct {
	Name    string `short:"n" long:"name" description:"Initialized project name"`
	ResLoad string `short:"r" long:"resload" description:"Resource loading method, 1: HTTP, 2: Local Load, default 1 HTTP"`
	IGo     bool
	INSIS   bool
	IUPX    bool
	IEnv    bool
	INpm    bool
}

type Build struct {
	Path    string `short:"p" long:"path" description:"Project path, default current path. Can be configured in energy.json" default:""`
	Upx     bool   `short:"u" long:"upx" description:"Set this parameter and install upx. Use upx to compress the execution file."`
	UpxFlag string `long:"upxFlag" description:"Upx command line parameters" default:""`
	Args    string `long:"args" description:"go build [args]" default:""`
	Libemfs bool   `long:"libemfs" description:"Built in dynamic libraries to executable files, Copy liblcl to the built-in directory every compilation"`
}

type Bindata struct {
	Debug      bool   `long:"debug" description:"Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk."`
	Dev        bool   `long:"dev" description:"Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package."`
	Tags       string `long:"tags" description:"Optional set of build tags to include." default:""`
	Prefix     string `long:"prefix" description:"Optional path prefix to strip off asset names." default:""`
	Package    string `long:"pkg" description:"Package name to use in the generated code." default:"main"`
	NoMemCopy  bool   `long:"nomemcopy" description:"Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries."`
	NoCompress bool   `long:"nocompress" description:"Assets will *not* be GZIP compressed when this flag is specified."`
	NoMetadata bool   `long:"nometadata" description:"Assets will not preserve size, mode, and modtime info."`
	FSSystem   bool   `long:"fs" description:"Whether generate instance http.FileSystem interface code."`
	Mode       uint   `long:"mode" description:"Optional file mode override for all files."`
	ModTime    int64  `long:"modtime" description:"Optional modification unix timestamp override for all files."`
	Output     string `long:"o" description:"Optional name of the output file to be generated." default:"./bindata.go"`
	Ignore     string `long:"Ignore" description:"Regex pattern to ignore." default:""`
	Paths      string `long:"paths" description:"Static resource directory, Multiple Catalogs: ./resource,./libs" default:""`
}

type EnergyConfig struct {
	Registry string `json:"registry"`
}

func (m OS) IsWindows() bool {
	return m == "windows"
}

func (m OS) IsLinux() bool {
	return m == "linux"
}

func (m OS) IsDarwin() bool {
	return m == "darwin"
}

func (m Arch) Is386() bool {
	return m == "386" || m == "32" // windows32
}

func (m Arch) IsAMD64() bool {
	return m == "amd64"
}

func (m Arch) IsARM64() bool {
	return m == "arm64"
}
