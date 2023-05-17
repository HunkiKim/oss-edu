# 파드를 안정적으로 유지하기
- 쿠버네티스에서 파드는 배포 가능한 기본 단위
- Pod의 상위 리소스가 없다면 파드를 수동으로 생성, 감독, 관리를 해야한다.
- k8s의 주요 이점 중 하나는 컨테이너 목록을 제공하면 클러스터에서 계속 실행상태를 유지해준다는 것이다.
    - 원하는 상태로 말이다.
- 파드내의 컨테이너가 죽으면 kubelet이 이를 확인하고 컨테이너를 실행하도록 할 것이다.
    - 특별한 작업을 하지 않아도 쿠버네티스에서 애플리케이션을 실행하는 것만으로 자동으로 치유를 하게된다.
- 하지만 Appliccation이 메모리가 가득차고, OutOfMemoryErrors를 발생시켜도, Container는 종료되지 않을 때가 잇다. 이 때에도 Application을 다시 시작하면 좋을것이다.
    - 오류를 캐치해서 프로세스를 종료 할 수 있지만 문제가 될 수 있다.
- 이를 해결하기 위해선 내부 기능에 의존하지 않고, 외부에서 애플리케이션의 상태를 체크해야 한다.
    - 라이브니스 프로브는 컨테이너가 살아있는지 확인한다.
        - GET 프로브, TCP 소켓 프로브,  Exec 프로브가 잇다.
            - 임의의 Get, TCP, 명령을 실행해 확인하는 방식이다.
        - manifest에 작성 가능하다.        
        ```yaml
        apiVersion: v1
        kind: Pod
        metadata:
        name: liveness-pod
        spec:
        containers:
        - name: my-app
            image: nginx
            livenessProbe:
            httpGet:
                path: /health
                port: 80
            initialDelaySeconds: 5
            periodSeconds: 10
        ```
        ```
        Liveness:       http-get http://:80/health delay=5s timeout=1s period=10s #success=1 #failure=3
        ```
        ```
          Type     Reason     Age                     From               Message
            ----     ------     ----                    ----               -------
            Normal   Scheduled  13m                     default-scheduler  Successfully assigned default/liveness-pod to lima-rancher-desktop
            Normal   Pulled     13m                     kubelet            Successfully pulled image "nginx" in 5.560608544s (5.560637336s including waiting)
            Normal   Pulled     12m                     kubelet            Successfully pulled image "nginx" in 5.492635503s (5.492659795s including waiting)
            Normal   Pulled     12m                     kubelet            Successfully pulled image "nginx" in 5.493649335s (5.493696919s including waiting)
            Normal   Created    12m (x3 over 13m)       kubelet            Created container my-app
            Normal   Started    12m (x3 over 13m)       kubelet            Started container my-app
            Warning  Unhealthy  11m (x9 over 13m)       kubelet            Liveness probe failed: HTTP probe failed with statuscode: 404
            Normal   Killing    11m (x3 over 13m)       kubelet            Container my-app failed liveness probe, will be restarted
            Normal   Pulling    11m (x4 over 13m)       kubelet            Pulling image "nginx"
            Warning  BackOff    3m41s (x26 over 9m41s)  kubelet            Back-off restarting failed container
            ☁  k8s 실습 [replicationcontroller-replicaset-daemonset] ⚡  k get pods
            NAME            READY   STATUS             RESTARTS     AGE
            annotated-pod   1/1     Running            0            44m
            liveness-pod    0/1     CrashLoopBackOff   7 (5m ago)   14m
        ```
        - 실제로 컨테이너는 올라갔지만, /health 응답이 없어 재시작 되는걸 볼 수 있다.
          - 하지만 3번이상 실패해 현재는 종료되었다.
        - 만약 pod의 Restart가 있을 때 describe po를 통해 확인해보면 재시작한 정보들을 확인할 수 있다.
        - 코드 137은 특별한 의미를 가진다.
            - 외부에서 종료됐음을 나타낸다.
        - 이외에도 추가적인 속성을 정의할 수 있다.
        - health를 만들어 내부 기능들을 확인하는 것도 방법이다.
            - 인증 확인. 인증있으면 무한정 컨테이너 실행할 수 있다.
    - 재시도 루프를 구현하지말자
    - 프로브는 자주 시작되기 때문에 가볍게 유지하자

# 레플리케이션컨트롤러

- 현재는 잘 쓰지 않으며, Deployment를 주로 쓴다.
- ReplicationController는 파드가 항상 실행되도록 보장
    - 사라진 파드 감지해 교체 파드를 생성
- 복제본을 통해 파드를 관리
    - 파드 템플릿을 통해 복제본을 만듬
- 레이블 셀렉터로 파드를 찾고, 매치되는 파드 수와 의도되는 파드 수를 찾아 다르면 삭제하거나 추가함
    - 이 때문에 특정 파드의 레이블을 변경하면 찾지 못해 새롭게 생성하고, 레이블이 변경된 파드는 좀비파드가 된다.
- 템플릿 변경이 가능하지만 다른 리소스를 볼때 보자.
- —cascade==false로 컨트롤러를 삭제하면 컨트롤러가 생성한 파드와 상관 없이 컨트롤러만 삭제한다.
    - 컨트롤러 변경시 유용

# 레플리카셋

- 초기엔 ReplicationController가 파드의 상위 리소스이며, 유일한 해결책이었다.
- 이후 ReplicaSet이라는 유사한 리소스가 도입됐다.
    - 차세대 ReplicationController라고 하며, 완전히 대체 될 것이라고 하는데 된것같다.
- 기능은 거의 동일하다
- 하지만 일반적으로 레플리카셋을 직접 생성하지 않고, 상위 수준의 디플로이먼트 리소스를 생성할 때 자동으로 생성되게 한다.

## 레플리카셋 vs 레플리케이션컨트롤러 비교

- 레플리카셋은 좀 더 풍부한 표현식을 사용하는 파드 셀렉터 보유
    - 레플리케이션컨트롤러의 레이블 셀렉터는 특정 레이블이 있는 파드만 매칭
      - matchLabels
    - 반면 레플리카셋의 셀렉터는 특정 레이블이 없는 파드나 특정 레이블 키를 갖는 파드를 매칭시킬 수 있다(값상관x).
      - env=*
    - 하지만 레플리카셋은 하나의 레플리카셋으로 두 파드 세트를 모두 매칭시켜 하나의 그룹으로 취급 가능
      - matchExpression 도입

## 레플리카셋 정의

- 상위 리소스가 없는 파드는 조건에 맞으면 레플리카셋의 소유가 된다.
- 레플리카셋은 apiVersion이 v1이 아니기 때문에 다른걸 선택해야 한다.
    - 예제에선 apps/v1beta2를 썼다.
    - Selector는 matchLabels Selector를 사용했으며, 레플리케이션컨트롤러와 유사하다.
        - 가장 단순한 방법이다
        - 기존 ReplicationController는 selector 바로 아래에 정의하지만 ReplicaSet은 selector.matchLabels에 정의한다.
    - 템플릿은 동일하다.
    
    ```json
    apiVersion: apps/v1beta2
    kind: ReplicaSet
    metadata:
    	name: kubia
    spec:
    	replicas: 3
    	selector: 
    		matchLabels:
    			app: kubia
    	template:
    			metadata:
    				labels:
    					app: kubia
    			spec:
    				containers:
    				- name: kubia
    					image: khk9346/kubia
    ```
    
- 레플리카셋의 조회는 rc가 아니라 rs이다.
    
    ```json
    kubectl get rs
    ```
    
- 위처럼 만들면 사실 ReplicaController와 똑같다. matchLabels 셀럭터는 ReplicaController와 똑같다.
    
    ```json
    apiVersion: apps/v1beta2
    kind: ReplicaSet
    metadata:
    	name: kubia
    spec:
    	replicas: 3
    	selector: 
    		matchExpressions: # 파드의 키가 app인 레이블을 포함하라는 의미
    		- key: app
    			operator: In
    			values:
    				- kubia
    	template:
    			metadata:
    				labels:
    					app: kubia
    			spec:
    				containers:
    				- name: kubia
    					image: khk9346/kubia
    ```
    
    - 위처럼 셀렉터에 표현식을 추가해 키, 연산자, 가능한 값이 포함되어야 한다.
    - 아래는 많이 쓰는 네 가지 유효 연산자이다.
        - In : 레이블의 값이 지정된 값 중 하나와 일치해야함
        - NotIn : 레이블의 값이 지정된 값과 일치하지 않아야함
        - Exists : 파드는 지정된 키를 가진 레이블이 포함
        - DoesNotExist : 파드는 지정된 키를 가진 레이블이 포함돼 잇지 않아야 함
    - matchLabels와 matchExpressions를 모두 지정할 수 있다.
        - 모든 레이블이 일치하고, 모든 표현식도 true인것맞 찾는다.

# 데몬셋

- 레플리카셋과 레플리케이션컨트롤러는 쿠버네티스 클러스터 내 어딘가에 지정된 수만큼의 파드를 실행하는 데 사용된다.
- 그러나 클러스터의 모든 노드에, 노드당 하나만 실행되길 원하는 경우가 있을 수 있다.
    - 시스템 수준의 작업을 수행하는 인프라 관련 파드가 그 예시다.
    - kube-proxy 프로세스도 좋은 예시다.
    - 모든 노드에서 실행되어야 한다.
- 레플리카셋은 클러스터 전체에서 무작위로 파드를 분산시키지만, 데몬셋은 각 노드에서 하나의 파드 복제본만 실행한다.

## 데몬셋으로 모든 노드에 파드 실행하기

![image](https://user-images.githubusercontent.com/66348135/238251746-bbc85f45-1b11-4f78-afbc-24282bea55f2.png)
- 모든 클러스터 노드마다 파드를 하나만 데실행시키려면 데몬셋 오브젝트를 생성하면 된다.
    - 사실 파드가 타깃 노드가 이미 있고, 스케줄러를 건너뛰는 것을 제외하면 기존의 레플리케이션컨트롤러 또는 레플리카셋과 매우 유사하다.
- 레플리카셋(컨트롤러)는 원하는 수의 파드 복제본이 존재하는지 확인하지만, 데몬셋은 복제본 수라는 개념이 없다.
    - 파드 셀렉터와 일치하는 파드 하나가 실행중인지 확인하는게 데몬셋의 역할이기 때문이다.
- 노드가 다운되어도 데몬셋은 다른 곳에 파드를 생성하지 않지만, 새 노드가 추가되면 해당 노드에 파드를 생성한다.
- 이전과 마찬가지로 파드 템플릿으로 파드를 생성한다.

## 데몬셋을 사용해 특정 노드에서만 파드를 실행하기

- 파드가 노드의 일부에서만 실행하려면 노드셀렉터를 이용해 특정 노드에만 파드를 배포할 수 있다.
- 데몬셋 YAML 정의를 해보자.

```json
apiVersion: apps/v1beta2
kind: DaemonSet
metadata:
  name: ssd-monitor
spec:
  selector: # matchLabels에 해당하는 파드를 관리함
    matchLabels:
      app: ssd-monitor
    template:
      metadata:
        labels:
          app: ssd-monitor
      spec:
        nodeSelector:
          disk: ssd
        containers:
        - name: main
        image: khk9346/ssd-monitor
```

- matchLabels
- selector
- ReplicaSet과 DaemonSet이 같은 matchLabels로 파드를 관리할 경우 어떻게 되는거지?
  - 아래처럼 정의하고 실행시켜보면 잘 실행된다.
  - 서로 다른 리소스 유형이므로 서로 간섭하거나 충돌하는 것은 아니다.
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: my-daemonset
spec:
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-container
        image: nginx
        command: ["sleep", "3600"]
```
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: my-replicaset
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-container
        image: nginx
```
```
my-replicaset-l5l2z   1/1     Running            0                3m23s
my-replicaset-z4kvr   1/1     Running            0                3m23s
my-replicaset-8jcdp   1/1     Running            0                3m23s
my-daemonset-7bv22    1/1     Running            0                43s
```
- 데몬셋은 약자가 ds이다.
- 만약 노드에 라벨을 추가하지 않았으면 아래를 통해 해보자.

```json
k label node [node-name | minikube] disk=ssd
```

- 제거는 아래와 같다.

```json
k label node [node-name] disk=hdd --overwrite
```

- 위처럼 label을 --overwrite하면 기존의 pod의 라벨이 바뀐다.
  - 아래처럼 daemonset의 소속이 더이상 아니게 된다.
  - 그러면 daemonset은 다시 node에 해당 매치 라벨이 없기 때문에 파드를 하나 더 생성한다.
  - 이후 daemonset을 삭제하면, --overwrite한 파드는 그대로 남아있고, 새로 생긴 파드가 삭제된다.
```
$ k label pod my-daemonset-7bv22 app=no-app --overwrite
pod/my-daemonset-7bv22 labeled
```
```
$ k delete ds my-daemonset
daemonset.apps "my-daemonset" deleted
```
```
$ k get pods
my-daemonset-7bv22    1/1     Running            0                7m19s
my-daemonset-pdrzp    1/1     Terminating        0                2m10s
```