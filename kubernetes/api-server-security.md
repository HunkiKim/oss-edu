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
  - 싱글 사인 온은 OAuth, OpenID Connect 등과 같이 사용자가 여러 애플리케이션 또는 시스템에 대해 단일 인증으로 접근할 수 있는 인증 메커니즘이다.
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
- 시크릿 볼륨으로 각 컨테이너의 파일시스템에 마운트된 /var/run/secrets/kubernetes.io/serviceaccount/token 파일의 내용을 전송해 파드를 인증할 수 있다.
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
  - 기본적으론 원하는 시크릿을 마운트 할 수 잇지만, 마운트 가능한 시크릿 목록에 잇는 시크릿만 마운트하도록 파드의 서비스어카운트를 설정할 수 있다.
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
- 주체(사람, 서비스어카운트, 또는 사용자나 서비스어카운트의 그룹일 수 있음)는 하나 이상의 롤과 연계돼 있으며 각 롤은 특정 리소스에 특정 동사를 수행할 수 있다.
- 여러 롤이 있다면 롤에서 허용하는 모든 작업을 수행할 수 있다.

## RBAC 리소스 소개
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/ae744a91-7dd0-42ee-8adf-e2907b6f8625)

- RBAC 인가 규칙은 네 개의 리소스로 구성되며 두 개의 그룹으로 분류할 수 있다.
  - 롤과 클러스터롤 : 리소스에서 수행할 수 있는 동사를 지정한다.
  - 롤바인딩과 클러스터롤바인딩 : 위의 롤을 특정 사용자, 그룹, 또는 서비스어카운트에 바인딩한다.
- 즉 롤은 수행할 수 있는 작업을 정의하고, 바인딩은 누가 이를 수행할 수 있는지 정의한다.
- 롤과 롤바인딩 그리고 클러스터롤과 클러스터롤바인딩의 차이점은 롤과 롤바인딩은 네임스페이스가 지정된 리소스이고 클러스터롤과 클러스터롤바인딩은 네임스페이스를 지정하지 않는 클러스터 수준의 리소스라는 것이다.

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/0ec58222-ad5e-4c28-b5cd-0114501d1747)
### 중간 정리

- Role: 특정 네임스페이스 내에서 사용되는 역할(Role)을 정의합니다. Role은 해당 네임스페이스의 리소스에 대한 권한을 명시적으로 지정하는 데 사용됩니다.

- RoleBinding: RoleBinding은 네임스페이스 내에서 Role을 참조하여 특정 주체(사용자, 그룹, 서비스 계정 등)에게 권한을 부여합니다. RoleBinding을 사용하여 Role에 정의된 권한을 특정 주체에게 할당합니다.

- ServiceAccount: ServiceAccount는 Kubernetes에서 실행되는 애플리케이션에 할당되는 서비스 계정을 나타냅니다. ServiceAccount는 네임스페이스 내에서 리소스에 대한 액세스 권한을 부여하는 데 사용됩니다. 각 Pod는 실행되는 동안에는 해당 네임스페이스의 기본 ServiceAccount를 가집니다.

- ServiceAccountToken: ServiceAccountToken은 서비스 계정 인증을 위해 사용되는 토큰입니다. Kubernetes 클러스터 내에서 각 서비스 계정은 고유한 토큰을 가지며, 이 토큰을 사용하여 인증 및 권한 부여를 수행합니다.

- 요약하면, Role은 권한을 정의하고, RoleBinding은 Role의 권한을 특정 주체에게 할당합니다. ServiceAccount는 애플리케이션에 할당되는 서비스 계정을 나타내며, ServiceAccountToken은 서비스 계정 인증을 위한 토큰입니다.


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
- RoleBinding : dev01 ServiceAccount에게 "default" namespace의 pod-reader role의 권한을 준다.
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

- 서비스 어카운트 토큰을 생성을 먼저하자.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: dev01-token
  annotations:
    kubernetes.io/service-account.name: dev01
type: kubernetes.io/service-account-token
```
```shell
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
```
- 이제 롤과 롤바인딩을 생성한다.
```shell

$ k apply -f role.yaml

$ k apply -f role-binding.yaml
```

- 이제 아래의 명령어를 통해 사용자 생성을 진행한다.
```shell
$ k get secret dev01-token -o jsonpath='{.data.token}' | base64 --decode
eyJhbGciOiJSUzI1NiIsImtpZCI6IkJKQS1fbDNJaGpyREh1cFo1d0ZKRVU3bTNMTU9ZOUd4R29yNUxQMUxfVzgifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRldjAxLXRva2VuIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRldjAxIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiYzkzNWQxZGQtZTFiNS00NTk4LTgwY2ItYzlhNDg5NjNiODZiIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGV2MDEifQ.HIGGuZP9NrL4kRLdM7WfJ9FYFr5f-cSYkG-CDfR9VgpB4iCb3T_oWxNJBlLvUahDJBhHYal8sIn_FpnYDNFIJsELClc31HX_A3wKNfWIADSzD91kElGvZtTjkxiBGkVWR9jYqGu3GVhmJnXFlWIFy3hYWFNXHesNMrjs5myZhuCXVfUHIFbtiKNDCZTTiUNLgrJB8_Q4P1gYizmLflOJFSyvAC0cDzKqXUz7gh-YmY_eDTTg3rQU-H4Wx7Rghx13EtmldS0QdSjcqghDysU86nkwBoTcJR3URy2KlngljDeA51ctiwL5PeYPmeQEl_kjL9ykO3nzebalrJ71QEZ2gA%

$ k config set-credentials hunki --token={decoded-token}
User "hunki" set.

$ k config get-clusters
NAME
rancher-desktop

$ k config set-context test-context --cluster=rancher-desktop --user=hunki
Context "test-context" created.

$ k config use-context test-context
Switched to context "test-context".
# 참고 : 삭제 명령은 kubectl config delete-context [context명]
```

- 위 까지 하면 해당 콘텍스트의 유저권한 테스트는 잘되지만 서비스 어카운트는 파드에 권한을 주기위해 주로 사용한다. 따라서 Pod에 서비스 어카운트를 주는걸 실습해보자.
- 유저같은 경우는 외부에서 인증서 형식을 사용한다.


## 클러스터롤과 클러스터롤바인딩 사용하기
- 롤과 롤바인딩은 네임스페이스가 지정된 리소스로, 하나의 네임스페이스상에 상주하며 해당 네임스페이스의 리소스에 적용된다는 것을 의미하지만 롤바인딩은 다른 네임스페이스의 서비스어카운트도 참조할 수 있다.
- 하지만 매번 네임스페이스가 생길때마다 네임스페이스에 같은 권한에 대해 Role, RoleBinding을 만들어 서비스 어카운트에 권하을 부여하는것은 귀찮은 일이다.
- 클러스터롤과 클러스터롤바인딩은 네임스페이스를 지정하지 않는 클러스터 수준의 RBAC 리소스이다. 따라서 모든 네임스페이스에서 사용이 가능하다.
- 또한 네임스페이스를 지정하지 않는 리소스가 있다. 여기에 룰을 적용하기 위해선 클러스터수준의 RBAC가 필요하다.
  - PV, Node, Namespace 등이 있다.
  - 그리고 API 서버는 리소스를 나타내지 않는 일부 URL 경로를 노출한다. (/healthz) 일반적인 롤로는 이런 리소스나 리소스가 아닌 URL에 관한 엑세스 권한을 부여할 수 없지만 클러스터롤은 가능하다.

## 클러스터 수준 리소스에 엑세스 허용
- 파드가 pv를 나열할 수 있는 권한을 이용해보자.
  - pv는 클러스터수준의 리소스이기 때문에 적합하다.
```shell
$ k create clusterrole pv-reader --verb=get,list --resource=persistentvolumes
```

이제 해당 context로 가서 get pv를 해보자.
```
$ k get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
pvc-fd5153fb-f5ad-40ec-9882-d4c5322730f1   1Gi        RWO            Delete           Bound    default/mysql-pvc   local-path              45h
```
- 권한을 삭제하면?


```
$ k get pv
Error from server (Forbidden): persistentvolumes is forbidden: User "system:serviceaccount:default:dev01" cannot list resource "persistentvolumes" in API group "" at the cluster scope
```

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/369150b2-a728-432c-b165-cd92f0ec24fe)
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/60ecb45c-3b44-400c-840d-cdbfa8acc3ac)


## 리소스가 아닌 URL에 엑세스 허용하기
- API서버는 리소스가 아닌 URL도 노출한다.
```yaml
$ k get clusterrole system:discovery -o yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
...
rules:
- nonResourceURLs:
  - /api
  - /api/*
  - /apis
  - /apis/*
  - /healthz
  - /livez
  - /openapi
  - /openapi/*
  - /readyz
  - /version
  - /version/
  verbs:
  - get
```

- 해당 클러스터롤이 get http method의 url들을 참조한다. 즉 권한을 준다.
- URL의 클러스터롤은 클러스터롤바인딩으로 바인딩돼야 한다.

```yaml
$ k get clusterrolebinding system:discovery -o yaml

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  creationTimestamp: "2023-05-22T08:20:05Z"
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: system:discovery
  resourceVersion: "135"
  uid: ed3fe50b-9893-4158-91a4-ef17fac7fb85
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:discovery
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:authenticated
```

- 인증된 사용자만 subjects에 있어 인증된 사용자만 클러스터롤에 나열된 URL에 엑세스 할 수 있다.

## 특정 네임스페이스의 리소스에 액세스 권한을 부여하기 위해 클러스터롤 사용하기
- 클러스터롤이 항상 클러스터 수준 클러스터롤바인딩과 바인딩될 필요 없다. 네임스페이스를 갖는 일반적인 롤 바인딩과 바인딩될 수도 있다.
- view를 보면 많인 규칙이 있지만 ConfigMaps, Endpoints, PersistentVolumeClaim 등과 같은 리소스를 가져오고(get) 나열(list)하고 볼 수 있게 허용된다.
```yaml
$ k get clusterrole view -o yaml

aggregationRule:
  clusterRoleSelectors:
  - matchLabels:
      rbac.authorization.k8s.io/aggregate-to-view: "true"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  creationTimestamp: "2023-05-22T08:20:05Z"
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
  name: view
  resourceVersion: "474"
  uid: 0dbecb6e-eca9-4e42-9c14-2cf8eae45c3b
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - persistentvolumeclaims
  - persistentvolumeclaims/status
  - pods
  - replicationcontrollers
  ...
```

- 위의 리소스들은 네임스페이스가 지정된 리소스이지만, 현재 보는건 클러스터롤이다.
- 클러스터롤은 클러스터롤바인딩과 롤바인딩 중 어디에 바인딩되느냐에 따라 다르다.
  - 클러스터롤바인딩을 생성하고 클러스터롤을 참조하면, 바인딩에 나열된 주체는 모든 네임스페이스에 있는 지정된 리소스를 볼 수 있다.
  - 반면 롤바인딩을 만들면 바인딩에 나열된 주체가 롤바인딩의 네임스페이스에 있는 리소스만 볼 수 있다.


- 이제 아무 권한이 없는 유저에 대해 클러스터롤바인딩을 부여하면 어떻게 되는지 보려한다.
```shell
$ k create clusterrolebinding view-test --clusterrole=view --serviceaccount=default:dev01
```
그러면 아래와 같은 결과를 얻는다.
```
$ k get pods
NAME                                READY   STATUS    RESTARTS   AGE
ubuntu-deployment-f78c99cc7-rf5xs   1/1     Running   0          171m
ubuntu-deployment-f78c99cc7-zxmc5   1/1     Running   0          171m
ubuntu-deployment-f78c99cc7-4xf2f   1/1     Running   0          171m
```

- default 뿐만 아니라 모든 네임스페이스의 파드를 나열할 수 있다.

```shell
$ k get pods -A

NAMESPACE     NAME                                      READY   STATUS      RESTARTS   AGE
kube-system   local-path-provisioner-79f67d76f8-8mlfb   1/1     Running     0          3d20h
kube-system   coredns-597584b69b-sfstb                  1/1     Running     0          3d20h
kube-system   metrics-server-5f9f776df5-c89lx           1/1     Running     0          3d20h
kube-system   helm-install-traefik-crd-frcb2            0/1     Completed   0          3d20h
kube-system   helm-install-traefik-pgxvf                0/1     Completed   0          3d20h
kube-system   svclb-traefik-574b90f0-l5lrf              2/2     Running     0          3d20h
kube-system   svclb-traefik-574b90f0-jtcwl              2/2     Running     0          3d20h
kube-system   svclb-traefik-574b90f0-6wkkj              2/2     Running     0          3d20h
kube-system   svclb-traefik-574b90f0-v9zcf              2/2     Running     0          3d20h
kube-system   traefik-66c46d954f-gwkb2                  1/1     Running     0          3d20h
default       ubuntu-deployment-f78c99cc7-rf5xs         1/1     Running     0          172m
default       ubuntu-deployment-f78c99cc7-zxmc5         1/1     Running     0          172m
default       ubuntu-deployment-f78c99cc7-4xf2f         1/1     Running     0          172m
```

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/90953ed8-c103-48eb-952d-0863bdd1aa61)

- 위의 사진처럼 특정 네임스페이스의 서비스 어카운트가 클러스터롤수준의 롤을 적용시켜 다른 네임스페이스의 파드까지 볼 수 있게 되었다.

이제 클러스터 롤 바인딩을 롤 바인딩으로 교체해보자.

```shell
$ k delete clusterrolebindings.rbac.authorization.k8s.io view-test
clusterrolebinding.rbac.authorization.k8s.io "view-test" deleted

$ k create rolebinding view-test --clusterrole=view --serviceaccount=default:dev01
rolebinding.rbac.authorization.k8s.io/view-test created
```

```shell
$ k get pods -n default
NAME                                READY   STATUS    RESTARTS   AGE
ubuntu-deployment-f78c99cc7-rf5xs   1/1     Running   0          3h2m
ubuntu-deployment-f78c99cc7-zxmc5   1/1     Running   0          3h2m
ubuntu-deployment-f78c99cc7-4xf2f   1/1     Running   0          3h2m
```

```shell
$ k get pods -n foo
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:default:dev01" cannot list resource "pods" in API group "" in the namespace "foo"

$ k get pods -A
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:default:dev01" cannot list resource "pods" in API group "" at the cluster scope
```

롤바인딩에 속한 , 즉 서비스어카운트가 속한 네임스페이스에서의 리소스 권한만 가지고 있다. 아래 이미지를 참고하자.

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/7004b328-e812-4795-9d1c-0f9cd2c435f3)

## 롤, 클러스터롤, 롤바인딩과 클러스터롤바인딩 조합의 관한 요약
요약하자면 아래와 같다.
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/628d72a8-736f-4ddf-9962-547bb81480a9)

# 디폴트 클러스터롤과 클러스터롤바인딩의 이해
- 쿠버네티스는 API 서버가 시작될 때마다 업데이트되는 클러스터롤과 클러스터롤바인딩의 디폴트 세트를 제공한다.
  - 이렇게 하면 실수로 삭제하거나 최신 버전의 쿠버네티스가 클러스터 롤과 바인딩을 다르게 설정해 사용하더라도 모든 디폴트 롤과 바인딩을 다시 생성되게 한다.
  - *k get clusterrolebindings*, *k get clusterrole* 명령어를 통해 확인 가능하다.
- 가장 중요한 롤은 view, edit, admin 그리고 cluster-admin 클러스터롤이다. 이들은 사용자 정의 파드에서 사용하는 서비스어카운트에 바인딩되기 위한 것이다.

## view 클러스터롤을 사용해 리소스에 읽기 전용 엑세스 허용하기
- 이전 예제에서 이미 default view clusterrole을 사용했다. 롤, 롤바인딩, 시크릿을 제외한 네임스페이스 내의 거의 모든 리소스를 읽을 수 있다. 
  - 시크릿은 view보다 더 큰 권한을 갖는 인증 토큰이 포함될 수 있으며, 사용자가 다른 사용자로 가정해 추가 권한(권한 에스컬레이션)을 얻을 수 있기 때문이다.

## edit 클러스터롤을 사용해 리소스에 변경 허용하기
- 다음으로 edit 클러스터롤은 네임스페이스 내의 리소스를 수정할 수 있을 뿐만 아니라 시크릿을 읽고 수정할 수도 있다. 
  - 그러나 롤 또는 롤바인딩을 보거나 수정하는 것은 허용되지 않는다. 이것 또한 권한 상승을 방지하기 위한 것이다.

## admin 클러스터롤을 사용해 네임스페이스에 제어 권한 허용하기
- 네임스페이스에 완전한 제어 권한이 admin 클러스터롤에 부여된다.
- 이 클러스터롤을 가진 주체는 리소스쿼터와 네임스페이스 리소스 자체를 제외한 네임스페이스 내의 모든 리소스를 읽고 수정할 수 있다.
- edit과 admin의 주요 차이는 롤과 롤바인딩을 수정할 수 있다는 점이다.
> 리소스 쿼터(Resource Quota)는 Kubernetes에서 네임스페이스별 리소스 사용량을 제한하는 정책입니다. 

## cluster-admin 클러스터롤을 사용해 완전한 제어 허용하기 
- 쿠버네티스 클러스터를 완전하게 제어하려면 cluster-admin 클러스터롤을 주체에 할당하면 된다. 
- 롤바인딩을 생성해 할당해주면 생성된 네임스페이스의 모든 측면을 완전하게 제어 가능하다.
  - 모든 네임스페이스를 와전하게 제어하려면 클러스터롤바인딩에서 cluster-admin롤을 참조하면 된다.

## 그 밖의 디폴트클러스터롤 이해하기
- 디폴트 클러스터롤 목록에는 접두사 system:으로 시작하는 클러스터롤이 있다.
  - 이들은 다양한 쿠버네티스의 구성요소에서 사용된다.
  - 그중 스케줄러에서 사용되는 system:kube-scheduler와 Kubelet에서 사용되는 system:node 등이 있다.
- 이외에도 컨트롤러 매니저 등과 같은 시스템 컴포넌트들의 클러스터롤들이 존재한다.

# 인가 권한을 현명하게 부여하기
- 기본적으로 네임 스페이스의 디폴트 서비스어카운트에는 인증되지 않은 사용자의 권한 이외에는 어떤 권한도 없다. 
- 따라서 기본적으로 파드는 클러스터 상태조차 볼 수 없다. 이에 대한 권한을 부여하는 것은 사용자의 몫이다.
- 하지만 보안이 그렇듯, 자신에 일에 꼭 필요한 권한만 부여해 한 가지 이상의 권한을 주지 않는 것이 가장 좋다.(최소 권한 법칙)

## 각 파드에 특정 서비스어카운트 생성
- 각 파드(복제본 세트)를 위한 특정 서비스어카운트를 생성한 다음 롤바인딩으로 맞춤형 롤(또는 클러스터롤)과 연계하는 것이 바람직한 접근 방법이다.
- 초반부에서 처럼 파드를 읽는 서비스어카운트, 수정하는 서비스 어카운트를 만들고, 각 파드 스펙의 serviceAccountName 속성에 이 서비스어카운트를 사용하는것이 예시이다.
  - 양쪽 파드에 필요한 모든 권한을 네임스페이스의 디폴트 서비스어카운트에 추가하지말자.

