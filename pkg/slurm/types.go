package slurm

// InterLinkConfig holds the whole configuration
type SlurmConfig struct {
	VKConfigPath      string `yaml:"VKConfigPath"`
	Sbatchpath        string `yaml:"SbatchPath"`
	Scancelpath       string `yaml:"ScancelPath"`
	Squeuepath        string `yaml:"SqueuePath"`
	Sidecarport       string `yaml:"SidecarPort"`
	ExportPodData     bool   `yaml:"ExportPodData"`
	Commandprefix     string `yaml:"CommandPrefix"`
	DataRootFolder    string `yaml:"DataRootFolder"`
	Namespace         string `yaml:"Namespace"`
	Tsocks            bool   `yaml:"Tsocks"`
	Tsockspath        string `yaml:"TsocksPath"`
	Tsockslogin       string `yaml:"TsocksLoginNode"`
	BashPath          string `yaml:"BashPath"`
	VerboseLogging    bool   `yaml:"VerboseLogging"`
	ErrorsOnlyLogging bool   `yaml:"ErrorsOnlyLogging"`
	SingularityPrefix string `yaml:"SingularityPrefix"`
	set               bool
}

type CreateStruct struct {
	PodUID string `json:"PodUID"`
	PodJID string `json:"PodJID"`
}
