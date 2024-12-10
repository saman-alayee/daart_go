package adsize

import "errors"

// AdSize struct holds the width and height for a given ad size.
type AdSize struct {
    Width  int
    Height int
}

// Define a map for converting integer ad sizes to their respective width and height.
var adSizesMap = map[int]AdSize{
    0:  {192, 192},    // SIZE_ICON
    1:  {720, 480},    // SIZE1
    2:  {728, 90},     // SIZE2
    3:  {480, 320},    // SIZE3
    4:  {492, 328},    // SIZE4
    5:  {468, 60},     // SIZE5
    6:  {360, 240},    // SIZE6
    7:  {320, 100},    // SIZE7
    8:  {320, 50},     // SIZE8
    9:  {300, 250},    // SIZE9
    10: {295, 98},     // SIZE10
    11: {160, 600},    // SIZE11
    12: {720, 360},    // SIZE12
    13: {1200, 1200},  // SIZE_INTERSTITIAL
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
