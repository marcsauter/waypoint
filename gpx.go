package waypoint

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
)

// GPX represents the element gpx according to GPX 1.1
type GPX struct {
	XMLName   xml.Name `xml:"gpx"`
	Version   string   `xml:"version,attr"`
	Creator   string   `xml:"creator,attr"`
	Waypoints []WPT    `xml:"wpt"`
}

// New returns GPX build from a Reader
func New(r io.Reader) (*GPX, error) {
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

// Closest returns the closest waypoint to the given latitude and longitude
func (g *GPX) Closest(lat, lon float64) (string, int) {
	closest := g.Waypoints[0]
	distance := g.Waypoints[0].Distance(lat, lon)
	for _, wp := range g.Waypoints[1:] {
		d := wp.Distance(lat, lon)
		if d < distance {
			closest = wp
			distance = d
		}
	}
	return closest.Name, int(distance)
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
