package ec2

import (
	// "bytes"
	"fmt"
	"time"
)

const (
	// instanceType
	T1_micro    = "t1.micro"
	M1_small    = "m1.small"
	M1_medium   = "m1.medium"
	M1_large    = "m1.large"
	M1_xlarge   = "m1.xlarge"
	M3_xlarge   = "m3.xlarge"
	M3_2xlarge  = "m3.2xlarge"
	C1_medium   = "c1.medium"
	C1_xlarge   = "c1.xlarge"
	M2_xlarge   = "m2.xlarge"
	M2_2xlarge  = "m2.2xlarge"
	M2_4xlarge  = "m2.4xlarge"
	CR1_8xlarge = "cr1.8xlarge"
	CC1_4xlarge = "cc1.4xlarge"
	CC2_8xlarge = "cc2.8xlarge"
	CG1_4xlarge = "cg1.4xlarge"
	// productDescription
	LINUX_UNIX     = "Linux/UNIX"
	LINUX_UNIX_VPC = "Linux/UNIX (Amazon VPC)"
	SUSE_LINUX     = "SUSE Linux"
	SUSE_LINUX_VPC = "SUSE Linux (Amazon VPC)"
	WINDOWS        = "Windows"
	WINDOWS_VPC    = "Windows (Amazon VPC)"
)

var InstanceTypes []string
var ProductDescriptions []string

func init() {
	InstanceTypes = []string{
		T1_micro,
		M1_small,
		M1_medium,
		M1_large,
		M1_xlarge,
		M3_xlarge,
		M3_2xlarge,
		C1_medium,
		C1_xlarge,
		M2_xlarge,
		M2_2xlarge,
		M2_4xlarge,
		CR1_8xlarge,
		CC1_4xlarge,
		CC2_8xlarge,
		CG1_4xlarge,
	}

	ProductDescriptions = []string{
		LINUX_UNIX,
		LINUX_UNIX_VPC,
		SUSE_LINUX,
		SUSE_LINUX_VPC,
		WINDOWS,
		WINDOWS_VPC,
	}
}

type SpotPriceRequest struct {
	// 2010-08-16T05:06:11.000Z
	StartTime time.Time `xml:"StartTime"`
	EndTime   time.Time `xml:"EndTime"`
	// t1.micro | m1.small | m1.medium |
	// m1.large | m1.xlarge | m3.xlarge |
	// m3.2xlarge | c1.medium | c1.xlarge |
	// m2.xlarge | m2.2xlarge | m2.4xlarge |
	// cr1.8xlarge | cc1.4xlarge |
	// cc2.8xlarge | cg1.4xlarge
	// http://goo.gl/Nk2JJ0
	// Required: NO
	InstanceType string `xml:"InstanceType"`
	// Linux/UNIX | SUSE Linux | Windows |
	// Linux/UNIX (Amazon VPC) |
	// SUSE Linux (Amazon VPC) |
	// Windows (Amazon VPC)
	// Required: NO
	ProductDescription string `xml:"ProductDescription"`
	// us-east-1a, etc.
	// Required: NO
	AvailabilityZone string `xml:"AvailabilityZone"`
}

// Response to DescribeSpotPriceHistory
type SpotPriceResponse struct {
	RequestId      string           `xml:"requestId"`
	SpotHistorySet []*SpotPriceItem `xml:"spotPriceHistorySet>item"`
	NextToken      string           `xml:"nextToken"`
}

// SpotPriceHistorySetItemType
type SpotPriceItem struct {
	InstanceType       string    `xml:"instanceType"`
	ProductDescription string    `xml:"productDescription"`
	SpotPrice          float64   `xml:"spotPrice"`
	Timestamp          time.Time `xml:"timestamp"`
	AvailabilityZone   string    `xml:"availabilityZone"`
}

func (s *SpotPriceItem) Key() string {
	return fmt.Sprintf("%s.%s.%s", s.AvailabilityZone, s.InstanceType, s.ProductDescription)
}

func (ec2 *EC2) SpotPriceHistory(o *SpotPriceRequest, filter *Filter) ([]*SpotPriceItem, error) {

	// make sure that startTime != endTime and startTime is before endTime
	if o.EndTime.Before(o.StartTime) || o.StartTime.Equal(o.EndTime) {
		return nil, fmt.Errorf("'startTime' must be before 'endTime'. startTime: %v, endTime: %v", o.StartTime, o.EndTime)
	}

	// define API call
	params := makeParams("DescribeSpotPriceHistory")

	// convert time objects into RFC3339
	params["StartTime"] = o.StartTime.In(time.UTC).Format(time.RFC3339)
	params["EndTime"] = o.EndTime.In(time.UTC).Format(time.RFC3339)

	// filter on instanceType
	if o.InstanceType != "" {
		params["InstanceType"] = o.InstanceType
	}

	// filter on productDesc
	if o.ProductDescription != "" {
		params["ProductDescription"] = o.ProductDescription
	}

	// filter on availabilityZone
	if o.AvailabilityZone != "" {
		params["AvailabilityZone"] = o.AvailabilityZone
	}

	// add filter parameters
	filter.addParams(params)

	resp := &SpotPriceResponse{}
	err := ec2.query(params, resp)
	if err != nil {
		return nil, err
	}

	return resp.SpotHistorySet, nil
}
