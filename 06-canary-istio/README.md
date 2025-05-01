# Lesson 05.2 - Canary Deployment with Istio

In this lesson, we perform a **Canary Deployment using Istio**, where a small portion of traffic is gradually routed to a new version (`v2`) before going fully live.
This allows safe testing in production with the ability to roll back instantly.

---

## 📘 What We'll Do
- Start with 100% traffic to `v1`
- Introduce `v2` with a small weight (canary)
- Gradually shift traffic to `v2`
- Final step: 100% to `v2` (full rollout)

---

## 🛠️ Step 1: Start with 100% to `v1`

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: sample-app-canary
spec:
  hosts:
  - "*"
  gateways:
  - sample-app-gateway
  http:
  - route:
    - destination:
        host: sample-app
        subset: v1
        port:
          number: 80
      weight: 100
```

✅ All external traffic goes to `v1`

---

## 🛠️ Step 2: Introduce Canary (10% to `v2`)

Update the `VirtualService` with traffic split:

```yaml
  http:
  - route:
    - destination:
        host: sample-app
        subset: v1
        port:
          number: 80
      weight: 90
    - destination:
        host: sample-app
        subset: v2
        port:
          number: 80
      weight: 10
```

✅ Now `v2` receives 10% of the traffic.

---

## 🔁 Step 3: Gradual rollout (50/50)

```yaml
      weight: 50  # v1
```
```yaml
      weight: 50  # v2
```

Apply the update:
```bash
kubectl apply -f virtual-service-canary.yaml -n istio-lab
```

---

## 🚀 Step 4: Full Rollout (100% to `v2`)

```yaml
  http:
  - route:
    - destination:
        host: sample-app
        subset: v2
        port:
          number: 80
      weight: 100
```

✅ Canary becomes the new production.

---

## 🧪 How to Observe

### Option 1: Curl loop inside cluster
```bash
kubectl run curl-loop --rm -it --image=curlimages/curl --restart=Never -n istio-lab -- sh
while true; do curl -s http://sample-app; sleep 1; done
```

### Option 2: Access via browser (if using NodePort)
```bash
http://<minikube-ip>:<nodePort>
```
Repeat refresh and observe `v1`/`v2` responses.

---

## ✅ Summary

| Stage | v1 Weight | v2 Weight | Behavior |
|--|--|--|--|
| Start | 100% | 0% | Safe base |
| Canary | 90% | 10% | Testing new version |
| Mid-Roll | 50% | 50% | Shared load |
| Full Rollout | 0% | 100% | `v2` becomes prod |

---

Next: Automating canary with metrics + rollback triggers (advanced) 🚦
