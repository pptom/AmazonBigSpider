package log

//打印的level等级

const (
	NOSET  = 0
	NO     = 110
	PANIC  = 120
	ERROR  = 130
	WARN   = 140
	NOTICE = 150
	LOG    = 160
	DEBUG  = 170
	ALL    = 180
)

var LogLevelMap map[string]int = map[string]int{
	"NO":     NO,
	"DEBUG":  DEBUG,
	"WARN":   WARN,
	"NOTICE": NOTICE,
	"LOG":    LOG,
	"ERROR":  ERROR,
	"PANIC":  PANIC,
	"ALL":    ALL,
}
var logLevelStringMap map[int]string = map[int]string{
	NO:     "NO",
	DEBUG:  "DEBUG",
	WARN:   "WARN",
	NOTICE: "NOTICE",
	LOG:    "LOG",
	ERROR:  "ERROR",
	PANIC:  "PANIC",
	ALL:    "ALL",
}
