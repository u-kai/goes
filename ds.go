package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type EKSClusterName string

func (e EKSClusterName) ToKey() string {
	return string(e)
}

type EKSNodeGroupName string

type EKSNodeGroupId struct {
	ClusterName   EKSClusterName
	NodeGroupName EKSNodeGroupName
}

func (e EKSNodeGroupId) ToKey() string {
	return string(e.ClusterName) + ":" + string(e.NodeGroupName)
}
func FromKeyEKSNodeGroupId(key string) (EKSNodeGroupId, error) {
	split := strings.Split(key, ":")
	if len(split) != 2 {
		return EKSNodeGroupId{}, fmt.Errorf("invalid key: %s", key)
	}
	return EKSNodeGroupId{
		ClusterName:   EKSClusterName(split[0]),
		NodeGroupName: EKSNodeGroupName(split[1]),
	}, nil
}

type EKSFargateProfileName string

type EKSFargateProfileId struct {
	ClusterName        EKSClusterName
	FargateProfileName EKSFargateProfileName
}

func (e EKSFargateProfileId) ToKey() string {
	return string(e.ClusterName) + ":" + string(e.FargateProfileName)
}

func FromKeyEKSFargateProfileId(key string) (EKSFargateProfileId, error) {
	split := strings.Split(key, ":")
	if len(split) != 2 {
		return EKSFargateProfileId{}, fmt.Errorf("invalid key: %s", key)
	}
	return EKSFargateProfileId{
		ClusterName:        EKSClusterName(split[0]),
		FargateProfileName: EKSFargateProfileName(split[1]),
	}, nil
}

func DeletePublicEKSFargateProfileActionByClient(client *eks.Client) DeletePublicResourceAction[EKSFargateProfileId] {
	return func(resources []PublicResourceReference[EKSFargateProfileId]) error {
		for _, resource := range resources {
			clusterName := string(resource.Id.ClusterName)
			fargateProfileName := string(resource.Id.FargateProfileName)
			_, err := client.DeleteFargateProfile(context.Background(), &eks.DeleteFargateProfileInput{
				ClusterName:        &clusterName,
				FargateProfileName: &fargateProfileName,
			})
			if err != nil {
				if strings.Contains(err.Error(), "NotFound") {
					continue
				}
				return err
			}
		}
		return nil
	}
}
func SearchPublicSubnetsEKSFargateActionByClient(client *eks.Client) SearchPublicSubnets[EKSFargateProfileId] {
	return func(subnets []PublicSubnetId) ([]PublicResourceReference[EKSFargateProfileId], error) {
		var resources []PublicResourceReference[EKSFargateProfileId]
		output, err := client.ListClusters(context.Background(), &eks.ListClustersInput{})
		if err != nil {
			return nil, err
		}
		for _, cluster := range output.Clusters {
			clusterInfo, err := client.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
				Name: &cluster,
			})
			if err != nil {
				return nil, err
			}
			fargateProfiles, err := client.ListFargateProfiles(context.Background(), &eks.ListFargateProfilesInput{
				ClusterName: &cluster,
			})
			if err != nil {
				return nil, err
			}
			for _, fargateProfile := range fargateProfiles.FargateProfileNames {
				fargateProfileInfo, err := client.DescribeFargateProfile(context.Background(), &eks.DescribeFargateProfileInput{
					ClusterName:        &cluster,
					FargateProfileName: &fargateProfile,
				})
				if err != nil {
					return nil, err
				}
				pubSubnetIds := FilterPublicSubnets(subnets, fargateProfileInfo.FargateProfile.Subnets)
				if len(pubSubnetIds) == 0 {
					continue
				}
				resources = append(resources, PublicResourceReference[EKSFargateProfileId]{

					ResourceType:    EKSFargateProfile,
					VpcId:           VpcId(*clusterInfo.Cluster.ResourcesVpcConfig.VpcId),
					PublicSubnetIds: pubSubnetIds,
					Id:              EKSFargateProfileId{ClusterName: EKSClusterName(cluster), FargateProfileName: EKSFargateProfileName(fargateProfile)},
				})
			}
		}
		return resources, nil
	}
}

func DeletePublicEKSNodeGroupActionByClient(client *eks.Client) DeletePublicResourceAction[EKSNodeGroupId] {
	return func(resources []PublicResourceReference[EKSNodeGroupId]) error {
		for _, resource := range resources {
			clusterName := string(resource.Id.ClusterName)
			nodeGroupName := string(resource.Id.NodeGroupName)

			_, err := client.DeleteNodegroup(context.Background(), &eks.DeleteNodegroupInput{
				ClusterName:   &clusterName,
				NodegroupName: &nodeGroupName,
			})
			if err != nil {
				if strings.Contains(err.Error(), "NotFound") {
					continue
				}
				return err
			}
		}
		return nil
	}
}

func SearchPublicSubnetsEKSNodeGroupActionByClient(client *eks.Client) SearchPublicSubnets[EKSNodeGroupId] {
	return func(subnets []PublicSubnetId) ([]PublicResourceReference[EKSNodeGroupId], error) {
		var resources []PublicResourceReference[EKSNodeGroupId]
		output, err := client.ListClusters(context.Background(), &eks.ListClustersInput{})
		if err != nil {
			return nil, err
		}
		for _, cluster := range output.Clusters {
			clusterInfo, err := client.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
				Name: &cluster,
			})
			if err != nil {
				return nil, err
			}
			nodeGroups, err := client.ListNodegroups(context.Background(), &eks.ListNodegroupsInput{
				ClusterName: &cluster,
			})
			if err != nil {
				return nil, err
			}
			for _, nodeGroup := range nodeGroups.Nodegroups {
				nodeGroupInfo, err := client.DescribeNodegroup(context.Background(), &eks.DescribeNodegroupInput{
					ClusterName:   &cluster,
					NodegroupName: &nodeGroup,
				})
				if err != nil {
					return nil, err
				}
				publicSubnets := FilterPublicSubnets(subnets, nodeGroupInfo.Nodegroup.Subnets)
				if len(publicSubnets) == 0 {
					continue
				}

				resources = append(resources, PublicResourceReference[EKSNodeGroupId]{
					ResourceType:    EKSNodeGroup,
					VpcId:           VpcId(*clusterInfo.Cluster.ResourcesVpcConfig.VpcId),
					PublicSubnetIds: publicSubnets,
					Id: EKSNodeGroupId{
						ClusterName:   EKSClusterName(cluster),
						NodeGroupName: EKSNodeGroupName(nodeGroup),
					},
				})
			}
		}
		return resources, nil
	}
}

func DeletePublicEKSClusterActionByClient(client *eks.Client) DeletePublicResourceAction[EKSClusterName] {
	return func(resources []PublicResourceReference[EKSClusterName]) error {
		for _, resource := range resources {
			name := resource.Id.ToKey()
			_, err := client.DeleteCluster(context.Background(), &eks.DeleteClusterInput{
				Name: &name,
			})
			if err != nil {
				if strings.Contains(err.Error(), "NotFound") {
					continue
				}
				return err
			}
		}
		return nil
	}
}

func SearchPublicSubnetsEKSClusterByClient(client *eks.Client) SearchPublicSubnets[EKSClusterName] {
	return func(subnets []PublicSubnetId) ([]PublicResourceReference[EKSClusterName], error) {
		var resources []PublicResourceReference[EKSClusterName]
		output, err := client.ListClusters(context.Background(), &eks.ListClustersInput{})
		if err != nil {
			return nil, err
		}
		for _, cluster := range output.Clusters {
			clusterInfo, err := client.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
				Name: &cluster,
			})
			if err != nil {
				return nil, err
			}
			publicSubnets := FilterPublicSubnets(subnets, clusterInfo.Cluster.ResourcesVpcConfig.SubnetIds)
			if len(publicSubnets) == 0 {
				continue
			}
			resources = append(resources, PublicResourceReference[EKSClusterName]{
				ResourceType:    EKSCluster,
				VpcId:           VpcId(*clusterInfo.Cluster.ResourcesVpcConfig.VpcId),
				PublicSubnetIds: publicSubnets,
				Id:              EKSClusterName(cluster),
			})
		}
		return resources, nil
	}
}
