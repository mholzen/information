package transforms

import (
	"log"
	"os"
	"testing"

	"github.com/mholzen/information/ilog"
)

func TestMain(m *testing.M) {
	ilog.Init()
	log.SetFlags(0) // Remove timestamps
	os.Exit(m.Run())
}
