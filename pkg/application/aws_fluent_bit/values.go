package aws_fluent_bit

const valuesTemplate = `---
image:
  repository: public.ecr.aws/aws-observability/aws-for-fluent-bit
  tag: {{ .Version }}
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  create: true
  name: {{ .ServiceAccount }}
rbac:
  nodeAccess: true
hostNetwork: true
dnsPolicy: ClusterFirstWithHostNet

## Merger of AWS best practice configuration and Helm chart values
## Adds "Flush 5", "Grace 30" and "storage" settings from AWS manifest
## Adds "Health_Check On" from Helm chart values
## Removes "Parsers_File custom_parsers.conf" from Helm chart values
##
## AWS best practices from Container Insights manifest:
## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#L46
##
## Helm values:
## https://github.com/fluent/helm-charts/blob/main/charts/fluent-bit/values.yaml#L345
##
## Documentation:
## https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/classic-mode/configuration-file
config:
  service: |
    [SERVICE]
        Flush                     5
        Grace                     30
        Log_Level                 {{"{{ .Values.logLevel }}"}}
        Daemon                    off
        Parsers_File              parsers.conf
        HTTP_Server               On
        HTTP_Listen               0.0.0.0
        HTTP_Port                 {{"{{ .Values.metricsPort }}"}}
        Health_Check              On
        storage.path              /var/fluent-bit/state/flb-storage/
        storage.sync              normal
        storage.checksum          off
        storage.backlog.mem_limit 5M

  ## AWS best practices from Container Insights manifest for Application INPUT
  ## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#L65
  ##
  ## AWS best practices from Container Insights manifest for Dataplane INPUT
  ## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#L126
  ##
  ## Documentation:
  ## https://docs.fluentbit.io/manual/pipeline/inputs
  inputs: |
    [INPUT]
        Name                tail
        Tag                 application.*
        Exclude_Path        /var/log/containers/cloudwatch-agent*, /var/log/containers/fluent-bit*, /var/log/containers/aws-node*, /var/log/containers/kube-proxy*
        Path                /var/log/containers/*.log
        multiline.parser    docker, cri
        DB                  /var/fluent-bit/state/flb_container.db
        Mem_Buf_Limit       50MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Rotate_Wait         30
        storage.type        filesystem
        Read_from_Head      {{ not .ReadFromTail }}

    [INPUT]
        Name                tail
        Tag                 application.*
        Path                /var/log/containers/fluent-bit*
        multiline.parser    docker, cri
        DB                  /var/fluent-bit/state/flb_log.db
        Mem_Buf_Limit       5MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Read_from_Head      {{ not .ReadFromTail }}

    [INPUT]
        Name                systemd
        Tag                 dataplane.systemd.*
        Systemd_Filter      _SYSTEMD_UNIT=containerd.service
        Systemd_Filter      _SYSTEMD_UNIT=docker.service
        Systemd_Filter      _SYSTEMD_UNIT=kubelet.service
        DB                  /var/fluent-bit/state/systemd.db
        Path                /var/log/journal
        Read_From_Tail      {{ .ReadFromTail }}

    [INPUT]
        Name                tail
        Tag                 dataplane.tail.*
        Path                /var/log/containers/aws-node*, /var/log/containers/kube-proxy*
        multiline.parser    docker, cri
        DB                  /var/fluent-bit/state/flb_dataplane_tail.db
        Mem_Buf_Limit       50MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Rotate_Wait         30
        storage.type        filesystem
        Read_from_Head      {{ not .ReadFromTail }}

  ## AWS best practices from Container Insights manifest for Application FILTER
  ## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#102
  ##
  ## AWS best practices from Container Insights manifest for Dataplane FILTER
  ## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#L149
  ##
  ## Modifications:
  ## 1. Perform the Kubernetes filter on dataplane logs for aws-node* and kube-proxy*
  ## 2. Add rewrite_tag filter for sane log names
  ## 2a. Application pods -- pod.<namespace>.<pod-name>
  ## 2b. Dataplane pods -- pod.<namespace>.<pod-name>.<hostname>
  ## 2c. Dataplane systemd -- systemd.<systemd-unit>.<hostname>
  ##
  ## Documentation:
  ## https://docs.fluentbit.io/manual/pipeline/filters
  filters: |
    [FILTER]
        Name                kubernetes
        Match               application.*
        Kube_Tag_Prefix     application.var.log.containers.
        Merge_Log           On
        Merge_Log_Key       log_processed
        K8S-Logging.Parser  On
        K8S-Logging.Exclude Off
        Labels              Off
        Annotations         Off
        Use_Kubelet         On
        Kubelet_Port        10250
        Buffer_Size         0

    [FILTER]
        Name                modify
        Match               dataplane.systemd.*
        Rename              _HOSTNAME                   hostname
        Rename              _SYSTEMD_UNIT               systemd_unit
        Rename              MESSAGE                     message
        Remove_regex        ^((?!hostname|systemd_unit|message).)*$

    [FILTER]
        Name                aws
        Match               dataplane.*
        imds_version        v1

    [FILTER]
        Name                kubernetes
        Match               dataplane.tail.*
        Kube_Tag_Prefix     dataplane.tail.var.log.containers.
        Merge_Log           On
        Merge_Log_Key       log_processed
        K8S-Logging.Parser  On
        K8S-Logging.Exclude Off
        Labels              Off
        Annotations         Off
        Use_Kubelet         On
        Kubelet_Port        10250
        Buffer_Size         0

    [FILTER]
        Name                rewrite_tag
        Match               application.*
        Rule                $log .* od.$kubernetes['namespace_name'].$kubernetes['pod_name'] false

    [FILTER]
        Name                rewrite_tag
        Match               dataplane.tail.*
        Rule                $log .* d.$kubernetes['namespace_name'].$kubernetes['pod_name'].$kubernetes['host'] false   

    [FILTER]
        Name                rewrite_tag
        Match               dataplane.systemd.*
        Rule                $message .* ystemd.$systemd_unit[0].$hostname false


  ## AWS best practices from Container Insights manifest for Application OUTPUT
  ## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#117
  ##
  ## AWS best practices from Container Insights manifest for Dataplane OUTPUT
  ## https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/main/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml#162
  ##
  ## Modifications:
  ## 1. Log group name changed prefix from /aws/containerinsights to /aws/fluentbit
  ## 2. Log stream prefix removed hostname, starts with either "pod" or "systemd", see rewrite_tag notes in filters section above
  ## 3. Remove "extra_user_agent container-insights"
  ##
  ## Documentation:
  ## https://docs.fluentbit.io/manual/pipeline/outputs
  outputs: |
    [OUTPUT]
        Name                cloudwatch_logs
        Match               od.*
        region              {{ .Region }}
        log_group_name      /aws/fluentbit/{{ .ClusterName }}/application
        log_stream_prefix   p
        auto_create_group   true

    [OUTPUT]
        Name                cloudwatch_logs
        Match               ystemd.*
        region              {{ .Region }}
        log_group_name      /aws/fluentbit/{{ .ClusterName }}/dataplane
        log_stream_prefix   s
        auto_create_group   true

    [OUTPUT]
        Name                cloudwatch_logs
        Match               d.*
        region              {{ .Region }}
        log_group_name      /aws/fluentbit/{{ .ClusterName }}/dataplane
        log_stream_prefix   po
        auto_create_group   true
`
