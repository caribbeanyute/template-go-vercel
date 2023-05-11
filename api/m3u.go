package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
	ID                     string        `json:"_id"`
	Title                  string        `json:"title"`
	Keywords               []string      `json:"keywords"`
	Description            interface{}   `json:"description"`
	MediaType              string        `json:"mediaType"`
	AdPolicyID             interface{}   `json:"adPolicyId"`
	Countries              interface{}   `json:"countries"`
	AllowedCountries       bool          `json:"allowedCountries"`
	ExcludeCountries       interface{}   `json:"excludeCountries"`
	Ads                    []interface{} `json:"ads"`
	ChannelLogoSmartphones struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"ChannelLogoSmartphones"`
	HLSBlockedStream struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"HLSBlockedStream"`
	ChannelLogoLarge struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"ChannelLogoLarge"`
	ChannelLogoSmall struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"ChannelLogoSmall"`
	ChannelLogoTablets struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"ChannelLogoTablets"`
	AndroidStream struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"AndroidStream"`
	AndroidBlockedStream struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"AndroidBlockedStream"`
	HLSStream struct {
		DownloadURL  string      `json:"downloadUrl"`
		DaiAssetID   interface{} `json:"daiAssetId"`
		StreamingURL string      `json:"streamingUrl"`
		URL          string      `json:"url"`
	} `json:"HLSStream"`
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
			Keywords:      channel.Keywords,
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
			Logo:    channel.ChannelLogoSmall.DownloadURL,
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
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

	req, _ := getJSON(os.Getenv("MEDIA_URL"))

	channels, err := jsonToChannels(req)
	//channels
	if err != nil {
		fmt.Fprintf(w, "Error")
	}

	extInfList := channels.StreamListToEXTINF("TVJ")

	popfd := &M3UData{extInfList}

	w.Write(popfd.M3UData())
}
