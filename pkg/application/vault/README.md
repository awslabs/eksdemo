# `Hashicorp Vault` - Swiss Army Knife of key-value Stores

## Documentation
[Kubernetes Deployment Guide](https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide)
[EKS-Specific Guide](https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-amazon-eks)

## Install with EKSDemo
`eksdemo install vault`

## Installation Arguments
- `enable-tls` -> Enable TLS for end-to-end encrypted transport
  - Diabled by default
- `replicas` -> Configure Number of replicas (3 recommended)
  - 1 replica by default

## Installation Validation

Following successful installation, confirm that the Pods have been created (note that 0/1 for Pod(s) `vault-*` is expected)
```
$ kubectl --namespace vault get pods
NAME                                   READY   STATUS    RESTARTS   AGE
vault-0                                0/1     Running   0          2m12s
vault-agent-injector-95c7b5566-f8trr   1/1     Running   0          2m12s
```

The expected output confirms that Vault is not initialized and is sealed:
```
$ kubectl --namespace vault exec vault-0 -- vault status
Key                Value
---                -----
Seal Type          shamir
Initialized        false
Sealed             true
Total Shares       0
Threshold          0
Unseal Progress    0/0
Unseal Nonce       n/a
Version            1.16.1
Build Date         2024-04-03T12:35:53Z
Storage Type       file
HA Enabled         false
command terminated with exit code 2
```

Initialize and unseal Vault:
```
$ kubectl --namespace vault exec vault-0 -- vault operator init \
   -key-shares=1 \
   -key-threshold=1 \
   -format=json > cluster-keys.json
```

> Note: This command can be executed as-is for any number of replicas

> Note: This command captures the root token (along with other data) into a local file named `cluster-keys.json`
```
{
  "unseal_keys_b64": [
    "6onhnY+ENwYCw2Bd9ts0SGUZac6YJ2go41wERaRr2jg="
  ],
  "unseal_keys_hex": [
    "ea89e19d8f84370602c3605df6db3448651969ce98276828e35c0445a46bda38"
  ],
  "unseal_shares": 1,
  "unseal_threshold": 1,
  "recovery_keys_b64": [],
  "recovery_keys_hex": [],
  "recovery_keys_shares": 0,
  "recovery_keys_threshold": 0,
  "root_token": "hvs.uruz2wt9NhbE1eXrvPU1ci62"
}
```

Vault is now initialized but still sealed (confirm by re-executing `kubectl --namespace vault exec vault-0 -- vault status`) - in order to unseal, use the `unseal_key` value from the initialization:
```
$ VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]") && \
  kubectl --namespace vault exec vault-0 -- vault operator unseal $VAULT_UNSEAL_KEY
```

Use the `root_token` to verify that all is well by logging into Vault:
```
$ CLUSTER_ROOT_TOKEN=$(cat cluster-keys.json | jq -r ".root_token") && \
  kubectl --namespace vault exec vault-0 -- vault login $CLUSTER_ROOT_TOKEN
```

Expected output is similar to:
```
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                  Value
---                  -----
token                hvs.uruz2wt9NhbE1eXrvPU1ci62
token_accessor       CviOdFCUsR41ZhBZSw8J0zzA
token_duration       ∞
token_renewable      false
token_policies       ["root"]
identity_policies    []
policies             ["root"]
```

Congratulations on a working Vault!


## Dependencies
- The Helm install will include a PersistentVolume for Vault that must succeed, otherwise the vault.vault-0 Pod will not successfully schedule or launch; there is nothing particularly unique about the PV but its creation success will be fully depedent on a working storage layer.
  CSI management is out of the scope of this application.
