package converters

import "myapp/models"

import mongoModels "myapp/mongo_models"



func ToDomainLeadSlice(mongoLeads []mongoModels.Lead) []models.Lead {
	leads := make([]models.Lead, len(mongoLeads))
	for i, mongoLead := range mongoLeads {
		leads[i] = ToDomainLead(mongoLead)
	}
	return leads
}

func ToDomainLead(mongoLead mongoModels.Lead) models.Lead {
	return models.Lead{
		ID:           mongoLead.ID.Hex(),
		Name:         mongoLead.Name,
		Phone:        mongoLead.Phone,
		Requirement:  mongoLead.Requirement,
		AadharNumber: mongoLead.AadharNumber,
		AadharPhoto:  mongoLead.AadharPhoto,
	}
}

