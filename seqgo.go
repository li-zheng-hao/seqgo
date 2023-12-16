package seqgo

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// SeqHook sends logs to Seq via HTTP.
type SeqHook struct {
	endpoint string
	apiKey   string
	levels   []logrus.Level
}

func (hook *SeqHook) Levels() []logrus.Level {
	return hook.levels
}

var SeqHookOption *SeqHookOptions = &SeqHookOptions{
	levels: []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	},
	period:    2,
	batchSize: 10,
}

func NewSeqHook(configure func(*SeqHookOptions)) *SeqHook {

	configure(SeqHookOption)

	endpoint := fmt.Sprintf("%v/api/events/raw", SeqHookOption.endpoint)

	SeqHookOption.endpoint = endpoint

	go ScheduleSend()

	return &SeqHook{
		endpoint: endpoint,
		apiKey:   SeqHookOption.apiKey,
		levels:   SeqHookOption.levels,
	}
}

// Fire sends a log entry to Seq.
func (hook *SeqHook) Fire(entry *logrus.Entry) error {
	formatter := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "@mt",
			logrus.FieldKeyLevel: "@l",
			logrus.FieldKeyTime:  "@t",
		},
	}
	for k, v := range SeqHookOption.fields {
		entry.Data[k] = v
	}
	data, err := formatter.Format(entry)

	if err != nil {
		return err
	}

	Push(data)

	return nil
}

// SeqHookOptions collects non-default Seq hook options.
type SeqHookOptions struct {
	apiKey    string
	levels    []logrus.Level
	period    int
	fields    map[string]string
	batchSize int64
	endpoint  string
}