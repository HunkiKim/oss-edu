# Job

- 잡 리소스는 완료 가능한 단일 태스크를 수행하는 파드를 실행한다.
    - 즉 단일 태스크를 수행하고, 수행 완료 후 다시 생성되지 않는다.
- kube-controller-manager에 의해 관리된다.
- restartPolicy는 보통 Always로 설정되어 있기 때문에 Job으로 만들경우 항상 이를 **OnFailuer** or **Never**를 작성하여 재실행하지않게 한다.

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

- 위의 manifest를 제출하면 처음엔 DESIRED로 되어있다. 완료되면 Completed로 표시되고 -a를 통해서만 볼 수 있다.
    - 삭제하지 않으면 남아있다고 한다.
- 추가적으로 **completions** 방식과 **parallelism** 방식이 있다.
    - **completions** 방식은 순차적으로 Job 리소스를 생성한다.
    - **parallelism** 방식은 병렬적으로 같이 Job 생성된다.
        - Job 실행중 **parallelism 속성을 변경**할 수 있다.
        
        ```json
        k scale job multi-completion-batch-job --replicas 3
        ```
        
        - 위처럼 변경하면 즉시 다른 파드가 기동돼, 세 개의 파드가 실행된다. 이를 **Job Scaling**이라고 한다.
    - 만약 둘 다 같이 정의되어 있다면 parallelism으로 실행된 Job들중 하나가 Completed가 된다면 이후로 completions의 잡들이 순차적으로 실행된다.
    
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
    
- **activeDeadSeconds** 속성을 통해 파드가 일정시간 실행되면, 잡을 실패하는 것으로 간주하게 할 수 있다.

# 잡을 주기적으로 또는 한 번 실행되도록 스케줄링

- 크론잡을 쓰면 된다.
    - linux에도 비슷한게 있다.
- CronJob이라는 resource가 있다.
- **Job과 마찬가지로 kube-controller-manager**에 있으며, **시간대 또한 kube-controller-manager 기준**으로 작동된다.

```json
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
  spec:
    schedule: "*/1 * * * *" # 크론 표현식
    jobTemplate:
      spec:
        template:
          spec:
            containers:
              - name: hello
                image: busybox
                args:
                  - /bin/sh
                  - -c
                  - date; echo Hello from the Kubernetes cluster
                  restartPolicy: OnFailure
```

- 위처럼 schedule에 **cron expression**을 통해 작성 가능하며 실행할 jobTemplate를 작성한다.
- spec의 컨테이너는 여러 개 정의 가능
- **restartDeadlineSeconds** 속석이 있다.
    - ?초안에 시작안하면 실패한다는 표시로 간주한다.
    
    ```json
    metadata:
      name: hello
      spec:
        schedule: "*/1 * * * *" # 크론 표현식
    		startingDeadlineSeconds: 15
    ```
    
- 크론잡은 두 개의 잡이 생성되거나 전혀 생성이 안될 수 잇다.
    - 첫 번째 문제인 두 개의 잡이 생성될 때의 해결 법은 멱등성을 갖는 방법이다.
        - 여러 번 실행되어도 같은 결과가 나온다면 두 번 실행되는게 문제될게 없다.
    - 두 번째 문제인 전혀 생성되지 않는 경우는, 여러 이유가 있을 수 있다.
        - 설정,리소스,로그,이벤트 등 다양한 이유를 찾아보자.