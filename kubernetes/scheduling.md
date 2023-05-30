# Scheduling
> 이미지 출처 : https://blog.kubecost.com/
## Affinity(관계성)
- nodeSelector는 파드를 특정 레이블이 있는 노드로 제한하는 가장 간단한 방법이다.
- Affinity와 Anti-Affinity 기능은 표현할 수 있는 제약 종류를 크게 확장한다.
- Pod와 Node간의 관계를 정의하는 기능입니다. 쉽게말하면 특정 노드에 파드를 스케줄링 하려면 Node Affinity, 특정 파드가 다른 파드와 같은 노드에 스케줄링 되도록 할 땐 Pod Affinity를 사용합니다.
- Affinity는 Pod의 sec 섹션에 설정됩니다.
### Node Affinity
<img width="804" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/f969f493-4f9a-41ce-9571-fea7bb884170">

- 노드와 관련된 Affinity를 설정합니다. 노드의 Label과 매치되는 조건을 지정해 파드를 특정 노드에 바인딩 할 수 있습니다.
- Node Affinity는 다음과 같은 세 가지 옵션을 제공합니다.
  - requiredDuringSchedulingIgnoredDuringExecution : 파드가 특정 노드와 매치되어아먄 스케줄링 됩니다. nodeSelector와 유사하지만, 좀 더 표현적인 문법을 제공합니다.
  - preferredDuringSchedulingIgnoredDuringExecution : 파드가 특정 노드와 가장 잘 매치되는 경우 스케줄링 됩니다. 해당 노드가 없어도, 스케줄러는 여전히 파드를 스케줄링합니다.
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
- 이제 agent1에 label-2:key-2 라벨을 삭제하고 다시 테스트하면 아래와 같이 설정된다.

```shell
$ kubectl get pods -o=custom-columns='POD:metadata.name,NODE:spec.nodeName'

POD                                 NODE
with-affinity-anti-affinity         k3d-k3s-default-agent-1
with-affinity-anti-affinity2        k3d-k3s-default-agent-0
```

- 위의 템플릿 두 개의 pod를 설정하여 파드를 생성하면 위와 같은 결과가 나오며, weight가 더 큰 쪽으로 pod가 할당된다.
  - 라벨 모두가 존재할 땐 weight가 agent1이 51, agent0이 1로 agent1이 더 높으며, 라벨을 삭제한 경우 agent1이 0, agent0이 1로 더 높다.
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

<img width="905" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/8322bb93-fa1c-4855-850c-d7438ee4f695">

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
- topologyKey 같은 경우엔, topology를 사용하는 이유에 대해 알면 좋은데, 기본적으로 각 존마다 균일하게 스케줄되도록 하며, 클러스터에 문제가 생길경우 스스로 치유하도록 설정하기 위해 존재한다.
  - 즉 topologyKey의 키와 동일한 값을 가진 레이블이 있는 노드는 동일한 토폴로지에 있는 것으로 간주한다.
  - 여기선 topologyKey에 기반하여 스케줄링 또는 분산 배치를 수행하는데 사용된다.
  - topology.kubernetes.io/zone는 클라우드 환경에서 노드가 속한 물리적인 존 또는 리전을 나타냅니다.
    - 즉 이 물리적인 존 또는 리전을 기반으로 분산 배치를 해준다는 의미입니다.
- podAntiAffinity같은 경우엔 label-2=key-2가 있는 레이블에 대해 가중치 50을 label-1=key-1은 1을 부여합니다.
- 다른점이 있다면 노드 어피니티는 노드의 레이블을 기반으로, 파드 어피니티는 노드에 속한 파드를 기반으로 선택됩니다.
- Anti가 붙어있으면 weight의 역순을 우선적으로 고른다.

우선 아래의 경우는 topologyzone이 agent0=zone-0, agent1=zone-1 agent2=zone-2가 적용되어 있다. 또한 pod1은 zone-0에 pod2는 zone-1을 가진 노드에 배치하였다.
```shell
k get pods -o wide
NAME            READY   STATUS              RESTARTS   AGE   IP           NODE                      NOMINATED NODE   READINESS GATES
pod-affinity    1/1     Running             0          22s   10.42.3.41   k3d-k3s-default-agent-0   <none>           <none>
pod-affinity2   0/1     ContainerCreating   0          3s    <none>       k3d-k3s-default-agent-1   <none>           <none>
```
이제 이 상태에서 위의 정의된 with-pod-affinity를 하면, zone-2로 설정된 agent2로 가야한다. weight가 가장 낮은 0이기 때문이다.
```shell
k get pods -o wide
NAME                READY   STATUS              RESTARTS   AGE     IP           NODE                      NOMINATED NODE   READINESS GATES
pod-affinity        1/1     Running             0          7m54s   10.42.3.41   k3d-k3s-default-agent-0   <none>           <none>
pod-affinity2       1/1     Running             0          7m35s   10.42.1.28   k3d-k3s-default-agent-1   <none>           <none>
with-pod-affinity   0/1     ContainerCreating   0          2s      <none>       k3d-k3s-default-agent-2   <none>           <none>
```
마스터 노드도 동일한 가중치이기 때문에 랜덤으로 선택된다.

# 테인트(Taints)와 톨러레이션(Tolerations)

<img width="900" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/04370f6e-8508-4efd-be08-83ab15fe002a">

- 노드 어피니티는 특정노드에 배치하거나, 배치하고싶지 않을 때 사용했다.
  - required나 prefered를 통해
  - 노드에 파드를 끌어들이는 속성
- 반면 테인트는 특정 노드에 파드의 스케줄링을 제한하기 위해 사용
  - 노드가 파드셋을 제외시킨다.
- 톨러레이션은 파드에 적용된다.
  - 톨러레이션을 통해 스케줄러는 그와 일치하는 테인트가 있는 파드를 스케줄할 수 있다.
- 테인트와 톨러레이션은 함께 작동하여 파드가 부적절한 노드에 스케줄되지 않게 한다.

## 테인트
- key:value:effect 형식으로 구성
- 노드에 조건을 설정하여 파드의 배치를 제한하는 역할을 한다.
- 일반적으로 NoSchedule, PreferNoSchedule 효과를 가지는 테인트를 사용
  - NoSchedule: 파드가 해당 테인트와 일치하는 노드에 배치되지 않도록 합니다. 즉, 파드가 해당 노드에 스케줄링되지 않습니다.
  - PreferNoSchedule: 파드가 해당 테인트와 일치하는 노드에 배치되지 않는 것을 선호하지만, 다른 조건을 충족하는 경우에는 배치될 수 있습니다.
  - NoExecute: Toleration이 없는 모든 파드들을 **즉시** 노드로부터 제외한다. tolerations 내의 tolerationSeconds 필드를 통해 테인트와 일치하는 파드가 노드안에서 얼마나 머물수 있게 할지 설정할 수 있다.
    - 만약 테인트가 tolerationSeconds안에 삭제되면 파드는 노드에서 제외되지 않는다.
    - 만약 파드에 tolerationSeconds가 없다면, 수동으로 제거하지 않는 한 노드에서 계속 유지된다.

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

<img width="911" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/6a75204c-7d5a-4758-9977-5e9449837966">

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