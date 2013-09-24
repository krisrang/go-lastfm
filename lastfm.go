package lastfm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	user string
	key  string

	apiRoot = "http://ws.audioscrobbler.com/2.0"
)

type User struct {
	User UserInfo
}

type UserInfo struct {
	Name       string
	Realname   string
	URL        string
	PlayCount  string
	Country    string
	Image      []Image
	Registered Date
}

func (t UserInfo) GetImage() string {
	image := t.Image[len(t.Image)-1]
	return image.URL
}

type Tracks struct {
	Tracks TrackList `json:"recenttracks"`
}

type TrackList struct {
	Tracks []Track `json:"track"`
}

type Track struct {
	Artist Artist
	Name   string
	URL    string
	MBID   string
	NP     NowPlaying `json:"@attr"`
	Image  []Image
	Date   Date
}

func (t Track) IsNowPlaying() bool {
	return t.NP.NowPlaying == "true"
}

type NowPlaying struct {
	NowPlaying string
}

func (t Track) GetImage() string {
	image := t.Image[len(t.Image)-1]
	return image.URL
}

type Artist struct {
	Name string `json:"#text"`
	MBID string
}

type Image struct {
	URL  string `json:"#text"`
	Size string
}

type Date struct {
	Text string `json:"#text"`
	UTS  string `json:"uts,unixtime"`
}

func (d Date) ParseDate() (time.Time, error) {
	date, err := time.Parse("2006-01-02 15:04", d.Text)
	if err != nil {
		date, err = time.Parse("02 Jan 2006, 15:04", d.Text)
		if err != nil {
			return time.Time{}, err
		}
	}

	return date, nil
}

func (d Date) ShortDate() string {
	date, err := d.ParseDate()
	if err != nil {
		return ""
	}

	return (string)(date.Format("2 Jan 2006"))
}

func (d Date) RelativeDate() string {
	date, err := d.ParseDate()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	s := time.Now().Sub(date)

	days := int(s / (24 * time.Hour))
	if days > 1 {
		return fmt.Sprintf("%v days ago", days)
	} else if days == 1 {
		return fmt.Sprintf("%v day ago", days)
	}

	hours := int(s / time.Hour)
	if hours > 1 {
		return fmt.Sprintf("%v hours ago", hours)
	}

	minutes := int(s / time.Minute)
	if minutes > 2 {
		return fmt.Sprintf("%v minutes ago", minutes)
	} else {
		return "Just now"
	}
}

// PUBLIC

func SetConfig(u, k string) {
	user = u
	key = k
}

func GetUser() *UserInfo {
	userdata := &User{}
	getData("user.getinfo", userdata)
	return &userdata.User
}

func GetTracks(limit int) *[]Track {
	trackdata := &Tracks{}
	getData("user.getrecenttracks", trackdata)
	tracks := trackdata.Tracks.Tracks

	if len(tracks) > limit && limit > 0 {
		tracks = tracks[:limit]
	}

	return &tracks
}

//  PRIVATE

func getData(method string, i interface{}) {
	uri := apiRoot + "?method=" + method + "&format=json&user=" + user + "&api_key=" + key
	data := getRequest(uri)
	jsonUnmarshal(data, i)
}

func getRequest(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func jsonUnmarshal(b []byte, i interface{}) {
	err := json.Unmarshal(b, i)
	if err != nil {
		log.Fatal(err)
	}
}
