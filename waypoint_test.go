package waypoint

import (
	"strings"
	"testing"
)

var testGPX = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="Flyland">
	<metadata>
		<name>Waypointfile by Flyland</name>
	</metadata>
	<wpt lon="8.74905846" lat="46.95258549">
		<ele>1620</ele>
		<name>Achslen</name>
	</wpt>
	<wpt lon="7.77341942" lat="46.70996411">
		<ele>1920</ele>
		<name>Niederhorn 1</name>
	</wpt>
	<wpt lon="7.77669719" lat="46.71130378">
		<ele>1920</ele>
		<name>Niederhorn 2</name>
	</wpt>
	<wpt lon="7.64958013" lat="46.64414295">
		<ele>2280</ele>
		<name>Niesen</name>
	</wpt>
	<wpt lon="8.29266929" lat="46.49895437">
		<ele>1350</ele>
		<name>Ulrichen</name>
	</wpt>
	<wpt lon="7.82421503" lat="46.6810178">
		<ele>570</ele>
		<name>Unterseen Lehn</name>
	</wpt>
	<wpt lon="9.306128" lat="47.195195">
		<ele>908</ele>
		<name>Unterwasser</name>
	</wpt>
</gpx>`

var testDistance = []struct {
	latitude, longitude float64
	distance            int
}{
	{46.72847214, 7.5186329, 0.0},
	{46.95258549, 8.74905846, 96959},
	{46.70996411, 7.77341942, 19554},
	{46.71130378, 7.77669719, 19788},
	{46.64414295, 7.64958013, 13716},
	{46.49895437, 8.29266929, 64469},
	{46.6810178, 7.82421503, 23919},
	{47.195195, 9.306128, 145404},
}

func TestGPX(t *testing.T) {
	gpx, err := New(strings.NewReader(testGPX))
	if err != nil {
		t.Error(err)
	}

	t.Run("Nearest", func(t *testing.T) {
		var name string
		name, _ = gpx.Closest(46.70996411, 7.77341942)
		if name != "Niederhorn 1" {
			t.Errorf("should be \"Niederhorn 1\" - got %s\n", name)
		}
		name, _ = gpx.Closest(46.6810178, 7.82421503)
		if name != "Unterseen Lehn" {
			t.Errorf("should be \"Unterseen Lehn\" - got %s\n", name)
		}
	})

	t.Run("Distance", func(t *testing.T) {
		g := WPT{
			Latitude:  testDistance[0].latitude,
			Longitude: testDistance[0].longitude,
		}
		for _, p := range testDistance[1:] {
			d := int(g.Distance(p.latitude, p.longitude) + 0.5)
			if d != p.distance {
				t.Errorf("distance between %f/%f should be %d not %d", p.latitude, p.longitude, p.distance, d)
			}
		}
	})
}
