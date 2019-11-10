package waypoint

import (
	"encoding/xml"
	"fmt"
	"math"
	"strings"
)

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
