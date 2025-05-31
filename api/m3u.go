package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type URL struct {
	url.URL
	Raw string
}

func (u *URL) String() string {
	return u.Raw
}
func (u *URL) Set(s string) (err error) {
	parsed, err := url.Parse(s)
	u.URL = *parsed
	u.Raw = s
	return
}

type EXTINF struct {
	Id        string `extinf:"tvg-id"`
	Name      string `extinf:"tvg-name"`
	Logo      string `extinf:"tvg-logo"`
	Group     string `extinf:"group-title"`
	Number    int    `extinf:"tvg-chno"`
	Title     string
	Url       string
	SD        bool
	HD        bool
	FHD       bool
	Prefix    string
	NewName   string
	MatchName string
}
type M3UData struct {
	List []*EXTINF
}

func (m3u *M3UData) M3UData() []byte {
	var stringSlice []string

	for _, inf := range m3u.List {
		if inf == nil {
			continue
		}
		/*suffix, prefix := "", ""
		if inf.HD {
			suffix = " HD"
		}
		if inf.FHD {
			suffix = " FHD"
		}
		if inf.Prefix != "" {
			prefix = inf.Prefix + ": "
		}*/

		//name := prefix + inf.NewName + suffix
		name := inf.Title

		//fmt.Printf("#EXTINF:-1 tvg-chno=\"%d\" tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\", %s \n%s\n",
		//    inf.Number, inf.Id, name, inf.Logo, inf.Group, name, inf.Url)
		stringSlice = append(stringSlice, fmt.Sprintf(
			"#EXTINF:-1 tvg-chno=\"%d\" tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\", %s \n%s\n",
			inf.Number,
			inf.Id,
			name,
			inf.Logo,
			inf.Group,
			name,
			inf.Url,
		))
	}
	stringByte := strings.Join(stringSlice, "\x0a")
	return []byte(stringByte)
}

type Channel struct {
	//Keywords         string        `json:"keywords"`
	VodCategory      []interface{} `json:"vod_category"`
	Categories       []interface{} `json:"categories"`
	ID               string        `json:"_id"`
	Title            string        `json:"title"`
	SeriesID         string        `json:"series_id"`
	AiredDate        int64         `json:"aired_date"`
	AllowedCountries interface{}   `json:"allowedCountries"`
	AdPolicyID       interface{}   `json:"adPolicyId"`
	Epg              struct {
		Events []struct {
			Title  string    `json:"title"`
			Start  time.Time `json:"start"`
			End    time.Time `json:"end"`
			Custom struct {
				Duration int    `json:"duration"`
				Rating   string `json:"rating"`
				Image    struct {
					Width       string `json:"width"`
					Height      string `json:"height"`
					DownloadURL string `json:"downloadUrl"`
				} `json:"image"`
			} `json:"custom"`
		} `json:"events"`
	} `json:"epg"`
	Rating    string `json:"rating"`
	MediaType string `json:"mediaType"`
	Order     int    `json:"order"`
	HLSStream struct {
		DownloadURL  string `json:"downloadUrl"`
		StreamingURL string `json:"streamingUrl"`
		URL          string `json:"url"`
	} `json:"HLSStream"`
	CommerceType            string   `json:"commerceType,omitempty"`
	PaidType                string   `json:"paidType,omitempty"`
	SubscriptionsCategories []string `json:"subscriptionsCategories,omitempty"`
	PosterH                 struct {
		DownloadURL string `json:"downloadUrl"`
	} `json:"PosterH"`
	AndroidStream struct {
		DownloadURL  string `json:"downloadUrl"`
		StreamingURL string `json:"streamingUrl"`
		URL          string `json:"url"`
	} `json:"AndroidStream"`
	AndroidBlockedStream struct {
		DownloadURL  string `json:"downloadUrl"`
		StreamingURL string `json:"streamingUrl"`
		URL          string `json:"url"`
	} `json:"AndroidBlockedStream"`
	HLSBlockedStream struct {
		DownloadURL  string `json:"downloadUrl"`
		StreamingURL string `json:"streamingUrl"`
		URL          string `json:"url"`
	} `json:"HLSBlockedStream"`
	LogoLarge        string `json:"logoLarge"`
	ChannelLogoLarge struct {
		DownloadURL  string `json:"downloadUrl"`
		StreamingURL string `json:"streamingUrl"`
		URL          string `json:"url"`
	} `json:"ChannelLogoLarge"`
	ChannelLogoTablets struct {
		DownloadURL  string `json:"downloadUrl"`
		StreamingURL string `json:"streamingUrl"`
		URL          string `json:"url"`
	} `json:"ChannelLogoTablets"`
	PosterF struct {
		DownloadURL string `json:"downloadUrl"`
	} `json:"PosterF,omitempty"`
}

type Channels []Channel

type response struct {
	Title         string   `json:"title"`
	ChannelImgURL string   `json:"channel_img_url"`
	HLSStreamURL  string   `json:"hls_stream_url"`
	Keywords      []string `json:"keywords"`
}

func (u Channels) StreamList() []response {
	var list []response
	for _, channel := range u {

		list = append(list, response{
			Title:         channel.Title,
			ChannelImgURL: channel.ChannelLogoTablets.StreamingURL,
			HLSStreamURL:  channel.HLSBlockedStream.StreamingURL,
			Keywords:      []string{}, //channel.Keywords,
		})
	}
	return list
}

func (u Channels) StreamListToEXTINF(group string) []*EXTINF {
	var list []*EXTINF
	for inx, channel := range u {
		list = append(list, &EXTINF{
			Id:      channel.ID,
			Name:    channel.Title,
			NewName: channel.Title,
			Logo:    channel.ChannelLogoTablets.DownloadURL,
			Url:     channel.HLSBlockedStream.StreamingURL,
			Group:   group,
			Number:  inx,
			Title:   channel.Title,
			FHD:     true,
		})

	}
	return list
}

func getJSON(url string) ([]byte, error) {
	client := &http.Client{}
	fmt.Println("Attempting to fetch for url: ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return bodyBytes, nil
}

func jsonToChannels(bytes []byte) (Channels, error) {
	var response Channels
	err := json.Unmarshal(
		bytes,
		&response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func M3u(w http.ResponseWriter, r *http.Request) {

	
	fmt.Print(r.Header.Get("User-Agent"))
	req, _ := getJSON(os.Getenv("MEDIA_URL"))

	channels, err := jsonToChannels(req)
	//channels
	if err != nil {
		fmt.Fprintf(w, "Error")
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	extInfList := channels.StreamListToEXTINF("TVJ")

	popfd := &M3UData{extInfList}

	w.Write(popfd.M3UData())
}
