# Volume
- 컨테이너 내의 파일들은 임시적이다.
  - 컨테이너 문제생기면 파일들이 손상된다.
- 컨테이너간 데이터 공유가 힘들다.
- 위의 이유 외에도 컨테이너의 데이터를 관리하기 위한 볼륨이라는 개념이 등장했다.
## Docker에서의 Volume
- 다소 느슨하고 덜 관리되는 볼륨이라는 개념이 있다.
- 다른 컨테이너에 있거나, 디스크에 디렉토리이다.
- 볼륨 드라이버를 제공하지만, 제한적이다.
- 기본적으로 Docker는 다음과 같은 볼륨 드라이버를 제공합니다:

  - local: 로컬 호스트 파일 시스템에 마운트되는 볼륨을 생성합니다.
  - nfs: NFS(Network File System)를 사용하여 네트워크 상의 공유 디렉터리를 마운트하는 볼륨을 생성합니다.
  - bind: 호스트의 특정 디렉터리를 컨테이너에 마운트하는 볼륨을 생성합니다.
  - tmpfs: 컨테이너의 임시 파일 시스템을 생성하는 볼륨을 생성합니다.
- 그러나 이러한 기본 볼륨 드라이버는 일부 제한 사항이 있습니다. 
  - 예를 들어, 볼륨 드라이버 간 데이터 공유, 데이터 복제 또는 분산 저장소와의 통합과 같은 고급 기능을 지원하지 않습니다. 

## k8s의 Volume
- 파드내의 컨테이너들의 각 경로를 볼륨에 마운트 한다면 컨테이너간 데이터 공유가 가능하다.
- 볼륨을 채우거나 마운트하는 프로세스는 파드의 컨테이너가 시작되기 전에 수행된다.
- 볼륨이 파드의 라이프사이클에 바인딩되면 파드 존재까진 유지도리 수 있지만 유형에 따라 파드와 볼륨이 사라진 후에도 볼륨의 파일이 유지돼 새로운 볼륨으로 마운트 될 수 있다.
- 쿠버네티스는 다양한 유형의 볼륨을 지원하며 여러 볼륨 유형을 동시에 사용할 수 있다. 

## 사용 가능한 볼륨 유형 소개
- emptyDir: 일시적인 데이터를 저장하는 데 사용되는 간단한 빈 디렉터리다.
- hostPath: 워커 노드의 파일시스템을 파드의 디렉터리로 마운트하는 데 사용한다.
- gitRepo: git 레포지토리 콘텐츠를 체크아웃해 초기화한 볼륨이다.
- nfs: NFS 공유를 파드에 마운트한다.
- cloudDisk: 클라우드 제공자의 전용 스토리지를 마운트하는 데 사용한다.
- cinder, cephfs, iscsi, flocker, glusterfs, quobyte, rbd, flexVolume, vsphere Volume, photonPersistentDisk, scaleIO : 다른 유형의 네트워크 스토리지를 마운트하는 데 사용한다.
- configMap, secret, downwardAPI: 쿠버네티스 리소스나 클러스터 정보를 파드에 노출하는 데 사용되는 특별한 유형의 볼륨이다.
- persistentVolumeClaim: 사전에 혹은 동적으로 프로비저닝된 퍼시스턴트 스토리지를 사용하는 방법이다.

# 볼륨을 사용한 컨테이너 간 데이터 공유
- 단일 컨테이너에서도 볼륨은 유용하지만 하나의 파드에 있는 여러 컨테이너에서 데이터를 공유하는 방법을 한 번 보자.

## emptyDir 볼륨 사용
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/c1b3b773-47e9-4bdb-93de-68ae085da160)

- 이름처럼 볼륨이 빈 디렉터리로 시작된다.
- 볼륨의 라이프사이클이 파드에 묶여있어, 파드 삭제시 콘텐츠도 사라진다.
- emptyDir 볼륨은 동일 파드에서 실행 중인 컨테이너 간 파일을 공유할 때 유용하다.
  - 단일 컨테이너에서도 컨테이너의 파일시스템이 쓰기가 불가능한 경우, 마운트된 볼륨에 쓰는 것이 유일한 옵션일 수 있다.

### 파드에 emptyDir 볼륨 사용
- 웹 서버, 콘텐츠 에이전트, 로그 순환기가 두 개의 볼륨을 공유한다고 해보자.
- Nginx를 웹 서버로 사용하고 유닉수 fortune 명령으로 HTML 콘텐츠를 생성한다.
- Fortune 명령은 실행할 대마다 임의의 인용문을 출력한다. 매 10초마다 fortune 명령을 실행하고 출력을 index.html에 저장하는 스크립트를 생성한다.

- fortuneloop.sh
  - 실행권한 확인하자
```shell
#!/bin/bash
trap "exit" SIGINT
while :
do
    echo $(date) Writing fortune to /var/htdocs/idnex.html
    /usr/games/fortune > /var/htdocs/index.html
    sleep 10
done
```
- Dockerfile
```docker
FROM ubuntu:latest
Run apt-get update; apt-get -y install fortune
ADD fortuneloop.sh /bin/fortuneloop.sh
ENTRYPOINT /bin/fortuneloop.sh # 쉘파일이 시작되어야 함을 명시
```
- image build
```shell
$ docker build -t khk9346/fortune .
$ docker push khk9346/fortune
```
- 파드 생성
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: fortune
spec:
  containers:
  - image: khk9346/fortune # html 생성 컨테이너 
    name: html-generator
    volumeMounts:
    - name: html
      mountPath: /var/htdocs
  - image: nginx:alpine # nginx 읽기 전용 마운트
    name: web-server
    volumeMounts:
    - name: html
      mountPath: /usr/share/nginx/html
      readOnly: true
    ports:
    - containerPort: 80
      protocol: TCP
  volumes: # html 단일 emptyDir 볼륨을 위의 컨테이너 두 개에 마운트
  - name: html
    emptyDir: {}
```
- 위의 파드는 컨테이너 두 개와 각 컨테이너에 각기 다른 경로로 마운트된 단일 볼륨을 갖는다. html-generator는 10초마다 fortune명령의 결과를 쓰기 시작하고, 이는 볼륨의 /var/htdocs에 마운트 된다.
- nginx는 html 파일을 서비스하기 시작한다.
- 위의 파드는 로컬 머신의 포트를 파드로 포워딩하면 된다.
```shell
$ k port-forward fortune 8080:80
```
- 실행 결과는 아래와 같다.
```shell
k port-forward fortune 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
Handling connection for 8080
Handling connection for 8080
```
```shell
☁  fortune  curl localhost:8080
Your heart is pure, and your mind clear, and your soul devout.
☁  fortune  curl localhost:8080
Q:	"What is the burning question on the mind of every dyslexic
	existentialist?"
A:	"Is there a dog?"
```

### emptyDir을 사용하기 위한 매체 지정하기
- 볼륨으로 사용한 emptyDir은 파드를 호스팅하는 워커 노드의 실제 디스크에 생성되므로 디스크의 유형에 따라 성능이 달라진다. 
- 하지만 k8s의 emptyDir은 디스크가 아닌 메모리를 사용하는 tmpfs 파일시스템으로 생성하도록 요청 가능하다.
```yaml
volumes:
  - name: html
    emptyDir: 
      medium: memory # Volume의 저장소 유형을 선택하는 필드
```
## git repository를 볼륨으로 사용하기
- gitRepo 볼륨은 기본적으로 emptyDir 볼륨이며 파드가 시작되면 깃 리포지터리를 복제하고 특정 리비전을 체크아웃해 데이터로 채운다.
- 하지만 gitRepo에 변경을 푸시할 때마다 웹사이트의 새 버전을 서비스하기 위해선 파드를 삭제해줘야 한다.
  - 볼륨이 항상 깃 레포지터리와 동기화하도록 추가 프로세스가 가능하다.
    - 사이드카 컨테이너: 파드 내의 도움 컨테이너인 사이드카 컨테이너를 실행하여 주 컨테이너의 동작을 보완할 수 있다. 사이드카 컨테이너가 주기적으로 gitRepo를 확인해 볼륨을 마운트 할 수 있다.

# 워커 노드 파일시스템의 파일 접근
- 대부분의 파드는 호스트 노드를 인식하지 못하므로 노드의 파일시스템에 있는 어떤 파일에도 접근하면 안 된다.
- 그러나 특정 시스템 레벨의 파드(Daemonset pod)는 노드의 파일을 읽거나 파일시스템을 통해 노드 디바이스를 접근하기 위해 파일시스템을 사용해야 한다.
  - hostPath 볼륨으로 가능하다.

## hostPath 볼륨 소개
- hostPath 볼륨은 노드 파일시스템의 특정 파일이나 디렉터리를 가리킨다.
- 동일 노드에 실행 중인 파드가 hostPath 볼륨의 동일 경로를 사용 중이면 동일한 파일이 표시된다.
- PV의 처음 소개 유형이다.
  - hostPath는 파드가 종료되어도 삭제되지 않기 때문이다.
- hostPath는 노드의 파일시스템을 사용하기 때문에 신중하게 써야한다.
  - 실제 Production에서는 잘 쓰지 않으며, 파드가 어떤 노드로 갈지 모르기 때문에도 민감하기 때문에 테스트용으로 주로 쓴다고 한다.

## hostPath 볼륨을 사용하는 시스템 파드 검사하기
- hostPath를 적절하게 사용하는 방법을 보기 전에 이 볼륨 유형을 사용하는 전역파드를 확인하자.
- 아래 명령어는 클러스터 내의 kube-system 네임스페이스에 모든 파드를 가져오는 명령어이며, kube-system은 k8s 시스템 컴포넌트 및 서비스가 배포되는 기본 네임스페이스이다.
```shell
$ k get pods --namespace kube-system
```
```shell
k get pods --namespace kube-system
NAME                                      READY   STATUS      RESTARTS   AGE
local-path-provisioner-79f67d76f8-xr5hm   1/1     Running     0          3d11h
coredns-597584b69b-rl6zs                  1/1     Running     0          3d11h
metrics-server-5f9f776df5-tvs78           1/1     Running     0          3d11h
helm-install-traefik-crd-dl4h5            0/1     Completed   0          3d11h
helm-install-traefik-zgkz9                0/1     Completed   2          3d11h
svclb-traefik-95360656-x9qd5              2/2     Running     0          3d11h
svclb-traefik-95360656-79n2v              2/2     Running     0          3d11h
svclb-traefik-95360656-2nkvq              2/2     Running     0          3d11h
svclb-traefik-95360656-hfngw              2/2     Running     0          3d11h
svclb-traefik-95360656-7tkzx              2/2     Running     0          3d11h
traefik-66c46d954f-qtdnc                  1/1     Running     0          3d11h
```
- 위의 pod들의 경우엔 configmap이 주된 볼륨을 사용한다.
  - 만약 hostPath가 있더라도 자체 데이터를 저장하기 위해 사용하는 경우는 없다.
  - Minikube와 같은 단일 클러스터, 단일 노드 인 경우엔 어쩔수없이 테스트용으로 종종 사용한다.

# 퍼시스턴트 스토리지 사용
- 파드에서 실행 중인 애플리케이션이 디스크에 데이터를 유지해야 하고 파드가 다른 노드로 재스케줄링된 경우에도 동일한 데이터를 사용해야 한다면 지금까지 사용한 hostPath, emptyDir은 사용할 수 없다.
  - 어떤 클러스터의 노드들에서든 접근이 가능해야 하기 때문이다.

## GCE Persistent Disk를 파드 볼륨으로 사용하기
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/3310cb54-1601-4de2-9015-18ddebabdb6e)
- 여기서 세번째 파드인 mysql 컨테이너가 사용하는 방식이다.
- mongodb라는 이름의 GCE Persistent Disk를 생성
```shell
$ gcloud compute disks create --szie=10GiB --zone=us-central1-a mongodb
```
- gcePersistentDisk 볼륨을 사용하는 파드
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mongodb
spec:
  volumes:
    - name: mongodb-data
      gcePersistentDisk: #볼륨의 유형은 GCE Persistent Disk이다.
        pdName: mongodb 
        fsType: ext4 # 파일시스템 유형은 EXT4(리눅스 파일시스템 유형 중 하나)이다.
  containers:
  - image: mongo
    name: mongodb
    volumeMounts:
      - name: mongodb-data
        mountPath: /data/db #MongoDB가 데이터를 저장할 경로
    ports:
    - containerPort: 27017
      protocol: TCP
```
- 위는 GCE의 같은 리전일 경우만 실행되며, 다른 리전이면 실행되지 않는다.

## 기반 퍼시스턴트 스토리지로 다른 유형의 볼륨 사용하기
- 위의 GCE 퍼시스턴트 디스크 볼륨은 GCE Cluster에서만 가능하다.
  - 위의 방식일 경우엔 zone이 같아야 하기 때문에 어쩔 수 없다.
- 즉 기반 퍼시스턴트 스토리지를 각각 사용해야 한다는 것이다.
  - aws면 EBS 써야한다.
- 하지만 파드정의에 이런 스토리지를 명시하는 것은 클러스터마다 파드 정의가 달라진다는 것이다.
- 개선해보자.

# 기반 스토리지 기술과 파드 분리
- 지금까지 살펴본 모든 퍼시스턴트 볼륨 유형은 파드 개발자가 실제 네트워크 스토리지 인프라스터럭처에 관한 지식이 필요하다.
- 애플리케이션 개발자로부터 실제 인프라스터럭처를 숨기는 쿠버네티스 기본 아이디어에 위배된다.
- 개발자는 애플리케이션을 위해 일정량의 퍼시스턴트 스토리지를 필요로 하면 쿠버네티스에 요청할 수 있어야 하고, 동일한 방식으로 파드 생성 시 CPU, 메모리와 다른 리소스를 요청할 수 있어야 한다.

## PV(Persistent Volume)와 PVC(Persistent Volume Claim)
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/a3c475ac-92c4-494c-8896-ea84bcd79925)

- 인프라스터럭처에 세부 사항을 처리하지 않고 k8s 클러스터에 스토리지를 요청하기 위해 새로운 리소스인 PV와 PVC가 도입됐다.
- 여태까지 했던 Persistent Storage도 영구적으로 저장 가능하다. 용어가 비슷하니 헷갈리지 말자.
- 즉 관리자가 네트워크 스토리지를 설정하고, PV manifest를 게시해 PV를 생성한다.
- 이후에 사용자는 PVC를 생성하면, k8s가 적정 크기와 접근 모드의 PV를 찾고 PVC를 PV에 바인딩 한다.
- 사용자는 파드를 생성할 때 PVC를 참조하는 볼륨을 가진 파드를 생성한다.
- 아래는 로컬 스토리지를 이용한 PV 예시 파일이다.
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mongodb-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  local:
    path: /Users/mantech/Desktop/k8s/fortune
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - k3d-k3s-default-agent-0
```
- PV를 생성할 때 관리자는 쿠버네티스에게 용량이 얼마인지, 단일 노드나 동시에 다수 노드에 읽기나 쓰기가 가능한지 알려야 한다. 
- 또한 PV가 해제되면 어떤 동작을 해야할지 알려야 한다.
- PV는 노드와 같은 클러스터 수준 리소스이다.
  - PV는 클러스터 수준 리소스이기 때문에 특정 네임스페이스에 속하지 않는다. 

## PVC 생성을 통한 PV 요청
- PV가 필요한 파드를 배포해보자.
- 파드에 직접 사용할 수 없고 클레임을 해야한다.
- 파드가 재스케줄링 되어도 동일한 PVC가 사용 가능한 상태로 유지되기를 원하므로 PV에 대한 클레임은 파드를 생성하는 것과 별개의 프로세스다.
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mongodb-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
  storageClassName: local-storage # 현재 로컬 스토리지를 사용해서 했지만 나중에 동적프로비저닝때 사용
```
 - pvc는 생성되면 적절한 PV를 찾고 클레임에 바인딩한다.
```shell
NAME         CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                 STORAGECLASS    REASON   AGE
mongodb-pv   2Gi        RWO            Retain           Bound    default/mongodb-pvc   local-storage            35m
```
```shell
NAME          STATUS   VOLUME       CAPACITY   ACCESS MODES   STORAGECLASS    AGE
mongodb-pvc   Bound    mongodb-pv   2Gi        RWO            local-storage   90s
```

## 파드에서 PVC 사용하기
- PV는 사용중에 있다. 볼륨을 해제할 때까지 다른 사용자는 동일한 볼륨에 클레임 할 수 없다.
- PV가 아니라 Pod는 PVC를 참조한다.
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mongodb
spec:
  volumes:
    - name: mongodb-data
      persistentVolumeClaim:
        claimName: mongodb-pvc
  containers:
    - name: mongodb
      image: mongo
      volumeMounts:
        - name: mongodb-data
          mountPath: /data/db
      resources:
        limits:
          memory: "128Mi"
          cpu: "500m"
      ports:
        - containerPort: 27017
          protocol: TCP
```
```shell
Volumes:
  mongodb-data:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  mongodb-pvc
    ReadOnly:   false
```
- 추가적인 절차가 생긴건 맞지만, volume에 대해서 특별한 설정없이 클레임만 지정해주면 내부를 알 필요없이 개발자가 볼륨을 사용할 수 있게 된다.

### 정책들
- PV의 reclaim 정책: PV에는 reclaim 정책이 지정되어 있습니다. reclaim 정책에는 Retain, Delete, Recycle 등이 있습니다. Retain 정책은 PV를 삭제해도 해당 스토리지를 유지하고 보존하는 것을 의미합니다. Delete 정책은 PV를 삭제하면 해당 스토리지를 삭제합니다. Recycle 정책은 PV를 삭제한 후 해당 스토리지를 재사용하기 위해 데이터를 삭제하거나 초기화합니다. PV의 reclaim 정책에 따라 종료하기 전에 스토리지의 보존 또는 삭제 여부를 결정해야 합니다.
  - persistentVolumeReclaimPolicy: Retain 우리의 PV는 이렇게 설정되어 있기 때문에 PV가 삭제되어도 스토리지는 유지합니다.
- PV와 PVC의 연결: PV와 PVC는 일대일 매핑입니다. PV를 종료하면 해당 PV와 연결된 PVC가 있을 수 있습니다. 이 경우 PVC를 사용하는 Pod에서 스토리지에 접근할 수 없게 됩니다. 따라서 PV를 종료하기 전에 해당 PV와 연결된 PVC를 관리하고 해당 PVC를 사용하는 Pod를 수정하거나 삭제해야 합니다.
- Pod의 안정성: PV를 종료하기 전에 해당 PV에 연결된 Pod가 안정적인 상태인지 확인해야 합니다. Pod에서 스토리지를 사용하는 중인 경우 PV를 종료하면 Pod의 동작이 영향을 받을 수 있습니다. 따라서 Pod를 안정적인 상태로 이전하거나 삭제한 후에 PV를 종료해야 합니다.
- pvc와 pv가 1:1 매핑되어 있으며 pvc가 pv에 의존하고 있기 때문에 pvc를 삭제하기 전에 pv를 삭제하면 아래와 같이 삭제가 되지 않는다.
```shell
k delete pv mongodb-pv
persistentvolume "mongodb-pv" deleted

k get pv
NAME         CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS        CLAIM                 STORAGECLASS    REASON   AGE
mongodb-pv   2Gi        RWO,ROX        Retain           Terminating   default/mongodb-pvc   local-storage            117m
```
- 종료가 되지 않는다.
- 여기서 또 파드와 pvc가 바인딩 되어 있으면 pvc도 꺼지지 않는다. 파드가 pvc를 참조하기 때문이다.
- 따라서 pod->pvc->pv 순으로 종료시켜줘야 종료가 성공적으로 이루어진다.

# PV의 동적 프로비저닝
- 지금까지 본걸로 많이 편해졌지만 클러스터 관리자는 실제 스토리지를 미리 프로비저닝해둬야 개발자가 쓸 수 있다.
  - 프로비저닝은 스토리지를 동적으로 할당하고 제공하는 과정
- 하지만 이 과정까지 자동으로 처리할 수 있다.
- 퍼시스턴트볼륨을 생성하는 대신 퍼시스턴트볼륨 프로비저너를 배포하고 사용자가 선택 가능한 퍼시스턴트 볼륨의 타입을 하나 이상의 스토리지클래스 오브젝트로 정의할 수 있다.
- 사용자가 퍼시스턴트볼륨클레임에서 스토리지클래스를 참조하면 프로비저너가 퍼시스턴트 스토리지를 프로비저닝할 때 이를 처리한다.

## 스토리지클래스 리소스를 통한 사용 가능한 스토리지 유형 정의
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/5ca1f277-36df-4b19-a09f-4c0cd59b163f)

- storageclass
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
provisioner: kubernetes.io/gce-pd #PV 프로비저닝을 위해 사용되는 볼륨 플러그인
parameters:
  type: pd-ssd
  zone: us-central1
```
- pvc
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mongodb-pvc
spec:
  storageClassName: fast
  resources:
    requests:
      storage: 100Mi
  accessModes:
    - ReadWriteOnce
```

- PV를 정의하지 않는 이유는 스토리지를 요구할 때 storage 사양에 따라 자동으로 PV를 생성하고 프로비저닝하기 때문입니다.
- 실제로 실행하면 결과는 아래와 같습니다.
- pv
```shell
$ k get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                 STORAGECLASS   REASON   AGE
pvc-c4c2bd8d-66fd-478c-ad44-25536098475f   1Gi        RWO            Delete           Bound    default/mongodb-pvc   fast                    25s
```
- pvc
```shell
$ k get pvc
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
mongodb-pvc   Bound    pvc-c4c2bd8d-66fd-478c-ad44-25536098475f   1Gi        RWO            fast           2m23s
```
- storageclass
```shell
NAME                        PROVISIONER                    RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
fast                        kubernetes.io/gce-pd           Delete          Immediate              false                  3m48s
```
- 여기서 RECLAIMPOLICY에서 DELETE는 PVC 삭제시 PV도 삭제됨을 의미한다. 
- 또 주목할 점은 pd 또한 자동으로 생성해줬다는 것이다.
  - kubernetes.io/gce-pd 프로비저너 때문에 가능하다.
```shell
NAME                                      LOCATION       LOCATION_SCOPE  SIZE_GB  TYPE         STATUS
pvc-c4c2bd8d-66fd-478c-ad44-25536098475f  us-central1-a  zone            1        pd-ssd       READY
```
