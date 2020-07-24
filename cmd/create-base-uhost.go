package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xiayinchang/uk8sctl/config"
	"github.com/xiayinchang/uk8sctl/pkg/util"

	"github.com/spf13/cobra"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
	"github.com/ucloud/ucloud-sdk-go/services/vpc"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

var uhostConfig = &uhost.CreateUHostInstanceRequest{}
var uhostClient *uhost.UHostClient
var unetClient *unet.UNetClient
var vpcClient *vpc.VPCClient
var ucloudConfig = ucloud.NewConfig()

const defaultImage = "centos"

func newCreateBaseUHostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-base-uhost",
		Short: "create base uhost for to make a uk8s image",
		Run: func(cmd *cobra.Command, args []string) {
			setDefaultConfig()
			createUHost()
		},
	}

	uhostConfig.CPU = cmd.Flags().Int("cpu", 2, "Optional. cpu count for uhost")
	uhostConfig.Memory = cmd.Flags().Int("memory", 4096, "Optional. memory size for uhost")
	uhostConfig.Name = cmd.Flags().String("name", "based-on-uk8s-image-for-creating-user-image", "Optional. name for uhost")
	uhostConfig.ChargeType = cmd.Flags().String("charge-type", "Month", "Optional. chargetype for uhost")
	isBoot := "true"
	uhostConfig.Disks = append(uhostConfig.Disks, uhost.UHostDisk{
		IsBoot: &isBoot,
	})
	uhostConfig.Disks[0].Type = cmd.Flags().String("disk-type", "CLOUD_SSD", "Optional. system disk type for uhost")
	uhostConfig.Disks[0].Size = cmd.Flags().Int("disk-size", 40, "Optional. system disk size for uhost")
	uhostConfig.LoginMode = ucloud.String("Password")
	uhostConfig.Password = cmd.Flags().String("password", "", "Required. login password for uhost")
	uhostConfig.MachineType = cmd.Flags().String("machine-type", "N", "Optional. machine type of uhost")
	uhostConfig.MinimalCpuPlatform = cmd.Flags().String("cpu-platform", "Intel/Auto", "Optional. cpu platform of uhost")
	uhostConfig.Zone = cmd.Flags().String("zone", "", "Optional. zone for uhost. if not specified, randomly choose one")
	uhostConfig.VPCId = cmd.Flags().String("vpcid", "", "Optional. vpcid for uhost. if not specified, randomly choose one")
	uhostConfig.SubnetId = cmd.Flags().String("subnetid", "", "Optional. subnetid for uhost. if not specified, randomly choose one")
	uhostConfig.Quantity = cmd.Flags().Int("quantity", 1, "Optional. ")
	cmd.MarkFlagRequired("password")

	return cmd
}

func getUHostClient() *uhost.UHostClient {
	if uhostClient == nil {
		uhostClient = uhost.NewClient(&ucloudConfig, &globalConfig.Credential)
	}
	return uhostClient
}

func getUNetClient() *unet.UNetClient {
	if unetClient == nil {
		unetClient = unet.NewClient(&ucloudConfig, &globalConfig.Credential)
	}
	return unetClient
}

func getVPCClient() *vpc.VPCClient {
	if vpcClient == nil {
		vpcClient = vpc.NewClient(&ucloudConfig, &globalConfig.Credential)
	}
	return vpcClient
}

func setDefaultConfig() {
	uhostConfig.Region = ucloud.String(globalConfig.Region)
	uhostConfig.ProjectId = ucloud.String(globalConfig.ProjectId)
	if *uhostConfig.Zone == "" {
		uhostConfig.Zone = ucloud.String(randomGetZone())
	}
	if *uhostConfig.VPCId == "" {
		uhostConfig.VPCId = ucloud.String(randomGetVPCId())
	}
	if *uhostConfig.SubnetId == "" {
		uhostConfig.SubnetId = ucloud.String(randomGetSubnetId())
	}
	uhostConfig.ImageId = ucloud.String(getUK8SImageId())
}

func randomGetZone() string {
	zones, ok := config.RegionZoneMap[globalConfig.Region]
	if !ok {
		log.Fatalf("regions %s is not supported", globalConfig.Region)
	}
	return zones[0]
}

func randomGetVPCId() string {
	vpcClient := getVPCClient()
	req := vpcClient.NewDescribeVPCRequest()
	req.Region = ucloud.String(globalConfig.Region)
	req.ProjectId = ucloud.String(globalConfig.ProjectId)
	vpcs, err := vpcClient.DescribeVPC(req)
	if err != nil {
		log.Fatalf("describe vpc failed: %s", err.Error())
	}
	if len(vpcs.DataSet) == 0 {
		log.Fatal("can not find available vpc")
	}
	return vpcs.DataSet[0].VPCId
}

func randomGetSubnetId() string {
	vpcClient := getVPCClient()
	req := vpcClient.NewDescribeSubnetRequest()
	req.Region = ucloud.String(globalConfig.Region)
	req.ProjectId = ucloud.String(globalConfig.ProjectId)
	req.VPCId = uhostConfig.VPCId
	subnets, err := vpcClient.DescribeSubnet(req)
	if err != nil {
		log.Fatalf("describe subnet failed: %s", err.Error())
	}

	if len(subnets.DataSet) == 0 {
		log.Fatal("can not find available subnet")
	}
	return subnets.DataSet[0].SubnetId
}

func getUK8SImageId() string {
	req := &DescribeUK8SImageRequest{}
	req.SetAction("DescribeUK8SImage")
	req.SetRequestTime(time.Now())
	req.Region = ucloud.String(globalConfig.Region)
	req.ProjectId = ucloud.String(globalConfig.ProjectId)
	req.Zone = uhostConfig.Zone

	var res DescribeUK8SImageResponse
	err := util.DoRequest(req, &res, &globalConfig.Credential)
	if err != nil {
		log.Fatalf("describe uk8s image failed: %s", err)
	}
	if len(res.ImageSet) == 0 {
		log.Fatalf("no image found")
	}
	for _, imageInfo := range res.ImageSet {
		if strings.Contains(strings.ToLower(imageInfo.ImageName), defaultImage) {
			return imageInfo.ImageId
		}
	}
	return ""
}

type DescribeUK8SImageRequest struct {
	request.CommonBase
}

type DescribeUK8SImageResponse struct {
	response.CommonBase
	ImageSet      []ImageInfo
	PHostImageSet []ImageInfo
}

type ImageInfo struct {
	ZoneId        int
	ImageId       string
	ImageName     string
	NotSupportGPU bool
}

func createUHost() {
	uhostClient = getUHostClient()
	newUHost, err := uhostClient.CreateUHostInstance(uhostConfig)
	if err != nil {
		log.Fatalf("create instance failed: %s", err)
	}

	reqDescribe := uhostClient.NewDescribeUHostInstanceRequest()
	reqDescribe.Region = uhostConfig.Region
	reqDescribe.ProjectId = uhostConfig.ProjectId
	reqDescribe.Zone = uhostConfig.Zone
	reqDescribe.UHostIds = newUHost.UHostIds
	hosts, err := uhostClient.DescribeUHostInstance(reqDescribe)
	if err != nil {
		log.Fatalf("describe uhost failed: %s", err)
	}
	finalHost := hosts.UHostSet[0]
	var buf []byte
	bufWriter := bytes.NewBuffer(buf)
	bufWriter.WriteString("The following uhost instance is created based on uk8s image:\n")
	bufWriter.WriteString(fmt.Sprintf("UhostId:   %s\n", finalHost.UHostId))
	bufWriter.WriteString(fmt.Sprintf("Name:      %s\n", finalHost.Name))
	bufWriter.WriteString(fmt.Sprintf("ProjectId: %s\n", globalConfig.ProjectId))
	bufWriter.WriteString(fmt.Sprintf("Region:    %s\n", globalConfig.Region))
	bufWriter.WriteString(fmt.Sprintf("Zone:      %s\n", finalHost.Zone))
	bufWriter.WriteString(fmt.Sprintf("VPCId:     %s\n", *uhostConfig.VPCId))
	bufWriter.WriteString(fmt.Sprintf("SubnetId:  %s\n", *uhostConfig.SubnetId))
	bufWriter.WriteString(fmt.Sprintf("ImageId:   %s\n", *uhostConfig.ImageId))
	bufWriter.WriteString(fmt.Sprintf("OsName:    %s\n", finalHost.OsName))
	io.Copy(os.Stdout, bufWriter)
}
