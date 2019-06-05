package code

// App time:            09-12-2019 15:02
// Reftime:             Mon Jan 2 15:04:05 -0700 MST 2006

import (
	"fmt"
	"time"
)

const (
	reftime                = "Mon Jan 2 15:04:05 -0700 MST 2006"
	reftimeNow             = "2006-01-02 15:04:05.000000 -0700 MST m=+0.000000000"
	reftimemmddyyyyhhcolmm = "01-02-2006 15:04"
)


func TranslateUTCandFormat(t time.Time) string {
	// Convert to UTC
	tutc := t.UTC()
	return fmt.Sprintf("%s", tutc.Format(reftimemmddyyyyhhcolmm))
}
