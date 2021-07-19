package dlog

type DriverInfo struct {
	Name   string  `json:"name"`
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
}

type DriverManager struct {
	drivers map[string]ConfigDriver
	infos   []DriverInfo
}

func NewDriverManager(drivers []ConfigDriver, ignoreField ...string) *DriverManager {

	driverMap := make(map[string]ConfigDriver)
	infos := make([]DriverInfo, 0, len(drivers))
	for _, d := range drivers {
		driverMap[d.Name()] = d

		infos = append(infos, DriverInfo{
			Name:   d.Name(),
			Title:  d.Title(),
			Fields: d.ConfigFields(ignoreField...),
		})
	}

	return &DriverManager{
		drivers: driverMap,
		infos:   infos,
	}
}

func (m *DriverManager) Get(driver string) (ConfigDriver, bool) {
	d, has := m.drivers[driver]
	return d, has
}
func (m *DriverManager) Infos() []DriverInfo {
	return m.infos
}
func (m *DriverManager) Title(name string) string {
	d, has := m.drivers[name]
	if has {
		return d.Title()
	}
	return name
}
