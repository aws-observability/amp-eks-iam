/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	defaultNamespace          = "prometheus"
	defaultIAMRole            = "EKS-AMP-ServiceAccount"
	defaultIAMRoleDescription = "IAM role to be used by a K8s service account with write access to AMP"
	ampRemoteWritePolicy      = "arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess"
)

var (
	account        string // AWS account ID
	cluster        string // EKS cluster
	region         string // AWS region, optional
	namespace      string // Kubernetes namespace, optional
	serviceAccount string // Kubernetes service account, optional
	role           string // IAM role name to create, optional
)

func main() {
	flag.StringVar(&account, "account", "", "AWS account ID")
	flag.StringVar(&cluster, "cluster", "", "EKS cluster name")
	flag.StringVar(&region, "region", "", "EKS cluster's region")
	flag.StringVar(&namespace, "namespace", defaultNamespace, "EKS namespace to restrict the IAM policy for")
	flag.StringVar(&serviceAccount, "service-account", "", "EKS service account")
	flag.StringVar(&role, "role", "", "IAM role to be created or updated")
	flag.Usage = func() {
		usageText(0)
	}
	flag.Parse()

	if account == "" || cluster == "" {
		usageText(1)
	}
	if region == "" {
		reg, err := defaultRegion()
		if err != nil {
			log.Fatalf("Cannot identify the default region: %v", err)
		}
		region = reg
	}
	if serviceAccount == "" {
		serviceAccount = namespace
	}
	if role == "" {
		role = fmt.Sprintf("%s-%s-%s-%s-%s",
			defaultIAMRole, region, cluster, namespace, serviceAccount)
	}

	cfg := &aws.Config{}
	if region != "" {
		cfg.Region = aws.String(region)
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		log.Fatalf("Cannot create session: %v", err)
	}

	if err := createRole(sess); err != nil {
		log.Fatalf("Cannot create IAM role: %v", err)
	}
	log.Printf("Role %q is created.", role)
}

func createRole(sess *session.Session) error {
	iamSvc := iam.New(sess)
	eksSvc := eks.New(sess)
	clusterOut, err := eksSvc.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(cluster),
	})
	if err != nil {
		return fmt.Errorf("failed to find the EKS cluster: %v", err)
	}

	trustDoc := bytes.NewBuffer(nil)
	if err := trustDocTmpl.Execute(trustDoc, &trustDocVars{
		Account:        account,
		Namespace:      namespace,
		ServiceAccount: serviceAccount,
		OIDC:           strings.ReplaceAll(*clusterOut.Cluster.Identity.Oidc.Issuer, "https://", ""),
	}); err != nil {
		return fmt.Errorf("failed to generate the trust relationship document: %v", err)
	}

	_, err = iamSvc.CreateRole(&iam.CreateRoleInput{
		RoleName:                 aws.String(role),
		AssumeRolePolicyDocument: aws.String(trustDoc.String()),
		Description:              aws.String(defaultIAMRoleDescription),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != iam.ErrCodeEntityAlreadyExistsException {
				return fmt.Errorf("failed to create the IAM role: %v", err)
			}
			// TODO(jbd): Instead of returning an error, validate the document.
			return fmt.Errorf("role %q already exists, delete it manually to recreate", role)
		}
	}

	_, err = iamSvc.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: aws.String(ampRemoteWritePolicy),
		RoleName:  aws.String(role),
	})
	if err != nil {
		return fmt.Errorf("failed to attach the policy to the role: %v", err)
	}

	return nil
}

func defaultRegion() (string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	defaultRegion := sess.Config.Region
	if defaultRegion == nil {
		return "", errors.New("no default region is set")
	}
	return *defaultRegion, nil
}

var trustDocTmpl = template.Must(template.New("trust-doc").Parse(`{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect":"Allow",
			"Principal":{
				"Federated": "arn:aws:iam::{{.Account}}:oidc-provider/{{.OIDC}}"
			},
			"Action":"sts:AssumeRoleWithWebIdentity",
			"Condition":{
				"StringEquals": {
					"{{.OIDC}}:sub": "system:serviceaccount:{{.Namespace}}:{{.ServiceAccount}}"
				}
			}
		}
	]
}`))

type trustDocVars struct {
	Account        string
	Namespace      string
	ServiceAccount string
	OIDC           string
}
