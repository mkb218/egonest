package echonest

//import "fmt"
import "encoding/json"

type Status struct {
	version string
	code int
	message string
}

type AudioSummary_t struct {
	Duration float64
	Energy float64
	Mode bool
	AnalysisUrl string `json:"analysis_url"`
	TimeSignature int `json:"time_signature"`
	Key int
	Loudness float64
	AudioMd5 string `json:"audio_md5"`
	Tempo float64
	Danceability float64
}

type Track struct {
	ForeignReleaseId string `json:"foreign_release_id"`
	PreviewUrl string `json:"preview_url"`
	ReleaseImage string `json:"release_image"`
	Catalog string
	ForeignId string `json:"foreign_release_id"`
}

type Song struct {
	AudioMd5 string `json:"audio_md5"`
	Id string 
	Title string
	ArtistName string `json:"artist_name"`
	AudioSummary AudioSummary_t `json:"audio_summary"`
	ArtistId string `json:"artist_id"`
	Tracks []Track
}
	
func (e *Echonest) SongSearch(args []*Arg) (songs []Song, err error) {
	r, err := e.GetCall("song", "search", args)
	if err != nil {
		return
	}
	defer r.Body.Close()
	m := make(map[string]interface{})
	// m["response"] = make(map[string]interface{})
	// m["response"]["status"] = new(Status)
	// m["response"]["songs"] = make([]Song, 0)
	jd := json.NewDecoder(r.Body)
	err = jd.Decode(&m)
	return
}