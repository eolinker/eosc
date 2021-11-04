package eosc

type IExtendersSetting interface {
	SetExtender(group, project, version string) error
	DelExtender(group, project string) (string, bool)
	GetExtenderVersion(group, project string) (string, bool)
	//ALl()map[string]string
}
