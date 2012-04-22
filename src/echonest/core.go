package echonest

import "encoding/json"
import "log"
import "bytes"
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

func (e *Echonest) Upload(filetype string, data []byte) (id, analysis_url string, err error) {
	args := make([]string, 3)
	args[0] = e.keyArg()
	args[1] = "bucket=audio_summary"
	args[2] = (&Arg{"filetype", filetype}).Joined()
	req, err := http.NewRequest("POST", strings.Join([]string{"http:/", e.Host, basepath, "track", "upload"}, "/") + "?" + strings.Join(args, "&"), bytes.NewReader(data))
	if err != nil {
		log.Println("NewReq error")
		return
	}
	req.ContentLength = int64(len(data))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("post error")
		return
	}
	defer resp.Body.Close()
	j := json.NewDecoder(resp.Body)
	var tmpA struct {
		Response struct {
			Track struct {
				Id string `json:"id"`
				Audio_summary AudioSummary_t `json:"audio_summary"`
			} `json:"track"`
		} `json:"response"`
	}
		
	err = j.Decode(&tmpA)
	if err != nil {
		log.Println("decode error")
		return
	}
	id = tmpA.Response.Track.Id
	analysis_url = tmpA.Response.Track.Audio_summary.AnalysisUrl
	return	
}


func (e *Echonest) Analyze(id string) (analysis_url string, err error) {
	resp, err := http.PostForm(strings.Join([]string{"http:/", e.Host, basepath, "track", "analyze"}, "/"),
		url.Values{"api_key": {e.Key}, "bucket": {"audio_summary"}, "id": {id}})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	j := json.NewDecoder(resp.Body)
	var tmpA struct {
		Response struct {
			Track struct {
				Id string `json:"id"`
				Audio_summary AudioSummary_t `json:"audio_summary"`
			} `json:"track"`
		} `json:"response"`
	}
		
	err = j.Decode(&tmpA)
	if err != nil {
		return
	}
	analysis_url = tmpA.Response.Track.Audio_summary.AnalysisUrl
	return	
}

