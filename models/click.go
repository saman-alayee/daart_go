package models

// Click represents a click event in the system.
type Click struct {
	ViewID      string `bson:"view_ID" json:"view_ID"`
	PublisherID int    `bson:"publisher_ID" json:"publisher_ID"`
	CampaignID  int    `bson:"campaign_ID" json:"campaign_ID"`
	Requested   int64  `bson:"requested" json:"requested"` // Unix timestamp when the click was requested
	Counted     bool   `bson:"counted" json:"counted"`     // Initially set to 0, can be updated later
	IP          string `bson:"ip" json:"ip"`               // Client IP address
}
