# Request and Validate a Public Certificate with AWS Certificate Manager (ACM)

Many `eksdemo` application installs have an option to expose the application via Kubernetes Ingress. The default Ingress option is with an Application Load Balancer (ALB) using the the AWS Load Balancer Controller. `eksdemo` requires every Ingress to have a fully qualified domain name (FQDN) and use TLS to secure the connection.

ALB integrates natively with AWS Certificate Manager (ACM) to configure TLS for incoming connections. `eksdemo` makes it easy to request and validate an ACM certificate with a single command.

1. [Prerequisites](#prerequisites)
2. [Create an ACM Certificate](#create-an-acm-certificate)
3. [(Optional) Inspect the ACM Certificate](#optional-inspect-the-acm-certificate)
4. [(Optional) Inspect the Hosted Zone Records](#optional-inspect-the-hosted-zone-records)

## Prerequisites

You need an domain that you own configured as a hosted zone in Route 53. If you have a domain but it's not yet configured as a hosted zone, follow the following steps.

1. [Create a Hosted Zone in Route 53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/migrate-dns-domain-inactive.html#migrate-dns-create-hosted-zone-domain-inactive)
2. [Update the domain registration to use Amazon Route 53 name servers](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/migrate-dns-domain-inactive.html#migrate-dns-update-domain-inactive)

## Create an ACM Certificate

The instructions below and in the following sections will refer to the domain as `<example.com>`. Please replace all instances of `example.com` with your Route 53 Hosted Zone.

Use the **`eksdemo create acm-certificate`** command to request an ACM certificate and automatically validate the certificate by adding a CNAME record to the Route 53 hosted zone. In this tutorial we will create a wildcard certificate that can be used across many application installs. **Replace `<example.com>` with your Route 53 Hosted Zone.**

```
» eksdemo create acm-certificate "*.<example.com>"
Creating ACM Certificate request for: *.example.com......done
Created ACM Certificate Id: 2df5cee5-cf53-44db-8857-e99ace4884f8
Validating domain "*.example.com" using hosted zone "example.com"
Waiting for certificate to be issued........done
```

`eksdemo` requests the ACM certificate, adds the necessary CNAME record entries to your Hosted Zone to validate the certificate and waits for the certificate to be issued.

## (Optional) Inspect the ACM Certificate

Use the **`eksdemo get acm-certificate`** command to view your ACM certificates.

```
» eksdemo get acm-certificate
+-----------+--------------------------------------+----------------+--------+--------+
|    Age    |                  Id                  |      Name      | Status | In Use |
+-----------+--------------------------------------+----------------+--------+--------+
| 43 weeks  | c67f9618-9302-4ae5-a93e-dc641903d411 | blue.eks.dev   | ISSUED | Yes    |
| 2 minutes | 2df5cee5-cf53-44db-8857-e99ace4884f8 | *.example.com  | ISSUED | No     |
+-----------+--------------------------------------+----------------+--------+--------+
```

To view the details of the ACM certificate, use the `-o yaml` output option. Remember to replace `<example.com>` with your Route 53 Hosted Zone.

```
» eksdemo get acm "*.<example.com>" -o yaml
- CertificateArn: arn:aws:acm:us-west-2:123456789012:certificate/2df5cee5-cf53-44db-8857-e99ace4884f8
  CertificateAuthorityArn: null
  CreatedAt: "2023-01-31T03:31:52.809Z"
  DomainName: '*.example.com'
  DomainValidationOptions:
  - DomainName: '*.example.com'
    ResourceRecord:
      Name: _d10591305d885fa85023860ca8f99d07.example.com.
      Type: CNAME
      Value: _354518f41374633f455edd1a64448c41.ndlxkpgcgs.acm-validations.aws.
    ValidationDomain: '*.example.com'
    ValidationEmails: null
    ValidationMethod: DNS
    ValidationStatus: SUCCESS
  ExtendedKeyUsages:
  - Name: TLS_WEB_SERVER_AUTHENTICATION
    OID: 1.3.6.1.5.5.7.3.1
  - Name: TLS_WEB_CLIENT_AUTHENTICATION
    OID: 1.3.6.1.5.5.7.3.2
  FailureReason: ""
  ImportedAt: null
  InUseBy: []
  IssuedAt: "2023-01-31T03:32:14.016Z"
  Issuer: Amazon
  KeyAlgorithm: RSA-2048
  KeyUsages:
  - Name: DIGITAL_SIGNATURE
  - Name: KEY_ENCIPHERMENT
  NotAfter: "2024-03-01T23:59:59Z"
  NotBefore: "2023-01-31T00:00:00Z"
  Options:
    CertificateTransparencyLoggingPreference: ENABLED
  RenewalEligibility: INELIGIBLE
  RenewalSummary: null
  RevocationReason: ""
  RevokedAt: null
  Serial: 0d:3e:51:36:a1:8b:2c:9e:95:65:13:d7:fe:2c:97:7c
  SignatureAlgorithm: SHA256WITHRSA
  Status: ISSUED
  Subject: CN=*.example.com
  SubjectAlternativeNames:
  - '*.example.com'
  Type: AMAZON_ISSUED
```

## (Optional) Inspect the Hosted Zone Records

Use the **`eksdemo get hosted-zone`** command to view your Route 53 hosted zones.

```
» eksdemo get hosted-zone
+------------------------+--------+---------+-----------------------+
|          Name          |  Type  | Records |        Zone Id        |
+------------------------+--------+---------+-----------------------+
| example.com            | Public |       3 | Z00613681GX1IL06L0N2S |
| eks.dev                | Public |       6 | Z02488393MP7RTN49WZYP |
| myalias.people.aws.dev | Public |       3 | Z04452933E7K9NRVVP5UC |
+------------------------+--------+---------+-----------------------+
```

To view the DNS records in the hosted zone, use the **`eksdemo get dns-record -z <hosted-zone>`** command. Remember to replace `<example.com>` with your Route 53 Hosted Zone.

```
» eksdemo get dns-records -z <example.com>
+-----------------------------------------------+-------+----------------------------------------------+
|                     Name                      | Type  |                    Value                     |
+-----------------------------------------------+-------+----------------------------------------------+
| example.com                                   | NS    | ns-1234.awsdns-98.co.uk.                     |
|                                               |       | ns-5678.awsdns-76.org.                       |
|                                               |       | ns-123.awsdns-45.net.                        |
|                                               |       | ns-45.awsdns-67.com.                         |
| example.com                                   | SOA   | ns-1234.awsdns-98.co.uk.                     |
|                                               |       | awsdns-hostmaster.amazon.com.                |
|                                               |       | 1 7200 900 1209600 86400                     |
| _d10591305d885fa85023860ca8f99d07.example.com | CNAME | _354518f41374633f455edd1a64448c41.ndlxkpg... |
+-----------------------------------------------+-------+----------------------------------------------+
```

To view details and the raw API output you can re-run the command with the `-o yaml` output flag. You can also choose to delete the CNAME validation record with the **`eksdemo delete dns-record <name>`** command. Note this will prevent the certificate from being automatically renewed.