# Lesson 04.2 - Deploy v1 and v2 Apps with Istio Sidecars

In this lesson, we deploy two versions of a sample app (`v1` and `v2`) into our Istio-enabled namespace `istio-lab`.

We will:
- Deploy `v1` and `v2` pods
- Confirm Envoy sidecars are injected
- Expose both using a shared Kubernetes Service
- Prepare for traffic shifting in the next lesson (VirtualService)

---

## 🛠️ Step 1: Deploy v1

```bash
kubectl apply -f sample-app-v1.yaml -n istio-lab
```

---

## 🛠️ Step 2: Deploy v2

```bash
kubectl apply -f sample-app-v2.yaml -n istio-lab
```

---

## 🛠️ Step 3: Create a shared Service

```bash
kubectl apply -f sample-app-service.yaml -n istio-lab
```

This service will forward traffic to any pod with:

```yaml
labels:
  app: sample-app
```

So both `v1` and `v2` will receive traffic!

---

## 🧪 Step 4: Confirm Sidecars Are Injected

```bash
kubectl get pods -n istio-lab
```

You should see:

```
sample-app-v1-xxxxx   2/2   Running
sample-app-v2-xxxxx   2/2   Running
```

✅ 2/2 = app container + istio-proxy sidecar

---

## 🧪 Step 5: Curl from a test client

```bash
kubectl run test-client --image=alpine -it --rm --restart=Never -n istio-lab -- sh
```

Inside the shell:

```sh
apk add curl
curl http://sample-app
```

Through Browser
``` minikube tunnel
> minikube tunnel

Status: 
        machine: minikube
        pid: 8901
        route: 10.96.0.0/12 -> 192.168.106.2
        minikube: Running
        services: [sample-app, test-nginx-svc, istio-ingressgateway]
    errors: 
                minikube: no errors
                router: no errors
                loadbalancer emulator: no errors

```
Get the
IP : 192.168.106.2

```
> kubectl get svc -n istio-lab

NAME             TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
sample-app       LoadBalancer   10.103.179.193   10.103.179.193   80:32625/TCP   4m16s

```
Get the Port: 32625

Then you can access the app using the following URL:


On Browser
```
http://192.168.106.2:32625
```

You will randomly hit either `v1` or `v2`.

---

## 📘 What We’ve Learned

| Concept | Explanation |
|:--|:--|
| Multiple deployments with same app label | Enables service to route to both |
| Istio injection enabled via namespace label | Each pod has Envoy proxy |
| Simple curl shows traffic load-balancing | (for now, round-robin) |

---

## ⏭️ Next Lesson: 04.3 - Traffic Shifting with VirtualService

In the next lesson, we’ll define an **Istio VirtualService** and **DestinationRule** to:

- Send 100% to v1
- Then shift 90/10
- Then 50/50
- Then 100% to v2

🚀 This is where Istio really shines.

---
