# Job

- 잡 리소스는 완료 가능한 단일 태스크를 수행하는 파드를 실행한다.
    - 즉 단일 태스크를 수행하고, 수행 완료 후 다시 생성되지 않는다.
- kube-controller-manager에 의해 관리된다.
- restartPolicy는 default가 Always로 설정되어 있기 때문에 Job으로 만들경우 항상 이를 **OnFailuer** or **Never**를 작성하여 재실행하지않게 한다.
  - 실제로 job 같은 경우 restartPolicy를 설정하지 않으면 아래처럼 출력이 되며, deployment, replicaset, pod등은 확인해보면 Always로 설정되어 있다.
  - restartPolicy를 job에서 사용할 수 없는 이유는 job 리소스의 경우 완료 가능한 단일 태스크를 수행하는 파드를 실행하기 때문이다. 일회성 리소스에 Awlays 설정이 들어갈 경우 중복 작업이 발생할 수 있기 때문이다. 
```shell
$ k apply -f job.yaml
The Job "my-job" is invalid: spec.template.spec.restartPolicy: Required value: valid values: "OnFailure", "Never"
```
```shell
$ k get rs kubia-deployment-67b44c64cc -o yaml | grep restartPolicy
      restartPolicy: Always
```

```json
apiVersion: batch/v1
kind: Job
metadata:
  name: my-job
spec:
  template:
    spec:
      restartPolicy: OnFailure
      containers:
      - name: my-container
        image: my-image
        command: ["my-command"]
```
- 실제로 실행하면 아래처럼 나온다.
  - job이 완료되면 컨테이너가 내려가고, job 파드는 그대로 남아있는 모습이 보인다.
  - Status는 Completed로 job이 잘 끝난걸 확인할 수 있으며, Completions는 1/1로 1번의 잡 시도중 1번 성공한 모습이 보인다. 
```shell
$ k get pod
NAME                                READY   STATUS      RESTARTS   AGE
my-job-mcj78                        0/1     Completed   0          24s
$ k get jobs
NAME     COMPLETIONS   DURATION   AGE
my-job   1/1           22s        27s
```
제가 이전에 설명한 내용에 오해가 있을 수 있습니다. 제가 언급한 내용이 일반적인 Job의 동작 방식과 일치하지 않았을 것입니다. 죄송합니다.

실제로, completions 속성은 Job의 완료를 지정하는 것이 아니라, 생성될 파드의 총 수를 지정합니다. Job은 해당 수의 파드를 생성하고 작업을 실행하며, 모든 파드의 작업이 완료되면 Job은 완료됩니다.

parallelism 속성은 동시에 실행될 수 있는 파드의 최대 수를 지정합니다. 이 값은 completions 값보다 작거나 같아야 합니다. parallelism 값이 1보다 큰 경우, Job은 동시에 여러 파드를 생성하고 병렬로 작업을 실행합니다. 작업이 완료되지 않은 파드는 Job이 재시작되거나 다른 파드가 완료될 때까지 대기합니다.

따라서, Job은 completions 속성에 지정된 수의 파드를 생성하고, parallelism 속성에 지정된 수 이하로 동시에 실행합니다. 모든 파드의 작업이 완료되면 Job은 완료됩니다.
- backoff 실패 정책
  - 구성에 논리적 오류가 포함되어 있어 몇 회의 재시도 이후에 잡이 실패되도록 만들어야 하는 경우가 있다.
  - .spec.backoffLimit의 값에 재시도 횟수를 설정 할 수 있다
  - backoffLimit에 도달하면, 잡이 실패한 것으로 간주된다.
- 추가적으로 **completions** 필드과 **parallelism** 필드이 있다.
    - **completions** 필드은 순차적으로 성공한 파드의 수를 지정합니다.
      - Job은 해당 수의 파드를 생성하고 작업을 실행하며, 모든 파드의 작업이 완료되면 Job은 완료됩니다.
      - 재시도가 backoffLimit에 도달하면 job은 실패로 처리됩니다.
    - **parallelism** 필드은 병렬적으로 실행될 수 있는 파드의 최대 수를 지정합니다.
      - completions 값보다 작거나 같아야 합니다.
      - Job 실행중 **parallelism 속성을 변경**할 수 있다.
  - completions 필드는 파드가 성공적으로 생성됐는지를 확인할 수 있습니다.
  ```
    $ k get jobs
    NAME     COMPLETIONS   DURATION   AGE
    my-job   1/15          8s         8s
  ```


- 아래처럼 작업의 복제본 수를 변경해 job을 확장할 수 있다.
```json
k scale job multi-completion-batch-job --replicas 3
```
- 현재는 job에서 scale은 deprecated 되어 spec.parallelism을 직접 변경해야 한다.
```
$ k  scale # 이후에 미리보기를 보면 아래와 같이 job은 없다.
--replicas             -- The new desired number of replicas. Required.
deployment             replicaset             replicationcontroller  statefulset
```
- 아래 처럼 manifest를 수정하면 parallelism이 적용된다.
  - 처음엔 3으로 하고, 이후 2로 변경하면 실제로 아래처럼 나온다.
```
k apply -f job-com-pa.yaml 
job.batch/my-job created # 이후에 변경
k apply -f job-com-pa.yaml
job.batch/my-job configured
```
```
spec:
  backoffLimit: 6
  completionMode: NonIndexed
  completions: 15
  parallelism: 2
```
- 아래는 completions 필드과 parallelism 의 예시이다.
```json
apiVersion: batch/v1
kind: Job
metadata:
    name: my-job
spec:
completions: 5
parallelism: 2
    template:
        spec:
        restartPolicy: OnFailure
        containers:
        - name: my-container
            image: my-image
            command: ["my-command"]
```

```shell
$ k get pods
NAME                                READY   STATUS              RESTARTS   AGE
my-job-mdgm9                        0/1     ContainerCreating   0          2s
my-job-2vmft                        0/1     ContainerCreating   0          2s
```
```shell
k get pods
NAME                                READY   STATUS              RESTARTS   AGE
my-job-8b5fd                        0/1     ContainerCreating   0          2s
my-job-2vmft                        0/1     Completed           0          11s
my-job-mdgm9                        0/1     Completed           0          11s
my-job-6vb2k                        0/1     ContainerCreating   0          2s
```
```shell
k get pods
NAME                                READY   STATUS              RESTARTS   AGE
my-job-2vmft                        0/1     Completed           0          20s
my-job-mdgm9                        0/1     Completed           0          20s
my-job-8b5fd                        0/1     Completed           0          11s
my-job-6vb2k                        0/1     Completed           0          11s
my-job-5vmp6                        0/1     ContainerCreating   0          1s
```
- **activeDeadlineSeconds** 필드를 통해 파드가 일정시간 실행되면, 잡을 실패하는 것으로 간주하게 할 수 있다.
  - 실제로 아래의 결과를 확인해보면 잡이 실행되지만, 일정시간 이후에 pod가 종료되며, job이 실패한것으로 나타난다.
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: my-job
spec:
  activeDeadlineSeconds: 10
  template:
    spec:
      restartPolicy: OnFailure
      containers:
      - name: my-job-container
        image: ubuntu:latest
        command: ["/bin/sh", "-c", "sleep 20"]  # 실행 시간이 activeDeadlineSeconds보다 긴 명령어로 변경
```
```shell
$ k get job
NAME     COMPLETIONS   DURATION   AGE
my-job   0/1           6s         6s
$ k get pods
NAME           READY   STATUS    RESTARTS   AGE
my-job-7s6wn   1/1     Running   0          9s
```
```
$ k get pods
No resources found in default namespace.
$ k get job
NAME     COMPLETIONS   DURATION   AGE
my-job   0/1           34s        34s
```

# 잡을 주기적으로 또는 한 번 실행되도록 스케줄링
- 배치 잡은 어떻게 할까? 운영체제에서는 크론이 있다.
    - 쿠버네티스도 이를 지원한다.
- 바로 CronJob이 있다.
    - kube-controller-manager의 시간대를 기준으로 동작한다.
    - 특정시간마다 “약” 한 번의 job 오브젝트를 생성한다.
        - 약인 이유는 두 번 혹은 0 번 실행될 수 있기 때문이다.
- CronJob이라는 resource가 있다.
- **Job과 마찬가지로 kube-controller-manager**에 있으며, **시간대 또한 kube-controller-manager 기준**으로 작동된다.

```json
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox:1.28
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```
```shell
$ k get cronjob
NAME    SCHEDULE    SUSPEND   ACTIVE   LAST SCHEDULE   AGE
hello   * * * * *   False     0        19s             116s
$ k get job
NAME             COMPLETIONS   DURATION   AGE
hello-28072845   1/1           21s        84s
hello-28072846   1/1           3s         24s
```
- 위처럼 schedule에 **cron expression**을 통해 작성 가능하며 실행할 jobTemplate를 작성한다.
- spec의 컨테이너는 여러 개 정의 가능
- **startingDeadlineSeconds** 필드 있다.
    - ?초안에 시작안하면 실패한다는 표시로 간주한다.
    - 실패한 Cronjob은 이후에 다시 스케줄되지 않습니다.
    ```json
    metadata:
      name: hello
      spec:
        schedule: "*/1 * * * *" # 크론 표현식
    		startingDeadlineSeconds: 15
    ```
    
- 위에서 크론잡은 두 개의 잡이 생성되거나 전혀 생성이 안될 수 있다고 하였다. 부가설명은 아래와 같다.
    - 첫 번째 문제인 두 개의 잡이 생성될 때의 해결 법은 멱등성을 갖는 방법이다.
        - 여러 번 실행되어도 같은 결과가 나온다면 두 번 실행되는게 문제될게 없다.
    - 두 번째 문제인 전혀 생성되지 않는 경우는, 여러 이유가 있을 수 있다.
        - 설정,리소스,로그,이벤트 등 다양한 이유를 찾아보자.
