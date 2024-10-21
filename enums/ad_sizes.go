package adsize

import "errors"

// AdSize struct holds the width and height for a given ad size.
type AdSize struct {
    Width  int
    Height int
}

// Define a map for converting integer ad sizes to their respective width and height.
var adSizesMap = map[int]AdSize{
    1:  {192, 192},    // SIZE_ICON
    2:  {720, 480},    // SIZE1
    3:  {728, 90},     // SIZE2
    4:  {480, 320},    // SIZE3
    5:  {492, 328},    // SIZE4
    6:  {468, 60},     // SIZE5
    7:  {360, 240},    // SIZE6
    8:  {320, 100},    // SIZE7
    9:  {320, 50},     // SIZE8
    10: {300, 250},    // SIZE9
    11: {295, 98},     // SIZE10
    12: {160, 600},    // SIZE11
    13: {720, 360},    // SIZE12
    14: {1200, 1200},  // SIZE_INTERSTITIAL
}

// GetAdSizeByID returns the corresponding AdSize based on the provided ID (adsize).
// It returns an error if the adsize is not found in the map.
func GetAdSizeByID(adsize int) (AdSize, error) {
    adSize, exists := adSizesMap[adsize]
    if !exists {
        return AdSize{}, errors.New("invalid adsize")
    }
    return adSize, nil
}
