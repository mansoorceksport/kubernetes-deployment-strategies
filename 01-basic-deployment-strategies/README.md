# Lesson 01 - Basic Kubernetes Deployment Strategies

In this lesson, we explore the two fundamental deployment strategies in Kubernetes:

---

## 1. RollingUpdate (default)

- Gradually updates Pods with the new version.
- Ensures minimum downtime.
- Old and new versions may temporarily run together.
- Good for stateless web apps and APIs.

**Visual:**

```
[ OLD ][ OLD ][ OLD ]

↓ rolling update starts

[ NEW ][ OLD ][ OLD ]
[ NEW ][ NEW ][ OLD ]
[ NEW ][ NEW ][ NEW ]
```
---

## 2. Recreate

- Deletes all existing Pods first.
- Then starts new Pods.
- Causes a brief downtime.
- Useful if old and new versions cannot run together (e.g., database migrations).

**Visual:**
```
[ OLD ][ OLD ][ OLD ] (all deleted)

(no pods running)

[ NEW ][ NEW ][ NEW ]
```

---

## ⚡ Pros and Cons

| Strategy | Pros | Cons |
|:---|:---|:---|
| RollingUpdate | Minimal downtime, safe | Old and new can overlap temporarily |
| Recreate | Clean switch, no overlap | Causes downtime during deployment |

---

## 📜 YAML Samples

### RollingUpdate Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rollingupdate-demo
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  selector:
    matchLabels:
      app: demo
  template:
    metadata:
      labels:
        app: demo
    spec:
      containers:
      - name: demo
        image: nginx:1.19
        ports:
        - containerPort: 80
```
### Recreate Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: recreate-demo
spec:
  replicas: 3
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: demo
  template:
    metadata:
      labels:
        app: demo
    spec:
      containers:
      - name: demo
        image: nginx:1.19
        ports:
        - containerPort: 80
```
* ✅ This understanding forms the base for understanding Blue-Green and Canary later. 
* ✅ Master this, and rolling updates will feel natural.

## 📈 Observing Rolling Updates in Action

Once you have deployed the RollingUpdate Deployment,  
you can **simulate an upgrade** and **observe** how Kubernetes gradually updates the Pods.

---

### 🛠 Step 1: Deploy the RollingUpdate Demo

```bash
kubectl apply -f deployment-rollingupdate.yaml
```

✅ This will create 3 Pods using the `nginx:1.19` image.

---

### 🛠 Step 2: Update the Deployment to a New Image Version

Trigger an upgrade by setting a new image:

```bash
kubectl set image deployment/rollingupdate-demo demo=nginx:1.21
```

✅ Kubernetes will now **start a rolling update** to replace old Pods (`nginx:1.19`) with new Pods (`nginx:1.21`).

---

### 🛠 Step 3: Watch the Rollout in Real Time

You can observe the update live using:

```bash
watch -n 2 kubectl get pods
```

✅ This command **refreshes the pod list every 2 seconds**.

✅ You will see:
- Some old Pods terminating
- New Pods being created
- Gradual replacement happening

**Example Output:**

```
NAME                                  READY   STATUS        RESTARTS   AGE
rollingupdate-demo-68c5b7d7c8-abc12   1/1     Running       0          30s
rollingupdate-demo-68c5b7d7c8-def34   0/1     Terminating   0          50s
rollingupdate-demo-68c5b7d7c8-ghi56   1/1     Running       0          1m
```

---

### 🛠 Step 4: Verify Rollout Status

Check if the rollout completed successfully:

```bash
kubectl rollout status deployment/rollingupdate-demo
```

✅ You should see:

> deployment "rollingupdate-demo" successfully rolled out

---

## 🔥 Why This Practice is Powerful

- You can **see** Kubernetes rolling updates happening live.
- You **understand** how downtime is minimized.
- You **build muscle memory** for real production upgrades.

✅ This small exercise gives you a **real feel** of how **RollingUpdate** strategy works behind the scenes.
