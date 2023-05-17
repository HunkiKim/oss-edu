- 애플리케이션은 결국 변경된다.
- 여기선 쿠버네티스 클러스터에서 실행되는 애플리케이션을 업데이트 하는 방법과 쿠버네티스가 어떻게 무중단 업데이트 프로세스로 전환하는 데 도움을 주는지 설명해준다.
- Deployment를 알아보자.

# 파드에서 실행 중인 애플리케이션 업데이트

- 지금까지 배운걸 토대로 보면 흐름은 아래와 같다.
    - 클라이언트 → 서비스 → 파드 ← 레플리카셋(컨트롤러)
- 기존 방식으로 애플리케이션을 업데이트 한다면 v1 태크가 지정된 기존 파드를 모두 제거하고, 새로운 v2 파드로 교체해야 한다.
- 또 다른 방법으로는 순차적으로 새 파드를 추가하고 기존 파드를 점진적으로 제거하는 방법이 있다.
- 둘의 차이는 아래와 같다.
    - 첫 번째 방식은 짧은 시간 애플리케이션 사용을 못한다.
    - 두 번째 방식은 동시에 두 가지 버전을 실행해야 한다.

## 오래된 파드를 삭제하고 새 파드로 교체

- 모든 파드 인스턴스를 새 버전의 파드로 교체하기 위해 레플리카셋을 사용하는 방법은 기존 레플리카셋의 라벨을 v2를 참조하도록 수정하고, 이전 파드 인스턴스를 삭제하는 것이다.
    - 이러면 v2로 바꿨기 떄문에 새로운 인스턴스를 desired state에서 정의한 개수만큼 실행될 것이다.
    - v1은 수동으로 삭제해준다.
    - v2는 자동으로 생성된다.
    - v1을 삭제하고 v2가 생성되는 사이의 시간동안 다운타임이 발생한다.

## 새 파드 기동과 이전 파드 삭제

- 만약 다운타임을 없애고 싶다면 새 파드들을 모두 띄우고, 이전 파드를 삭제하면 된다.
- 잠시 동안 두배의 파드가 실행되기 때문에 리소스가 많이 필요하다.

### 한 번에 이전 버전에서 새 버전으로 전환

- 파드의 앞쪽에는 일반적으로 서비스를 배치한다.
- 새 버전 실행동안 서비스는 이전 버전의 파드들로 포워딩이 되어있다.
- 새 파드가 모두 실행되면 서비스의 레이블 셀렉터를 변경하고 서비스를 새 파드로 전환할 수 있다.
- 이를 블루-그린 디플로이먼트라고 한다.
- 버전이 어느정도 안정화되면 이전 버전의 ReplicaSet을 삭제해 모든 파드와 상위리소스를 삭제할 수 있다.

### 롤링 업데이트 수행

- 한 번에 모든 것을 하는 대신 단계별로 교체하는 롤링 업데이트도 가능하다.
- 천천히 이전 버전 스케일 다운 및 천천히 새로운 버전 스케일 업을 하는 방식이다.
- 수동으로 롤링 업데이트는 힘들다.
    - 서비스에서 v1라벨과 v2라벨에 대해 모두 관리가 필요하기 때문이다.
- 쿠버네티스에선 자동으로 가능하게 해준다.

# 레플리카셋으로 자동 롤링 업데이트 수행

- kubectl을 통해 롤링 업데이트가 가능하지만 이 방식도 이젠 옛날 방식이다.

## 애플리케이션의 초기 버전 실행

- 먼저 배포된 애플리케이션이 필요한데 이전에 작성한 nodejs 애플리케이션 코드를 그대로 쓴다.

```json
const http = require('http');
const os = require('os');

console.log("test");

var handler = function(request, response) {
    console.log("a" + request.connection.remoteAddress);
    response.writeHead(200);
    response.end("response end " + os.hostname() + "\n");
}

var www = http.createServer(handler);
www.listen(8080);
```

- 일단 v1이라는 태그는 안달려있으니 이를 예시로 해본다.
    
    ```bash
    docker tag khk9346/kubia:latest khk9346/kubia:v1;
    docker push khk9346/kubia:v1;
    ```
    
    - 위의 명령어를 통해 docker-hub의 repository의 이미지 태그 변경
- v1의 yaml을 작성하자.

```json
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: kubia-v1
spec:
  replicas: 3
  selector:
    matchExpressions:
      - key: app
        operator: In
        values:
          - kubia
  template:
    metadata:
      name: kubia
      labels:
        app: kubia
    spec:
      containers:
        - name: nodejs
          image: khk9346/kubia:v1
--- # 대시 3개가 있는 줄로 여러 리소스 정의를 포함 가능
apiVersion: v1
kind: Service
metadata:
  name: kubia-service
spec:
  type: LoadBalancer
  selector:
    app: kubia
  ports:
  - port: 80
    targetPort: 8080
```

- 이제 curl로 한 번 요청해보자.

```bash
while true; do curl http://34.64.68.148; done;
```

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/d07a9ded-7448-4641-8126-58b9bfb321f5)

- 여기서 [http://34.64.68.148은](http://34.64.68.148은) 서비스의 external-ip이다.

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/3c4a3928-8721-48f0-9af9-92189c2fcdcc)

## kubectl을 이용한 롤링 업데이트

- application v2를 만들자. 출력을 v2라고 바꾸자.
- v2를 build하고 docker-hub에 배포한다.
- while문으로 요청을 계속 보내면서 아래의 rolling-update 명령을 수행해보자.
    - Deprecated 되었다.

```bash
k rolling-update kubia-v1 kubia-v2 --image=khk9346/kubia:v2
```

- v2 yaml을 새로 만들자

```json
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: kubia-v2
spec:
  replicas: 3
  selector:
    matchExpressions:
      - key: app
        operator: In
        values:
          - kubia
  template:
    metadata:
      name: kubia
      labels:
        app: kubia
    spec:
      containers:
        - name: nodejs
          image: khk9346/kubia:v2
--- # 대시 3개가 있는 줄로 여러 리소스 정의를 포함 가능
apiVersion: v1
kind: Service
metadata:
  name: kubia-service
spec:
  type: LoadBalancer
  selector:
    app: kubia
  ports:
  - port: 80
    targetPort: 8080
```

- 이미 deprecated 되었지만 원리만 간단하게 알아보자
    - 이미 실행중인 v1 pods에 deployment=v1 레이블을 추가한다.
    - scaling kubia-v2 up to 1
        - kubia-v2의 파드 라벨은 app:kubia, deployment: v2이며 selector도 이와 맞춘다.
    - scaling kubia-v1 down to 2
    - 이렇게 되면 점진적으로 교체가 가능하다. 서비스는 그대로 app=kubia에 대해 selector가 되어있어 v1,v2를 둘 다 요청한다.
        - 처음엔 v1의 비율이 높지만 점차 v2의 비율이 높아진다.
- 쓰지않는 이유는 아래와 같다.
    - 위에서 보듯이 오브젝트가 수정된다.
        - 레이블과 레이블 셀렉터가 수정되는데, 이건 저자가 예상못하기 때문이다.
    - 또한 롤링업데이트가 kubectl 클라이언트에 의해 이루어 진다는 것이다.
        - 업데이트중에 네트워크 끊어지면 업데이트 프로세스 중간에 중단되고, 중간상태에서 끝이나기 때문에 치명적이다.
    - 또한 쿠버네티스는 desired state가 알아서 유지되게 해야하는데 실제 명령으로 업데이트가 수행된다.
        - 쿠버네티스 철학에 맞지 않다.

# 애플리케이션을 선언적으로 업데이트하기 위한 디플로이먼트 사용하기

- 명령 말고 선언적으로 업데이트 하는게 낫다.
- 디플로이먼트는 낮은 수준(lower-level)의 개념으로 간주되는 replicaset을 통해 수행하는게 아니라, 선언적으로 업데이트 하기 위한 높은 수준(high-level)의 리소스다.
    - 더 추상적이다.
- 디플로이먼트를 생성하면 레플리카셋 리소스가 그 아래에 생성된다.
    - Deployment → ReplicaSet → Pods 순서로 관리된다.
    - 그러면 ReplicaSet은 deployment가 관리하는건가?
- 디플로이먼트를 통해 업데이트를 훨씬 쉽게 해보자.(위에 어려웠으니)

## 디플로이먼트 생성

- 디플로이먼트 생성은 레플리케이션컨트롤러 생성과 다르지 않다.
- 디플로이먼트는 레이블 셀렉터, 원하는 레플리카 수, 파드 템플릿으로 구성된다.
- manifest 작성해보자
    - 처음엔 v1 이다.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubia-deployment
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
        - name: nodejs
          image: khk9346/kubia:v1
---
apiVersion: v1
kind: Service
metadata:
  name: kubia-service
spec:
  type: LoadBalancer
  selector:
    app: kubia
  ports:
  - port: 80
    targetPort: 8080
```

## 배포

```bash
k apply -f deployment-v1.yaml
```

- 위 명령어를 치면 배포가 된다.
    - deployment, replicaset, pods 모두 잘 나온다.
- 이제 아래 명령어로 롤아웃 상태 출력을 해보자.
    - 롤아웃은 일반적으로 새로운 버전을 점진적으로 배포하는 프로세스를 의미

```bash
k rollout status deployment
```

- 아래와 같이 나온다.
    - 롤아웃은 성공적으로 실행 됐고, 파드 레플리카 세 개가 시작돼 실행중이라는 의미
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/bee4b40c-3cb7-4e79-afbb-6da360d926b8)

## 디플로이먼트가 레플리카셋 생성 방법, 레플리카셋이 파드를 생성하는 방식

- 이전 파드의 이름은 kubia-m3213 이런식으로 구성됐다. 하지만 디플로이먼트에서 생성한 파드 세 개에는 이름 중간에 값이 추가된다.

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/0d1091b1-1f43-4d67-9c15-690c8b544d5c)

- 디플로이먼트와 파드 템플릿의 해시값을 의미하며 레플리카셋이 파드를 관리함을 의미한다.
    - 즉 deployment의 name인 kubia-deployment 뒤의 숫자는 디플로이먼트의 해시값, 그리고 그 뒤는 파드 템플릿의 해시값을 의미한다.
        - 레플리카셋이 이러한 파드를 관리한다는 뜻이다.
        - 디플로이먼트는 파드를 직접 관리하지 않는다. 레플리카셋을 생성하고, 이들이 파드를 관리하도록 한다.
        - 마찬가지로 레플리카셋의 이름에도 해당 파드 템플릿의 해시값이 포함
        
        ![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/ee9ff7bc-2479-4e10-a344-b188c0e026d6)
        
    - 즉 디플로이먼트의 특정 해시값에 의해 관리되는 레플리카셋에 의해 관리되는 파드들 이라는 의미를 지닌다.
- 서비스파드는 그대로 한다.

## 업데이트

- 이제 업데이트를 해보자.
- 원래는 위의 과정을 거쳤다면 지금은 다르다.
- 아래의 명령어로 이미지를 바꿔보자

```bash
k set image deployment kubia-deployment nodejs=khk9346/kubia:v2
```

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/2506efd1-b067-4810-b3fe-fef4f56c17fc)

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/9b3fbc29-d9d4-4f17-a81a-17bd2263c22b)

- 위처럼 알아서 바뀐다

<aside>
💡 k edit은 기본 편집기로 오브젝트 메니페스트를 변경하며, k patch는 오브젝트의 개별 속성을 수정, k apply는 yaml 속성 값 적용, replace는 새 것으로 교체, set image는 컨테이너 이미지 변경이다.

</aside>

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/b1974de9-5c65-48a1-b9be-e500f6de609e)

- 보면 순차적으로 바뀐걸 알 수 있다
- 어쨋든 디플로이먼트가 알아서 내부에서 버전이 바뀐걸 감지하고, replicaset을 생성하고, pod를 줄이고 생성해줫다.

## 롤백

- 일단 망가지는 application을 작성해보자

```bash
const http = require('http');
const os = require('os');

console.log("Kuba server starting...");

var requestCount = 0;

var handler = function(request, response) {
    console.log("a" + request.connection.remoteAddress);
    if (++requestCount >= 5) {
        response.writeHead(500);
        response.end("Some internal error has occurred! This is pod " + os.hostname() + "\n");
        return;
    }
    response.writeHead(200);
    response.end("v2 response end " + os.hostname() + "\n");
}

var www = http.createServer(handler);
www.listen(8080);
```

- 이후 롤아웃을 통해 진행 상황을 확인하자

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/9c057d68-9cb5-45a8-bfc3-b2359947c808)

- 오류가 잘 난다

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/eee0255e-8adb-48b3-be82-e58da916f61f)

- 이제 롤아웃 돌려보자.

### 롤아웃 되돌리기

- 자동으로 롤아웃을 차단할 수 있지만, 일단 수동으로 먼저 해보자.

```bash
k rollout undo deployment kubia
```

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/62b02aaf-f32b-423a-b255-f412a5bb03a7)

- CHANGE-CAUSE가 나오는 이유는 —record 명령줄 옵션이 없기 때문이다.

### 특정 디플로이먼트 개정으로 롤백

- Revision 번호를 통해 롤백도 가능하다. 아래 명령어를 써보자

```bash
k rollout undo deployment kubia-deployment --to-revision=3
```

- editionHistoryLimit 속성으로 과거 레플리카셋 목록을 없앨 수 있다.

## 롤아웃 속도 제어

- maxSurge
    - 롤링 업데이트 중 동시에 생성할 수 있는 파드 수
    - 가용성 증가, 클러스터 리소스 사용량 증가
- maxUnavaillable
    - 롤링 업데이트 중 사용 못하게할 파드 수
    - 가용성 저하, 클러스터 리소스 사용량 감소
- 아래와 같이 두 가지 속성을 지정해 줄 수 있다.

```json
spec:
	strategy:
		rollingUpdate:
			maxSurge: 1
			maxUnavailable: 0
		type: RollingUpdate
```

- maxSurge → 1, maxUnavailable → 0, 그리고 레플리카셋이 쓰는게 3이라고 해보자.
    - 4개 → 삭제 → 3개 → 생성 → 4개 → 삭제 → 3개 → 생성 을 반복하여 롤링 업데이트를 진행한다.

## 롤아웃 프로세스 일시 중지

- 롤아웃 프로세스 중에 배포를 일시 중지해, 나머지 롤아웃 진행 전 새 버전으로 모든 것이 정상인지 확인 가능
    - 카나리 배포를 통해 소수의 사용자가 새 버전을 사용하게 하고, 제대로 작동하면 이후 롤아웃 전체 진행하는 방식이다.
- deployment중에 pause 하면 된다.

## 롤아웃 재개

- 롤아웃 프로세스를 일시 중지하면 클라이언트 요청 중 일부만파드 v4에 도달
- 대부분 파드 v3을 호출
- 새버전 작동 확신시 디플로이먼트를 다시 시작해 이전 파드를 모두 새 파드로 교체 가능하다.
- 하지만 수동으로 rollout pause 와 rollout resume을 통해 하는건 힘들다.

## 잘못된 버전의 롤아웃 방지

- minReadySeconds는 롤아웃 속도를 늦추는 것이다.
    - 이는 오작동 버전의 배포를 방지하기 위해 존재한다.
- minReadySeconds 속성은 파드를 사용 가능한 것으로 취급하기 전에 새로 만든 파드를 준비할 시간을 지정한다.
    - 모든 파드가 레디니스 프로브가 실패하기 시작하면 새 버전의 롤아웃이 효과적으로 차단한다.

### 버전 v3가 완전히 롤아웃되는 것을 방지하기 위한 레디니스 프로브 정의

- 이번엔 레디니스 프로브를 가진 디플로이먼트를 정의해보자.

```json
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubia-deployment
spec:
  replicas: 3
  minReadySeconds: 10
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0 # deployment가 파드를 하나씩 교체하도록 
  selector:
    matchLabels:
      app: kubia
  template:
    metadata:
      labels:
        app: kubia
    spec:
      containers:
        - image: khk9346/kubia:v3
          name: nodejs
          readinessProbe:
            periodSeconds: 1
            httpGet:
              path: /
              port: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: kubia-service
spec:
  type: LoadBalancer
  selector:
    app: kubia
  ports:
  - port: 80
    targetPort: 8080
```

- 위의 yaml파일을 deployment로 하고 롤아웃 상태를 확인해보면 아래와 같이 나오며, 파드가 새로 생성되지 않는다.

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/62b02aaf-f32b-423a-b255-f412a5bb03a7)

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/05a7668a-4f37-43dc-ba9f-1ed9f00438a5)

- maxUnavailable을 0으로 설정했기 때문에 기존의 파드는 삭제하지 않고, 파드를 생성하여 배포하는 방식인데 파드또한 생성되지 않기 때문에 배포가 중단된 것이다.
- 하지만 롤아웃이 계속 Waiting에 걸렸다.

### 롤아웃 데드라인 설정

- 기본은 10분동안 진행되지 않으면 롤아웃이 실패된 것으로 간주된다.
- progressDealineSeconds를 사용하면 데드라인 시간을 설정할 수 있으니 기억해두자.