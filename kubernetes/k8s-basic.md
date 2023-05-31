# Kubernetes 넓게 바라보기

- 컨테이너화 된 애플리케이션들을 쉽게 관리하고 운영할 수 있게 해주는 소프트웨어 시스템
- 개발자가 매니페스트(앱 디스크립터)를 마스터에 게시하면, 쿠버네티스가 워커 노드 클러스터에 배포
    - 개발자는 이 앱이 어떤 노드에 배포되는지는 중요하지 않다. 이는 설정값에 맞춰 자동적으로 된다.
- 쿠버네티스는 OS와 비슷한 구조로 형성이 된다.

# 쿠버네티스 클러스터 아키텍처 이해
![아키텍처이미지](https://user-images.githubusercontent.com/66348135/238534163-c6ace61b-50f1-4d78-88a6-2d6c31b713cc.png)
> 출처 : https://kubernetes.io/docs/concepts/architecture/
- 여러 노드들이 있지만 크게 마스터 노드와 워커 노드 두 가지 유형으로 나뉜다.
    - 마스터노드는 클러스터의 전체 쿠버네티스 시스템을 관리하고, 쿠버네티스 컨트롤 플레인을 실행한다
    - 워커노드는 실제 애플리케이션을 실행한다

# 마스터 노드
- 일반적으로 고가용성을 위해 최소 3개 이상의 마스터 노드이용한다.
  - 마스터 노드가 한 개로 운영될 경우 마스터 노드가 죽게되면 파드 스케쥴링 및 API 접근이 불가하기 때문에 안정적인 클러스터 운영이 불가능하다.
  - 노드가 복구되지 않으면 etcd 정보가 사라질 위험이 있지만 위와 운영하면 클러스터를 안정적으로 운영할 수 있게 된다.
- 여러개의 마스터 노드가 동작중일 경우, 로드밸런서를 통해 마스터노드에 접근한다.
  - 이는 트래픽 관리에도 유용하기 때문에 고가용성에 도움이 됩니다.
## 컨트롤 플레인

- 클러스터를 제어하고 작동시킨다.
- 구성요소론 다음과 같다.
    - API Server
        - 사실상 중앙에서 요청을 처리하는 역할
        - Kubernetes Cluster의 어느 요청이든 API Server를 통해 이뤄진다.
        - 시스템 컴포넌트들은 오직 API 서버와 통신하며, 모든 오브젝트가 사실상 API Server에서 이루어진다고 봐야합니다.
        - 모든 요청이 있기 때문에 Authentication, Authorization, Admission Control 또한 API Server에서 이뤄집니다.
            - 인증, 인가는 말 그대로 사용자 인증, 권한 확인의 과정이며 Admission Control은 Validation, 값 자동 생성, 기본값 추가와 같은 작업을 합니다.
    - Scheduler
        - Manifest의 리소스와 kubelet을 통해 어떤 노드로 파드를 배치할지 정하는 역할을 합니다.
        - 즉 애플리케이션을 실제로 배포합니다.
        - 최적의 노드 선택에는 두 가지 방법이 있습니다.
            - 필터링
            - 스코어링
    - etcd (distributed key-value store)
        - 분산 key-value 저장소로, k8s에선 사실상 모든 정보들이 저장되는 곳입니다.
        - 즉 클러스터의 모든 정보를 가지고 있는 곳입니다.
        - RAFT 알고리즘을 사용한다고 합니다.
            - RSM이라고도 불립니다.
              - 똑같은 데이터를 여러 서버에 계속하여 복제하는 것이며, 이 방법으로 사용하는 머신을 RSM이라고 한다.
              - 데이터를 복제한다고 모든 문제가 해결되는 것은 아니며, 오히려 복제를 하며 문제가 생길 수 있다. 따라서 컨센서스를 확보해야 하며 이는 아래의 4가지 속성을 만족한다는 것을 의미한다.
                - Safety : 항상 올바른 결과 리턴
                - Available : 서버 몇 대가 다운되어도 항상 올바른 결과 리턴
                - Independent from timing : 네트워크 지연이 발생해도 일관성 유지
                - Reactivity : 모든 서버에 복제되지 않아도 조건 만족시 빠르게 요청 응답
            - 리더,팔로워,후보자 세 가지 역할로 노드를 구성하고 리더 선출, 로그 복제 등의 과정을 통해 분산시스템의 일관성을 유지한다.
              - 리더 선출
                - 분산된 노드들 간의 election 과정을 통해 이루어진다.
                - 이런 리더는 클러스터의 상태를 변화시키고, 클라이언트 요청을 처리하고, 로그를 복제하는 역할을 한다.
                - 일반적으로 리더가 쓰기작업, 팔로워들은 리더에 쓰인 내용을 복제하고, 읽기작업을 한다. (무조건은 아니지만 일반적으로 이런식으로, master-slave와 유사하게 쓰인다.)
                - 리더가 문제가 생긴다면 팔로워들중 후보를 선출해 투표를 하고, 다시 리더를 정하게 된다.
              - 로그 복제
                - 리더를 중심으로 로그를 복제해 분산 시스템 일관성을 유지한다.
                - 일반적으로 기록은 리더가, 이를 팔로워들에게 복제한다.
                - 복제된 로그는 팔로워들도 동일한 순서로 저장해 일관성을 유지한다.
                - 팔로워들로 부터 충분한 응답을 받으면 로그를 커밋하고, 커밋된 로그는 상태변화로 적용된다.
        - API를 통해 쿠버네티스에서 사용되는 모든 리소스들을 저장합니다.
    - Controller Manager
        - 컨트롤러 오브젝트들이 존재합니다.
            - Controller Manager들이 관리하는 컨트롤러 오브젝트들은 Desired State를 가집니다. spec field가 바로 그것이며, 이 상태를 기반으로 모니터링하고, 현재 상태가 Desired State와 일치하지 않으면 최대한 가깝게 만듭니다.
            - Deployment, ReplicaSet, DamemonSet, Job 등의 리소스 타입을 관리합니다.

# 워커 노드

- 워커 노드는 컨테이너화된 애플리케이션을 실행하는 시스템입니다.
- 애플리케이션을 실행하고 모니터링합니다.
- 또한 다음 구성요소에 의해 애플리케이션에 서비스를 제공합니다.
    - 컨테이너 런타임
    - kubelet 컴포넌트
    - kube-proxy 컴포넌트

## Kubelet

- 워커 노드에 있는 컴포넌트이며, 클러스터의 각 노드마다 실행되는 에이전트입니다.
    - 여기서 에이전트는 노드가 다른 시스템, 네트워크와 상호 작용하는 작업을 수행해주는 소프트웨어 프로그램이라는 뜻입니다.
    - 즉 다른 노드나 API서버와 상호작용하는 컴포넌트입니다.
- 파드에서 컨테이너가 확실하게 동작하도록 관리
    - Pod의 목록을 주기적으로 가져와 Pod안에 있는 컨테이너들을 확인
    - 각 컨테이너 상태를 주기적으로 감시, 필요하면 컨테이너를 다시 시작하거나 중지
    - API 서버로 해당 정보도 저장
- 네트워크와 스토리지 또한 구성하고 관리
- kublet도 liveness probe를 실행하는 컴포넌트, 다운이 되면 Pod의 상태를 감시 못해 클러스터 전체의 안정성에 영향을 미친다.

## Kube-proxy

![kube-proxy](https://user-images.githubusercontent.com/66348135/238748068-f79af140-48c2-4b02-b87b-0eb0d5deb490.png)
> 출처 : https://kubernetes.io/ko/docs/concepts/services-networking/service/

- kube-proxy는 노드의 네트워크 규칙을 관리하며, 이 규칙을 통해 바깥에서 파드로 네트워크 통신이 가능
- 컨트롤 플레인의 서비스 및 엔드포인트 오브젝트의 추가와 제거를 감시
- usernamespace, iptables, IPVS를 통해 외부에서 Pod을 찾을 수 있습니다.
    - 각자 다른 특징과 장단점이 있습니다.
    - 보안적인 이유로 usernamespace 방식은 잘 사용되지 않습니다.
    - iptables 방식은 kube-proxy가 iptables를 사용해 가상 ip주소를 물리적인 pod ip 주소로 매핑하는 방식. 가장 많이 사용되며 빠르게 구성 가능
    - IPVS 방식은 Linux 커널에서 제공하는 IPVS를 사용해 로드밸런싱을 수행하는 방식. 더 빠르고 성능 좋고, 다양한 로드밸런싱 알고리즘이 있어 성능이 중요하면 더 적합
- 즉 순서는 이렇습니다.
    - 서비스의 Cluster IP 주소로 요청
    - kube-proxy는 etcd에 저장된 엔드포인트 정보를 기반으로 로드밸런싱을 결정하여 api 서버가 kube-proxy와 같은 클라이언트에게 제공
    - kube-proxy는 클라이언트의 요청을 적절한 백엔드 Pod로 라우팅 (iptables 룰이나 IPVS를 설저해야 이루어짐)
      - 서비스는 같은 레이블 셀렉터를 가진 Pod를 자동으로 발견하고 이를 엔드포인트로 사용
      - iptables룰을 통해, 즉 ClusterIP나 NodePort, LoadBalancer IP를 통해 클라이언트는 Pod에 접근

# 쿠버네티스 특징
- 서비스 디스커버리
  - 서비스를 통해 클러스터 내부의 애플리케이션 통신을 간단하게 함
- 로드밸런싱
  - 파드에 들어오는 여러 트래픽을 분산시킴
- 스토리지 오케스트레이션
  - 스토리지 또한 유기적으로 늘리고 줄일 수 있음
- 자동화된 롤아웃,롤백
  - 배포시 자동으로 롤링업데이트, 롤백 지원
- 자동화된 빈 패킹
  - 여러 파드를 적절한 노드에 자동으로 배치하여 자원활용 극대화
- Self-Healing
  - 애플리케이션, 노드 상태 지속적감시를 통해 문제 발생시 자동으로 복구
- Auto-Scailing
  - 부하 상황에 따라 자동으로 노드 크기 조절
- Job
  - 배치작업 간편화
- Secret과 구성 관리
  - Secret과 Config Map을 통해 보안정보, 설정정보를 안전하게 관리