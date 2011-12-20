package echonest

import "os"
import "net/url"
import "strings"
import "sync"
import "net/http"
import "fmt"

var defaultkey string
var defaultkeylock sync.RWMutex
var defaulthost string = "developer.echonest.com"
var defaulthostlock sync.RWMutex

var basepath string = "api/v4"

type Echonest struct {
	Key string
	Host string
}

type Arg struct {
	Key, Value string
}

// return url-escape key=value string
func (a *Arg) Joined() string {
	return fmt.Sprintf("%s=%s", a.Key, url.QueryEscape(a.Value))
}

func init() {
	for _, r := range os.Environ() {
		if strings.HasPrefix(r, "ECHO_NEST_API_KEY") {
			split := strings.SplitN(r, "=", 2)
			defaultkey = split[1]
			break
		}
	}
}

func SetDefaultKey(k string) {
	defaultkeylock.Lock()
	defaultkey = k
	defaultkeylock.Unlock()
}

// not required currently, but someday Echonest may have unexported fields
func New() *Echonest {
	defaultkeylock.RLock()
	defaulthostlock.RLock()
	o := &Echonest{ defaultkey, defaulthost }
	defaulthostlock.RUnlock()
	defaultkeylock.RUnlock()
	return o
}

func (e *Echonest) keyArg() string {
	return fmt.Sprintf("api_key=%s", e.Key)
}

func (e *Echonest) GetCall(path string, method string, inargs []*Arg) (*http.Response, error) {
	args := make([]string, 0, len(inargs) + 1) // +1 for api key
	for _, r := range inargs {
		args = append(args, r.Joined())
	}
	args = append(args, e.keyArg())
	return http.Get(strings.Join([]string{"http:/", e.Host, basepath, path, method}, "/") + "?" + strings.Join(args, "&"))
}