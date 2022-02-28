package gpx

import "time"

type XMLGPX struct {
	Tracks []XMLTrack `xml:"trk"`
}

type XMLTrack struct {
	Segments []XMLTrackSegment `xml:"trkseg"`
}

type XMLTrackSegment struct {
	Points []XMLTrackPoint `xml:"trkpt"`
}

type XMLTrackPoint struct {
	Latitude  float64   `xml:"lat,attr"`
	Longitude float64   `xml:"lon,attr"`
	Time      time.Time `xml:"time"`
	Elevation float64   `xml:"ele"`
}
