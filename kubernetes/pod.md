# Pod
- k8s의 가장 기본 단위
- 관련된 컨테이너들을(일반적으로 1개) 함께 그룹화하여 배포하고 실행하는 데 사용

## 컨테이너 재시작
- 컨테이너가 종료되었을 때 어떻게 처리할지 처리를 결정하는 정책입니다.
- 3가지 정책이 있습니다.
  - Always: 컨테이너가 종료되면 항상 재시작
  - OnFailure: 컨테이너가 실패한 경우에만 재시작
  - Never: 컨테이너 종료시 재시작 X
## imagePullPolicy
- 파드가 생성되고, 컨테이너를 생성할 때 이미지를 PULL할때의 정책을 정하는 속성이다.
- 만약 따로 설정하지 않는다면 
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
    