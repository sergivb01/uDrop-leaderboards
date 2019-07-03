package services

type config struct {
	MongoURI     string `yaml:"mongoURI"`
	Leaderboards []struct {
		Name      string `yaml:"name"`
		Shortname string `yaml:"shortName"`
	}
}
