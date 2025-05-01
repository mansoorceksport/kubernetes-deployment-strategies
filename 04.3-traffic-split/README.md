# Lesson 04.3 - Istio Traffic Splitting (VirtualService + DestinationRule)

In this lesson, we will:
- Define a `DestinationRule` for version subsets
- Use a `VirtualService` to control traffic split
- Observe traffic shift using curl
- Understand how `host` maps to Kubernetes service

---

## 🛠️ Step 1: Create DestinationRule

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: sample-app-destination
spec:
  host: sample-app
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

✅ This defines two named subsets based on pod labels: `version: v1` and `version: v2`.

🧠 The `host: sample-app` refers to the **Kubernetes service name** (`sample-app`). It's not a domain, just an internal DNS name inside the cluster.

---

## 🛠️ Step 2: Send 100% traffic to v1

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: sample-app-vs
spec:
  hosts:
    - sample-app
  http:
    - route:
        - destination:
            host: sample-app
            subset: v1
          weight: 100
```

✅ All traffic goes to v1.

🧠 `VirtualService` controls **how** traffic is split between subsets.
- It references the same `host: sample-app` (the service)
- And decides traffic routing using subset weights

---

## 🛠️ Step 3: Send 90% to v1, 10% to v2

Update the `VirtualService` to:

```yaml
  http:
    - route:
        - destination:
            host: sample-app
            subset: v1
          weight: 90
        - destination:
            host: sample-app
            subset: v2
          weight: 10
```

✅ Gradual rollout.

---

## 🛠️ Step 4: Send 50/50 to v1 and v2

```yaml
  http:
    - route:
        - destination:
            host: sample-app
            subset: v1
          weight: 50
        - destination:
            host: sample-app
            subset: v2
          weight: 50
```

---

## 🛠️ Step 5: Send 100% to v2

```yaml
  http:
    - route:
        - destination:
            host: sample-app
            subset: v2
          weight: 100
```

✅ Traffic cutover complete!

---

## 🧪 Step 6: Test with curl

```bash
kubectl run curl-client --rm -it --image=curlimages/curl --restart=Never -n istio-lab -- sh

while true; do curl -s http://sample-app; sleep 1; done
```

✅ This immediately issues a GET request to `sample-app` from inside the mesh.

Run the command multiple times to see which version responds (`v1` or `v2`).

---

## 🧼 Cleanup

```bash
kubectl delete virtualservice sample-app-vs
kubectl delete destinationrule sample-app-destination
```

---

## ✅ Summary

| Resource | Purpose |
|:--|:--|
| DestinationRule | Define subsets (v1, v2) by labels. These subsets are mapped to real Pod labels. |
| VirtualService | Route traffic to those subsets with weights, using internal service name as host. |
| host | Refers to the Kubernetes **Service name**, not a real domain. For example: `sample-app` refers to the service `sample-app` inside the cluster. |
| Curl Test | Observe live traffic behavior without redeploying apps |

---

## 🧠 FAQ: How Traffic Weights Work in Istio

| Scenario | Result |
|:--|:--|
| `100` + `100` | Normalized to 50% / 50% (equal split) |
| `0` + `0` | Invalid. Traffic has no destination — leads to 503 or routing failure |
| `1000` + `50` | Normalized to ~95.2% and ~4.8% respectively |

🧠 Istio normalizes all `weight` values to percentages based on their sum. You don’t need them to add up to 100 — they’re just relative proportions.

---

Next up: Canary & Blue/Green with Istio Gateways 🔥
