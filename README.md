## AMP-EKS-IAM tool

The amp-eks-iam tool helps reduce friction for Amazon Managed Service for Prometheus (AMP) customers by providing remote write privileges to authorized EKS accounts with an easier way to set up IAM.

amp-eks-iam creates the required IAM policies and roles to give remote write privileges to an EKS service account. 

Also if you are collecting Prometheus metrics on EKS and want to write them to Amazon Managed Service for Prometheus (AMP), you can use this tool to give minimal privileges to a Kubernetes namespace and service account to send data to AMP.

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.

