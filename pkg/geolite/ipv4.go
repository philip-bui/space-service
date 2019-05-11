package geolite

import (
	"os"

	"github.com/golang/geo/s2"
	"github.com/philip-bui/space-service/pkg/csv"
	"github.com/rs/zerolog/log"
)

// GeoliteIPv4Row is a struct matching a GeoLite2 CSV IPv4 row.
type GeoliteIPv4Row struct {
	IP        string  `csv:"network"`
	GeonameID int     `csv:"geoname_id"`
	Latitude  float64 `csv:"latitude"`
	Longitude float64 `csv:"longitude"`
}

// ReadIPv4 reads a Geolite2 CSV IPv4 file, parses and batches it into PostgreSQL.
func ReadIPv4(fName string) {
	f, err := os.Open(fName)
	defer f.Close()
	if err != nil {
		log.Fatal().Err(err).Str("fName", fName).Msg("error opening file")
	}
	row := &GeoliteIPv4Row{}
	u := csv.NewCsvUnmarshallerFromFile(f, row)
	batch, err := BeginIPBatch()
	if err != nil {
		log.Fatal().Err(err).Msg("error creating ipv4 batch")
	}
	for err = u.UnmarshalToStruct(row); err == nil; err = u.UnmarshalToStruct(row) {
		if err := batch.AddRow(row.IP, row.GeonameID, int64(s2.CellIDFromLatLng(s2.LatLngFromDegrees(row.Latitude, row.Longitude)).Parent(29)>>2)); err != nil {
			log.Fatal().Err(err).Msg("error inserting ipv4")
		}
	}
	if err := batch.ExecAndCommit(); err != nil {
		log.Fatal().Err(err).Msg("error commiting ipv4 batch")
	}
}
