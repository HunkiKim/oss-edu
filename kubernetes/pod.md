# Pod
- k8s의 가장 기본 단위
- 관련된 컨테이너들을(일반적으로 1개) 함께 그룹화하여 배포하고 실행하는 데 사용

## 컨테이너 재시작
- 컨테이너가 종료되었을 때 어떻게 처리할지 처리를 결정하는 정책입니다.
- 파드의 spec에 restartPolicy 필드가 있다.
- restartPolicy는 Pod의 모든 컨테이너에 적용된다.
- 3가지 정책이 있습니다.
  - Always: 컨테이너가 종료되면 항상 재시작
  - OnFailure: 컨테이너가 실패한 경우에만 재시작
  - Never: 컨테이너 종료시 재시작 X
- 기본값은 Alawys입니다.

## imagePullPolicy
- kubelet이 이미지를 pull 할 때 사용되는 속성이다.
- 3가지 종류가 있다.
  - IfNotPresent : 해당 이미지가 로컬에 없는 경우만 pull 작업을 한다.
  - Always : 매번 pull 작업을 합니다.
  - Never : 이미지를 pull하지 않으며, 로컬에 없으면 작업에 실패하게된다.
- 만약 따로 설정하지 않는다면 조건에 따라 기본값이 달라진다.
- 프로덕션 환경에선 :latest는 지양해야한다.
  - 이미지의 어떤 버전이 기동되고 있는지 추적이 어렵고, 롤백도 어렵다.
- 기본적으로 default도 3가지가 있다.
    - 컨테이너 이미지의 태그가 :latest이며, imagePullPolicy를 설정하지 않았다면 imagePullPolicy는 자동으로 Always로 설정된다.
    - 이미지 태그 명시하지 않고, imagePullPolicy를 설정하지 않았다면 자동으로 Always로 된다.
    - 태그가 latest가 아닌 태그가 있고, imagePullPolicy를 설정하지 않았다면 자동으로 IfNotPresent로 설정된다.
- 이미지가 없고 Never인 경우 실행 결과
```
$ k get pods
NAME     READY   STATUS              RESTARTS   AGE
my-pod   0/1     ErrImageNeverPull   0          2s
```
- 이미지가 없고 Always인 경우 실행 결과
```
$ k get pods
NAME     READY   STATUS             RESTARTS   AGE
my-pod   0/1     ImagePullBackOff   0          21s
```

## Affinity(관계성)
- nodeSelector는 파드를 특정 레이블이 있는 노드로 제한하는 가장 간단한 방법이다.
- Affinity와 Anti-Affinity 기능은 표현할 수 있는 제약 종류를 크게 확장한다.
- Pod와 Node간의 관계를 정의하는 기능입니다. 쉽게말하면 특정 노드에 파드를 스케줄링 하려면 Node Affinity, 특정 파드가 다른 파드와 같은 노드에 스케줄링 되도록 할 땐 Pod Affinity를 사용합니다.
- Affinity는 Pod의 sec 섹션에 설정됩니다.
### Node Affinity
- 노드와 관련된 Affinity를 설정합니다. 노드의 Label과 매치되는 조건을 지정해 파드를 특정 노드에 바인딩 할 수 있습니다.
- Node Affinity는 다음과 같은 세 가지 옵션을 제공합니다.
  - requiredDuringSchedulingIgnoredDuringExecution : 파드가 특정 노드와 매치되어아먄 스케줄링 됩니다. nodeSelector와 유사하지만, 좀 더 표현적인 문법을 제공합니다.
  - preferredDuringSchedulingIgnoredDuringExecution : 파드가 특정 노드와 가장 잘 매치되는 경우 스케줄링 됩니다. 해당 노드가 없어도, 스케줄러는 여전히 파드를 스케줄링합니다.
  - requiredDuringScheduingRequiredDuringExecution : 파드가 특정 노드와 매치되어야 스케줄링되고 실행 중인 동안에도 계속해서 해당 노드에 유지됩니다.
- Pod는 nodeSelectorTerms에 해당하는 노드에 스케줄링됩니다.
  - matchExpression : 매치 표현식을 사용하여 노드 레이블과 일치시킬 조건을 정의합니다. 노드는 하나 이상의 매치 표현식을 충족해야 합니다.
    - key : 매치 표현식의 키로 노드 레이블 키를 지정
    - operator : 매치 표현식에 사용할 연산자를 지정
      - 주로 In, NotIn, Exists, DoesNotExist, Gt, Lt 등이 사용됩니다.
      - NotIn과 DeosNotExist 연산자를 사용하면 노드 안티-어피니티 규칙을 정의한다고도 부른다.
    - values: 매치 표현식과 일치해야 하는 노드 레이블 값의 목록을 지정합니다.
  - matchFields : 매치 필드를 사용하여 노드의 필드와 일치시킬 조건을 정의합니다. 노드는 하나 이상의 매치 필드를 충족해야 합니다.
    - key: 매치 필드의 키로 노드의 필드를 지정합니다. 예를 들어 metadata.name, spec.providerID 등이 사용될 수 있습니다.
    - operator: 매치 필드에 사용할 연산자를 지정합니다. 
      - 주로 Equals, NotEquals, In, NotIn, Exists, DoesNotExist 등이 사용됩니다.
    - values: 매치 필드와 일치해야 하는 값의 목록을 지정합니다.
- weight는 노드 어피니티 가중치라고도 부르며 preferredDuringSchedulingIgnoredExecution 어피니티 타입 인스턴스에 대해 1-100 범위의 weight를 명시합니다.
  - 스케줄러가 다른 모든 파드 스케줄링 요구 사항을 만족하는 노드를 찾으면, 만족한 모든 선호 규칙에 대해 합계 계산을 위한 weight 값을 각각 추가합니다.
  - 스케줄러가 파드에 대한 스케줄링을 판단 할 때, 총 점수가 가장 높은 노드가 우선순위를 갖게됩니다.
> IgnoredDuringExecution는 쿠버네티스가 파드를 스케줄링한 뒤에 노드 레이블이 변경되어도 파드는 계속 해당 노드에서 실행됨을 의미한다.

### 실습
```json
apiVersion: v1
kind: Pod
metadata:
  name: with-affinity-anti-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/os
            operator: In
            values:
            - linux
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: label-1
            operator: In
            values:
            - key-1
      - weight: 50
        preference:
          matchExpressions:
          - key: label-2
            operator: In
            values:
            - key-2
  containers:
  - name: with-node-affinity
    image: registry.k8s.io/pause:2.0
```
- label-1:key-1은 agent0에 label-2:key-2는 agent1에 설정했다.
```
$ k describe pod with-affinity-anti-affinity
Name:             with-affinity-anti-affinity
Namespace:        default
Priority:         0
Service Account:  default
Node:             k3d-k3s-default-agent-1/172.20.0.4
```
```
$ k describe pod with-affinity-anti-affinity2
Name:             with-affinity-anti-affinity2
Namespace:        default
Priority:         0
Service Account:  default
Node:             k3d-k3s-default-agent-1/172.20.0.4
```
- 위의 템플릿 두 개의 pod를 설정하여 파드를 생성하면 위와 같은 결과가 나오며, weight가 더 큰 쪽으로 pod가 할당된다.

- Node Anti Affinity
  - NotIn이나 DoesNotExist를 통해 파드가 해당 노드에 할당되지 않게 할 때 사용한다.
  - 현재 두 파드 모두 agent1로 설정이 되었는데 약간의 설정 변경으로 다시 파드를 생성하면 아래와 같은 yaml파일이 나온다.
- weight가 아무리 높아도, requiredDuringSchedulingIgnoredDuringExecution을 충족시키지 못하면 해당 노드에 스케줄링 될 수 없다.
```json
apiVersion: v1
kind: Pod
metadata:
  name: with-affinity-anti-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/os
            operator: NotIn
            values:
            - linux
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: label-1
            operator: In
            values:
            - key-1
      - weight: 50
        preference:
          matchExpressions:
          - key: label-2
            operator: In
            values:
            - key-2
  containers:
  - name: with-node-affinity
    image: registry.k8s.io/pause:2.0
```
```shell
k get pods -o wide --show-labels
NAME                          READY   STATUS    RESTARTS   AGE   IP       NODE     NOMINATED NODE   READINESS GATES   LABELS
with-affinity-anti-affinity   0/1     Pending   0          8s    <none>   <none>   <none>           <none>            <none>
```
### Pod Affinity
- 파드 간의 Affinity를 설정합니다.
- 파드의 Label과 매치되는 조건을 지정하여 파드를 특정 파드와 같은 노드 또는 다른 노드에 스케줄링할 수 있습니다.
- Pod Affinity 다음과 같은 옵션을 제공합니다.
  - requiredDuringSchedulingIgnoredDuringExecution : 파드가 특정 파드와 매치되어야만 스케줄링됩니다.
  - preferredDuringSchedulingIgnoredDuringExecution : 파드가 특정 파드와 가장 잘 매치되는 경우 스케줄링됩니다.

### 실습
```
apiVersion: v1
kind: Pod
metadata:
  name: with-pod-affinity
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: label-1
            operator: In
            values:
            - key-1
          - key: label-2
            operator: In
            values:
              - key-2
        topologyKey: topology.kubernetes.io/zone
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 50
        podAffinityTerm:
          labelSelector:
            matchExpressions:
            - key: label-2
              operator: In
              values:
              - key-2
          topologyKey: topology.kubernetes.io/zone
      - weight: 1
        podAffinityTerm:
          labelSelector:
            matchExpressions:
            - key: label-1
              operator: In
              values:
              - key-1
          topologyKey: topology.kubernetes.io/zone
  containers:
  - name: with-pod-affinity
    image: registry.k8s.io/pause:2.0

```
- 위의 문서를 있는 그대로 해석해보면 파드가 label-1=key-1 인 레이블이 있는 파드의 노드를 선택하고, 그곳에 해당하는 topologyKey가 있는 zone을 선택해서 Pod를 스케줄링 한다.
- topologyKey 같은 경우엔, topology를 사용하는 이유에 대해 알면 좋은데, 기본적으로 각 존마다 균일하게 스케줄되도록 하며, 클러스터에 문제가 생길경우 스스로 치유하도록 설정하기 위해 존재한다.
  - 즉 topologyKey의 키와 동일한 값을 가진 레이블이 있는 노드는 동일한 토폴로지에 있는 것으로 간주한다.
  - 여기선 topologyKey에 기반하여 스케줄링 또는 분산 배치를 수행하는데 사용된다.
  - topology.kubernetes.io/zone는 클라우드 환경에서 노드가 속한 물리적인 존 또는 리전을 나타냅니다.
    - 즉 이 물리적인 존 또는 리전을 기반으로 분산 배치를 해준다는 의미입니다.
- podAntiAffinity같은 경우엔 label-2=key-2가 있는 레이블에 대해 가중치 50을 label-1=key-1은 1을 부여합니다.
- 다른점이 있다면 노드 어피니티는 노드의 레이블을 기반으로, 파드 어피니티는 노드에 속한 파드를 기반으로 선택됩니다.
- Anti가 붙어있으면 weight의 역순을 우선적으로 고른다.
```shell
$ k get pods -o wide --show-labels
NAME                READY   STATUS             RESTARTS     AGE    IP          NODE                      NOMINATED NODE   READINESS GATES   LABELS
pod-affinity        1/1     Running            0            2m2s   10.42.1.4   k3d-k3s-default-agent-3   <none>           <none>            label-1=key-1
pod-affinity2       1/1     Running            0            118s   10.42.4.4   k3d-k3s-default-agent-0   <none>           <none>            label-1=key-1,label-2=key-2
with-pod-affinity   0/1     CrashLoopBackOff   1 (6s ago)   9s     10.42.0.4   k3d-k3s-default-agent-1   <none>           <none>            <none>
```
- 0번 노드는 현재 weight 51, 3번 노드는 weight는 50이다.
- 처음엔 가장 weight가 낮은 1번 노드를 고른다. 다시 삭제하고 실행하면 아래와 같이 나온다.
```
$ k get pods -o wide --show-labels
NAME                READY   STATUS              RESTARTS   AGE     IP          NODE                      NOMINATED NODE   READINESS GATES   LABELS
pod-affinity        1/1     Running             0          2m26s   10.42.1.4   k3d-k3s-default-agent-3   <none>           <none>            label-1=key-1
pod-affinity2       1/1     Running             0          2m22s   10.42.4.4   k3d-k3s-default-agent-0   <none>           <none>            label-1=key-1,label-2=key-2
with-pod-affinity   0/1     ContainerCreating   0          2s      <none>      k3d-k3s-default-agent-2   <none>           <none>            <none>
```
- 이번엔 2번 노드를 고른다. 같은 우선순위인 1,2번 노드에 대해서 랜덤으로 파드가 생성되는 것을 볼 수 있다.

# 테인트(Taints)와 톨러레이션(Tolerations)
- 노드 어피니티는 특정노드에 배치하거나, 배치하고싶지 않을 때 사용했다.
  - required나 prefered를 통해
  - 노드에 파드를 끌어들이는 속성
- 반면 테인트는 특정 노드에 파드의 스케줄링을 제한하기 위해 사용
  - 노드가 파드셋을 제외시킨다.
- 톨러레이션은 파드에 적용된다.
  - 톨러레이션을 통해 스케줄러는 그와 일치하는 테인트가 있는 파드를 스케줄할 수 있다.
- 테인트와 톨러레이션은 함께 작동하여 파드가 부적절한 노드에 스케줄되지 않게 한다.

## 테인트
- 라벨과 효과(Effect)로 구성
- 노드에 조건을 설정하여 파드의 배치를 제한하는 역할을 한다.
- 일반적으로 NoSchedule, PreferNoSchedule 효과를 가지는 테인트를 사용
  - NoSchedule: 파드가 해당 테인트와 일치하는 노드에 배치되지 않도록 합니다. 즉, 파드가 해당 노드에 스케줄링되지 않습니다.
  - PreferNoSchedule: 파드가 해당 테인트와 일치하는 노드에 배치되지 않는 것을 선호하지만, 다른 조건을 충족하는 경우에는 배치될 수 있습니다.
- 쉽게 말하면 라벨에 효과를 다는 것이다.

```shell
k taint nodes k3d-k3s-default-agent-1 label-2=key-2:NoSchedule
```
```shell
  taints:
  - effect: NoSchedule
    key: label-2
    value: key-2
```

- taint를 설정하고 아래와 같은 affinity가 있는 파드를 생성해보자.
```json
apiVersion: v1
kind: Pod
metadata:
  name: with-affinity-anti-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: label-2
            operator: In
            values:
            - key-2
  containers:
  - name: with-node-affinity
    image: registry.k8s.io/pause:2.0
```
- 이론상 nodeAffinity는 label-2=key-2인 곳만 생성될 수 있지만, taints가 걸려있어 어떤 노드에도 할당되면 안된다.
  - 참고로 label-2=key-2 라벨은 agent-1에만 있다.
```shell
$k get pods -o wide --show-labels
NAME                          READY   STATUS    RESTARTS   AGE   IP       NODE     NOMINATED NODE   READINESS GATES   LABELS
with-affinity-anti-affinity   0/1     Pending   0          7s    <none>   <none>   <none>           <none>            <none>
```
- 결과가 예상대로 나오는 모습을 볼 수 있다.
- taints가 없을땐 당연히 아래와 같았다.
```shell
$ k get pods -o wide --show-labels
NAME                          READY   STATUS             RESTARTS     AGE   IP          NODE                      NOMINATED NODE   READINESS GATES   LABELS
with-affinity-anti-affinity   0/1     CrashLoopBackOff   1 (3s ago)   4s    10.42.0.7   k3d-k3s-default-agent-1   <none>
```

## 톨러레이션
- 특정 테인트(Taint)를 가진 노드에서도 파드(Pod)를 실행할 수 있도록 허용하는 설정입니다.
- 파드가 배치될 때 톨러레이션은 다음과 같은 순서로 동작합니다:
  - 파드가 스케줄러에 의해 노드에 배치되기 전에 톨러레이션 규칙을 검사합니다.
  - 파드의 tolerations 필드에 지정된 톨러레이션과 노드의 테인트를 비교합니다.
  - 톨러레이션 규칙과 노드의 테인트가 일치하는 경우, 해당 테인트를 가진 노드에도 파드를 스케줄링합니다.
- 톨러레인트는 다음과 같은 속성들이 있습니다.
  - key: 테인트의 키를 지정합니다.
  - operator: 테인트 키와 톨러레이션 규칙의 키가 일치해야 하는지 확인하는 연산자를 지정합니다. 일치하는 연산자로는 Equal, Exists, NotEqual 등이 있습니다.
    - equal: taint의 키와 값이 정확히 일치해야 해야 스케줄링 가능하다.
    - exists: key만 일치하면 된다.
    - NotEqual: taint의 키와 값이 일치하지 않는 곳에만 스케줄링한다.
  - value: 테인트의 값(Value)을 지정합니다. 톨러레이션 규칙의 값과 일치해야 합니다.
  - effect: 톨러레이션 규칙이 적용되는 효과(Effect)를 지정합니다. 일반적으로 NoSchedule을 사용하여 특정 테인트를 가진 노드에서도 파드를 스케줄링할 수 있도록 허용합니다.
```json
apiVersion: v1
kind: Pod
metadata:
  name: with-affinity-anti-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: label-2
            operator: In
            values:
            - key-2
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: label-1
            operator: In
            values:
            - key-1
      - weight: 50
        preference:
          matchExpressions:
          - key: label-2
            operator: In
            values:
            - key-2
  tolerations:
  - key: "label-2"
    operator: "Equal"
    value: "key-2"
    effect: "NoSchedule"
  containers:
  - name: with-node-affinity
    image: registry.k8s.io/pause:2.0
```
- 위를 실행하면 아까 할당되지 못한 agent-1에 파드가 배치되어야 한다.
```shell
$ k get pods -o wide --show-labels
NAME                          READY   STATUS             RESTARTS     AGE   IP          NODE                      NOMINATED NODE   READINESS GATES   LABELS
with-affinity-anti-affinity   0/1     CrashLoopBackOff   1 (2s ago)   2s    10.42.0.8   k3d-k3s-default-agent-1   <none>           <none>            <none>
```
- 예상대로 되는 모습을 볼 수 있다.

## Pod Network
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/24ad72ea-6904-4ec5-9e47-263a81f58c90)
- 컨테이너는 서로 완전 격리, 하지만 파드안에 모든 컨테이너가 자체 네임스페이스가 아니라 동일한 리눅스 네임스페이스를 공유하도록 도커를 설정
    - 네임스페이스는 리눅스 커널의 자원을 격리하는 기능
        - IPC, ID, PID 등
    - 동일한 네임스페이스를 공유해 Pod 안의 컨테이너들 서로 간 통신이 필요할 때, [localhost](http://localhost)(루프백)로 통신 가능
        - 하지만 각각의 컨테이너들은 자체 PID, USER, MNT 네임스페이스를 가져 내부에선 독립적으로 동작 가능
        - 즉 Pod 내부는 하나의 네트워크로 묶입니다.


# YAML 또는 JSON 디스크립터로 파드 생성

- Object를 만들때 [https://kubernetes.io/docs/reference/](https://kubernetes.io/docs/reference/를) 를 참고하라고 한다.
    - 하지만 레거시이다.
- 이미 생성된 pod같은 경우는 아래와 같은 명령어로 yaml 정보를 읽을 수 있다.
    
    ```go
    k get po(d) [pod-name] -o yaml
    ```
    
    <aside>
    💡 po 같은 경우는 pod의 약어이며, -o는 출력 포맷을 나타낸다.
    
    </aside>
    
- pod는 크게 metadata, spec, status가 중요하다.
    - metadata: 이름, 네임스페이스, 레이블 및 파드에 관한 기타 정보를 포함한다.
    - spec: 파드 컨테이너, 볼륨, 기타 데이터 등 파드 자체에 관한 실제 명세를 가진다.
    - status: 파드 상태, 각 컨테이너 설명과 상태, 파드 내부 IP, 기타 기본 정보 등 현재 실행 중인 파드에 관한 현재 정보를 포함한다.
- 이를 토대로 간단한 YAML을 작성해보자.
    
    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
    name: kubia-manual
    spec:
    containers:
    - name: kubia-container
        image: khk9346/kubia
        ports:
        - containerPort: 8080
        protocol: TCP
    ```
    
- status가 없는 이유는 아직 생성되지 않았기 때문입니다.
  - API Server에서 Pending상태로 처음 생성됩니다.
    - 정확히는 etcd에 정보를 저장하고, pod manifest 기반으로 스케줄링을 수행하여 kubelet에게 전달합니다.
  - 이후 Kubelet에서 스케줄링된 Pod manifest를 받아 노드에서 실행하며, 이 시점에 Status가 다시 ContainerCreating으로 변경됩니다.
  - 성공적으로 컨테이너가 실행되면 Running으로 업데이트합니다.
- status 이외에도 다양한 yaml 정보들이 생성됩니다. 
  ```shell
  k get pod kubia-manual -o yaml
  ```
  ```yaml
  apiVersion: v1
    kind: Pod
    metadata:
    annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{},"name":"kubia-manual","namespace":"default"},"spec":{"containers":[{"image":"khk9346/kubia","name":"kubia-container","ports":[{"containerPort":8080,"protocol":"TCP"}]}]}}
    creationTimestamp: "2023-05-17T12:40:04Z"
    name: kubia-manual
    namespace: default
    resourceVersion: "782"
    uid: a6c0deb5-e668-4982-9971-f3bbc6ea24d5
    spec:
    containers:
    - image: khk9346/kubia
        imagePullPolicy: Always
        name: kubia-container
        ports:
        - containerPort: 8080
        protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
        name: kube-api-access-c88r4
        readOnly: true
    dnsPolicy: ClusterFirst
    enableServiceLinks: true
    nodeName: lima-rancher-desktop
    preemptionPolicy: PreemptLowerPriority
    priority: 0
    restartPolicy: Always
    schedulerName: default-scheduler
    securityContext: {}
    serviceAccount: default
    serviceAccountName: default
    terminationGracePeriodSeconds: 30
    tolerations:
    - effect: NoExecute
        key: node.kubernetes.io/not-ready
        operator: Exists
        tolerationSeconds: 300
    - effect: NoExecute
        key: node.kubernetes.io/unreachable
        operator: Exists
        tolerationSeconds: 300
    volumes:
    - name: kube-api-access-c88r4
        projected:
        defaultMode: 420
        sources:
        - serviceAccountToken:
            expirationSeconds: 3607
            path: token
        - configMap:
            items:
            - key: ca.crt
                path: ca.crt
            name: kube-root-ca.crt
        - downwardAPI:
            items:
            - fieldRef:
                apiVersion: v1
                fieldPath: metadata.namespace
                path: namespace
    status:
    conditions:
    - lastProbeTime: null
        lastTransitionTime: "2023-05-17T12:40:04Z"
        status: "True"
        type: Initialized
    - lastProbeTime: null
        lastTransitionTime: "2023-05-17T12:40:47Z"
        status: "True"
        type: Ready
    - lastProbeTime: null
        lastTransitionTime: "2023-05-17T12:40:47Z"
        status: "True"
        type: ContainersReady
    - lastProbeTime: null
        lastTransitionTime: "2023-05-17T12:40:04Z"
        status: "True"
        type: PodScheduled
    containerStatuses:
    - containerID: containerd://aba0b87ad26079ab210008c81f094291aa764bd5c5efedd1cc9022af4aca47f3
        image: docker.io/khk9346/kubia:latest
        imageID: docker.io/khk9346/kubia@sha256:6a53af0ff1cfe885e062a94f11b7b76fa84984c19064717a41e132e0cdd632a6
        lastState: {}
        name: kubia-container
        ready: true
        restartCount: 0
        started: true
        state:
        running:
            startedAt: "2023-05-17T12:40:47Z"
    hostIP: 192.168.5.15
    phase: Running
    podIP: 10.42.0.9
    podIPs:
    - ip: 10.42.0.9
    qosClass: BestEffort
    startTime: "2023-05-17T12:40:04Z"
  ```
    ```shell
    NAME           READY   STATUS              RESTARTS   AGE
    kubia-manual   0/1     ContainerCreating   0          12s
    ```
    ```shell
    NAME           READY   STATUS    RESTARTS   AGE
    kubia-manual   1/1     Running   0          **97s**
    ```

# Pod Life Cycle
## 용어
- Pending : Pod이 생성되고 스케줄링을 기다리는 단계입니다. 필요한 리소스가 할당되지 않거나 다른 Pod이 사용 중인 리소스 때문에 대기하는 동안 이 상태가 될 수 있습니다.
- Running : Pod의 모든 컨테이너가 실행 중이며, 최소한 하나의 컨테이너가 동작 중인 상태입니다. 이 상태에서 Pod은 애플리케이션을 처리하고 외부 요청에 응답할 수 있습니다.
- Succeeded : Pod의 모든 컨테이너가 성공적으로 종료된 상태입니다. 이 상태에서 Pod은 더 이상 실행되지 않으며, 작업이 완료된 것으로 간주됩니다. 일회성 작업이나 배치 작업을 수행한 후 작업이 성공적으로 완료되었을 때 이 상태가 될 수 있습니다.
- Failed : Pod의 모든 컨테이너가 종료되었지만, 적어도 하나의 컨테이너가 실패한 상태입니다. 컨테이너가 오류로 인해 비정상적으로 종료되거나 컨테이너 실행이 실패한 경우 이 상태가 될 수 있습니다.

## 파드의 단계
<img width="783" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/0dd671b8-307f-4cd0-bd85-e264f1dcf38f">

1. Pod을 스케줄링할 때까지 Pending상태입니다.
2. 이후 Running은 Pod가 스케줄링되고 컨테이너가 시작되어 정상적으로 실행 중인 단계입니다.
3. Succeeded는 Pod가 정상적으로 실행되고 작업이 완료되었을때의 단계이며, Pod가 필요없어, 종료될 수 있습니다. 
4. 반면 Failed는 하나 이상의 컨테이너가 실행 중에 오류가 발생해 작업이 실패한 단계입니다.

## 파드의 조건 (status.conditions)
- 리소스의 현재 상태를 설명하는 여러 조건을 포함하는 배열입니다.
- 일반적으로 아래의 속성들을 갖을 수 있습니다.
  - type: 조건의 유형을 식별하는 문자열입니다. 예를 들어, "Ready", "Initialized", "PodScheduled" 등이 될 수 있습니다.
  - status: 조건의 상태를 나타내는 문자열입니다. "True", "False", "Unknown" 중 하나의 값일 수 있습니다.
    - True: 조건이 만족되었거나 성공적으로 완료되었음을 나타냅니다.
    - False: 조건이 만족되지 않았거나 실패했음을 나타냅니다.
    - Unknown: 조건의 상태를 확인할 수 없음을 나타냅니다.
  - lastProbeTime: 조건을 마지막으로 확인한 시간입니다.
  - lastTransitionTime: 조건이 마지막으로 변경된 시간입니다.
  - reason: 조건의 변경 또는 상태를 설명하는 문자열입니다.
  - message: 조건에 대한 상세한 메시지를 포함하는 문자열입니다.
- 하지만 리소스마다 다르게 구성될 수 있기 떄문에, 리소스의 API 문서를 잘 참고해야 합니다.
### 컨테이너 라이프사이클 훅 (spec.containers[*].lifecycle)

- 간단하게 컨테이너 라이프사이클 사이에 훅을 넣는다는 말
- 컨테이너의 노출되는 훅은 두 가지 있습니다.
  - PostStart
    - 컨테이너가 생성된 직후에 실행된다. 
    - 하지만 훅이 컨테이너 엔트리포인트에 앞서서 실행된다는 보장은 없다.
    - 파라미터는 핸들러에 전달되지 않는다.
  - PresStop
    - API 요청 or liveness probe 실패, 선점, 자원 경합 등의 관리 이벤트로 인해 컨테이너가 종료되기 직전에 호출된다.
    - 이미 컨테이너가 termineated or completed 상태인 경우 PreStop 훅 요청이 실패한다.
    - 컨테이너를 중지하기 위한 TERM 신호가 보내지기 전에 완료해야 한다.
- 아래는 예제
  - exec 또는 httpGet을 통해 작업을 정의할 수 있다.
  - 사용 예시론 데이터 정리, 로그 기록 등의 작업을 처리할 때가 있다.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
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
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "echo 'Stopping...' && sleep 5"]

```

## 초기화 컨테이너 (spec.initContainers)

- 초기화 컨테이너는 항상 완료를 목표로 실행
- 각 초기화 컨테이너는 다음 초기화 컨테이너 시작 전에 성공적으로 완료해야함
- Docker Compose에서 depends_on으로 의존성 명시한것에 더해 순서들까지 만들었다고 생각하면 편안
    - 예로 db + 스키마를 선택한다든지 할 수 있다.
- 프로브 지원 X
    - 프로브는 파드가 실행 중일 때 컨테이너 상태를 확인하는데 사용되지만, 초기화 컨테이너는 파드 실행 전이라 그렇다고 합니다.

## Probe

- kubelet을 통해 관리
- Probe는 세 가지 종류가 있으며, Manifest에 작성할 수 있고 셋 다 작성 가능합니다. Spec에 Container마다 각각 다른 종류의 Probe를 둘 수 있다.
    - livenessProbe
        - 컨테이너가 살아있는지 여부를 결정하는 데 사용
    - readinessProbe
        - 컨테이너가 클라이언트 요청을 수신할 수 있는지 여부
    - startupProbe
        - 컨테이너가 시작되고 초기화 걸리는 시간을 고려해 초기화 전에 요청 전송 방지
        - startup probe가 있다면, 성공하기 전에 나머지 프로브는 활성화 되지 않음
        - 없으면 기본 상태 Success
    