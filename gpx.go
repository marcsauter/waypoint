package waypoint

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
)

// NewFromGPX ...
func NewFromGPX(r io.Reader) (*GPX, error) {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	g := GPX{}
	if err := xml.Unmarshal(d, &g); err != nil {
		return nil, err
	}
	return &g, err
}

// WPT represents the wptType according to GPX 1.1
type WPT struct {
	XMLName   xml.Name `xml:"wpt"`
	Name      string   `xml:"name"`
	Latitude  float64  `xml:"lat,attr"`
	Longitude float64  `xml:"lon,attr"`
	Elevation int      `xml:"ele"`
	Comment   string   `xml:"cmt"`
}

// Distance returns the distance to given point
// https://www.mkompf.com/gps/distcalc.html
func (w *WPT) Distance(lat, lon float64) float64 {
	lat1 := lat * math.Pi / 180
	lon1 := lon * math.Pi / 180
	lat2 := w.Latitude * math.Pi / 180
	lon2 := w.Longitude * math.Pi / 180
	return 6378.388 * math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)) * 1000
}

// CompeGPS returns the waypoint in CompeGPS waypoint format
func (w *WPT) CompeGPS() string {
	name := strings.Replace(w.Name, " ", "_", -1)
	// add suffix for latitude
	var lat string
	switch {
	case w.Latitude > 0:
		lat = fmt.Sprintf("%fN", w.Latitude)
	case w.Latitude < 0:
		lat = fmt.Sprintf("%fS", w.Latitude)
	default:
		lat = fmt.Sprintf("%f", w.Latitude)
	}
	// add suffix for longitude
	var lon string
	switch {
	case w.Longitude > 0:
		lat = fmt.Sprintf("%fN", w.Longitude)
	case w.Longitude < 0:
		lat = fmt.Sprintf("%fS", w.Longitude)
	default:
		lat = fmt.Sprintf("%f", w.Longitude)
	}
	return fmt.Sprintf("W %s A %s %s 27-MAR-62 00:00:00 %d %s\n", name, lat, lon, w.Elevation, w.Comment)
}

// GPX represents the element gpx according to GPX 1.1
type GPX struct {
	XMLName   xml.Name `xml:"gpx"`
	Version   string   `xml:"version,attr"`
	Creator   string   `xml:"creator,attr"`
	Waypoints []WPT    `xml:"wpt"`
}

// Nearest returns the nearest waypoint to the given coordinates
func (g *GPX) Nearest(lat, lon float64) (string, int) {
	nearest := g.Waypoints[0]
	distance := g.Waypoints[0].Distance(lat, lon)
	for _, wp := range g.Waypoints[1:] {
		d := wp.Distance(lat, lon)
		if d < distance {
			nearest = wp
			distance = d
		}
	}
	return nearest.Name, int(distance)
}

// CompeGPS returns the waypoints in CompeGPS format according to
// http://www.compegps.com/desarrollo/CompegpsAPI/CompeGPS%20Formats_UK%20v102.pdf
func (g *GPX) CompeGPS() (io.Reader, error) {
	var b bytes.Buffer
	b.WriteString("G WGS 84\n")
	b.WriteString("U 1\n")
	for _, p := range g.Waypoints {
		b.WriteString(p.CompeGPS())
	}
	return bytes.NewReader(b.Bytes()), nil
}
