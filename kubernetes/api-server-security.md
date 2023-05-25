# Access to the Kubernetes API
- 쿠버네티스엔 여러 접근 방법이 있다.
  - kubectl을 통해 쿠버네티스 클러스터에 직접적으로 접근하는 방법
  - UserAccount가 Kubernetes API Server에 접근하는 방법
  - ServiceAccount가 Pod로 접근하는 방법

## API 서버의 기능
- k8s api server는 다른 모든 구성요소와 kubectl 같은 클라이언트에서 사용하는 중심 구성 요소다.
- 클러스터 상태를 조회하고 변경하기 위해 RESTful API로 CRUD 인터페이스를 제공.
  - 상태는 etcd안에 저장
- 오브젝트를 etcd에 저장하는 일관된 방법을 제공하는 것뿐만 아니라, 오브젝트 유효성 검사 잠업도 수행하기 때문에 잘못된 오브젝트를 저장할 수 없다.
- 유효성 검사와 함께 낙관적 잠금도 처리하기 때문에 동시에 업데이터가 발생하더라도 다른 클라이언트에 의해 변경 사항이 재정의되지 않는다.
- API 서버의 클라이언트 중에 하나는 책의 시작 부분에서 사용한 kubectl 명령줄 도구다.
  - 예를들면 JSON 파일에서 리소스를 생성할 때 kubectl은 파일의 내용을 API 서버에 HTTP POST 요청으로 전달한다.

아래는 API 서버가 요청을 받을 때 내부에서 일어나는 일을 보여준다.
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/e6f62a23-f443-40aa-973e-444c99bdc957)

### 인증 플러그인으로 클라이언트 인증
- API서버는 먼저 요청을 보낸 클라이언트를 인증해야 한다.
  - API 서버는 누가 요청을 보낸 것인지 결정할 수 있을 때 까지 이 플러그인을 차례로 호출한다.
- 인증 방법에 따라 클라이언트 인증서 혹은 HTTP 헤더에서 가져온다.
- 플러그인은 클라이언트의 사용자 이름, 사용자 ID, 속해 있는 그룹 정보를 추출한다.
  - 이 데이터는 다음 단계인 인가 단계에서 사용된다.

### 인가 플러그인을 통한 클라이언트 인가
- API 서버는 인증 플러그인 외에도 하나 이상의 인가 플러그인을 사용하도록 설정돼 있다.
- 인가 플러그인을 통해 사용자가 해당 동작을 할 수 있는지 판단하고, 가능하다면 다음 단계로 넘어간다.

### 어드미션 컨트롤 플러그인으로 요청된 리소스 확인과 수정
- 리소스를 생성, 수정, 삭제하려는 요청인 경우에 해당 요청은 어드미션 컨트롤(Admission Control)로 보내진다.
- 서버는 여러 어드미션 컨트롤 플러그인을 사용하도록 설정돼 있다.
- 이 플러그인은 리소스를 여러 이유로 수정할 수 있다.
  - 리소스 정의에서 누락된 필드를 기본값으로 초기화하거나 재정의할 수 있다.
  - 요청에 없는 관계된 리소스를 수정하거나 어떤 이유로든 요청을 거부할 수도 있다.
  - 해당 리소스는 모든 어드미션 컨트롤 플러그인을 통과한다.

> 데이터를 읽는 요청은 어드미션 컨트롤을 거치지 않는다.

아래는 어드미션 컨트롤 플러그인의 예시이다.
- AlwaysPullImages: 파드의 imagePullPolicy를 Always로 변경해 파드가 배포될때마다 이미지를 항상 강제로 가져오도록 재정의한다.
- ServiceAccount: 명시적으로 지정하지 않을 경우 default 서비스 어카운트를 적용한다.
- NamespaceLifecycle: 삭제되는 과정에 잇는 네임스페이스와 존재하지 않는 네임스페이스 안에 파드가 생성되는 것을 방지한다.
- ResourceQuota: 특정 네임스페이스 안에 있는 파드가 해당 네임스페이스에 할당된 CPU와 메모리만을 사용하도록 강제한다.

이외에도 더 많은 어드미션 컨트롤러는 https://kubernetes.io/docs/admin/admission-controllers/ 에서 확인 가능하다.

### 리소스 유효성 검사
- 요청이 모든 어드미션 컨트롤 플러그인을 통과하면, API 서버는 오브젝트의 유효성을 검증하고 etcd에 저장한다. 그리고 클라이언트에 응답을 반환한다.

# 인증 이해
- API서버가 요청을 받으면 인증 플러그인 목록을 거치면서 요청이 전달되고, 각각의 인증 플러그인이 요청을 검사해 보낸 사람이 누구인가를 밝혀내려 시도한다.
  - 요청에서 해당 정보를 처음 추출해낸 플러그인은 사용자 이름, 사용자 ID와 클라이언트가 속한 그룹을 API 서버 코어에 반환한다.
  - API 서버는 나머지 인증 플러그인의 호출을 중지하고, 계속해서 인가 단계를 진행한다.
- API서버는 아래의 인증 플러그인들을 사용해 클라이언트의 아이덴티티를 얻는다.
  - 클라이언트 인증서
  - HTTP 헤더로 전달된 인증 토큰
  - 기본 HTTP 인증
  - 기타
- 인증 플러그인은 API 서버를 시작할 때 명령행 옵션을 통해 활성화 가능하다.

## 사용자와 그룹
- 인증 플러그인은 인증된 사용자의 사용자 이름과 그룹을 반환한다.
- 쿠버네티스는 해당 정보를 어디에도 저장하지 않는다. 
- 이를 사용해 사용자가 작업을 수행할 권한이 있는지 여부를 확인한다.

### 사용자
쿠버네티스는 API 서버에 접속하는 두 종류의 클라이언트를 구분한다.
- 실제 사람(사용자)
- 파드(더 구체적으로, 파드 내부에서 실행되는 애플리케이션)

이 두 가지 유형의 클라이언트는 모두 위에서 언급한 인증 플러그인을 사용해 인증된다.

- 사용자는 싱글 사인 온(SSO, Single Sign On)과 같은 외부 시스템에 의해 관리돼야 하지만 파드는 서비스 어카운트(service account)라는 메커니즘을 사용하며, 클러스터에 서비스 어카운트 리소스로 생성되고 저장된다.
  - 사용자 계정을 나태내는 자원은 없고, API 서버를 통해 사용자를 관리할 수 없다는 의미이다.

서비스 어카운트는 파드를 실행하는 데 필수적이므로 먼저 공부한다.

### 그룹
- 휴먼 사용자와 서비스어카운트는 하나 이상의 그룹에 속할 수 있다.
- 인증 플러그인이 반환하는 그룹은 임의의 그룹 이름을 나타내는 문자열일 뿐이지만, built-in 그룹은 특별한 의미를 갖는다.
  - system: unauthenticated 그룹은 어떤 인증 플러그인에서도 클라이언트를 인증할 수 없는 요청에 사용된다.
  - system: authenticated 그룹은 성공적으로 인증된 사용자에게 자동으로 할당된다.
  - system: serviceaccounts 그룹은 시스템의 모든 서비스어카운트를 포함한다.
  - system: serviceaccounts:<namespace>는 특정 네임스페이스의 모든 서비스어카운트를 포함한다.

## 서비스어카운트 소개
- 클라이언트는 API 서버에서 작업을 수행하기 전에 자신을 인증해야 한다.
- 시크릿 볼륨으로 각 컨테이너의 파일시스템에 마운트된 /var/run/secrets/kubernetes.io/serviceaccount/token 파일의 내용을 전송해 파다를 인증할 수 있다.
  - 모든 파드는 파드에서 실행 중인 애플리케이션의 아이덴티티를 나타내는 서비스어카운트와 연계돼있다.
  - 이 토큰 파일은 서비스어카운트의 인증 토큰을 갖고 있다.
  - 애플리케이션이 이 토큰을 사용해 API 서버에 접속하면 인증 플러그인이 서비스어카운트를 인증하고 서비스어카운트의 사용자 이름을 API 서버 코어로 전달한다.
  - 이후 인가 플러그인이 수행여부를 판단한다.


## 서비스 어카운트 리소스
서비스 어카운트는 파드, 시크릿, 컨피그맵 등과 같은 리소스이며 개별 네임스페이스로 범위가 지정된다. 각 네임스페이스마다 default 서비스어카운트가 자동으로 생성된다. (그동안 파드가 사용한 서비스어카운트)
- 아래의 명령어로 서비스어카운트 나열이 가능하다.
```shell
$ k get sa
NAME      SECRETS   AGE
default   0         2d20h
```
- 같은 네임스페이스의 서비스어카운트만 사용할 수 있다.
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/c09f9222-9f78-46a8-b469-201662ed2419)

## 서비스어카운트가 인가와 어떻게 밀접하게 연계돼 있는지 이해하기
- 파드 매니페스트에 서비스어카운트의 이름을 지정해 파드에 서비스어카운트를 할당할 수 있다.
  - 안하면 default 할당
- 파드에 서로 다른 서비스어카운트를 할당하면 각 파드가 액세스할 수 있는 리소스를 제어 가능하다.
- API 서버가 인증 토큰이 있는 요청을 수신하면, API 서버는 토큰을 사용해 요청을 보낸 클라이언트를 인증한 다음 관련 서비스어카운트가 요청된 작업을 수행할 수 있는지 여부를 결정한다.
- 인가 플러그인에서 수행하며, 그중 하나는 RBAC(역할 기반 엑세스 제어)플러그인이다.
  - RBAC는 대부분의 클러스터에서 사용해야 하는 플러그인이다.

## 서비스어카운트 생성
- 서비스어카운트를 생성하는 이유는 클러스터 보안 때문이다.
  - 클러스터의 메타데이터를 읽을 필요가 없는 파드는 클러스터에 배포된 리소스를 검색하거나 수정할 수 없는 제한된 계정으로 실행해야 한다.
  - 검색,수정 등 필요한 권한만 주어야 한다.

```shell
$ k create serviceaccount foo
serviceaccount/foo created
```

어떤 정보가 있는지 살펴보자.
```yaml
Name:                foo
Namespace:           default
Labels:              <none>
Annotations:         <none> 
Image pull secrets:  <none> # 이 서비스어카운트를 사용하는 파드에 이 필드의 값이 자동으로 추가된다.
Mountable secrets:   <none> # 마운트 가능한 시크릿이 강제화된 경우 이 서비스 어카운트를 사용하는 파드만 해당 시크릿을 마운트할 수 있다.
Tokens:              <none> # 인증 토큰, 첫 번째 토큰이 컨테이너에 마운트된다.
Events:              <none>
```

### 서비스어카운트의 마운트 가능한 시크릿 이해
- k describe를 사용해 서비스어카운트를 검사하면 토큰이 Mountable secrets 목록에 표시된다.
  - 기본적으론 원하는 시크릿을 마운트 할 수 잇지만, 마운트 가능한 시크릿 목록에 잇는 시크릿만 마운트하드록 파드의 서비스어카운트를 설정할 수 있다.
  - 이 기능을 사용하기 위해선 *kubernetes.io/enforce-mountable-secrets="trun ".* 어노테이션을 포함해야한다.

### 서비스어카운트의 이미지 풀 시크릿 이해
- 컨테이너 이미지를 가져오는 데 필요한 자격증명을 가지고 있는 시크릿이다.
- 예를들면 아래와 갖다.
```shell
$ k get secrets
NAME                TYPE                             DATA   AGE
fortune-https       Opaque                           1      2d5h
mydockerhubsecret   kubernetes.io/dockerconfigjson   1      2d4h
```

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account
imagePullSecrets:
  - name: mydockerhubsecret
```

- 서비스어카운트의 이미지 풀 시크릿은 마운트 가능한 시크릿과 다르게 동작한다.
  - 마운트 가능한 시크릿과 달리 각각의 파드가 어떤 이미지 풀 시크릿을 사용할 수 있는지 결정하는 것이 아니라 서비스어카운트를 사용해 모든 파드에 특정 이미지 풀 시크릿을 자동으로 추가해준다.
- 즉 서비스어카운트에 이미지 풀 시크릿을 추가하면 각 파드에 개별적으로 추가할 필요 없다.

## 파드에 서비스어카운트 할당
- 추가 서비스어카운트를 만든 후엔 이를 파드에 할당해야 한다.
  - spec.serviceAccountName 필드에서 서비스어카운트 이름을 설정하면 된다.

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: local-path
  resources:
    requests:
      storage: 1Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: khk9346/private-repo:v1
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: "1234"
          ports:
            - containerPort: 3306
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql
      volumes:
        - name: mysql-data
          persistentVolumeClaim:
            claimName: mysql-pvc
```
- 예를들어 서비스 어카운트 없이 docker private repository의 이미지를 접근해보자.
```
$ k get pods
NAME                                READY   STATUS         RESTARTS   AGE
mysql-deployment-6878f4c9df-9npvc   0/1     ErrImagePull   0          10s
mysql-deployment-6878f4c9df-8cpmp   0/1     ErrImagePull   0          10s
mysql-deployment-6878f4c9df-56ffn   0/1     ErrImagePull   0          10s
```
이제 아까 만든 서비스어카운트를 설정하자.
```
spec:
    serviceAccountName: my-service-account
    containers:
    - name: mysql
```
이제 확인해보면 아래와 같다.
```shell
$ k describe sa my-service-account
Name:                my-service-account
Namespace:           default
Labels:              <none>
Annotations:         <none>
Image pull secrets:  mydockerhubsecret
Mountable secrets:   <none>
Tokens:              <none>
Events:              <none>

$ k get pods
NAME                               READY   STATUS    RESTARTS   AGE
mysql-deployment-cc6d96fc7-pw7ph   1/1     Running   0          72s
mysql-deployment-cc6d96fc7-qv6qw   1/1     Running   0          55s
mysql-deployment-cc6d96fc7-npqc6   1/1     Rk unning   0          54s
```

# 역할 기반 엑세스 제어로 클러스터 보안 (RBAC Cluster Security)
- 쿠버네티스 버전 1.6.0부터 클러스터 보안이 크게 강화됐다.
- 1.8.0부턴 RBAC 인가 플러그인이 GA로 승격됐고, 이제 많은 클러스터에서 기본적으로 활성화돼 잇다.
- RBAC는 권한이 없는 사용자가 클러스터 상태를 보거나 수정하지 못하게 한다.
- RBAC외에도 속성 기반 액세스 제어(ABAC, Attribute-Based Access Control) 플러그인, 웹훅(Web Hook) 플러그인, 사용자 정의 플러그인 구현과 같은 여러 인가 플러그인이 포함돼 있다. 하지만 RBAC가 표준이다.

## RBAC 인가 플러그인 소개
- API 서버가 REST 인터페이스를 통해 사용자의 HTTP 요청 액션을 인가 플러그인을 통해 점검 가능하다.

### 액션 이해하기
- REST 클라이언트는 GET, POST, PUT, DELETE 및 기타 유형의 HTTP 요청을 특정 REST 리소스를 나타내는 특정 URL 경로로 보낸다.
- 쿠버네티스에서는 이러한 리소스에는 파드, 서비스, 시크릿 등이 있다.
  - 파드 가져오기 (Get)
  - 서비스 생성하기 (Create)
  - 시크릿 업데이트 (Update)
  - 기타 등등

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/746155f4-33e3-43aa-a7be-5dd49723c1c0)

- API 서버 내에서 실행되는 RBAC와 같은 인가 플러그인은 클라이언트가 요청한 자원에서 요청한 동사를 수행할 수 있는 지 판별한다.
  - 명사는 쿠버네티스 리소스에 정확하게 매핑되며, 동사는 클라이언트가 수행한 HTTP 메서드에 매핑된다.
- 전체 리소스 유형에 보안 권한을 적용하는 것 외에도 RBAC 규칙은 특정 리소스 인스턴스에도 적용할 수 있다.
  - 그리고 URL 경로에도 권한을 적용할 수 있다. 이는 API 서버가 노출하는 모든 경로가 리소스를 매핑하는 것은 아니기 때문이다.
    - /api 경로나 서버의 상태 정보를 갖는 /healthz


### RBAC 플러그인 이해
- 이름처럼 RBAC 인가 플러그인은 사용자가 액션을 수행할 수 있는지 여부를 결정하는 핵심 요소로 사용자 롤(user role)을 사용한다.
- 주체(사람, 서비스어카운트, 또는 사용자나 서비스어카운트의 그룹일 수 ㅐㅔ있음)는 하나 이상의 롤과 연계돼 있으며 각 롤은 특정 리소스에 특정 동사를 수행할 수 있다.
- 여러 롤이 있다면 롤에서 허용하는 모든 작업을 수행할 수 있다.

## RBAC 리소스 소개
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/ae744a91-7dd0-42ee-8adf-e2907b6f8625)

- RBAC 인가 규칙은 네 개의 리소스로 구성되며 두 개의 그룹으로 분류할 수 있다.
  - 롤과 클러스터롤 : 리소스에서 수행할 수 있는 동사를 지정한다.
  - 롤바인딩과 클러스터롤바인딩 : 위의 롤을 특정 사용자, 그룹, 또는 서비스어카운트에 바인딩한다.
- 즉 롤은 수행할 수 있는 작업을 정의하고, 바인딩은 누가 이를 수행할 수 있는지 정의한다.
- 롤과 롤바인딩 그리고 클러스터롤과 클러스터롤바인딩의 차이점은 롤과 롤바인딩은 네임스페이스가 지정된 리소스이고 클러스터롤과 클러스터롤바인딩은 네임스페이스를 지정하지 않는 클러스터 수준의 리소스라는 것이다.

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/0ec58222-ad5e-4c28-b5cd-0114501d1747)

### 실습을 위한 설정
- RBAC 리소스가 API 서버로 수행할 수 있는 작업에 어떤 영향을 주는지 살펴보기 전에 클러스터에 RBAC가 활성화돼 있는지 확인해야 한다.
- 버전 1.6 이상에, 인가 플러그인으로 RBAC 플러그인만 구성돼 있는지 확인해야 한다.
- 여러 개의 플러그인이 병렬로 활성화돼 있을 수 있으며 그중 하나라도 액션을 수행하도록 허용하는 경우 액션이 허용된다.

### 네임스페이스 생성과 파드 실행
- Role : pod에 읽기 권한을 부여해주는 default namespace의 role
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""] # core API 그룹
  resources: ["pods"] # 리소스
  verbs: ["get", "list", "watch"] # 허용되는 동작

```
- RoleBinding : jane에게 "default" namespace의 pod-reader role의 권한을 준다. 이는 jane이 "default" namespace의 파드들을 읽는 권한을 준다.
  - 또 추가적으로 ClusterRole을 참조할 수 있다.
```yaml
apiVersion: rbac.authorization.k8s.io/v1
# This role binding allows "jane" to read pods in the "default" namespace.
# You need to already have a Role named "pod-reader" in that namespace.
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
subjects:
# You can specify more than one "subject"
- kind: ServiceAccount
  name: dev01
  apiGroup: ""
roleRef:
# "roleRef" specifies the binding to a Role / ClusterRole
  kind: Role #this must be Role or ClusterRole
  name: pod-reader # this must match the name of the Role or ClusterRole you wish to bind to
  apiGroup: rbac.authorization.k8s.io
```
- ClusterRole : 모든 네임스페이스의 secret의 읽기 권한이 주어진다.
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClutserRole
metadata:
  name: secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
```
- ClusterRoleBinding : manager 그룹에게 어떤 네임스페이스안에 있는 시크릿을 읽을 수 있는 권한을 줍니다.
```yaml
apiVersion: rbac.authorization.k8s.io/v1
# This cluster role binding allows anyone in the "manager" group to read secrets in any namespace.
kind: ClusterRoleBinding
metadata:
  name: read-secrets-global
subjects:
- kind: Group
  name: manager # Name is case sensitive
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

이제 Role을 추가하면 아래와 같이 나온다.

서비스 어카운트에 롤을 바인딩

우선 서비스 어카운트 토큰을 생성해야 한다.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: dev01-token
  annotations:
    kubernetes.io/service-account.name: dev01
type: kubernetes.io/service-account-token
```
$ k apply -f secret-serviceaccount.yaml
secret/dev01-token created

$ k describe secret dev01-token
Name:         dev01-token
Namespace:    default
Labels:       <none>
Annotations:  kubernetes.io/service-account.name: dev01
              kubernetes.io/service-account.uid: 5da8dfba-20db-45a5-adbb-e85fc2790db5

Type:  kubernetes.io/service-account-token

Data
====
ca.crt:     570 bytes
namespace:  7 bytes
token:      eyJhbGciOiJSUzI1NiIsImtpZCI6IjNmWUdpcTNDdHQ3akszVVpDVFpxTnlxaldERkJKU3pGbUFFTmJMbDZCRFkifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRldjAxLXRva2VuIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRldjAxIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiNWRhOGRmYmEtMjBkYi00NWE1LWFkYmItZTg1ZmMyNzkwZGI1Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGV2MDEifQ.TK2iUWT0l1lDG5H8rqUidvQwAn_073hcpPTkbr79PD0kaCHCrwZ9EwN8F-QKV9hmSD7la0b2s83JSET696-lUOGRTbZ4qH5ZFUxt80nKfHnxMaR_7gypPGefd7dlpBWpxYH3Wc1dv_tD9R3H8-CYKPpagptd0-vXHyU7mQbR2MR9gpD00Rmm01atL4ySco7dmWW-ciaqQTG4qpWQFzTr8s_WfWBTGq8BDzItGSv8Gn_JbVRptfdgRmRudgaRwya71PK8iIzu-4956d00W0r2MwSTSb2e_6lX87cRx5PfWWf9SZoaZYpD47uNYimqFqcFNoerS3G-Xl-Ix_lhomiopg
```shell


$ openssl rand -base64 32
PiQKy8D4XX4cNu/apsKcWfRoJH9cU+63+nntDsqhIIE=

$ k apply -f role.yaml

$ k apply -f role-binding.yaml
```

- 이제 아래의 명령어를 통해 사용자 생성을 진행한다.
```shell
$ k config set-credentials hunki --token=dev01
User "hunki" set.

$ k config get-clusters
NAME
rancher-desktop

$ k config set-context test-context --cluster=rancher-desktop --user=hunki
Context "test-context" created.

$ k apply -f role-binding.yaml
rolebinding.rbac.authorization.k8s.io/read-pods created

$ k config use-context test-context
Switched to context "test-context".
# 참고 : 삭제 명령은 kubectl config delete-context [context명]
```