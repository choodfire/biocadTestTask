package data

type LogRow struct {
	N         int
	Mqtt      string
	Invid     string
	Unit_guid string // globally unique identifier
	Msg_id    string
	Text      string
	Context   string
	Class     string
	Level     int
	Area      string
	Addr      string
	Block     string
	Typee     string // can't use word 'type' in Go
	Bit       string
}
