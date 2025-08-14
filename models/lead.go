package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct{
	ID               primitive.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	DealerID         primitive.ObjectID      `json:"dealer_id" bson:"dealer_id"`
	Name             string                  `json:"name" bson:"name"`
	Area             string                  `json:"area" bson:"area"`
	Requirement      string                  `json:"requirement" bson:"requirement"`
	Status           string                  `json:"status" bson:"status"`
	CreatedOn        time.Time               `json:"created_on" bson:"created_on"`
}

const (
	LeadStatusNew              = "new"
	LeadStatusMatchFound       = "match_found"
	LeadStatusDealerContacted  = "dealer_contacted"
	LeadStatusDealInProgress   = "deal_in_progress"
	LeadStatusSuccess          = "success"
)

var AllLeadStatuses = []string{
	LeadStatusNew,
	LeadStatusMatchFound,
	LeadStatusDealerContacted,
	LeadStatusDealInProgress,
	LeadStatusSuccess,
}
