package models

import (
)

type AdsRequest struct {
	Web        bool   `bson:"web" json:"web"`
	CampaignID int    `bson:"campaign_id" json:"campaign_id"`
	PublisherID int   `bson:"publisher_id" json:"publisher_id"`
	Created    int64  `bson:"created" json:"created"`
	Origin     string `bson:"origin" json:"origin"`
}
