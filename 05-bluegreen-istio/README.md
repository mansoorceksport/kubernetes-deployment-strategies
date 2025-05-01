# Lesson 05.1 - Blue/Green Deployment with Istio

In this lesson, we implement a real **Blue-Green Deployment strategy using Istio**.
We will use the knowledge from previous lessons (Gateway, VirtualService, DestinationRule) to switch traffic from `v1` (blue) to `v2` (green) **with zero downtime**.

---

## 📘 What We'll Do
- Use existing Gateway + sample-app Service
- Keep both `v1` and `v2` deployments running
- Route all external traffic to `v1` (blue)
- Then update VirtualService to send traffic to `v2` (green)

---

## 🛠️ Step 1: Ensure both Deployments are running

You should already have these from earlier lessons:
- `sample-app-v1` deployment
- `sample-app-v2` deployment
- Labels: `app: sample-app`, `version: v1` / `v2`

✅ If not, re-apply from lesson `04.2-deploy-v1-v2`

---

## 🛠️ Step 2: VirtualService routing 100% to v1 (Blue)

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: sample-app-bluegreen
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

✅ This will serve `v1` only (blue)

---

## 🧪 Step 3: Access in browser

Access via browser using your previous method:

```bash
minikube ip  # get IP
```
Visit:
```
http://<minikube-ip>:<nodePort>
```
You should see `v1` consistently.

---

## 🔁 Step 4: Switch to `v2` (Green)

Update the VirtualService route:

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

Apply it:
```bash
kubectl apply -f virtual-service-bluegreen.yaml -n istio-lab
```
✅ All traffic now goes to `v2`.

---

## ✅ Summary

| Step | Traffic | Visual |
|--|--|--|
| Initial | 100% to `v1` (Blue) | 🟦🟦🟦 |
| After switch | 100% to `v2` (Green) | 🟩🟩🟩 |

✅ Instant switch without restarting pods
✅ Safe rollback by re-pointing VirtualService to `v1`

---

Next: Canary deployment with traffic split + gradual rollout 🐟
