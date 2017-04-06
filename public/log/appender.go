package log

import (
	"fmt"
	golog "log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	DefaultLowCallpath      = 0
	DefaultAppenderCallpath = 2
	DefaultLoggerCallpath   = 3
)
var UseShortFile bool

//日志的输出接口
type Appender interface {
	Log(extendCallpath int, level string, args ...interface{})
	Logln(extendCallpath int, level string, args ...interface{})
	Logf(extendCallpath int, level string, format string, args ...interface{})
	SetCallpath(int)
}

type baseAppender struct {
	*golog.Logger
	Name     string
	Callpath int
}

func (l *baseAppender) SetCallpath(level int) {
	l.Callpath = level
}

func (l *baseAppender) log(extendCallpath int, level string, fmtFunc func(...interface{}) string, args ...interface{}) {
	v := make([]interface{}, 1, len(args)+1)
	v[0] = "[" + level + "] "
	v = append(v, args...)
	if l.Callpath == 0 {
		l.Callpath = DefaultAppenderCallpath
	}
	// fmt.Println("=============================")
	// fmt.Println(DefaultLowCallpath, l.Callpath, extendCallpath)
	// for i := 0; i < l.Callpath+extendCallpath+2; i++ {
	// 	_, name, line, _ := runtime.Caller(i)
	// 	fmt.Println(name, line)
	// }
	// fmt.Println("=============================")

	l.Output(DefaultLowCallpath+l.Callpath+extendCallpath, fmtFunc(v...))
}

func (l *baseAppender) logf(extendCallpath int, level string, fmtFunc func(string, ...interface{}) string, format string, args ...interface{}) {
	format = "[" + level + "] " + format
	if l.Callpath == 0 {
		l.Callpath = DefaultAppenderCallpath
	}
	l.Output(DefaultLowCallpath+l.Callpath+extendCallpath, fmtFunc(format, args...))
}

func (l *baseAppender) Log(extendCallpath int, level string, args ...interface{}) {
	l.log(extendCallpath, level, fmt.Sprint, args...)
}

func (l *baseAppender) Logln(extendCallpath int, level string, args ...interface{}) {
	l.log(extendCallpath, level, fmt.Sprintln, args...)
}

func (l *baseAppender) Logf(extendCallpath int, level string, format string, args ...interface{}) {
	l.logf(extendCallpath, level, fmt.Sprintf, format, args...)
}

type FileAppender struct {
	*baseAppender
	fileName string
}

func NewFileAppender(name string, fileName string) *FileAppender {
	l := &FileAppender{
		baseAppender: &baseAppender{
			Name:     name,
			Callpath: 2,
		},
		fileName: fileName,
	}
	return l
}

func (l *FileAppender) Log(extendCallpath int, level string, args ...interface{}) {
	l.lazyNewLogger()
	l.log(extendCallpath, level, fmt.Sprint, args...)
}

func (l *FileAppender) Logln(extendCallpath int, level string, args ...interface{}) {
	l.lazyNewLogger()
	l.log(extendCallpath, level, fmt.Sprintln, args...)
}

func (l *FileAppender) Logf(extendCallpath int, level string, format string, args ...interface{}) {
	l.lazyNewLogger()
	l.logf(extendCallpath, level, fmt.Sprintf, format, args...)
}

func (l *FileAppender) lazyNewLogger() {
	if l.baseAppender.Logger == nil {
		os.MkdirAll(path.Dir(l.fileName), 0640)
		logFile, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			golog.Fatalln("log conf error:", err.Error())
			return
		}
		var tag int
		if UseShortFile {
			tag = golog.LstdFlags | golog.Lshortfile
		} else {
			tag = golog.LstdFlags | golog.Llongfile
		}

		defaultLogger := golog.New(logFile, "", tag)
		l.baseAppender.Logger = defaultLogger
	}
}

type DailyAppender struct {
	*FileAppender
	today     string //当天日期
	fileName  string //文件
	extension string //后缀名
	lock      sync.Mutex
}

func NewDailyAppenderEx(name, fileName, extension string) *DailyAppender {
	var fname string = fileName
	if strings.HasSuffix(fileName, extension) {
		fname = fileName[:len(fileName)-4]
	}
	var appender = &DailyAppender{
		fileName:  fname,
		extension: extension,
	}
	appender.setLogger(name, time.Now().Format("20060102"))

	//更新文件
	go func() {
		for {
			now := time.Now()
			h, m, s := now.Clock()
			leave := 86400 - (h*60+m)*60 + s
			select {
			case <-time.After(time.Duration(leave) * time.Second):
				appender.setLogger(name, time.Now().Format("20060102"))
			}
		}
	}()
	return appender
}

func NewDailyAppender(name, fileName string) *DailyAppender {
	return NewDailyAppenderEx(name, fileName, ".log")
}

func (self *DailyAppender) setLogger(name string, day string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.today == day {
		return nil
	}
	self.today = day

	self.FileAppender = NewFileAppender(name, self.fileName+"_"+day+self.extension)
	return nil
}

type ConsoleAppender struct {
	*baseAppender
}

func NewConsoleAppender(name string) *ConsoleAppender {
	var tag int
	if UseShortFile {
		tag = golog.LstdFlags | golog.Lshortfile
	} else {
		tag = golog.LstdFlags | golog.Llongfile
	}

	a := &ConsoleAppender{
		baseAppender: &baseAppender{
			Logger:   golog.New(os.Stdout, "", tag),
			Name:     name,
			Callpath: 2,
		},
	}
	return a
}

//把不同等级的信息打印到不同的Appender中
type LevelSeparationAppender struct {
	Name      string
	appenders map[string]Appender
}

func NewLevelSeparationAppender(name string) *LevelSeparationAppender {
	return &LevelSeparationAppender{
		Name:      name,
		appenders: make(map[string]Appender),
	}
}

func (this *LevelSeparationAppender) SetLevelAppender(level string, appender Appender) {
	this.appenders[level] = appender
}
func (this *LevelSeparationAppender) Log(extendCallpath int, level string, args ...interface{}) {
	if l, ok := this.appenders[level]; ok {
		l.Log(extendCallpath+1, level, args...)
	}
}

func (this *LevelSeparationAppender) Logln(extendCallpath int, level string, args ...interface{}) {
	if l, ok := this.appenders[level]; ok {
		l.Logln(extendCallpath+1, level, args...)
	}
}

func (this *LevelSeparationAppender) Logf(extendCallpath int, level string, format string, args ...interface{}) {
	if l, ok := this.appenders[level]; ok {
		l.Logf(extendCallpath+1, level, format, args...)
	}
}

func (this *LevelSeparationAppender) SetCallpath(level int) {
	for _, ap := range this.appenders {
		ap.SetCallpath(level)
	}
}

func NewLevelSeparationDailyAppender(name string, fileName string) *LevelSeparationAppender {
	l := NewLevelSeparationAppender(name)

	fname := fileName
	if strings.HasSuffix(fileName, ".log") {
		fname = fileName[:len(fileName)-4]
	}
	for _, level := range logLevelStringMap {
		if level == "ALL" || level == "NO" {
			continue
		}
		levelAppender := NewDailyAppenderEx(name+"_"+level, fname, "."+strings.ToLower(level))
		l.SetLevelAppender(level, levelAppender)
	}
	return l
}
