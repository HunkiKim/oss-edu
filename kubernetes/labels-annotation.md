# 레이블

- k8s의 모든 리소스에 붙일 수 있다.
- Labels을 통해 파드와 기타 다른 쿠버네티스 오브젝트의 조직화가 이루어진다.
- 레이블은 리소스에 첨부하는 키-값 쌍으로, 이 쌍은 레이블 셀렉터를 사용해 리소스를 선택할 때 활용
    - 예를들면 app=web 과 같이 부여 가능
    - key를 이용해 오브젝트를 분류
    - value는 구체적인 값으로 애플리케이션의 Pod을 식별하거나 Service 오브젝트가 Pod을 연결할 때, 혹은 로그 수집, 스케일링 작업 등에서 쓰일 수 있습니다.
- 레이블 키가 리소스 내에 고유하다면, 하나 이상 원하는 만큼 레이블을 가질 수 있음
- 일반적으론 생성할때 붙이지만, 나중에 추가하거나 수정 가능

## 레이블 셀렉터
- 실습을 해보자. 아래의 yaml파일은 각각 다른 라벨이 있는 pod이다.
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-1
  labels:
    app: my-app
    environment: production
spec:
  containers:
  - name: nginx
    image: nginx

---

apiVersion: v1
kind: Pod
metadata:
  name: pod-2
  labels:
    app: my-app
    environment: staging
spec:
  containers:
  - name: nginx
    image: nginx

---

apiVersion: v1
kind: Pod
metadata:
  name: pod-3
  labels:
    app: other-app
    environment: production
spec:
  containers:
  - name: nginx
    image: nginx
```
```shell
NAME    READY   STATUS    RESTARTS   AGE
pod-2   1/1     Running   0          38s
pod-1   1/1     Running   0          38s
pod-3   1/1     Running   0          38s
```
```shell
$ k get pods -l app=my-app
NAME    READY   STATUS    RESTARTS   AGE
pod-2   1/1     Running   0          67s
pod-1   1/1     Running   0          67s
```
```shell
$ k get pods -l environment=production
NAME    READY   STATUS    RESTARTS   AGE
pod-1   1/1     Running   0          2m11s
pod-3   1/1     Running   0          2m11s
```
```shell
$ k get pods -l app=my-app,environment=production
NAME    READY   STATUS    RESTARTS   AGE
pod-1   1/1     Running   0          2m54s
```
- 쿠버네티스 객체를 선택할 수 있는 쿼리 기능
- 셀렉터를 사용해 특정 레이블, 레이블 집합 객체를 선택할 수 있음
    - 셀렉터는 쿠버네티스 오브젝트를 검색하거나 묶을 때 사용하는 기능

# 어노테이션

- 파드 및 다른 오브젝트들은 레이블 외에 어노테이션을 가질 수 있다.
- 키-값 쌍으로 레이블과 비슷하지만 식별정보가 없고, 레이블보다 더 큰 값을 기록할 때 많이쓴다.
    - 자동으로 추가되는 경우도 있고, 수동으로도 가능하다
- 쿠버네티스 새로운 기능 추가할 때 주로 쓴다고 한다.
    - Javadocs 같은 느낌인것같다.
    - 하지만 일반적으론 CRD(Custom Resource Definition)을 주로 쓴다고 한다.
- 주로 사용되는 경우는 파드나 다른 API 오브젝트에 설명을 쓸 때 사용해, 개별 오브젝트 관한 정보를 신속하게 찾을 수 있다고 한다.
    - 예를 들어 클러스터 내에 오브젝트를 만든 사람 이름을 어노테이션으로 지정하면, 협업이 용이해진다.

## 오브젝트의 어노테이션 조회

- 실습을 해보자

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: annotated-pod
  annotations:
    description: "This is an annotated pod."
spec:
  containers:
  - name: nginx
    image: nginx
```
```shell
$ k describe pod annotated-pod
Name:             annotated-pod
Namespace:        default
Priority:         0
Service Account:  default
Node:             lima-rancher-desktop/192.168.5.15
Start Time:       Wed, 17 May 2023 23:40:32 +0900
Labels:           <none>
Annotations:      description: This is an annotated pod.
```
```shell

```
- 위를 보면 어노테이션이 들어가 있는걸 볼 수 있다.
- 레이블과 달리 어노테이션은 256KB까지 긴 텍스트 값을 넣을 수 있다.

## 어노테이션 추가 및 수정

레이블을 만들 때와 같은 방법으로 파드를 생성할 때 어노테이션을 추가할 수 있다.

```shell
$ kubectl annotate pod kubia-manual mycompany.com/someannotation="foo bar"
```

위와 같이 annotate 명령어를 쓰면 된다.