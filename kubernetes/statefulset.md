# 스테이트풀 소개
- 디플로이먼트와 유사하게, 스테이트풀셋은 동일한 컨테이너 스펙을 기반으로 둔 파드들을 관리한다. 
- 디플로이먼트와는 다르게, 스테이트풀셋은 각 파드의 독자성을 유지한다. 
  - 이 파드들은 동일한 스펙으로 생성되었지만, 서로 교체는 불가능하다. 다시 말해, 각각은 재스케줄링 간에도 지속적으로 유지되는 식별자를 가진다.
- 스토리지 볼륨을 사용해서 워크로드에 지속성을 제공하려는 경우, 솔루션의 일부로 스테이트풀셋을 사용할 수 있다.
- 스테이트풀셋의 개별 파드는 장애에 취약하지만, 퍼시스턴트 파드 식별자는 기존 볼륨을 실패한 볼륨을 대체하는 새 파드에 더 쉽게 일치시킬 수 있다.

# 스테이트풀 파드 복제하기
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/7eaa3341-e598-4e97-a9b0-d9881663cf11)

- 레플리카셋은 하나의 파드 템플릿에서 여러 개의 파드 레플리카를 생성한다.
- 레플리카는 이름과 IP 주소를 제외하면 서로 동일하다.
- 파드 템플릿이 특정 PVC을 참조하는 볼륨을  포함하면 레플리카셋의 모든 레플리카는 정확히 동일한 PVC을 사용할 것이고 클레임에 바인딩된 동일한 PV을 사용하게 된다.
- 여러 개의 파드 레플리카를 복제하는 데 사용하는 파드 템플릿에는 클레임 참조가 있기 때문에 각 레플리카가 별도의 PVC을 사용하도록 만들 수 없다.
  - 즉 분산 데이터 저장소를 사용할 때 레플리카셋은 사용할 수 없다.

## 개별 스토리지를 갖는 레플리카 여러 개 실행하기
- 그러면 개별 스토리지를 갖는 레플리카를 여러개 실행하면 되지 않을까?
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/3f5d54e3-eab7-4610-8b4d-5b18d4e8bf8f)
- 수동으로 파드를 생성해 다른 PVC를 갖게하는 방법이 있다.
  - 하지만 레플리카셋이 파드를 감시하지 않기 때문에 파드 관리가 어렵다.
- 그러면 파드 인스턴스별로 하나의 레플리카셋을 사용하는건 어떨까?
  - 가능은 하다. 하지만 단일 레플리카셋보다 설정이 번거롭다. 

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/4bccb231-2426-49df-8309-d7df4c28ee3c)
- 방법은 동일 볼륨을 여러 개의 디렉토리로 사용하는 방법이다.
  - 각각의 레플리카셋은 같은 PVC를 참조하지만, 실제 디스크에서 다른 디렉토리에 데이터를 분산해 저장할 수 있다.
  - 하지만 이 방법도 어떤 디렉토리를 각 파드가 선택해야할지, 공유 스토리지 볼륨에 병목현상 등의 문제가 있다.

## 각 파드에 안정적인 아이덴티티 
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/d8287a39-9f7a-4b9a-b520-f580f0660532)

- 스토리지 외에도 특정 클러스터 애플리케이션은 각 인스턴스에서 장시간 지속되는 안정적인 아이덴티티들 필요로 한다.   
  - 파드는 계속 삭제되고 생성되고를 반복한다. 
  - 레플리카셋이 파드를 교체하면 새 파드가 갖는 스토리지 볼륨은 동일해도 완전히 새로운 호스트 이름과 IP를 가진다.
  - 특정 애플리케이션에선 이전 인스턴스의 데이터를 가지고 시작할 때 새로운 네트워크 아이덴티티로 인해 문제가 발생할 수 있다.
- 특정 애플리케이션이 안정적인 네트워크 아이덴티티를 요구하는 이유는 분산 스테이트풀 애플리케이션에서 필요하기 때문이다.
- 이 문제를 해결하기 위해선 각 개별 멤버에게 전용 서비스를 생성해 안정적인 네트워크 주소를 제공해야 한다.

# 스테이트풀셋 이해하기
- 위의 복잡한 해결책에 의존하지 않고 해결할 수 있는 방법이 있다. 바로 레플리카셋을 사용하는 대신 StatefulSet 리소스를 생성하는 것이다. 
- 스테이트풀셋은 애플리케이션의 인스턴스가 각각 안정적인 이름과 상태를 가지며 개별적으로 취급돼야 하는 애플리케이션에 알맞게 만들어졌다.

## 스테이트풀셋과 레플리카셋 비교하기
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/0964256d-66e5-4908-9652-3b0ffbdba8d1)
- 스테이트풀셋의 목적을 이해하는 방법은 비교를 하는게 가장 빠르다.
- 가축과 애완동물의 비유가 가장 적절하다.
  - 각각의 동물에 이름을 부여하고 개별적으로 관리하는 경우는 스테이트풀셋과 유사하다. 만약 애완동물이 죽는다면 바로 살 수 없고, 다른 사람들도 바로 알아차린다. 다시 구하려고 해도 이전과 완전히 같은 상태를 찾아야 한다. 스테이트풀셋도 새 인스턴스가 이전 인스턴스와 완전히 같은 상태와 아이덴티티를 가진다.
    - 같은 이름, 네트워크 아이덴티티, 상태로 다른 노드에서 다시 생성되어야 한다. 즉 파드가 아이덴티티와 상태를 유지하면서 다시 스케줄링되게 한다.
  - 하지만 가축처럼 따로 이름을 부여하지 않고, 병든 가축을 없애고 새롭게 다시 가축을 추가하는것처럼 관리도 가능하다. 레플리카셋이 이와 동일하다.
- 스테이트풀셋은 자체 볼륨 세트(퍼시스턴트 상태, 스토리지)를 가진다. 또한 파드 인스턴스가 예측가능한 아이덴티티를 가진다.

### 거버닝 서비스 소개
- 스테이트풀셋은 거버닝 헤드리스 서비스(governing headless service)를 생성해 각 파드에게 실제 네트워크 엔티티를 제공한다. 
  - 이 서비스를 통해 각 파드는 자체 DNS 엔트리를 가지며 클러스터의 피어 혹은 클러스터의 다른 클라이언트가 호스트 이름을 통해 파드의 주소를 지정할 수 있다.
  - StatefulSet은 FQDN(Fully Qualified Domain Name, a-0.foo(namespace).default.svc.cluster.local)을 통해 접근 가능하다. ReplicaSet은 물가능하다.
  - 또한 FQDN 도메인의 SRV 레코드를 조회해 모든 스테이트풀셋의 파드 이름을 찾는 목적으로도 DNS를 사용할 수 있다.
    - SRV레코드는 서비스의 위치 정보를 제공하기위해 사용되는 DNS 리소스 레코드(도메인 이름과 관련된 정보)이다.


### 스테이트풀셋 교체
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/1ab2d384-ee8e-40d9-8632-6673f828279e)
위처럼 파드가 교체될 때 새로운 인스터스로 교체되는 것은 동일하지만, 레플리카셋과 달리 교체된 파드는 사라진 파드와 동일한 이름을과 호스트 이름을 갖는다.

### 스테이트풀셋 스케일링
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/bd4c0011-67c5-4900-ab4c-8e5d6d6b6e3b)

- 레플리카셋처럼 랜덤으로 파드가 스케일다운하는 것이 아닌 항상 가장 높은 서수 인덱스를 먼저 제가한다. 따라서 스케일 다운의 영향을 예측할 수 있다.
- 특정 스테이트풀 애플리케이션은 빠른 스케일 다운을 잘처리하지 못하기 때문에 스테이트풀셋은 한 시점에 하나의 파드 인스턴스만 스케일  다운한다. 분산 데이터 스토어를 예를 들면 여러 개 노드가 동시에 다운되면 데이터를 잃을 수 있다.


## 각 스테이트풀 인스턴스에 안정적인 전용 스토리지 제공하기
- 아이덴티티를 갖는 예제들은 보았지만 스토리지는?
  - 각 스테이트풀 파드 인스턴스는 자체 스토리지를 사용할 필요가 있고 스테이트풀 파드가 다시 스케줄링되면, 새로운 인스턴스는 동일한 스토리지에 연결돼야 한다.

### 볼륨 클레임 템플릿과 파드 템플릿을 같이 구성
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/5e8e7ed2-88db-43c3-9393-44222bd5741f)

위의 그림처럼 파드템플릿과 볼륨클레임템플릿을 정의하면 각 파드마다 PVC도 같이 셍성한다. 이런 이유로 각 파드와 함께하는 PVC을 복제하는 하나 이상의 볼륨 클레임 템플릿을 가질 수 있다.

- 미리 관리자가 프로비저닝 하거나 동적 프로비저닝 하여 PV를 생성할 수 있다.

### PVC 생성과 삭젱의 이해
- 스테이트풀셋을 하나 스케일 업하면 두 개 이상의 API 오브젝트(Pod와 PVC)가 생성된다.
- 하지만 스케일 다운할 때는 파드만 삭제하고 PVC는 남겨둔다.
  - 바인딩된 PV가 재활용되거나 삭제돼 콘텐츠가 손실되기 때문이다.
- 위와같은 이유로 기반 PV을 해제하려면 PVC을 수동으로 삭제해야 한다.

### 동일 파드의 새 인스턴스에 PVC 다시 붙이기
스케일 다운 이후 PVC가 남아 있다는 사실은 이후에 스케일 업 할때 PV에 바인딩된 동일한 PVC를 다시 연결할 수 있고, 새 파드에 해당 콘텐츠가 연결됨을 의미한다. 즉 파드를 롤백 할 수 있다는 것을 말한다.

## 스테이트풀셋 보장 이해하기
![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/bda5f278-f761-4f48-a2bb-3869adb7a36f)

지금까지는 아이덴티티와 스토리지가 보장됨을 말했는데 또 한가지 보장하는 것이 있다.

### 안정된 아이덴티티와 스토리지의 의미
지금까진 쿠버네티스가 파드 상태를 확신 할 수 있을때를 기준으로 하였다. 하지만 쿠버네티스가 동일한 아이덴티티를 가지는 교체 파드를 생성하면 두 개 인스턴스가 동일한 아이덴티티로 시스템에서 실행할 수 있다. 그러면 두 인스턴스가 같은 스토리지에 바인딩되고 두 프로세스가 동일한 아이덴티티로 같은 파일을 쓰려 한다.

**즉 스테이트풀셋은 동일한 아이덴티티로 실행되지 않고 동일한 PVC에 바인딩되지 않도록 보장한다.** 이는 최대 하나의 의미(at-most-one semantics)를 보장해야 한다는 것이다.

즉, 스테이트풀셋은 교체 파드를 생성하기 전에 파드가 더 이상 실행 중이지 않는다는 점을 절대적으로 확신해야 한다.

# 스테이트풀셋 사용하기
- POST 요청마다 요청의 body로 쓴 데이터를 /var/data/kubia.txt파일에 쓰는 app.js가 있다고 해보자.
- GET 요청은 호스트 이름과 저장된 데이터를 반환한다.

## 스테이트풀셋을 통한 애플리케이션 배포
- 애플리케이션 배포를 위해 추가적인 오브젝트가 더 필요하다.
  - 데이터파일 저장을 위한 PV
  - 스테이트풀셋에 필요한 거버닝 서비스
  - 스테이트 풀셋 자체

### PV 생성하기
- 예제는 PV를 생성하도록 되어있지만 현재 GCE를 동적 프로비저닝을 사용할 수 있다. storageclass(storage-provisioner)는 rancher에 있는걸 사용합니다.

- 거버닝 서비스는 아래와 같다.
  - clusterIP 필드를 None으로 하면 헤드리스 서비스가 된다. 
```yaml
apiVersion: v1
kind: Service
metadata:
  name: kubia
spec:
  clusterIP: None
selector:
  app: kubia
ports:
  - name: http
    port: 80
```
- 스테이트풀셋 매니페스트는 아래와 같다.
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kubia
spec:
  serviceName: kubia
  replicas: 2
  selector:
    matchLabels:
      app: kubia # has to match .spec.template.metadata.labels
  template:
    metadata:
      labels:
        app: kubia
    spec:
      containers:
      - name: kubia
        image: luksa/kubia-pet
        ports:
        - name: http
          containerPort: 8080
        volumeMounts:
        - name: data
          mountPath: /var/data
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      resources:
        requests:
          storage: 1Mi
      accessModes:
      - ReadWriteOnce
```

### 스테이트풀셋 생성하기
```yaml
$ k apply -f kubia-statefulset.yaml
statefulset.apps/kubia created
$ k get pods
NAME            READY   STATUS    RESTARTS   AGE
fortune-https   2/2     Running   0          133m
kubia-0         0/1     Pending   0          6s
```

- 파드를 조회하면 이전처럼 동시에 생성하는 것이 아닌, 하나씩 생성하는걸 확인할 수 있다. 
  - 첫 번째 파드가 생성되고 준비가 완료되면 두 번째 파드가 생성된다. 
  - 특정 클러스터된 스테이트풀 애플리케이션은 두 개 이상의 멤버가 동시에 생성되면 레이스 컨디션에 빠질 가능성이 있기 때문에 스테이트풀셋은 이와같이 동작한다.
- describe를 통해 확인하면 manifest에 명시된대로 동작하는것을 확인할 수 있다.
```yaml
Mounts:
      /var/data from data (rw)
Volumes:
  data:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  data-kubia-0
```

- 현재 노드들은 실행중이지만 헤드리스 서비스를 생성했으므로 서비스를 통해 파드와 통신할 수 없다.
  - headless service는 k8s에서 제공하는 서비스 유형 중 하나입니다.
  - clusterIP: None으로 설정해 클러스터 IP를 할당받지 않습니다.
  - 이로 인해 외부에선 직접 접근할 수 없으며, 클러스터 내부의 파드의 IP 주소로 직접 통신 가능합니다.
- 개별 파드와 직접 연결해야 한다
  - kubectl proxy를 통해 API 서버와 통신하여 localhost:8001에 요청을 하자.
```shell
$ curl localhost:8001/api/v1/namespaces/default/pods/kubia-0/proxy/
You've hit kubia-0
Data stored on this pod: No data posted yet
```
잘 나온다. 이제 POST 요청을 통해 데이터를 볼륨에 저장하고, 다시 get요청으로 볼륨에 잘 저장되었는지 확인해보자.
```shell
$ curl -X POST -d "hello world" localhost:8001/api/v1/namespaces/default/pods/kubia-0/proxy/
Data stored on pod kubia-0

$ curl localhost:8001/api/v1/namespaces/default/pods/kubia-0/proxy/
You've hit kubia-0
Data stored on this pod: hello world
```
- 잘 되는것을 확인할 수 있다.

### 스테이트풀 파드를 헤드리스가 아닌 일반적인 서비스로 노출하기
10장 마지막 부분으로 넘어가기 전에 보통 클라이언트는 파드를 직접 연결하는 것보다 서비스를 통해 연결하므로 non-headless 서비스를 파드 앞에 추가해보자.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: kubia-public
spec:
  selector:
    app: kubia
  ports:
  - port: 80
    targetPort: 8080
```
```shell
$ curl localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
You've hit kubia-0
Data stored on this pod: hello world
$ curl localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
You've hit kubia-1
Data stored on this pod: No data posted yet
```
- 서비스는 랜덤으로 요청을 한다.

# 스테이트풀셋의 피어 디스커버리
- 클러스터된 애플리케이션의 중요한 요구 사항은 피어 디스커버리(클러스터의 다른 멤버를 찾는 기능)다.
- 스테이트풀셋의 각 멤버는 모든 다른 멤버를 쉽게 찾을 수 있어야 한다.
  - API서버와 통신해 찾을 수 있지만, 쿠버네티스의 목표는 애플리케이션을 쿠버네티스 독립적으로 유지하여 기능을 노출하는 것이다.
  - 따라서 대화없이 찾아야 한다
- DNS의 SRV 레코드를 사용하면 된다

### SRV 레코드 소개
- SRV 레코드는 특정 서비스를 제공하는 서버의 호스트 이름과 포트를 가리키는 데 사용된다.
- 쿠버네티스는 헤드리스 서비스를 뒷받침하는 파드의 호스트 이름을 가리키도록 SRV 레코드를 생성한다.
- dig를 실행해 스테이트풀 파드의 SRV 레코드를 조회할 수 있다.
```shell
$ k run -it srvlookup --image=tutum/dnsutils --rm --restart=Never -- dig SRV kubia.default.svc.cluster.local
```
- 1회용 파드이다.
- ANSWER SECTION에는 헤드리스 서비스를 뒷받침하는 두 개의 파드를 가리키는 두 개의 SRV 레코드를 보여준다. 또한 각 파드는 ADDITIONAL SECTION에 표시된 것처럼 자체 A 레코드를 가진다.
  - A레코드는 도메인 이름을 IPv4주소로 매핑하는데 사용된다.
    - 예를들어 A레코드가 192.0.2.1이면 도메인을 192.0.2.1로 연결 할 수 있게 해준다.
```shell
;; ANSWER SECTION:
kubia.default.svc.cluster.local. 5 IN	SRV	0 50 80 kubia-0.kubia.default.svc.cluster.local.
kubia.default.svc.cluster.local. 5 IN	SRV	0 50 80 kubia-1.kubia.default.svc.cluster.local.

;; ADDITIONAL SECTION:
kubia-1.kubia.default.svc.cluster.local. 5 IN A	10.42.3.9
kubia-0.kubia.default.svc.cluster.local. 5 IN A	10.42.0.9
```

# DNS를 통한 피어 디스커버리
- kubia-public 서비스를 통해 데이터 저장소 클러스터로 연결한 클라이언트가 게시한 데이터는 임의의 클러스터 노드에 전달된다. 
  - 즉 지금 랜덤요청이라는 뜻이다.
  - 스테이트풀셋과 SRV 레코드를 사용해 문제를 해결해보자.
```node
const http = require('http');
const os = require('os');
const fs = require('fs');
const dns = require('dns');

const dataFile = "/var/data/kubia.txt";
const serviceName = "kubia.default.svc.cluster.local";
const port = 8080;


function fileExists(file) {
  try {
    fs.statSync(file);
    return true;
  } catch (e) {
    return false;
  }
}

function httpGet(reqOptions, callback) {
  return http.get(reqOptions, function(response) {
    var body = '';
    response.on('data', function(d) { body += d; });
    response.on('end', function() { callback(body); });
  }).on('error', function(e) {
    callback("Error: " + e.message);
  });
}

var handler = function(request, response) {
  if (request.method == 'POST') {
    var file = fs.createWriteStream(dataFile);
    file.on('open', function (fd) {
      request.pipe(file);
      response.writeHead(200);
      response.end("Data stored on pod " + os.hostname() + "\n");
    });
  } else {
    response.writeHead(200);
    if (request.url == '/data') {
      var data = fileExists(dataFile) ? fs.readFileSync(dataFile, 'utf8') : "No data posted yet";
      response.end(data);
    } else {
      response.write("You've hit " + os.hostname() + "\n");
      response.write("Data stored in the cluster:\n");
      dns.resolveSrv(serviceName, function (err, addresses) { # SRV 레코드 확인
        if (err) {
          response.end("Could not look up DNS SRV records: " + err);
          return;
        }
        var numResponses = 0;
        if (addresses.length == 0) {
          response.end("No peers discovered.");
        } else {
          addresses.forEach(function (item) {
            var requestOptions = {
              host: item.name,
              port: port,
              path: '/data'
            };
            httpGet(requestOptions, function (returnedData) {
              numResponses++;
              response.write("- " + item.name + ": " + returnedData + "\n");
              if (numResponses == addresses.length) {
                response.end();
              }
            });
          });
        }
      });
    }
  }
};

var www = http.createServer(handler);
www.listen(port);
```

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/8e8c3ffe-caef-4138-9667-5bff5a24aa76)

```shell
$ k edit statefulset kubia
```

- 위 명령어를 통해 레플리카 개수를 3으로 늘려본다. 그러면 새로운 파드가 생성되지만 이전에 있던 파드는 업데이트 되지 않는다.
  - 이는 StatefulSet이 ReplicaSet과 유사하여 새로운 복제본을 기반으로 새로운 파드를 생성하고 롤아웃은 수행하지 않기 때문이다.

# 데이터스토어에 생성해보기
```shell
$ curl -X POST -d "The sun is shining" localhost:8001/api/v1/namespace/default/services/kubia-public/proxy/
```

```
$ curl -X POST -d "The sun is shining" localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
Data stored on pod kubia-0
$ curl -X POST -d "The sun is shining2" localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
Data stored on pod kubia-2
$ curl -X POST -d "The sun is shining3" localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
Data stored on pod kubia-0
$ curl -X POST -d "The sun is shining3" localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
Data stored on pod kubia-0
$ curl -X POST -d "The sun is shining3" localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
Data stored on pod kubia-0
$ curl -X POST -d "The sun is shining3" localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
Data stored on pod kubia-1
```

```
curl localhost:8001/api/v1/namespaces/default/services/kubia-public/proxy/
You've hit kubia-2
Data stored in the cluster:
- kubia-0.kubia.default.svc.cluster.local: The sun is shining3
- kubia-2.kubia.default.svc.cluster.local: The sun is shining2
- kubia-1.kubia.default.svc.cluster.local: The sun is shining3
```
모든 피어들을 찾았다. 즉 애플리케이션 인스턴스가 피어를 디스커버리하고 수평 확장을 쉽게 처리할 수 있는지 보여줬다.