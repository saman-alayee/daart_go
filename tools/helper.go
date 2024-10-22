package tools

import (
    "fmt"
    "math/rand"
    "time"
)

// UUID generates a version 4 UUID
func UUID() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%04x%04x-%04x-%04x-%04x-%04x%04x%04x",
        rand.Intn(0x10000), rand.Intn(0x10000), rand.Intn(0x10000),
        rand.Intn(0x1000)|0x4000, rand.Intn(0x4000)|0x8000,
        rand.Intn(0x10000), rand.Intn(0x10000), rand.Intn(0x10000),
    )
}
