# Lesson 04.4 - Exposing Istio App with Gateway (Browser Access)

In this lesson, we expose our app (`sample-app`) to the outside world using an **Istio Gateway**.
We will **continue from the previous lesson**, where we set up traffic splitting (v1 and v2).

---

## 📘 What We'll Do
- Reuse the `DestinationRule` from Lesson 04.3 (required for versioned traffic split)
- Create an Istio `Gateway` resource to listen for external HTTP traffic
- Update the `VirtualService` to bind to the Gateway **and split traffic externally**
- Use `minikube tunnel` to get a real external IP
- Access your app via a browser (e.g., http://<external-ip>)

---

> 🧠 **Note:** `DestinationRule` is only required if you want to use traffic splitting (e.g., v1/v2 subsets).
> In this case, since we're continuing from the previous lesson where we already implemented traffic splitting, we will include it.

---

## 🛠️ Step 1: Apply the Gateway

```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: sample-app-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      hosts:
        - "*"
```

Save as `gateway.yaml` and apply:

```bash
kubectl apply -f gateway.yaml -n istio-lab
```

---

## 🛠️ Step 2: Create a VirtualService bound to the Gateway (with traffic split)

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: sample-app-vs-gateway
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
          weight: 50
        - destination:
            host: sample-app
            subset: v2
            port:
              number: 80
          weight: 50
```

Save as `virtual-service-gateway.yaml` and apply:

```bash
kubectl apply -f virtual-service-gateway.yaml -n istio-lab
```

---

## 🛠️ Step 3: Start minikube tunnel

```bash
sudo minikube tunnel
```

Keep it running in a separate terminal.

---

## 🧪 Step 4: Get External IP

```bash
kubectl get svc istio-ingressgateway -n istio-system
```

Look for `EXTERNAL-IP`, e.g., `192.168.49.2`

---

## 🌐 Step 5: Access in Browser

⚠️ **Important:** The IP shown in the `EXTERNAL-IP` column from `kubectl get svc` may still be an internal cluster IP (e.g., `10.x.x.x`).
This will not work in your browser directly.

✅ Instead, use the IP exposed by `minikube tunnel`. This is often:
- `127.0.0.1`
- or `localhost`
- or a VM IP like `192.168.x.x`

Test in your browser:
```
http://127.0.0.1
```
OR
```
http://localhost
```
If those don’t work, try:
```
curl http://localhost
```
to verify from terminal first.

You should see alternating responses (`v1` and `v2`) based on the 50/50 traffic split.
Try refreshing the page repeatedly!

---

## 🧯 Troubleshooting: minikube tunnel not working (Mac / vfkit issue)

If you're using **Minikube with vfkit on macOS**, you may encounter a situation where:
- `minikube tunnel` runs, but
- The EXTERNAL-IP of `istio-ingressgateway` is stuck as a `10.x.x.x` internal IP
- Browser access via `localhost` / `127.0.0.1` / `192.168.x.x` doesn't work

✅ This is a known issue. Reference:
https://github.com/kubernetes/minikube/issues/10085

### 🛠️ Workaround: Use NodePort instead of LoadBalancer

Temporarily expose Istio Ingress Gateway via NodePort:

```bash
kubectl patch svc istio-ingressgateway -n istio-system -p '{"spec":{"type":"NodePort"}}'
```

Get the exposed port:

```bash
kubectl get svc istio-ingressgateway -n istio-system
```
Look for something like:
```
80:31560/TCP
```

Then get the Minikube VM IP:

```bash
minikube ip
```

Example:
```
192.168.49.2
```

✅ Now access in browser:
```
http://192.168.49.2:31560
```

When done, you can **revert the Istio Ingress Gateway service back to `LoadBalancer`** using:

```bash
kubectl patch svc istio-ingressgateway -n istio-system -p '{"spec":{"type":"LoadBalancer"}}'
```

---

## ✅ Summary

| Resource | Purpose |
|--|--|
| Gateway | Opens port 80 to the world (via Istio Ingress) |
| VirtualService | Binds your host/path to internal service and applies traffic split |
| DestinationRule | Enables use of `subsets` like v1/v2 for routing decisions |
| minikube tunnel | Bridges LoadBalancer to local network |
| External IP | Access app in browser! |

---

Next: Add Hostnames and TLS with Ingress Gateway 🚀
