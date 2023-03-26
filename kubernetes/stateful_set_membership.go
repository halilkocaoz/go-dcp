package kubernetes

import (
	"fmt"

	"github.com/Trendyol/go-dcp-client/helpers"
	"github.com/Trendyol/go-dcp-client/logger"
	"github.com/Trendyol/go-dcp-client/membership"
)

type statefulSetMembership struct {
	info *membership.Model
}

func (s *statefulSetMembership) GetInfo() *membership.Model {
	return s.info
}

func (s *statefulSetMembership) Close() {
}

func NewStatefulSetMembership(config *helpers.Config) membership.Membership {
	statefulSetInfo, err := NewStatefulSetInfoFromHostname()
	if err != nil {
		logger.ErrorLog.Printf("error while creating statefulSet membership: %v", err)
		panic(err)
	}

	memberNumber := statefulSetInfo.PodOrdinal + 1

	if memberNumber > config.Dcp.Group.Membership.TotalMembers {
		err := fmt.Errorf("memberNumber is greater than totalMembers")
		logger.ErrorLog.Printf("memberNumber: %v, totalMembers: %v, err: %v", memberNumber, config.Dcp.Group.Membership.TotalMembers, err)
		panic(err)
	}

	return &statefulSetMembership{
		info: &membership.Model{
			MemberNumber: memberNumber,
			TotalMembers: config.Dcp.Group.Membership.TotalMembers,
		},
	}
}