# Kmesh Bypass

## Introduction
The Bypass feature allows you to temporarily remove a pod from the Kmesh/Istio mesh traffic path. This is particularly useful for troubleshooting, as it helps you determine whether a communication issue is caused by the mesh service (Kmesh/Sidecar) or by the application logic itself.

When bypass is enabled for a pod, its inbound and outbound traffic will skip the mesh redirection rules and communicate directly with the peer.

## Prerequisites
- Kmesh is installed and running in your cluster.
- The bypass controller is enabled (it is disabled by default in some versions, check `enable-bypass` flag in Kmesh options).

## How to Enable Bypass
To enable the bypass function for a specific pod, add the label `kmesh.net/bypass=enabled` to the pod:

```bash
kubectl label pod <pod_name> kmesh.net/bypass=enabled
```

Once this label is added, the Kmesh daemon will detect it and add short-circuit rules (iptables `RETURN`) to the pod's network namespace.

## How to Disable Bypass
To restore normal mesh operation and remove the bypass, simply remove the label or set it to a different value:

```bash
kubectl label pod <pod_name> kmesh.net/bypass-
```

## How it Works
The Kmesh bypass controller watches for pods with the `kmesh.net/bypass=enabled` label. When it detects such a pod:
1. It enters the pod's network namespace.
2. It inserts `iptables -t nat -I PREROUTING 1 -j RETURN` and `iptables -t nat -I OUTPUT 1 -j RETURN` rules.
3. These rules ensure that any subsequent mesh redirection rules (added by Istio or Kmesh) are skipped.

For Kmesh-native mode (eBPF based redirection), the bypass also ensures that the eBPF programs do not intercept the traffic (currently implemented via similar short-circuiting logic).

## Use Cases
- **Troubleshooting**: If you experience connectivity issues between services, bypass one service to see if the issue persists. If it works without the mesh, the issue might be in how the mesh handles that traffic.
- **Performance Testing**: Compare performance with and without mesh interception for a specific pod without affecting the entire cluster.
