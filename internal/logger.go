package internal

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

// The Logger is roughly based on https://github.com/kpurdon/zlg with a number
// of features and config options tweaked or removed.

// We use this sync.Once to make sure the global zerolog configuration only
// happens once.
var once sync.Once

type LoggerConfig struct {
	Level LogLevel `envconfig:"level"`
}

type LogLevel zerolog.Level

// Ensure LogLevel implements flag.Value, because that's what we use it for.
var _ flag.Value = (*LogLevel)(nil)

func (ll *LogLevel) Set(val string) error {
	level, err := zerolog.ParseLevel(strings.ToLower(val))
	if err == nil {
		*ll = LogLevel(level)
	}
	return err
}

func (ll *LogLevel) String() string {
	return zerolog.Level(*ll).String()
}

// Logger provides limited structured logging capabilities backed by a well
// configured zerolog.Logger.
type Logger struct {
	logger zerolog.Logger
}

func NewLogger() (*Logger, error) {
	once.Do(func() {
		// The default in zerolog is 2, but we want to ignore this package as
		// well, so we set it to 3.
		zerolog.CallerSkipFrameCount = 3

		// The default TimestampFunc is time.Now without the UTC, so we update
		// it here.
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}
	})

	var config LoggerConfig
	var logger zerolog.Logger

	// Yes, there are instances where you could have a terminal and be in prod
	// (or not have a terminal and be in dev mode), but because of how this is
	// set up, that shouldn't be an issue.
	if Env() == EnvDev {
		logger = zerolog.New(zerolog.NewConsoleWriter())
	} else {
		logger = zerolog.New(os.Stdout)
	}

	err := envconfig.Process("LOG", &config)

	logger = logger.
		Level(zerolog.Level(config.Level)).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{
		logger: logger,
	}, err
}

func (l *Logger) WithLevel(lvl zerolog.Level) *Logger {
	l.logger = l.logger.Level(lvl)
	return l
}

// With embeds the given key/value in the Logger returning the new Logger.
func (l *Logger) With(k string, v interface{}) *Logger {
	if v == nil {
		return l
	}

	return &Logger{
		logger: loggerWithValue(l.logger, k, v),
	}
}

func (l *Logger) UpdateWith(k string, v interface{}) {
	l.logger = loggerWithValue(l.logger, k, v)
}

func loggerWithValue(logger zerolog.Logger, k string, v interface{}) zerolog.Logger {
	switch vv := v.(type) {
	case bool:
		return logger.With().Bool(k, vv).Logger()
	case string:
		return logger.With().Str(k, vv).Logger()
	case int:
		return logger.With().Int(k, vv).Logger()
	case uint:
		return logger.With().Uint(k, vv).Logger()
	case int64:
		return logger.With().Int64(k, vv).Logger()
	case uint64:
		return logger.With().Uint64(k, vv).Logger()
	case float64:
		return logger.With().Float64(k, vv).Logger()
	case time.Duration:
		return logger.With().Dur(k, vv).Logger()
	case time.Time:
		return logger.With().Time(k, vv).Logger()
	default:
		return logger.With().Interface(k, vv).Logger()
	}
}

// WithError embeds a given error in the Logger.
// Prefer to use the *Logger.Error method when logging at the error level.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With().Stack().Err(err).Logger(),
	}
}

// Debug writes the given message at the debug level.
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf writes the given formatted message at the debug level. Note that even
// though this function is being left in for debugging messages, formatted
// message functions will not be added for other log levels.
func (l *Logger) Debugf(format string, params ...interface{}) {
	l.logger.Debug().Msgf(format, params...)
}

// Info writes the given message at the info level.
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Warn writes the given message at the warning level.
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Error writes the given message at the error level.
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Panic writes the given message at error level and then calls panic.
// Generally this will be used along with an error Logger.WithError(err).Panic("something").
func (l *Logger) Panic(msg string) {
	l.logger.Panic().Msg(msg)
}

func LoggerMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return contextValueMiddleware(LoggerContextKey, logger)
}
