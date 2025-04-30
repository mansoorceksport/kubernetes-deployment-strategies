# Lesson 04.1 - Installing Istio on Minikube 🚀

In this lesson, we install Istio on our local Minikube cluster using `istioctl`.  
We’ll do everything in **baby steps** to ensure you understand what’s going on.

---

## 🧰 Pre-requisites

- Minikube is already running
- You have `kubectl` installed and connected to Minikube
- Internet access (to download Istio CLI and images)

---

## 🛠️ Step 1: Download Istio CLI

```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-*
export PATH="$PWD/bin:$PATH"
```

✅ This downloads Istio CLI and adds `istioctl` to your path temporarily.

---

## 🛠️ Step 2: Install Istio on Minikube

Run the default profile install:

```bash
istioctl install --set profile=demo -y
```

✅ This will install:
- Istiod (the control plane)
- Envoy sidecar injection support
- Istio gateways and CRDs

---

## 🛠️ Step 3: Label the Namespace for Auto Sidecar Injection

We’ll use a new namespace to deploy our Istio apps later:

```bash
kubectl create namespace istio-lab
kubectl label namespace istio-lab istio-injection=enabled
```

check if the label is set:
```bash
kubectl get ns istio-lab -o yaml
```
✅ You should see the label `istio-injection=enabled` in the output.

✅ Now any Pod created in `istio-lab` will have the Envoy proxy automatically injected.

---

## 🛠️ Step 4: Verify the Installation

### 🔎 Check if Istio components are running:

```bash
kubectl get pods -n istio-system
```

✅ You should see Pods like:
- istiod
- istio-ingressgateway
- istio-egressgateway

---

### 🔎 Check if injection works:

Create a dummy pod:

```bash
kubectl run testpod --image=nginx -n istio-lab
kubectl get pod testpod -n istio-lab -o json | jq '.spec.containers[].name'
```

✅ You should see:
```json
"nginx"
"istio-proxy"
```

That means Istio successfully injected the Envoy sidecar! 🎉

---

## 🧼 Step 5: Cleanup (Optional)

If you want to delete everything:

```bash
istioctl x uninstall --purge
kubectl delete namespace istio-system
kubectl delete namespace istio-lab
```

---

## ✅ Summary

| Task | Status |
|:---|:---|
| Download Istio CLI | ✅ |
| Install Istio (demo profile) | ✅ |
| Enable auto-injection on namespace | ✅ |
| Verify components and injection | ✅ |

---

## 🚀 Next Steps

Now you are ready for:

> **Lesson 04.2 — Deploy v1 and v2 Apps with Istio Sidecars**

Where you’ll:
- Deploy your sample-app
- Split traffic using VirtualService + DestinationRule

Let's go deeper into Service Mesh!

---
