package models

// Click_mobile represents a click_mobile event in the system.
type Click_mobile struct {
	PublisherID string `bson:"publisher_ID" json:"publisher_ID"`
	CampaignID  string `bson:"campaign_ID" json:"campaign_ID"`
	Requested   int64  `bson:"requested" json:"requested"` // Unix timestamp when the click_mobile was requested
	Counted     bool    `bson:"counted" json:"counted"`     // Initially set to 0, can be updated later
	IP          string `bson:"ip" json:"ip"`               // Client IP address
}
