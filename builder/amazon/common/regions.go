package common

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
)

func listEC2Regions(ec2conn ec2iface.EC2API) ([]string, error) {
	var regions []string
	resultRegions, err := ec2conn.DescribeRegions(nil)
	if err != nil {
		log.Printf("DescribeRegions: %v", err)
		return nil, fmt.Errorf("listEC2Regions: %v", err)
	}
	for _, region := range resultRegions.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

// ValidateRegion returns an error if the region name is valid
// and exists; otherwise nil.
// ValidateRegion calls ec2conn.DescribeRegions to get the list of
// regions available to this account, a DescribeRegions error
// could be returned
func ValidateRegion(region string, ec2conn ec2iface.EC2API) error {
	regions, err := listEC2Regions(ec2conn)
	if err != nil {
		return err
	}
	for _, valid := range regions {
		if region == valid {
			return nil
		}
	}
	return fmt.Errorf("Invalid region %s, available regions: %v", region, regions)
}
