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
	"fmt"
	"os"
)

const helpText = `amp-eks-iam <cluster flags> [options...]

amp-eks-iam creates the required IAM policies and roles to give
remote write privileges to an EKS service account.

Example:
$ amp-eks-iam \
   -region=us-east-1 -cluster=eks-cluster

Cluster flags:
-cluster         EKS cluster name.

Options:
-namespace       Kubernetes namespace to apply the policy to. By default, "prometheus".
-service-account Kubernetes service account to apply the policy to. By default, the namespace.
-role            IAM role name to create, default is
                 "EKS-AMP-ServiceAccount-{region}-{cluster}-{namespace}-{sa}".
-region          AWS region of the EKS cluster.`

func usageText(exit int) {
	fmt.Println(helpText)
	os.Exit(exit)
}
