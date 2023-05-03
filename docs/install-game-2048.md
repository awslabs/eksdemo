# Install Game 2048 Example Application with TLS using an ACM Certificate

`eksdemo` makes it extremely easy to install applications from it’s extensive application catalog in your EKS clusters. In this section we will walk through the installation of the Game 2048 example application with TLS using an AWS Certificate Manager (ACM) certificate.

1. [Prerequisites](#prerequisites)
2. [Install Game 2048 Example Application](#install-game-2048-example-application) — Will use optional configuration flags to specify an Ingress with TLS
3. [(Optional) Game 2048 Installation Configurations](#Optional-game-2048-installation-configurations) — How to deploy without a Hosted Zone using CLB or NLB

## Prerequisites

Click on the name of the prerequisites in the list below for a link to detailed instructions.

* [ACM Certificate](/docs/create-acm-cert.md) — A publicly trusted certificate is required to access the Game 2048 example application securely over TLS easily with your web browser of choice. If you don’t have a Hosted Zone configured in Route 53, you can skip this section and deploy the application insecurely over HTTP.
* [AWS Load Balancer Controller](/docs/install-awslb.md) — The Game 2048 example application includes an `Ingress` resource that instructs the AWS Load Balancer Controller to provision an ALB that will enable access to the application over the Internet.
* [ExternalDNS](/docs/install-edns.md) — ExternalDNS will add a DNS record to your Route 53 Hosted Zone for the Game 2048 application. 

### Install Game 2048 Example Application

The [Game 2048](https://play2048.co/) example application is included as [part of the EKS documentation](https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html#application-load-balancer-sample-application) to test and validate the successful deployment of the AWS Load Balancer Controller. The install of the Game 2048 example application includes an `Ingress` resources that instructs the AWS Load Balancer Controller to provision an ALB that will enable access to the application over the Internet.

In this section we will walk through the process of installing the Game 2048 example application. The command for performing the installation is **`eksdemo install example-game-2048 -c <cluster-name>`**

Let’s learn a bit more about the command and it’s options before we continue by using the `-h` help shorthand flag.

```
» eksdemo install example-game-2048 -h
Install example-game-2048

Usage:
  eksdemo install example-game-2048 [flags]

Aliases:
  example-game-2048, example-game2048, example-2048

Flags:
  -c, --cluster string         cluster to install application (required)
      --dry-run                don't install, just print out all installation steps
  -h, --help                   help for example-game-2048
      --ingress-class string   name of IngressClass (default "alb")
  -I, --ingress-host string    hostname for Ingress with TLS (default is Service of type LoadBalancer)
  -n, --namespace string       namespace to install (default "game-2048")
  -X, --nginx-pass string      basic auth password for admin user (only valid with --ingress-class=nginx)
      --nlb                    use NLB instead of CLB (when not using Ingress)
      --replicas int           number of replicas for the deployment (default 1)
      --target-type string     target type when deploying NLB or ALB Ingress (default "ip")
      --use-previous           use previous working chart/app versions (""/"latest")
  -v, --version string         application version (default "latest")

Global Flags:
      --profile string   use the specific profile from your credential file
      --region string    the region to use, overrides config/env settings
```

You’ll notice above there is an optional `--ingress-host` flag with a `-I` shorthand version of the flag. For this application and others that have external access, `eksdemo` defaults to using a Service of type `LoadBalancer` without any encryption (HTTPS). If you have a Hosted Zone configured in Route 53, you will include the Ingress Host flag with the fully qualified domain name for the application, like `-I game2048.example.com`.

Since Game 2048 is included in the EKS documentation as a manifest file, let’s use the the `--dry-run` flag to understand how the application will be installed. **Replace `example.com` with your Hosted Zone.**

```
» eksdemo install example-game-2048 -c blue -I game2048.example.com --dry-run

Manifest Installer Dry Run:
---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: game-2048
  name: deployment-2048
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: app-2048
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: app-2048
    spec:
      containers:
      - image: public.ecr.aws/l6m2t8p7/docker-2048:latest
        imagePullPolicy: Always
        name: app-2048
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  namespace: game-2048
  name: service-2048
  annotations:
    {}
spec:
  ports:
    - port: 80
      targetPort: 80
      protocol: TCP
  type: ClusterIP
  selector:
    app.kubernetes.io/name: app-2048
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: game-2048
  name: ingress-2048
  annotations:
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    alb.ingress.kubernetes.io/target-type: ip
spec:
  ingressClassName: alb
  rules:
    - host: game2048.example.com
      http:
        paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: service-2048
              port:
                number: 80
  tls:
  - hosts:
    - game2048.example.com
```

The dry run output is different from the previous example and simply displays the manifest to be installed. When a Helm chart is not available for an application, the manifest is included in the EKS binary. The manifest is stored as a template and is rendered dynamically each time `eksdemo` is run and can change significantly depending on the flags used. You can run the command again without the `-I` flag to see how the Service object type is changed to `LoadBalancer` and the Ingress resource is removed.

One of the benefits of using a Helm chart is that applications can be easily managed and uninstalled. A powerful feature of `eksdemo` is that all applications are installed as Helm charts even if the application is only available as a manifest. Since `eksdemo` bundles Helm as a library, it dynamically generate a Helm chart in memory from the manifest files before deploying the application.

Now that you know how a manifest install works, let’s install the Game 2048 example application. **Replace `example.com` with your Hosted Zone.**

```
» eksdemo install example-game-2048 -c blue -I game2048.example.com
Helm installing...
2022/11/15 19:45:20 creating 1 resource(s)
2022/11/15 19:45:20 creating 3 resource(s)
Using chart version "n/a", installed "example-game-2048" version "latest" in namespace "game-2048"
```

Let’s check the status of all three of our installed applications, understanding that they are all installed as Helm charts.

```
» eksdemo get application -c blue
+-------------------+--------------+---------+----------+--------+
|       Name        |  Namespace   | Version |  Status  | Chart  |
+-------------------+--------------+---------+----------+--------+
| aws-lb-controller | awslb        | v2.4.4  | deployed | 1.4.5  |
| example-game-2048 | game-2048    | latest  | deployed | n/a    |
| external-dns      | external-dns | v0.12.2 | deployed | 1.11.0 |
+-------------------+--------------+---------+----------+--------+
```

The Ingress resource that is created as part the Game 2048 example application install will trigger the AWS Load Balancer Controller to create an ALB. This will take a few minutes to provision. You can check on the status of the ALB by using the **`eksdemo get load-balancer`** command. For this command, the `-c <cluster-name>` flag is optional, and if used it will filter the query to ELB’s in the same VPC as the `blue` EKS cluster.

```
» eksdemo get load-balancer -c blue
+-----------+--------+----------------------------------+------+-------+-----+-----+
|    Age    | State  |               Name               | Type | Stack | AZs | SGs |
+-----------+--------+----------------------------------+------+-------+-----+-----+
| 3 minutes | active | k8s-game2048-ingress2-0d50dcef8e | ALB  | ipv4  |   3 |   2 |
+-----------+--------+----------------------------------+------+-------+-----+-----+
* Indicates internal load balancer
```

If the state shows as `provisioning`, wait a moment and run the command again.

Next let’s confirm that ExternalDNS has created a Route 53 resource record for our application. The command to query Route 53 records is **`eksdemo get dns-records --zone <zone-name>`.** `eksdemo` has a lot of shorthand aliases for commands and flags and you can discover these by using the `--help` flag on any command. For the `get dns-records` command we’ll use the command alias `dns` and for the `--zone` flag, we’ll use the shorthand `-z`.

**Replace `example.com` with your Hosted Zone.**

```
» eksdemo get dns -z example.com
+----------------------------+------+---------------------------------------------------------------------+
|          Name              | Type |                                Value                                |
+----------------------------+------+---------------------------------------------------------------------+
| example.com                | NS   | ns-1855.awsdns-39.co.uk.                                            |
|                            |      | ns-1452.awsdns-53.org.                                              |
|                            |      | ns-921.awsdns-51.net.                                               |
|                            |      | ns-35.awsdns-04.com.                                                |
| example.com                | SOA  | ns-1855.awsdns-39.co.uk.                                            |
|                            |      | awsdns-hostmaster.amazon.com.                                       |
|                            |      | 1 7200 900 1209600 86400                                            |
| cname-game2048.example.com | TXT  | "heritage=external-dns,external-dns/owner=blue,external-dns/reso... |
| game2048.example.com       | A    | k8s-game2048-ingress2-0d50dcef8e-334176506.us-west-2.elb.amazona... |
| game2048.example.com       | TXT  | "heritage=external-dns,external-dns/owner=blue,external-dns/reso... |
+----------------------------+------+---------------------------------------------------------------------+
```

We can see that an `A` record has been created for `game2048.example.com` that points to the DNS name of the ALB. Next open your web browser and enter `https://game2048.example.com` (**replace `example.com` with your Hosted Zone**) to load your Game 2048 example application!

![Game 2048 Screenshot](/docs/images/game-2048-screenshot.jpg "Game 2048 Screenshot")

Congratulations! You’ve successfully deployed the Game 2048 example application over HTTPS with a publicly trusted certificate!

NOTE: It’s possible you may have to wait for DNS to propagate. The time depends on your local ISP and operating system. If you get a DNS resolution error, you can wait and try again later. Or if you’d like to troubleshoot a bit further, A2 Hosting has a Knowledge base article [How to test DNS with dig and nslookup](https://www.a2hosting.com/kb/getting-started-guide/internet-and-networking/troubleshooting-dns-with-dig-and-nslookup).

Tips:

* Wait a minute or two after the Route 53 A record is created before querying on your computer. I’ve found that if I perform a lookup too fast before DNS has propagated, the operating system can cache the response for some time.
* On my Mac I’ve found that `dig` will directly query the local name servers and will have the latest data and `nslookup` will use the host cache that can have stale data.
* If you believe your DNS cache is to blame, consider this article [How to Flush DNS Cache: Windows and Mac](https://constellix.com/news/how-to-flush-dns-cache-windows-mac).

## (Optional) Game 2048 Installation Configurations

If you don’t have a Hosted Zone or want to deploy the Game 2048 example application unencrypted over HTTP, you can run the command without the `--ingress-host` flag or `-I` shorthand flag: **`eksdemo install example-game-2048 -c blue`**

By default, the application will deployed with a Service of type `LoadBalancer`, which will deploy a Classic Load Balancer (CLB). There are a number of flags that allow you to choose more options:

```
Flags:
      --ingress-class string   name of IngressClass (default "alb")
  -I, --ingress-host string    hostname for Ingress with TLS (default is Service of type LoadBalancer)
  -X, --nginx-pass string      basic auth password for admin user (only valid with --ingress-class=nginx)
      --nlb                    use NLB instead of CLB (when not using Ingress)
      --target-type string     target type when deploying NLB or ALB Ingress (default "ip")
```

To expose the application unencrypted as a Service using an NLB in Instance mode, the command is:
**`eksdemo install example-game-2048 -c blue --nlb --target-type instance`**

To expose the application encrypted as an Ingress using Nginx Ingress, the command is:
**`eksdemo install example-game-2048 -c blue -I game2048.example.com --ingress-class nginx`**

If exposing using a Service and NLB, you will need to have AWS Load Balancer Controller installed. If exposing using an Ingress, you will need to have the Ingress Controller and ExternalDNS installed. Also, if using an IngressClass other than `alb`, you will need to have cert-manager installed.
