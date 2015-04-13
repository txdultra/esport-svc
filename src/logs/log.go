package logs

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/astaxie/beego/toolbox"
	log "github.com/cihub/seelog"
)

const (
	CONFIG_XML = `
		<seelog type="asynctimer" asyncinterval="10000000">
			<outputs>
				<filter levels="trace">
					<file path="@trace" formatid="trace_log_fmt"/>
				</filter>
				<filter levels="debug">
					<file path="@debug" formatid="debug_log_fmt"/>
				</filter>
				<filter levels="error">
					<file path="@error" formatid="error_log_fmt"/>
				</filter>
				<filter levels="critical">
					<file path="@critical" formatid="critical_log_fmt"/>
				</filter>
			</outputs>
			<formats>
				<format id="trace_log_fmt" format="%Date(2006 Jan 02/3:04:05.000000000 PM MST) [%Level] %Msg%n"/>
				<format id="debug_log_fmt" format="%Date(2006 Jan 02/3:04:05.000000000 PM MST) [%Level] %Msg%n"/>
				<format id="error_log_fmt" format="%Date(2006 Jan 02/3:04:05.000000000 PM MST) [%Level] %Msg%n"/>
				<format id="critical_log_fmt" format="%Date(2006 Jan 02/3:04:05.000000000 PM MST) [%Level] %Msg%n"/>
			</formats>
		</seelog>
		`
	LOG_DIR string = "log_files"
)

func init() {
	config_xml := get_config()
	logger, err := log.LoggerFromConfigAsString(config_xml)
	if err != nil {
		fmt.Println("seelog fail:", err)
	}
	log.UseLogger(logger)
	fmt.Println("log config completed...")
	bg_task()
}

func get_config() string {
	t := time.Now()
	date := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	trace_path := filepath.Join(LOG_DIR, "trace_"+date+".log")
	debug_path := filepath.Join(LOG_DIR, "debug_"+date+".log")
	error_path := filepath.Join(LOG_DIR, "error_"+date+".log")
	critical_path := filepath.Join(LOG_DIR, "critical_"+date+".log")
	xml := strings.Replace(CONFIG_XML, "@trace", trace_path, -1)
	xml = strings.Replace(xml, "@debug", debug_path, -1)
	xml = strings.Replace(xml, "@error", error_path, -1)
	xml = strings.Replace(xml, "@critical", critical_path, -1)
	return xml
}

func bg_task() {
	spec := "0 */1 * * * *"
	task_name := "seelog_task"
	toolbox.AddTask(task_name, toolbox.NewTask(task_name, spec, func() error {
		fmt.Println("running seelog task...")
		cxml := get_config()
		logger, err := log.LoggerFromConfigAsString(cxml)
		if err != nil {
			fmt.Println("seelog fail:", err)
		}
		err = log.ReplaceLogger(logger)
		if err != nil {
			fmt.Println("seelog fail:", err)
		}
		return nil
	}))
}

func Errorf(format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, filename := path.Split(file)
		str := fmt.Sprintf(format, v...)
		log.Errorf("file:%s line=%d %s", filename, line, str)
	}
}

func Error(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, filename := path.Split(file)
		str := fmt.Sprintf("file:%s line=%d ", filename, line)
		log.Errorf(str, v...)
	}
}
