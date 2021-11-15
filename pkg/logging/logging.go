package logging

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type CreateIndexFunc func(entry *logrus.Entry, l *LoggerModel) error

type LoggerService interface {
	CreateLog() *logrus.Logger
}

type loggerService struct {
	logger *logrus.Logger
}

func (l loggerService) CreateLog() *logrus.Logger{
	log := logrus.New()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200/"),
		elastic.SetSniff(false))
	if err != nil {
		log.Panic(err)
	}

	createLogOperation, err := PrepareLogging(client, "localhost", logrus.DebugLevel, "product_api_log", CreateIndexAsync)
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(createLogOperation)

	return log
}

func NewLoggerService (logger *logrus.Logger) LoggerService{
	return &loggerService{logger: logger}
}

func CreateIndexAsync(entry *logrus.Entry, l *LoggerModel) error {
	go CreateIndex(entry, l)
	return nil
}

func PrepareLogging(client *elastic.Client, host string, level logrus.Level, elasticIndex string, createIndexAsync CreateIndexFunc) (*LoggerModel, error){
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())

	exists, err := client.IndexExists(elasticIndex).Do(ctx)
	if err != nil {
		cancel()
		return nil, err
	}

	if !exists {
		createIndex, err := client.CreateIndex(elasticIndex).Do(ctx)
		if err != nil {
			cancel()
			return nil, err
		}
		if !createIndex.Acknowledged {
			cancel()
			return nil, fmt.Errorf("Index cannot be create!")
		}
	}


	return &LoggerModel{
		client:    client,
		host:      host,
		levels:    levels,
		ctx: ctx,
		createIndexAsync: createIndexAsync,
	}, nil
}

func (l LoggerModel) Levels() []logrus.Level {
	return l.levels
}

func (l *LoggerModel) Fire(entry *logrus.Entry) error {
	return l.createIndexAsync(entry, l)
}

func CreateIndex(entry *logrus.Entry, l *LoggerModel) error{
	_, err := l.client.Index().Index("product_api_log").BodyJson(*PrepareMessage(entry, l)).Do(l.ctx)
	return err
}

func PrepareMessage(entry *logrus.Entry, l *LoggerModel) *LoggerMessage{
	lwl := entry.Level.String()
	if e, ok := entry.Data[logrus.ErrorKey]; ok && e != nil {
		if err, ok := e.(error); ok {
			entry.Data[logrus.ErrorKey] = err.Error()
		}
	}

	var file string
	var function string
	if entry.HasCaller() {
		file = entry.Caller.File
		function = entry.Caller.Function
	}

	return &LoggerMessage{
		l.host,
		entry.Time.UTC().Format(time.RFC3339Nano),
		file,
		function,
		entry.Message,
		entry.Data,
		strings.ToUpper(lwl),
	}
}

type LoggerModel struct {
	client *elastic.Client
	host string
	levels    []logrus.Level
	ctx context.Context
	createIndexAsync CreateIndexFunc
}

type LoggerMessage struct {
	Host      string `json:"Host,omitempty"`
	Timestamp string `json:"@timestamp"`
	File      string `json:"File,omitempty"`
	Func      string `json:"Func,omitempty"`
	Message   string `json:"Message,omitempty"`
	Data      logrus.Fields
	Level     string `json:"Level,omitempty"`
}
