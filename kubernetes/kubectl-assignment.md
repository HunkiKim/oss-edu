# Kubectl 과제


## 사전 조건

`nginx` 디플로이먼트 배포
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

## Viewing and finding

1. 디플로이먼트 json 형식으로 출력
```shell
$ kubectl get deployment nginx-deployment -o json
```
2. 디플로이먼트 `image` 필드 값 출력(jsonpath)
```shell
$ kubectl get deployment nginx-deployment -o=jsonpath='{.spec.template.spec.containers[*].image}'
```
- container는 0이면 첫 번째 컨테이너를 가리킵니다.
3. `app` 라벨 셀렉터를 이용한 파드 목록 조회
```shell
$ kubectl get pods -l app=nginx
```
4. 디플로이먼트 kubectl 조회 시 라벨 목록 포함하여 조회
```shell
$ kubectl get deployment -o wide
```
5. 동일한 결과 출력
```
PodName                             PodUID                                 Image
nginx-deployment-7fb96c846b-g6pms   b11b5344-caf0-4993-b13a-52ede4111bec   nginx:1.14.2
nginx-deployment-7fb96c846b-jz7dg   a97e6512-7e06-4d34-9d58-f496f08a131d   nginx:1.14.2
nginx-deployment-7fb96c846b-xnrk9   9dc294e7-fac5-4721-b63e-efd4e761bb8e   nginx:1.14.2
```
```shell
$ kubectl get pods -o=custom-columns='PodName:metadata.name,PodUID:metadata.uid,Image:spec.containers[*].image'
```