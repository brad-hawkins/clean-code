package types

var (
	SettingWritebackDisabled = "writeback_disabled"
)

type Settings interface {
	GetRaw(name string) string
}
