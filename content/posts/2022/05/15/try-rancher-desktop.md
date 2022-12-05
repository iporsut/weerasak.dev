---
title: "ลองใช้ Rancher Desktop"
date: 2022-05-15T16:47:34+07:00
draft: false
---

[Rancher Desktop](https://rancherdesktop.io/) เป็นโปรแกรมที่ช่วยจัดการสร้าง [Kubernetes](https://kubernetes.io/) cluster บนเครื่อง Desktop (local machine) ให้เรา พร้อมลงเครื่องมาต่างๆให้เราเสร็จสรรพ หรือถ้าใครจะลง Rancher Desktop แล้วใช่แค่ [Docker](https://www.docker.com/) ไม่ใช่ Kubernetes ก็ยังได้ เราสามารถลง Rancher Desktop แทนการลง Docker Desktop ได้เลย

<!--more-->

Rancher Desktop รองรับระบบปฏิบัติการทั้ง macOS (Apple Silicon / M1), macOS (Intel), Windows และ Linux

Rancher Desktop มันจะลง Kubernetes โดยใช้ [k3s](https://k3s.io/) แทน ซึ่งเป็น Kubernetes distribution นึงที่มีขนาดเล็ก โดยบน macOS จะทำงานภายใต้ VM ที่รัน Linux ผ่าน [Lima](https://github.com/lima-vm/) ส่วน Windows จะรันผ่าน WSL (Windows Subsystem for Linux v2)

Rancher Desktop ให้เราเลือก version ของ Kubernetes ได้ แล้วเลือก container runtime ได้ด้วยว่าจะใช้ [containerd](https://containerd.io/) หรือจะใช้ dockerd

ส่วน tools commandline ก็ลงมาให้แล้วทั้ง
- docker
- docker-buildx
- kubectl
- nerdctl
- rdctl

และได้ลง Rancher Dashboard สำหรับ local cluster เพื่อดู resources ต่างๆของ Kubernetes ได้อีกด้วย

หลังจากติดตั้งเสร็จ มันจะสร้าง cluster context มาให้ชื่อว่า `rancher-desktop` ถ้าเราใช้หลายๆ cluster ก่อนจะสั่งงาน kubectl ก็ให้คอนฟิคค่าเป็น context นี้ก่อน

## ลอง deploy nginx

```txt
$ kubectl run nginx --image=nginx:alpine --port=80
pod/nginx created
```

```txt
$ kubectl get pods
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          23s
```

```txt
$ kubectl port-forward pods/nginx 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
```

ลองยิง request ไปที่ localhost:8080

```txt
% curl localhost:8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

อยากรู้รายละเอียดเพิ่มเติมและลองเล่นเองบ้างลองเข้าไปดูต่อได้ที่นี่เลย [https://docs.rancherdesktop.io](https://docs.rancherdesktop.io)
