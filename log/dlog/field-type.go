package dlog

type FieldType string

const (
	String FieldType = "string"
	Number FieldType = "number"
	Array  FieldType = "array"
)

type InputType string

const (
	LineInput       InputType = "line"            // 单行文本
	MultiInput      InputType = "multi"           // 多行文本
	SelectInput     InputType = "select"          // 下拉框
	CheckboxInput   InputType = "checkbox"        // 多选
	RadioInput      InputType = "radio"           // 单选
	NumberInput     InputType = "number"          // 数值输入框
	Gradient_select InputType = "gradient_select" // 梯度选择
	TableInput      InputType = "table"           // 梯度选择
)
