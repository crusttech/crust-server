package incoming

type (
	ExecCommand struct {
		ChannelID uint64            `json:"channelId,string"`
		Command   string            `json:"command"`
		Params    map[string]string `json:"params"`
		Input     string            `json:"input"`
	}
)
