// converters/common.go
package converters

import "time"


func Now() time.Time {
    return time.Now()
}

func CreationTimestamps() (time.Time, time.Time) {
    now := time.Now()
    return now, now
}

func UpdateTimestamp() time.Time {
    return time.Now()
}