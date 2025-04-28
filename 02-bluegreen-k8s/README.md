# Lesson 02 - Blue-Green Deployment (Manual in Kubernetes)

In this lesson, we simulate a **Blue-Green Deployment** strategy manually using only Kubernetes Deployments and Services — **without Istio**.

---

## 🧠 What is Blue-Green Deployment?

- Deploy **two separate environments** (Blue and Green).
- Blue = current live version.
- Green = new version, prepared but NOT yet receiving traffic.
- After testing, **switch traffic** to Green.
- Minimal downtime if done carefully.

---

## 🖼️ Visual Diagram

```
Step 1: Initial State
Service --> Blue Pods (v1)

Step 2: Deploy Green quietly
Service --> still Blue Pods (v1)

Step 3: Switch
Service --> Green Pods (v2)
```

✅ No downtime if switch is fast.
✅ Easy rollback: switch back to Blue if Green has issues.

---

## 🛡️ Important Real-World Note

**⚡ Blue-Green is NOT a built-in Kubernetes feature.**  
It is a **manual trick** using:
- Deployments
- Services
- Label selectors

✅ Kubernetes does not "understand" Blue or Green.  
✅ It just **routes traffic based on Pod labels**.

✅ **YOU** are responsible to manage this switching carefully.

---

## 📚 Real-World Story Flavor

Imagine you are upgrading a payment service:

- You deploy Green (new version) quietly.
- You switch Service to Green.
- 🧨 You forgot to test Green properly!
- Customers start failing transactions.
- Now you must **rollback quickly** to Blue by switching Service back!

✅ **Lesson:** Always test the Green version **thoroughly** before switching traffic.

✅ Manual switching is **error-prone** if not careful!

✅ That's why in production, companies automate it — or later use Service Mesh tools like Istio.

---

## 📜 YAML Examples

> 🧠 **Important:**  
> The Service matches Pods **based on their labels**.  
> It does **NOT** match Deployments directly.
>
> In each Deployment below, pay attention to the:
>
> ```yaml
> spec.template.metadata.labels
> ```
>
> These labels are what the Service selector will use to route traffic.

---

### Blue Deployment (v1)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blue-deployment
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

### Green Deployment (v2)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: green-deployment
spec:
  replicas: 3
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

### Service (Initially pointing to Blue)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
    version: v1
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
```

---

## 🛠 Practical Steps

### Step 1: Deploy Blue (v1)

```bash
kubectl apply -f blue-deployment.yaml
kubectl apply -f service.yaml
```

✅ Service points to Blue Pods (`version=v1`).

---

### Step 2: Deploy Green (v2)

```bash
kubectl apply -f green-deployment.yaml
```

✅ Green Pods exist, but Service still sends traffic to Blue.

---

### Step 3: Switch Traffic to Green

Edit Service selector:

```bash
kubectl edit svc myapp-service
```
Change:

```yaml
selector:
  app: myapp
  version: v2
```

✅ Now traffic flows to Green Pods (`version=v2`).

---

### Step 4: Cleanup (Optional)

After confirming stability:

```bash
kubectl delete deployment blue-deployment
```

✅ Clean up old Blue deployment.

---

## 🎨 How Service Matches Pods (Not Deployments)

```
+---------------------+
|      Service        |
| selector:           |
|  app: myapp         |
|  version: v1        |
+----------+----------+
           |
           | Matches labels on Pods
           v
+----------------------+        +----------------------+
|       Pod A          |        |        Pod B          |
| labels:              |        | labels:               |
|  app: myapp          |        |  app: myapp           |
|  version: v1         |        |  version: v1          |
+----------------------+        +----------------------+

(no direct connection to Deployment)

[ Deployment manages Pods separately ]
```

✅ **Services match Pods directly based on labels — not Deployments!**


## ⚡ Key Lessons

- Blue-Green provides minimal downtime.
- Manual switching is powerful but **dangerous if careless**.
- Always **test Green environment before switching**.
- Manual method doesn't scale well for many microservices → **Service Mesh to the rescue!** 🚀

---

