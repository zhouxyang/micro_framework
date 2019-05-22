package caller

// refer to https://github.com/xdxiaodong/logrus-hook-caller
import (
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type CallerHook struct {
	CallerHookOptions *CallerHookOptions
}

// NewHook creates a new caller hook with options. If options are nil or unspecified, options.Field defaults to "src"
// and options.Flags defaults to log.Llongfile
func NewHook(options *CallerHookOptions) *CallerHook {
	// new
	if options.FileAlias == "" {
		options.FileAlias = "file"
	}
	if options.LineAlias == "" {
		options.LineAlias = "line"
	}
	// Set default caller flag to Std logger log.Llongfile
	if options.Flags == 0 {
		options.Flags = log.Lshortfile
	}
	return &CallerHook{options}
}

// CallerHookOptions stores caller hook options
type CallerHookOptions struct {
	FileAlias  string //default:file
	EnableFile bool
	LineAlias  string //default:line
	EnableLine bool
	// Stores the flags
	Flags int
}

// HasFlag returns true if the report caller options contains the specified flag
func (options *CallerHookOptions) HasFlag(flag int) bool {
	return options.Flags&flag != 0
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {

	file, line := getCallerIgnoringLogMulti(1)
	if hook.CallerHookOptions.HasFlag(log.Lshortfile) && !hook.CallerHookOptions.HasFlag(log.Llongfile) {
		file = path.Base(file)
	}
	if hook.CallerHookOptions.EnableFile {
		entry.Data[hook.CallerHookOptions.FileAlias] = file
	}
	if hook.CallerHookOptions.EnableLine {
		entry.Data[hook.CallerHookOptions.LineAlias] = line
	}
	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func getCaller(callDepth int, suffixesToIgnore ...string) (file string, line int) {
	// bump by 1 to ignore the getCaller (this) stackframe
	callDepth++
outer:
	for {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
			break
		}

		for _, s := range suffixesToIgnore {
			if strings.HasSuffix(file, s) {
				callDepth++
				continue outer
			}
		}
		break
	}
	return
}

//new
func getCallerIgnoringLogMulti(callDepth int) (string, int) {
	// the +1 is to ignore this (getCallerIgnoringLogMulti) frame
	return getCaller(callDepth+1, "logrus/hooks.go", "logrus/entry.go", "logrus/logger.go", "logrus/exported.go", "asm_amd64.s")
}
