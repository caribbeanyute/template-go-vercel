package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// UrlStruct struct for parsing and setting URLs
type UrlStruct struct {
	url.URL
	Raw string
}

func (u *UrlStruct) String() string {
	return u.Raw
}

func (u *UrlStruct) Set(s string) (err error) {
	parsed, err := url.Parse(s)
	u.URL = *parsed
	u.Raw = s
	return
}

// ExtinfEntry struct represents a single entry in an M3U playlist.
type ExtinfEntry struct {
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

// M3uPlaylistData struct holds a list of ExtinfEntry entries.
type M3uPlaylistData struct {
	List []*ExtinfEntry
}

// GenerateM3uData generates the M3U playlist content as a byte slice.
func (mpd *M3uPlaylistData) GenerateM3uData() []byte {
	var stringSlice []string

	for _, inf := range mpd.List {
		if inf == nil {
			continue
		}
		name := inf.Title

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

// MediaChannel struct represents the structure of a channel object from the external JSON API.
type MediaChannel struct {
	VodCategory []interface{} `json:"vod_category"`
	Categories  []interface{} `json:"categories"`
	ID          string        `json:"_id"`
	Title       string        `json:"title"`
	SeriesID    string        `json:"series_id"`
	AiredDate   int64         `json:"aired_date"`
	AllowedCountries interface{} `json:"allowedCountries"`
	AdPolicyID  interface{}   `json:"adPolicyId"`
	Epg         struct {
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
	Rating      string `json:"rating"`
	MediaType   string `json:"mediaType"`
	Order       int    `json:"order"`
	HLSStream   struct {
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

// MediaChannels is a slice of MediaChannel structs.
type MediaChannels []MediaChannel

// StreamResponse struct for a simplified stream list.
type StreamResponse struct {
	Title         string   `json:"title"`
	ChannelImgURL string   `json:"channel_img_url"`
	HLSStreamURL  string   `json:"hls_stream_url"`
	Keywords      []string `json:"keywords"`
}

// GetStreamList converts MediaChannels to a slice of simplified StreamResponse structs.
func (mc MediaChannels) GetStreamList() []StreamResponse {
	var list []StreamResponse
	for _, channel := range mc {
		list = append(list, StreamResponse{
			Title:         channel.Title,
			ChannelImgURL: channel.ChannelLogoTablets.StreamingURL,
			HLSStreamURL:  channel.AndroidStream.StreamingURL,
			Keywords:      []string{},
		})
	}
	return list
}

// ConvertToExtinfList converts MediaChannels to a slice of ExtinfEntry structs.
func (mc MediaChannels) ConvertToExtinfList(group string) []*ExtinfEntry {
	var list []*ExtinfEntry
	for inx, channel := range mc {
		list = append(list, &ExtinfEntry{
			Id:        channel.ID,
			Name:      channel.Title,
			NewName:   channel.Title,
			Logo:      channel.ChannelLogoTablets.DownloadURL,
			Url:       channel.AndroidStream.StreamingURL,
			Group:     group,
			Number:    inx,
			Title:     channel.Title,
			FHD:       true,
		})
	}
	return list
}

// fetchJSONData fetches JSON data from a given URL.
func fetchJSONData(url string) ([]byte, error) {
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

// unmarshalJSONToMediaChannels unmarshals JSON bytes into a MediaChannels slice.
func unmarshalJSONToMediaChannels(bytes []byte) (MediaChannels, error) {
	var response MediaChannels
	err := json.Unmarshal(
		bytes,
		&response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}


// --- XMLTV RELATED CODE ---

// Tv struct represents the root <tv> element in XMLTV.
type Tv struct {
	XMLName           xml.Name       `xml:"tv"`
	Date              string         `xml:"date,attr"`
	GeneratorInfoName string         `xml:"generator-info-name,attr"`
	SourceInfoName    string         `xml:"source-info-name,attr"`
	Channels          []XmltvChannel `xml:"channel"`
	Programmes        []Programme    `xml:"programme"`
}

// XmltvChannel represents a <channel> element in XMLTV.
type XmltvChannel struct {
	XMLName     xml.Name      `xml:"channel"`
	ID          string        `xml:"id,attr"`
	DisplayName []DisplayName `xml:"display-name"`
	Icon        *Icon         `xml:"icon,omitempty"` // Add icon for channel logo
	URL         string        `xml:"url"`
}

// Icon represents an <icon> element for channel logos.
type Icon struct {
	XMLName xml.Name `xml:"icon"`
	Src     string   `xml:"src,attr"`
}

// DisplayName represents a <display-name> element in XMLTV.
type DisplayName struct {
	XMLName xml.Name `xml:"display-name"`
	Lang    string   `xml:"lang,attr"`
	Text    string   `xml:",chardata"`
}

// Programme represents a <programme> element in XMLTV.
type Programme struct {
	XMLName xml.Name `xml:"programme"`
	Start   string   `xml:"start,attr"`
	Stop    string   `xml:"stop,attr"`
	Channel string   `xml:"channel,attr"`
	Title   []Title  `xml:"title"`
	Desc    []Desc   `xml:"desc"`
	Category []Category `xml:"category,omitempty"` // Add category
	Rating   *Rating    `xml:"rating,omitempty"`   // Add rating
}

// Title represents a <title> element in XMLTV.
type Title struct {
	XMLName xml.Name `xml:"title"`
	Lang    string   `xml:"lang,attr"`
	Text    string   `xml:",chardata"`
}

// Desc represents a <desc> element in XMLTV.
type Desc struct {
	XMLName xml.Name `xml:"desc"`
	Lang    string   `xml:"lang,attr"`
	Text    string   `xml:",chardata"`
}

// Category represents a <category> element in XMLTV.
type Category struct {
	XMLName xml.Name `xml:"category"`
	Lang    string   `xml:"lang,attr"`
	Text    string   `xml:",chardata"`
}

// Rating represents a <rating> element in XMLTV.
type Rating struct {
	XMLName xml.Name `xml:"rating"`
	System  string   `xml:"system,attr,omitempty"`
	Value   string   `xml:"value"`
}

// generateXMLTVData generates XMLTV data from a slice of Channel structs.
func generateXMLTVData(channels MediaChannels) ([]byte, error) { // Updated type to MediaChannels
	tv := Tv{
		Date:              time.Now().Format("20060102"),
		GeneratorInfoName: "MyGoEPGGenerator",
		SourceInfoName:    "EPG Data from Go Application",
	}

	for _, ch := range channels {
		xmltvChannel := XmltvChannel{
			ID: ch.ID,
			DisplayName: []DisplayName{
				{Lang: "en", Text: ch.Title},
			},
			URL: ch.HLSStream.StreamingURL,
		}
		// Add channel logo if available
		if ch.ChannelLogoTablets.DownloadURL != "" {
			xmltvChannel.Icon = &Icon{Src: ch.ChannelLogoTablets.DownloadURL}
		}
		tv.Channels = append(tv.Channels, xmltvChannel)

		if len(ch.Epg.Events) > 0 {
			for _, event := range ch.Epg.Events {
				startFormatted := event.Start.Format("20060102150405 -0700")
				stopFormatted := event.End.Format("20060102150405 -0700")

				programme := Programme{
					Start:   startFormatted,
					Stop:    stopFormatted,
					Channel: ch.ID,
					Title: []Title{
						{Lang: "en", Text: event.Title},
					},
					Desc: []Desc{
						{Lang: "en", Text: fmt.Sprintf("Duration: %d minutes. Rating: %s.", event.Custom.Duration, event.Custom.Rating)},
					},
				}

				// Add category if available (e.g., from VodCategory or Categories)
				if len(ch.VodCategory) > 0 {
					if catStr, ok := ch.VodCategory[0].(string); ok {
						programme.Category = []Category{{Lang: "en", Text: catStr}}
					}
				} else if len(ch.Categories) > 0 {
					if catStr, ok := ch.Categories[0].(string); ok {
						programme.Category = []Category{{Lang: "en", Text: catStr}}
					}
				}

				// Add rating if available
				if event.Custom.Rating != "" {
					programme.Rating = &Rating{System: "MPAA", Value: event.Custom.Rating}
				}

				tv.Programmes = append(tv.Programmes, programme)
			}
		} else {
			// If no EPG events are present, generate some dummy programs
			now := time.Now().UTC()
			programDuration := time.Hour

			for p_num := 0; p_num < 3; p_num++ {
				startTime := now.Add(time.Duration(p_num) * programDuration)
				stopTime := startTime.Add(programDuration)

				programme := Programme{
					Start:   startTime.Format("20060102150405 -0700"),
					Stop:    stopTime.Format("20060102150405 -0700"),
					Channel: ch.ID,
					Title: []Title{
						{Lang: "en", Text: fmt.Sprintf("Dummy Show %d on %s", p_num+1, ch.Title)},
					},
					Desc: []Desc{
						{Lang: "en", Text: fmt.Sprintf("This is a placeholder description for Dummy Show %d on %s.", p_num+1, ch.Title)},
					},
				}
				tv.Programmes = append(tv.Programmes, programme)
			}
		}
	}

	xmlBytes, err := xml.MarshalIndent(tv, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling XML: %w", err)
	}

	xmlDeclaration := []byte(xml.Header)
	doctypeDeclaration := []byte(`<!DOCTYPE tv SYSTEM "xmltv.dtd">` + "\n")

	finalXML := append(xmlDeclaration, doctypeDeclaration...)
	finalXML = append(finalXML, xmlBytes...)

	return finalXML, nil
}

// XMLTVHandler is the HTTP handler for fetching EPG data in XMLTV format.
// This function is exported and can be used as a Vercel handler.
func XMLTV(w http.ResponseWriter, r *http.Request) {
	//mediaURL := "https://1spotmedia.com/index.php/api/vod/get_live_streams" //os.Getenv("MEDIA_URL")
	// If you want to use an environment variable, uncomment the line below and comment out the hardcoded URL:
	mediaURL := os.Getenv("MEDIA_URL")

	if mediaURL == "" {
		http.Error(w, "MEDIA_URL environment variable is not set", http.StatusInternalServerError)
		return
	}

	reqBytes, err := fetchJSONData(mediaURL) // Renamed function call
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching media data from MEDIA_URL: %v", err), http.StatusInternalServerError)
		return
	}

	channels, err := unmarshalJSONToMediaChannels(reqBytes) // Renamed function call
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing channels JSON: %v", err), http.StatusInternalServerError)
		return
	}

	xmlData, err := generateXMLTVData(channels)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating XMLTV data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlData)
}
