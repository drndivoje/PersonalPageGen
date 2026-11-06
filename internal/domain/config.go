package domain

type ConfigFile struct {
	Main string `yaml:"site_name"`
	Menu []struct {
		Title string `yaml:"title"`
		Path  string `yaml:"path"`
	} `yaml:"menu"`
	Domain string `yaml:"domain"`
	Author string `yaml:"author"`
	Footer string `yaml:"footer"`
}

func (c ConfigFile) getDomainUrl() string {
	return "https://" + c.Domain
}

func (c ConfigFile) GetFooter() string {
	defaultFooter := "Built with <a href=\"https://github.com/drndivoje/PersonalPageGen\" target=\"_blank\" >Personal Page Gen</a>"
	if c.Footer == "" {
		return defaultFooter
	}
	return c.Footer + " - " + defaultFooter
}