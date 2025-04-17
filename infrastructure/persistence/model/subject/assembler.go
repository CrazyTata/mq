package subject

import (
	"mq/domain/subject"
)

// POTODOGetSubject converts persistence object to domain object
func POTODOGetSubject(res *Subject) *subject.Subject {
	return &subject.Subject{
		Id:           res.Id,
		AppId:        res.AppId,
		AppSecret:    res.AppSecret,
		MerchantId:   res.MerchantId,
		MerchantName: res.MerchantName,
		IsActive:     res.IsActive,
		ExpireAt:     res.ExpireAt,
		AllowedIps:   res.AllowedIps,
		AllowedPaths: res.AllowedPaths,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
}

// POTODOGetSubjects converts persistence objects to domain objects
func POTODOGetSubjects(res []Subject) []*subject.Subject {
	subjects := make([]*subject.Subject, len(res))
	for i, v := range res {
		subjects[i] = POTODOGetSubject(&v)
	}
	return subjects
}

// DOTOPOSubject converts domain object to persistence object
func DOTOPOSubject(res *subject.Subject) *Subject {
	return &Subject{
		Id:           res.Id,
		AppId:        res.AppId,
		AppSecret:    res.AppSecret,
		MerchantId:   res.MerchantId,
		MerchantName: res.MerchantName,
		IsActive:     res.IsActive,
		ExpireAt:     res.ExpireAt,
		AllowedIps:   res.AllowedIps,
		AllowedPaths: res.AllowedPaths,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
}

// POTODOGetAppSharings converts persistence objects to domain objects
func POTODOGetAppSharings(pos []AppSharing) []*subject.AppSharing {
	if len(pos) == 0 {
		return nil
	}
	result := make([]*subject.AppSharing, len(pos))
	for i, po := range pos {
		result[i] = POTODOGetAppSharing(&po)
	}
	return result
}

// POTODOGetAppSharing converts persistence object to domain object
func POTODOGetAppSharing(po *AppSharing) *subject.AppSharing {
	if po == nil {
		return nil
	}
	return &subject.AppSharing{
		Id:        po.Id,
		AppId:     po.AppId,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

// DOTOPOAppSharing converts domain object to persistence object
func DOTOPOAppSharing(do *subject.AppSharing) *AppSharing {
	if do == nil {
		return nil
	}
	return &AppSharing{
		Id:        do.Id,
		AppId:     do.AppId,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
	}
}
