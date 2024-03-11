# slog-exp
slog experimental features such individual log files for different levels, and other extensions.


![GitHub](https://img.shields.io/github/license/smallnest/slog-exp) ![GitHub Action](https://github.com/smallnest/slog-exp/actions/workflows/action.yaml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/slog-exp)](https://goreportcard.com/report/github.com/smallnest/slog-exp)  [![GoDoc](https://godoc.org/github.com/smallnest/slog-exp?status.png)](http://godoc.org/github.com/smallnest/slog-exp)  

**See also:**

- [slog-multi](https://github.com/samber/slog-multi): `slog.Handler` chaining, fanout, routing, failover, load balancing...
- [slog-formatter](https://github.com/samber/slog-formatter): `slog` attribute formatting
- [slog-sampling](https://github.com/samber/slog-sampling): `slog` sampling policy
- [slog-gin](https://github.com/samber/slog-gin): Gin middleware for `slog` logger
- [slog-echo](https://github.com/samber/slog-echo): Echo middleware for `slog` logger
- [slog-fiber](https://github.com/samber/slog-fiber): Fiber middleware for `slog` logger
- [slog-chi](https://github.com/samber/slog-chi): Chi middleware for `slog` logger
- [slog-datadog](https://github.com/samber/slog-datadog): A `slog` handler for `Datadog`
- [slog-rollbar](https://github.com/samber/slog-rollbar): A `slog` handler for `Rollbar`
- [slog-sentry](https://github.com/samber/slog-sentry): A `slog` handler for `Sentry`
- [slog-syslog](https://github.com/samber/slog-syslog): A `slog` handler for `Syslog`
- [slog-logstash](https://github.com/samber/slog-logstash): A `slog` handler for `Logstash`
- [slog-fluentd](https://github.com/samber/slog-fluentd): A `slog` handler for `Fluentd`
- [slog-graylog](https://github.com/samber/slog-graylog): A `slog` handler for `Graylog`
- [slog-loki](https://github.com/samber/slog-loki): A `slog` handler for `Loki`
- [slog-slack](https://github.com/samber/slog-slack): A `slog` handler for `Slack`
- [slog-telegram](https://github.com/samber/slog-telegram): A `slog` handler for `Telegram`
- [slog-mattermost](https://github.com/samber/slog-mattermost): A `slog` handler for `Mattermost`
- [slog-microsoft-teams](https://github.com/samber/slog-microsoft-teams): A `slog` handler for `Microsoft Teams`
- [slog-webhook](https://github.com/samber/slog-webhook): A `slog` handler for `Webhook`
- [slog-kafka](https://github.com/samber/slog-kafka): A `slog` handler for `Kafka`
- [slog-nats](https://github.com/samber/slog-nats): A `slog` handler for `NATS`
- [slog-parquet](https://github.com/samber/slog-parquet): A `slog` handler for `Parquet` + `Object Storage`
- [slog-zap](https://github.com/samber/slog-zap): A `slog` handler for `Zap`
- [slog-zerolog](https://github.com/samber/slog-zerolog): A `slog` handler for `Zerolog`
- [slog-logrus](https://github.com/samber/slog-logrus): A `slog` handler for `Logrus`
- [slog-channel](https://github.com/samber/slog-channel): A `slog` handler for Go channels
- [slog-clickhouse](https://github.com/smallnest/slog-clickhouse): A `slog` handler for ClickHouse

## Installation

```bash
go get github.com/smllnest/slog-clickhouse
```

**Compatibility**: `go >= 1.21`

## Usage

### LevelHandler

Set different handlers for different log levels.
For example, you can set different handlers to save files for different log levels, such as `info`, `warn`, and `error`.

```go
    infoFile, err := os.OpenFile("testdata/info.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
    if err != nil {
        panic(err)
    }
    defer infoFile.Close()

    warnFile, err := os.OpenFile("testdata/warn.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
    if err != nil {
        panic(err)
    }
    defer warnFile.Close()

    errorFile, err := os.OpenFile("testdata/error.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
    if err != nil {
        panic(err)
    }
    defer errorFile.Close()

    infoHandler := slog.NewTextHandler(infoFile, &slog.HandlerOptions{Level: slog.LevelInfo})
    warnHandler := slog.NewTextHandler(warnFile, &slog.HandlerOptions{Level: slog.LevelWarn})
    errorHandler := slog.NewTextHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError})

    handler := NewLevelHandler(map[slog.Level]slog.Handler{
        slog.LevelInfo:  infoHandler,
        slog.LevelWarn:  warnHandler,
        slog.LevelError: errorHandler,
    })

    logger := slog.New(handler)

    logger.Info("info text")
    logger.Warn("warn text")
    logger.Error("error text")
```

### customized datetime and source

use `ReplaceTimeAttr` and `ReplaceSourceAttr`.

```go
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: WrapReplaceAttrFunc(ReplaceTimeAttr(time.TimeOnly), ReplaceSourceAttr()),
		AddSource:   true,
	})
	logger := slog.New(handler)

	logger.Info("info text")
	logger.Warn("warn text")
	logger.Error("error text")
```