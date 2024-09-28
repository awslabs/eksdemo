package parameter

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2/recommended/image_id
const eksAL2AMI = "/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id"

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2-arm64/recommended/image_id
const eksAL2Arm64AMI = "/aws/service/eks/optimized-ami/%s/amazon-linux-2-arm64/recommended/image_id"

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2023/x86_64/standard/recommended/image_id
const eksAL2023AMI = "/aws/service/eks/optimized-ami/%s/amazon-linux-2023/x86_64/standard/recommended/image_id"

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2023/arm64/standard/recommended/image_id
const eksAL2023Arm64AMI = "/aws/service/eks/optimized-ami/%s/amazon-linux-2023/arm64/standard/recommended/image_id"

// /aws/service/bottlerocket/aws-k8s-<eks-version>/x86_64/latest/image_id
const bottlerocketAMI = "/aws/service/bottlerocket/aws-k8s-%s/x86_64/latest/image_id"

// /aws/service/bottlerocket/aws-k8s-<eks-version>/arm64/latest/image_id
const bottlerocketArm64AMI = "/aws/service/bottlerocket/aws-k8s-%s/arm64/latest/image_id"

func (g *Getter) GetEKSOptimizedAL2AMI(eksVersion string) (string, error) {
	return g.getEKSOptimizedAMI(eksAL2AMI, eksVersion)
}

func (g *Getter) GetEKSOptimizedAL2Arm64AMI(eksVersion string) (string, error) {
	return g.getEKSOptimizedAMI(eksAL2Arm64AMI, eksVersion)
}

func (g *Getter) GetEKSOptimizedAL2023AMI(eksVersion string) (string, error) {
	return g.getEKSOptimizedAMI(eksAL2023AMI, eksVersion)
}

func (g *Getter) GetEKSOptimizedAL2023Arm64AMI(eksVersion string) (string, error) {
	return g.getEKSOptimizedAMI(eksAL2023Arm64AMI, eksVersion)
}

func (g *Getter) GetBottlerocketAMI(eksVersion string) (string, error) {
	return g.getEKSOptimizedAMI(bottlerocketAMI, eksVersion)
}

func (g *Getter) GetBottlerocketArm64AMI(eksVersion string) (string, error) {
	return g.getEKSOptimizedAMI(bottlerocketArm64AMI, eksVersion)
}

func (g *Getter) getEKSOptimizedAMI(paramName, eksVersion string) (string, error) {
	param, err := g.ssmClient.GetParameter(fmt.Sprintf(paramName, eksVersion))
	if err != nil {
		return "", err
	}
	return aws.ToString(param.Value), nil
}
