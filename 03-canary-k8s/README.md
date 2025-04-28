# Lesson 03 - Canary Deployment (Manual in Kubernetes)

In this lesson, we simulate a **Canary Deployment** strategy manually using Kubernetes Deployments and Services — **without Istio**.

---

## 🧠 What is Canary Deployment?

- Deploy a **small number** of new version Pods alongside the stable Pods.
- Route **some traffic** (e.g., 10%) to the new version.
- Monitor carefully for any issues.
- If successful, gradually increase rollout.
- If errors, rollback without affecting all users.

---

## 🖼️ Visual Diagram

```
Step 1: Initial State
Service --> [ v1 Pod ][ v1 Pod ][ v1 Pod ]

Step 2: Add Canary (small v2 Pod)
Service --> [ v1 Pod ][ v1 Pod ][ v2 Pod ]

Step 3: Monitor traffic
- v2 gets a small percentage of traffic

Step 4: Success
- Deploy more v2 Pods
- Retire v1 Pods

OR

Step 4: Failure
- Remove v2 Pod
- Stick with v1 Pods
```

---

## 🛡️ Important Real-World Note

**⚡ Manual Canary in Kubernetes is a trick, not a native feature.**  
You create **two Deployments** (v1 and v2) and **let the Service load-balance between them**.

✅ Service doesn’t understand versions.  
✅ Service just load-balances blindly across all selected Pods.

---

## 📚 Real-World Story Flavor

Imagine you are releasing a new signup page:

- You deploy a Canary (v2) with only 1 Pod.
- v2 has a hidden bug causing signup form failure.
- Customers randomly hitting v2 experience failure.
- If **too many** customers hit the bad Pod, it **damages user trust**.
- You realize **manual canary is dangerous without good monitoring!**

✅ **Lesson:**
- Always monitor Canary Pods closely.
- Better use traffic control tools later (like Istio).

---

## 📜 YAML Examples

---

### v1 Deployment (stable version)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: v1-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
      version: v1
  template:
    metadata:
      labels:
        app: myapp
        version: v1
    spec:
      containers:
      - name: myapp
        image: nginx:1.19
        ports:
        - containerPort: 80
```

---

### v2 Deployment (canary version)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: v2-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
      version: v2
  template:
    metadata:
      labels:
        app: myapp
        version: v2
    spec:
      containers:
      - name: myapp
        image: nginx:1.21
        ports:
        - containerPort: 80
```

---

### Service (targets both v1 and v2)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
```

## 🎯 How Kubernetes Services Select Pods

In Kubernetes, a **Service** uses **selectors** to decide **which Pods to send traffic to**.

✅ A Service does **NOT care about the Deployment**.  
✅ A Service **ONLY cares about Pod labels**.

When you define a Service:

```yaml
spec:
  selector:
    app: myapp
```

It tells Kubernetes:

> "Find all running Pods that have the label `app: myapp` — and load-balance traffic to them."

---

### 🔥 In Our Canary Example:

| Object | Labels | Notes |
|:---|:---|:---|
| v1 Pods | `app: myapp, version: v1` | Selected ✅ |
| v2 Pods (canary) | `app: myapp, version: v2` | Selected ✅ |

✅ Because both v1 and v2 Pods have `app: myapp`,  
✅ The Service selects **both types of Pods**.

---

### 🧠 Important Tip:

- If you change only **app** label, Service will or will not select Pods.
- If you want to control *only v1* or *only v2*, you must create **different Services** or use **different selectors**.

---

### 🛠️ Watching the Selector in Real Time

You can **watch** which Pods are being targeted by the Service:

```bash
kubectl get endpoints myapp-service -o wide
```

✅ This will show the **IP addresses of Pods** currently selected by the Service.

✅ If you see both v1 and v2 Pod IPs listed —  
✅ Your manual Canary traffic split is working!

---

### 📦 Simple Diagram

```
[ Service selector: app=myapp ]
           |
           |--> Pod 1 (app=myapp, version=v1)
           |--> Pod 2 (app=myapp, version=v1)
           |--> Pod 3 (app=myapp, version=v1)
           |--> Pod 4 (app=myapp, version=v2)  (canary)
```

✅ All Pods that match `app=myapp` are **treated equally** for load-balancing.

---

# ⚡ Why Understanding This is Critical

- ✅ You fully control traffic flow with **labels and selectors**.
- ✅ If labels are wrong, traffic may go to wrong Pods.
- ✅ It explains why Canary traffic split is based on **Pod count** in manual method.
- ✅ Later when you use Istio, you will appreciate how **fine-grained traffic routing** (not based only on Pod count) becomes possible.

---



> 🧠 **Note:**  
> The Service selector is broad (`app: myapp`) — so it **selects both** v1 and v2 Pods.

✅ The v2 canary Pod will receive a **small portion** of total traffic because there are fewer v2 Pods.

---

## 🛠 Practical Steps

### Step 1: Deploy v1 (stable version)

```bash
kubectl apply -f v1-deployment.yaml
kubectl apply -f service.yaml
```

✅ All traffic goes to v1 Pods initially.

---

### Step 2: Deploy v2 (canary version)

```bash
kubectl apply -f v2-deployment.yaml
```

✅ Now traffic is split: 3 Pods (v1) + 1 Pod (v2).

---

### Step 3: Monitor

✅ Watch logs, metrics, alerts.

✅ Check if v2 Pod is causing any errors.

---

### Step 4: Decide

- ✅ If success: Gradually increase v2 replicas, retire v1.
- ❌ If failure: Delete v2 Deployment, stick with v1.

---

## ⚡ Key Lessons

- Manual canary gives **basic traffic split** based on Pod count.
- No fine-grained control over percentage (e.g., exactly 5% traffic).
- High risk without observability.
- In real production, better to use traffic shaping tools like Istio, Linkerd, etc.

---

