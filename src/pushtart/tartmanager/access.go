package tartmanager

import (
	"pushtart/config"
)

//Exists returns true if a tart with the given pushURL exists.
func Exists(pushURL string) bool {
	if config.All().Tarts == nil {
		config.All().Tarts = map[string]config.Tart{}
		return false
	}
	if _, ok := config.All().Tarts[pushURL]; ok {
		return true
	}
	return false
}

//Get returns the tart with the given pushURL. If one does not exist, it returns an empty tart object.
func Get(pushURL string) config.Tart {
	if config.All().Tarts == nil {
		return config.Tart{}
	}
	return config.All().Tarts[pushURL]
}

//Save writes the given tart to global configuration, then to disk.
func Save(pushURL string, tart config.Tart) {
	if config.All().Tarts == nil {
		config.All().Tarts = map[string]config.Tart{}
	}
	config.All().Tarts[pushURL] = tart
	config.Flush()
}

//New creates a new (empty) tart with the given pushURL and owner.
func New(pushURL, owner string) {
	Save(pushURL, config.Tart{
		Name:             pushURL,
		PushURL:          pushURL,
		IsRunning:        false,
		Owners:           []string{owner},
		PID:              -1,
		RestartDelaySecs: 30,
		RestartOnStop:    false,
		LogStdout:        true,
	})
}

//List returns a []string of all the tarts in the system.
func List() []string {
	var output []string
	for pushURL := range config.All().Tarts {
		output = append(output, pushURL)
	}
	return output
}

//UserHasTartOwnership returns true if the given user is an owner of the given tart.Owners field.
func UserHasTartOwnership(user string, owners []string) bool {
	for _, owner := range owners {
		if user == owner {
			return true
		}
	}
	return false
}
