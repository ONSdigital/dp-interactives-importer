package importer_test

import (
	"flag"
	"testing"

	"github.com/ONSdigital/dp-interactives-importer/importer"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	becnhmarkFlag = flag.Bool("benchmark", false, "perform benchmark tests")
)

func TestLargeArchive(t *testing.T) {
	if *becnhmarkFlag {
		Convey("Given a large zip file", t, func() {
			Convey("Then open should run successfully", func() {
				err := importer.Process("/Users/markryan/Postman/files/largetest.zip", importer.EmptyProcessor)
				So(err, ShouldBeNil)
			})
		})
	} else {
		t.Skip("benchmark flag required to run benchmark tests")
	}
}
