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

func init() {
	for _, r := range os.Environ() {
		if strings.HasPrefix(r, "ECHO_NEST_API_KEY") {
			split := strings.SplitN(r, "=", 2)
			key = split[1]
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

func (e Echonest) GetCall(path string, method string, stringargs map[string]string, floatargs map[string]float64, intargs map[string]int) http.Response, error {
	args := make([]string, 0, len(stringargs) + len(floatargs) + len(intargs) + 1) // +1 for api key
	for k, v := range stringargs {
		args = append(args, url.QueryEscape(k + "=" + v))
	}
	for k, v := range intargs {
		args = append(args, url.QueryEscape(fmt.Sprintf("%s=%d", v)))
	}
	for k, v := range floatargs {
		args = append(args, url.QueryEscape(fmt.Sprintf("%s=%f", v)))
	}
	
	return Get(strings.Join([]string{"http:/", e.Host, basepath, path, method}, "/") + "?" + strings.Join(args, "/"))
}