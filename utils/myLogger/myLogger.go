package myLogger

import (
	"genosha/utils/confs"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"runtime"
	"strconv"
	"time"
)

var Log *zap.Logger

func MyLogInit(path io.Writer) {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(path),
		zap.DebugLevel,
	)
	Log = zap.New(core, zap.AddCaller())
}

func SendErrorEmail(msg string) {
	_, file, line, _ := runtime.Caller(1)
	from := mail.NewEmail(confs.ConfigMap["senderName"], confs.ConfigMap["senderAddress"])
	to := mail.NewEmail(confs.ConfigMap["midEndName"], confs.ConfigMap["midEndEmail"])
	subject := "Genosha Error"
	plainTextContent := "-"
	htmlContent := `<p><b>Genosha Error</b></p>` + `<p>---------------------</p>` +
		`<p><b>Location: ` + file + `</b></p>` +
		`<p><b>Line: ` + strconv.Itoa(line) + `</b></p>` +
		`<p><b>Error: ` + msg + `</b></p>`
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(confs.ConfigMap["SENDGRID_API_KEY"])
	response, err := client.Send(message)
	if err != nil {
		Log.Error("error", zap.Any("err", err))
	} else {
		Log.Info("response.StatusCode", zap.Any("response.StatusCode", response.StatusCode))
		Log.Info("response.Body", zap.Any("response.Body", response.Body))
	}
}
