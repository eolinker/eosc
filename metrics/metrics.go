package metrics

import (
	"github.com/eolinker/eosc/common/pool"
	"strings"
)

const (
	startSign         = '#'
	metricsFormat     = "#{%s}"
	variableStartSign = '$'
	spaceEnd          = ' '
)

type LabelReader interface {
	ReadLabel(name string) string
}

type Metrics interface {
	Metrics(ctx LabelReader) string
	Key() string
}
type metricsReader interface {
	read(r LabelReader) string
}

var (
	_ Metrics = (*imlMetrics)(nil)
	
	bufferPool = pool.New(func() *strings.Builder {
		return &strings.Builder{}
	})
)

type imlMetrics struct {
	key     string
	readers []metricsReader
}

func (m *imlMetrics) Metrics(ctx LabelReader) string {
	if len(m.readers) == 0 {
		return ""
	}
	
	builder := bufferPool.Get()
	builder.Reset()
	defer bufferPool.PUT(builder)
	for _, reader := range m.readers {
		builder.WriteString(reader.read(ctx))
	}
	return builder.String()
}

func (m *imlMetrics) Key() string {
	return m.key
}

func Parse(str string) Metrics {
	// service-#{app} => service-app
	// ["service","{app}"] => service-#{app} => service-app
	m := &imlMetrics{key: strings.TrimSpace(str)}
	
	strReader := strings.NewReader(m.key)
	
	for {
		v, isContinue := readConst(strReader)
		if len(v) > 0 {
			m.readers = append(m.readers, metricsConst(v))
		}
		if !isContinue {
			break
		}
		v, isContinue = readReader(strReader)
		if len(v) > 0 {
			m.readers = append(m.readers, labelReader(v))
		}
		if !isContinue {
			break
		}
	}
	
	return m
}

func readConst(r *strings.Reader) (string, bool) {
	builder := bufferPool.Get()
	builder.Reset()
	defer bufferPool.PUT(builder)
	
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return builder.String(), false
		}
		if c != startSign {
			if c == variableStartSign {
				return builder.String(), true
			}
			builder.WriteRune(c)
			continue
		}
		
		n, _, errNext := r.ReadRune()
		if errNext != nil {
			builder.WriteRune(c)
			return builder.String(), false
		}
		if n == '{' {
			return builder.String(), true
		}
		
		builder.WriteRune(c)
		builder.WriteRune(n)
		
	}
}
func readReader(r *strings.Reader) (string, bool) {
	builder := bufferPool.Get()
	builder.Reset()
	defer bufferPool.PUT(builder)
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return builder.String(), false
		}
		if c == '}' || c == spaceEnd {
			return builder.String(), true
		}
		builder.WriteRune(c)
		
	}
}
