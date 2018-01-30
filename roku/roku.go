package roku

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/resty.v1"
)

var port = 8060

var keys = map[string]string{
	// Standard Keys
	"home":      "Home",
	"reverse":   "Rev",
	"forward":   "Fwd",
	"play":      "Play",
	"select":    "Select",
	"left":      "Left",
	"right":     "Right",
	"down":      "Down",
	"up":        "Up",
	"back":      "Back",
	"replay":    "InstantReplay",
	"info":      "Info",
	"backspace": "Backspace",
	"search":    "Search",
	"enter":     "Enter",
	"literal":   "Lit",

	// For devices that support "Find Remote"
	"find_remote": "FindRemote",

	// For Roku TV
	"volume_down": "VolumeDown",
	"volume_up":   "VolumeUp",
	"volume_mute": "VolumeMute",

	// For Roku TV while on TV tuner channel
	"channel_up":   "ChannelUp",
	"channel_down": "ChannelDown",

	// For Roku TV current input
	"input_tuner": "InputTuner",
	"input_hdmi1": "InputHDMI1",
	"input_hdmi2": "InputHDMI2",
	"input_hdmi3": "InputHDMI3",
	"input_hdmi4": "InputHDMI4",
	"input_av1":   "InputAV1",

	// For devices that support being turned on/off
	"power": "Power",
}

// App application
type App struct {
	Name    string `xml:",chardata"`
	ID      string `xml:"id,attr"`
	Type    string `xml:"type,attr"`
	SubType string `xml:"subtype,attr"`
	Version string `xml:"version,attr"`
}

// Apps applications
type Apps struct {
	Apps []App `xml:"app"`
}

func makeURL(ip string) string {
	return fmt.Sprintf("http://%v:%v", ip, port)
}

func getStrKeys(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

// GetApps gets a url
func QueryActiveApp(ip string) (*App, error) {
	a := new(Apps)
	_, err := resty.R().SetResult(a).Get(makeURL(ip) + "/query/active-app")
	if len(a.Apps) != 1 {
		return &App{}, err
	}
	return &a.Apps[0], err
}

func QueryApps(ip string) (*Apps, error) {
	a := new(Apps)
	_, err := resty.R().SetResult(a).Get(makeURL(ip) + "/query/apps")
	return a, err
}

// GetCommands commands
func GetCommands() []string {
	return getStrKeys(keys)
}

// KeyPress command
func KeyPress(ip string, k string) error {
	if kReal, ok := keys[k]; ok {
		_, err := resty.R().Post(makeURL(ip) + "/keypress/" + kReal)
		return err
	}

	return errors.New("invalid command")
}

// KeyPress command
func KeyDown(ip string, k string) error {
	if kReal, ok := keys[k]; ok {
		_, err := resty.R().Post(makeURL(ip) + "/keydown/" + kReal)
		return err
	}

	return errors.New("invalid command")
}

// KeyPress command
func KeyUp(ip string, k string) error {
	if kReal, ok := keys[k]; ok {
		_, err := resty.R().Post(makeURL(ip) + "/keyup/" + kReal)
		return err
	}

	return errors.New("invalid command")
}

func LaunchApp(ip string, id string) error {
	_, err := resty.R().Post(makeURL(ip) + "/launch/" + id)
	return err
}

func LaunchAppName(ip string, n string) error {
	as, err := QueryApps(ip)
	if err != nil {
		return err
	}
	for _, a := range as.Apps {
		if n == a.Name {
			return LaunchApp(ip, a.ID)
		}
	}
	return errors.New("app not found")
}

func LaunchAppNameMatch(ip string, m string) error {
	as, err := QueryApps(ip)
	if err != nil {
		return err
	}
	for _, a := range as.Apps {
		if strings.Contains(strings.ToLower(a.Name), strings.ToLower(m)) {
			return LaunchApp(ip, a.ID)
		}
	}
	return errors.New("app not found")
}
