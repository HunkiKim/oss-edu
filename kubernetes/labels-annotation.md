# 레이블

- k8s의 모든 리소스에 붙일 수 있다.
- Labels을 통해 파드와 기타 다른 쿠버네티스 오브젝트의 조직화가 이루어진다.
- 레이블은 리소스에 첨부하는 키-값 쌍으로, 이 쌍은 레이블 셀렉터를 사용해 리소스를 선택할 때 활용
    - 예를들면 app=web 과 같이 부여 가능
    - key를 이용해 오브젝트를 분류
    - value는 구체적인 값으로 애플리케이션의 Pod을 식별하거나 Service 오브젝트가 Pod을 연결할 때, 혹은 로그 수집, 스케일링 작업 등에서 쓰일 수 있습니다.
- 레이블 키가 리소스 내에 고유하다면, 하나 이상 원하는 만큼 레이블을 가질 수 있음
- 일반적으론 생성할때 붙이지만, 나중에 추가하거나 수정 가능

## 권장 레이블
- 권장 레이블은 지원 도구 외에도 쿼리하는 방식으로 애플리케이션을 식별하게 한다.
- 메타데이터는 애플리케이션의 개념을 중심으로 정리된다.
- 애플리케이션에 대해 공식적인 개념이 없거나 강요하지 않는다.
  - 대신 애플리케이션은 비공식적이며 메타데이터로 설명된다.
- 메타데이터들은 권장하는 레이블이다.
- 공용 레이블과 주석에는 공통 접두사인 *app.kubernetes.io* 가 있다.
  - 접두사가 없는 레이블은 사용자가 개인적으로 사용할 수 있다.
  - 공유 접두사는 공유 레이블이 사용자 정의 레이블을 방해하지 않도록 한다.
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

- 모든 리소스들은 레이블과 오너테이션을 가질 수 있다.
- 키-값 쌍으로 레이블과 비슷하지만 식별정보가 없고, 레이블보다 더 큰 값을 기록할 때 많이쓴다.
    - 자동으로 추가되는 경우도 있고, 수동으로도 가능하다
- 어노테이션의 기록 예제는 다음과 같다.
  - 빌드, 릴리스, 또는 타임스탬프, 릴리스 ID, git branch, PR 번호, 이미지 해시 및 레지스트리 주소와 같은 이미지 정보
  - 로깅, 모니터링 분석 또는 감사 리포지터리에 대한 포인터
  - 디버깅 목적으로 사용될 수 있는 클라이언트 라이브러리 또는 도구 정보
  - 다른 생태계 구성 요소의 관련 오브젝트 URL과 같은 사용자 또는 도구/시스템 출처 정보
  - 경량 롤아웃 도구 메타데이터
  - 책임자의 전화번호 또는 호출기 번호ㅡ 또는 팀 웹사이트 같은 해당 정보를 찾을 수 있는 디렉터리 진입점
  - 행동을 수정하거나 비표준 기능을 수행하기 위한 최종 사용자의 지시 사항

## 문법과 캐릭터 셋
- 어노테이션은 키/값 쌍이다.
  - 키에는 두 개의 세그먼트가 있다.
    - 선택적인 접두사와 이름이며, 슬래시(/)로 구분된다.
    - 이름 세그먼트는 필수이며, [a-z0-9A-Z] 63자 이하이며, 사이에 -,_,.이 들어갈 수 있다.
    - 접두사는 선택적이다.
    - 접두사는 DNS 서브도메인이어야 한다.
    - .으로 구분된 일련의 DNS 레이블은 253자를 넘지 않고, 뒤에 /가 붙는다.
  - 접두사가 생략되면, 어노테이션 키는 사용자에게 비공개로 간주한다.
    - kubernetes.io/ 와 k8s.io/ 접두사는 쿠버네티스 핵심 구성 요소를 위해 예약되어 있다.

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