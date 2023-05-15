# ReplicaSet

- 초기엔 ReplicationController지만 지금은 ReplicaSet을 대용으로 많이 쓰고 있다.
- 기능은 거의 동일하지만 추가된 느낌이다.
- 일반적으로 Deployment와 같이 쓰인다고한다.

## ReplicaSet vs ReplicationController

- 레이블 셀렉터 이외에도 레이블 표현식이 추가됨
- 레플리케이션컨트롤러와 다르게 여러 레이블을 가진 파드도 선택 가능
    - 레플리케이션컨트롤러는 equality-based 라는 뜻
    - ReplicaSet은 set-based selector이다.
    - 즉 또는이라는 선택지가 생긴다.
    - In을 통해 특정 value만 가진 파드도 선택 가능하다.
        - (Not)In : 레이블의 값을 통해 찾음
        - (DoesNot)Exist : 파드의 키를 통해 찾음
        
        ```yaml
        - {key: app, operator: In, values: [my-app, my-app-v2]}
        - {key: env, operator: NotIn, values: [dev]}
        ```
        
    - 위처럼 사용한다.
- 또한 값 이외에도 키만으로도 매칭 가능
- ReplicationController는 일괄적으로 파드를 생성하거나 제거하지만, ReplicaSet은 아니다.
- 

## 특징

- 컨트롤러가 없는 파드는 조건에 맞으면 레플리카셋의 소유가 된다.
- 레플리카셋은 apiVersion v1이 아니기 떄문에 다른걸 선택해야함
- Selector는 matchLabels를 동일하게 쓸 수도 있다.
    - 기능도 동일하긴 하다.
- matchExpressions를 통해 더 파드를 찾는 기준이 많아졌다.
- rs가 약자이다.

# DaemonSet

- 파드가 특정 노드에서만 실행 가능하도록 지정할 때 사용
- DaemonSet은 클러스터 내의 모든 파드를 배치하고 유지하기 위해 사용한다.
- 각 노드당 하나의 파드가 배치되는 것을 목표로 한다.
- 노드셀렉터를 통해 특정 레이블을 가진 노드에만 배포하도록 할 수 있다.
- ReplicaSet과 DaemonSet이 배포하는 파드가 같은 matchLabels로 관리하면 어떻게 돼지?
- 약자는 ds이다.
- node에 추가로 라벨을 넣고, 삭제하는건 가능하다.