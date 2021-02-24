# 🔑 amp-eks-iam

amp-eks-iam creates an IAM role to give
remote write privileges to an EKS service account. If you are
collecting Prometheus metrics on EKS and want to send them to [Amazon
Managed Service for Prometheus (AMP)](https://aws.amazon.com/prometheus/),
you can use this tool to give minimal privileges to your Kubernetes
namespace and service account.

## Installation
```
$ go get github.com/aws-observability/amp-eks-iam
```

## Usage

```
amp-eks-iam <cluster flags> [options...]

amp-eks-iam creates the required IAM policies and roles to give
remote write priviledges to an EKS service account.

Example:
$ amp-eks-iam \
   -account=999999999999 -region=us-east-1 -cluster=eks-cluster

Cluster flags:
-account         AWS account ID, for example 999999999999.
-cluster         EKS cluster name.

Options:
-namespace       Kubernetes namespace to apply the policy to. By default, "prometheus".
-service-account Kubernetes service account to apply the policy to. By default, the namespace.
-role            IAM role name to create, default is
                 "EKS-AMP-ServiceAccount-{region}-{cluster}-{namespace}-{sa}".
-region          AWS region of the EKS cluster.
```

By default, amp-eks-iam creates the role and the privileges for the
"prometheus" Kubernetes namespace and service account. You can specify your own
namespaces and service accounts. For example, if you are deploying Grafana Agent
as explained [in this article](https://aws.amazon.com/blogs/opensource/configuring-grafana-cloud-agent-for-amazon-managed-service-for-prometheus/),
use the following command:

```
$ amp-eks-iam \
   -account=999999999999 -region=us-east-1 -cluster=eks-cluster \
   -namespace=grafana-agent \
   -service-account=grafana-agent
```

## Troubleshooting

If you received an error telling "roleName" is above the character limits like below,

```
2021/02/20 09:46:27 Cannot create IAM role: failed to create the IAM role: ValidationError: 1 validation error detected: Value 'EKS-AMP-ServiceAccount-us-west-2-demo-prometheusdeployment-prometheusdeploymentaccount' at 'roleName' failed to satisfy constraint: Member must have length less than or equal to 64
```

You can set a custom role name with -role:

```
$ amp-eks-iam \
   -account=999999999999 -region=us-east-1 -cluster=eks-cluster \
   -role AMPIngestRole
```

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.
