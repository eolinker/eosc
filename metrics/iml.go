package metrics

var (
	_ metricsReader = metricsConst("")
	_ metricsReader = labelReader("")
)

type metricsConst string

func (m metricsConst) read(r LabelReader) string {
	return string(m)
}

type labelReader string

func (m labelReader) read(r LabelReader) string {
	return r.GetLabel(string(m))
}
