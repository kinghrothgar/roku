package roku

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cloudflare/cfssl/log"
	ssdp "github.com/koron/go-ssdp"
	"gopkg.in/resty.v1"
)

var (
	port    = 8060
	timeout = time.Duration(5 * time.Second)
)

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
	//"literal":   "Lit",

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

type Roku struct {
	Apps       []App
	IP         string
	RestClient *resty.Client
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

func New(ip string) *Roku {
	restC := resty.New()
	restC.SetTimeout(timeout)
	return &Roku{
		IP:         ip,
		RestClient: restC,
	}
}

// GetCommands commands
func GetCommands() []string {
	return getStrKeys(keys)
}

// GetCommands commands
func FindRoku() ([]*url.URL, error) {
	list, err := ssdp.Search("roku:ecp", 5, "")
	if err != nil {
		log.Info("%v", err)
		return []*url.URL{}, err
	}

	rokus := make([]*url.URL, len(list))
	log.Info(len(list))
	for _, srv := range list {
		u, err := url.Parse(srv.Location)
		if err != nil {
			log.Info("%v", err)
			return []*url.URL{}, err
		}
		rokus = append(rokus, u)
	}
	log.Info("%v", rokus)
	return rokus, nil
}

// GetApps gets a url
func (r *Roku) QueryActiveApp() (*App, error) {
	a := new(Apps)
	_, err := r.RestClient.R().SetResult(a).Get(makeURL(r.IP) + "/query/active-app")
	if len(a.Apps) != 1 {
		return &App{}, err
	}
	return &a.Apps[0], err
}

func (r *Roku) QueryApps() (*Apps, error) {
	a := new(Apps)
	_, err := r.RestClient.R().SetResult(a).Get(makeURL(r.IP) + "/query/apps")
	return a, err
}

// KeyPress command
func (r *Roku) KeyPress(k string) error {
	if kReal, ok := keys[k]; ok {
		_, err := r.RestClient.R().Post(makeURL(r.IP) + "/keypress/" + kReal)
		return err
	}

	return errors.New("invalid command")
}

func (r *Roku) Literal(str string) error {
	for _, c := range str {
		if _, err := r.RestClient.R().Post(makeURL(r.IP) + "/keypress/Lit_" + url.QueryEscape(string(c))); err != nil {
			return err
		}
	}
	return nil
}

// KeyPress command
func (r *Roku) KeyDown(k string) error {
	if kReal, ok := keys[k]; ok {
		_, err := r.RestClient.R().Post(makeURL(r.IP) + "/keydown/" + kReal)
		return err
	}

	return errors.New("invalid command")
}

// KeyPress command
func (r *Roku) KeyUp(k string) error {
	if kReal, ok := keys[k]; ok {
		_, err := r.RestClient.R().Post(makeURL(r.IP) + "/keyup/" + kReal)
		return err
	}

	return errors.New("invalid command")
}

func (r *Roku) LaunchApp(id string) error {
	_, err := r.RestClient.R().Post(makeURL(r.IP) + "/launch/" + id)
	return err
}

func (r *Roku) LaunchAppName(n string) (string, error) {
	as, err := r.QueryApps()
	if err != nil {
		return "", err
	}
	for _, a := range as.Apps {
		if n == a.Name {
			return a.Name, r.LaunchApp(a.ID)
		}
	}
	return "", errors.New("app not found")
}

func (r *Roku) LaunchAppNameMatch(m string) (string, error) {
	as, err := r.QueryApps()
	if err != nil {
		return "", err
	}
	for _, a := range as.Apps {
		if strings.Contains(strings.ToLower(a.Name), strings.ToLower(m)) {
			return a.Name, r.LaunchApp(a.ID)
		}
	}
	return "", errors.New("app not found")
}
