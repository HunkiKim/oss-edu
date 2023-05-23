# 컨테이너화된 애플리케이션 설정
- 애플리케이션에 설정을 넘겨주는 방법으로 환경변수가 널리 쓰인다.
- 마찬가지로 설정 데이터를 저장하는 쿠버네티스 리소스를 ConfigMap이라고 한다.
- 일반적으로 아래의 방법을 통해 애플리케이션 구성을 할 순 있다.
  - 컨테이너에 명령줄 인수 전달
  - 각 컨테이너를 위한 사용자 정의 환경변수 지정
  - 특수한 유형의 볼륨을 통해 설정 파일을 컨테이너에 마운트
- 하지만 민감정보도 있는데 이를 위해 Secret이라는 오브젝트도 존재한다.

# 컨테이너에 명령줄 인자 전달
- 아래처럼 쿠버네티스에선 명령과 인자를 지정하는 파드를 만들 수 있다.
```yaml
kind: Pod
spec:
  containers:
  - image: some/image
    command: ["/bin/command"]
    args: ["arg1", "arg2", "arg3"]
```
- command(ENTRYPOINT) : 컨테이너 안에서 실행되는 실행파일
- args(CMD) : 실행파일에 전달되는 인자

# 컨테이너의 환경변수 설정
- 환경변수로 fortune 이미지 안에 간격을 설정할 수 있도록 만들자.

**fortuneloop.sh**
```shell
#!/bin/bash
trap "exit" SIGINT
echo Configured to generate new fortune every $INTERVAL seconds
mkdir -p /var/htdocs
while :
do
    echo $(date) Writing fortune to /var/htdocs/idnex.html
    /usr/games/fortune > /var/htdocs/index.html
    sleep $INTERVAL
done
```
**fortune-pod.yaml**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: fortune
spec:
  containers:
  - image: khk9346/fortune # html 생성 컨테이너 
    env:
    - name: INTERVAL
      value: "20"
    name: html-generator
    volumeMounts:
    - name: html
      mountPath: /var/htdocs
  - image: nginx:alpine # nginx 읽기 전용 마운트
    name: web-server
    volumeMounts:
    - name: html
      mountPath: /usr/share/nginx/html
      readOnly: true
    ports:
    - containerPort: 80
      protocol: TCP
  volumes: # html 단일 emptyDir 볼륨을 위의 컨테이너 두 개에 마운트
  - name: html
    emptyDir: {}

```

- 컨테이너에 종속이 아닌 파드에 종속된 환경변수가 있을 필요가 있다.

# Configmap으로 설정 분리
<img width="310" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/59fdb543-5927-41eb-a46c-8693a936adcd">

- 애플리케이션 구성의 요점은 환경에 따라 다르거나 자주 변경되는 설정 옵션을 애플리케이션 소스 코드와 별도로 유지하느 ㄴ것이다.

## 컨피그맵 소개
- 설정 옵션을 별도 오브젝토로 분리 가능하다. 위의 사진처럼 말이다. 
- 컨피그맵은 짧은 문자열에서 전체 설정 파일에 이르는 값을 가진 키/쌍으로 구성된 맵이다.
- 애플리케이션은 컨피그맵을 직접 읽거나 존재도 몰라돋 ㅚㄴ다. 대신 맵의 내용은 컨테이너의 환경변수 또는 볼륨 파일로 전달된다.
  - $(ENV) 와 같은 형식으로 명령줄 인수에서 참조할 수 있기 때문에, 프로세스 명령줄 인자로도 참조 가능

<img width="581" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/e7e7b14b-1387-4799-9d23-ee65a3fb6311">

위처럼 네임스페이스마다 다른 컨피그맵을 구성해 애플리케이션에 주입할 수 있다.

## 컨피그맵 생성
- yaml 없이도 생성이 가능하다.
```shell
k create configmap fortune-config --from-literal=sleep-interval=25
```
```shell
$ k create configmap fortune-config --from-literal=sleep-interval=25
configmap/fortune-config created
$ k get configmap
NAME               DATA   AGE
kube-root-ca.crt   1      11h
fortune-config     1      7s
```

- 하지만 유지보수를 위해선 결국 yaml이 필요하다.
```yaml
apiVersion: v1
data:
  sleep-interval: "25"
kind: ConfigMap
metadata:
  creationTimestamp: "2023-05-22T19:30:30Z"
  name: fortune-config
  namespace: default
  resourceVersion: "5980"
  uid: 3b858240-cab2-4d3c-8dbb-98888eabe17a
```

## 컨피그맵 항목을 컨테이너에 전달
```yaml
spec:
  containers:
  - image: khk9346/fortune # html 생성 컨테이너 
    env:
    - name: INTERVAL
      valueFrom: # 고정값이 아닌 컨피그맵 키에서 값을 가져와 초기화
        configMapKeyRef: # fortune-config에서 sleep-interval 키의 값을 가져옴
          key: sleep-interval
          name: fortune-config
```

## 컨피그맵의 모든 항목을 한 번에 환경변수로 전달
- 만약 FOO, BAR, FOO-BAR라는 세 개의 키를 갖는 컨피그맵을 생각해보자.
```yaml
spec:
  containers:
  - image: some-image
    envFrom: # env대신 envFrom
    - prefix: CONFIG_ ## 모든 환경변수는 CONFIG_ 접두사
    configMapRef:
      name: my-config-map
```
- 위의 결과는 CONFIG_FOO와 COFIG_BAR가 존재한다.
- FOO-BAR는 대시를 가지고 있어 올바른 환경변수 이름이 아니라 변환하지 않는다.

## 컨피그맵 볼륨을 사용해 컨피그맵 항목을 파일로 노출
- 환경변수나 명령줄 인자로 설정 옵션을 전달하는건 짧은 변숫값에 대해 사용된다.
- 하지만 컨피그맵은 모든 설정 파일들을 포함할 수 있다.
  - 컨피그맵 볼륨음 통해 가능하다.
<img width="391" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/31c46bfa-9ba1-47d9-a6d7-96fc1d4481e0">
- 예를들어 아래의 nginx config 파일이 있다고 해보자.
```nginx
server {
    listen 80;
    server_name www.kubia-example.com;

    gzip on;
    gzip_types text/plain application/xml;

    location / {
        root /usr/share/nginx/html;
        index index.html index.html;
    }
}
```
이제 아래의 명령어로 컨피그맵을 실행해보자.
```
$ k create configmap fortune-config --from-file=configmap-files
```
- configmap-files라는 폴더엔 my-nginx-config.conf와 sleep-interval 파일이 있다.
  - sleep-interval 파일엔 252525252525252525252525252525252525252525252525
가 들어가 있다.

실행하면 아래와 같은 결과가 나온다.
```
apiVersion: v1
data:
  my-nginx-config.conf: |-
    server {
        listen 80;
        server_name www.kubia-example.com;

        gzip on;
        gzip_types text/plain application/xml;

        location / {
            root /usr/share/nginx/html;
            index index.html index.html;
        }
    }
  sleep-interval: |
    252525252525252525252525252525252525252525252525
kind: ConfigMap
metadata:
  creationTimestamp: "2023-05-22T19:53:12Z"
  name: fortune-config
  namespace: default
  resourceVersion: "6698"
  uid: 19393bc8-e116-47e3-a64e-7c8aaff7174f
```

### 볼륨안에 컨피그맵 사용하기
- Nginx는 기본적으로 /etc/nginx/nginx.conf파일을 읽는다.
  - 기본 설정 파일은 /etc/nginx/conf.d/ 디렉터리 안에 있는 모든 .conf 파일을 포함한다.
- 아래처럼 설정하면 configVolume을 마운트 할 수 있다.

**fortune-pod.yaml**
```
apiVersion: v1
kind: Pod
metadata:
  name: fortune
spec:
  containers:
  - image: khk9346/fortune # html 생성 컨테이너 
    env:
    - name: INTERVAL
      valueFrom: # 고정값이 아닌 컨피그맵 키에서 값을 가져와 초기화
        configMapKeyRef: # fortune-config에서 sleep-interval 키의 값을 가져옴
          key: sleep-interval
          name: fortune-config
    name: html-generator
    volumeMounts:
    - name: html
      mountPath: /var/htdocs
  - image: nginx:alpine # nginx 읽기 전용 마운트
    name: web-server
    volumeMounts:
    - name: config
      mountPath: /etc/nginx/conf.d # 컨피그맵 볼륨을 마운트하는 위치
      readOnly: true
    - name: html
      mountPath: /usr/share/nginx/html
      readOnly: true
    ports:
    - containerPort: 80
      protocol: TCP
  volumes: # html 단일 emptyDir 볼륨을 위의 컨테이너 두 개에 마운트
  - name: html
    emptyDir: {}
  - name: config
    configMap:
      name: fortune-config

```

결과는 아래처럼 잘 나온다.
```shell
curl -H "Accept-Encoding: gzip" -I localhost:8080
HTTP/1.1 200 OK
Server: nginx/1.23.4
Date: Mon, 22 May 2023 20:00:40 GMT
Content-Type: text/html
Last-Modified: Mon, 22 May 2023 20:00:31 GMT
Connection: keep-alive
ETag: W/"646bc9df-4f"
Content-Encoding: gzip
```

### 마운트된 컨피그맵 볼륨 내용 살펴보기
```
$ k exec fortune -c web-server ls /etc/nginx/conf.d
my-nginx-config.conf
sleep-interval
```
- sleep-interval도 포함되어 있어 이상하다.

### 볼륨에 특정 컨피그맵 항목 노출
아래와 같이 수정하면 볼륨엔 my-nginx-conf.conf만 존재하며, interval-sleep은 환경변수로 들어간다.
```yaml
volumes: # html 단일 emptyDir 볼륨을 위의 컨테이너 두 개에 마운트
  - name: html
    emptyDir: {}
  - name: config
    configMap:
      name: fortune-config
      items:
      - key: my-nginx-config.conf
        path: gzip.conf
      
```
- path의 파일이름으로 들어간다.
```shell
$ k exec fortune -c web-server ls /etc/nginx/conf.d
gzip.conf
```

- 위처럼 설정하면 컨테이너 이미지 자체에 있던 /etc/nginx/conf.d 디렉터리 안에 저장된 파일이 숨겨지고 설정한 파일만 들어가게 되어 문제가 생길 수 있다.
  - 경로가 etc였다면 치명적인 오류가 발생한다.

### 디렉터리 안에 다른 파일을 숨기지 않고 개별 컨피그맵 항목을 파일로 마운트
- 전체 볼륨을 마운트하는게 아니라 subPath를 쓰면 된다.
```yaml
volumeMounts:
    - name: config
      mountPath: /etc/someconfig.conf # 파일을 마운트
      subPath: myconfig.conf # 전체 볼륨이 아니라 myconfig.conf항목만 마운트
```
하지만 업데이트 관련 큰 결함이 있다.

### 컨피그맵 볼륨 안에 있는 파일 권한 설정
- 기본적으로 컨피그맵 볼륨의 모든 파일 권한은 644(-rw-r-r--)로 설정된다.
  - defaultMode를 정의해주면 모든 파일권한을 바꿀 수 있다.
```yaml
volumes:
- name: config
  configMap:
    name: fortune-config
    defaultMode: "6600"
```

## 애플리케이션 재시작 없이 설정 업데이트
- 컨피그맵을 볼륨으로 노출하면 파드를 다시 만들거나 컨테이너를 다시 시작할 필요 없이 설정을 업데이트 가능하다.
- 컨피그맵을 업데이트하면, 이를 참조하는 모든 볼륨의 파일이 업데이트된다.
```yaml
$ k edit configmap fortune-config
configmap/fortune-config edited
```
- 위의 명령어를 통해 gzip을 off 하고 편집기를 닫으면 볼륨의 실제 파일도 업데이트 된다.
- 하지만 nginx를 확인해보면 업데이트가 되지않는다.
- 아래의 명령어를 통해 설정파일을 다시 로드해야한다.
```
$ k exec fortune -c web-server -- nginx -s reload
```
그러면 이제 gzip encoding을 하지 않는다.
```shell
curl -H "Accept-Encoding: gzip" -I localhost:8080
HTTP/1.1 200 OK
Server: nginx/1.23.4
Date: Mon, 22 May 2023 23:32:44 GMT
Content-Type: text/html
Content-Length: 72
Last-Modified: Mon, 22 May 2023 23:32:37 GMT
Connection: keep-alive
ETag: "646bfb95-48"
Accept-Ranges: bytes
```

- 하지만 모든 인스턴스가 동기적으로 업데이트되지 않기 때문에 주의하자

# 시크릿으로 민감한 데이터를 컨테이너에 전달
지금까지 컨테이너에 전달한 정보는 민감하지 않은 데이터였다. 하지만 aws 정보같은 민감한 정보들은 보안이 유지돼야 한다.

## 시크릿 소개
- 키-값 쌍을 가진 맵으로 컨피그맵과 비슷하다.
- 또한 같은 방식으로 환경변수로 시크릿 항목을 컨테이너에 전달하거나 시크릿 항목을 볼륨 파일로 노출할 수 있다.
- 시크릿은 물리저장소가 아닌 항상 메모리에 둔다. 
  - 물리 저장소는 시크릿을 삭제한 후에도 디스크를 완전히 삭제하는 작업이 필요하기 때문이다.
- 마스터 노드(etcd)에는 시크릿을 암호화되지 않은 형식으로 저장하므로, 시크릿에 저장한 민감한 데이터를 보호하려면 마스터 노드를 보호해야 한다.
  - 그러기 위해선 권한 없는 사용자가 API 서버를 이용하지 못하게 하는 것도 포함된다.
  - 파드를 만들 수 있는 사람은 누구나 민감 데이터에 접근하는 것이 가능하기 때문이다.
- 1.7버전부턴 etcd가 시크릿을 암호화된 형태로 저장해 시스템을 좀 더 안전하게 만든다.
- 시크릿을 사용하는 기준은 아래 두가지가 잇다.
  - 민감하지 않고, 일반 설정 데이터는 컨피그맵을 사용한다.
  - 본질적으로 민감한 데이터는 시크릿을 사용해 키 아래에 보관하는 것이 필요하다. 만약 설정 파일이 둘 다 가지고 있다면 해당 파일을 시크릿 안에 저장해야 한다.

## 시크릿 생성
https 트래픽을 제공하게 개선해보자.
```shell
$ openssl genrsa -out https.key 2048
$ openssl req -new -x509 -key https.key -out https.cert -days 3650 -subj /CN=www.kubia-example.com
```
위는 인증서와 개인키를 만드는 명령어다.

이제 secret을 만든다.
```
$ k create secret generic fortune-https --from-file=https.key
secret/fortune-https created
```
- 컨피그맵을 작성하는 것과 크게 다르지 않다. fortune-https 이름을 가진 generic 시크릿을 생성한다.
  - 시크릿의 세 가지 유형에는 도커 레지스트리를 사용하기 위한 docker--registry, TLS 통신을 위한 tls, generic이 있다.
- --from-file=fortune-https 옵션을 이용해 디렉터리 전체를 포함할 수 있다.
- 시크릿과 컨피그맵은 매우 큰 차이가 잇다. yaml파일을 보자.
```yaml
data:
  https.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdDE3UG1INUNteURVak95dU1JNXkzQUJYdm1vRnRBQ3Q4OVQ3ait2enVST0NqTlFaCkRYNURXa3NNbHp0aWZHa0tHTnlwSDZTSFFWcTNNamZ5M2I4QUN6ZlFUTnFpRWs2YjRMYUF3V3NvZmE1SGQxRDEKRUVlU0VXTlQxSXpkN1VOUkV6ZVlMT081WXlVb2hHSmRkUCtDTllQem5leVBUN0dGc2lRSmdTM282UFRXWEV4bQpyOVBtcEVUOUFXN3FaY1U1dDZOQzhyQytRc1ZJcXhNbERQdGFLNzQ5Z...
```
- 시크릿 항목의 내용은 Base64 인코딩 문자열로 표시된다.
  - 스크릿 안의 manifest를 다룰땐 인코딩, 디코딩을 매번 읽을때마다 해줘야 한다.
- Base64를 사용한 이유는 간단하다. 시크릿 항목에 일반 텍스트뿐만 아니라 바이너리 값도 담을 수 잇기 때문이다.
- 하지만 모든 민감한 데이터가 바이너리 형태는 아니기 때문에 stringData 필드로 시크릿 값을 설정할 수 있다.
  - stringData 필드는 쓰기 전용이다. 값을 설정할때만 사용할 수 있다.

## 파드에 시크릿 사용
HTTPS를 사용하도록 fortune-config 컨피그맵을 수정해서 사용할 수 있게 해보자.
```shell
$ k edit configmap fortune-config
configmap/fortune-config edited
```
- 그리고 마운트를 하자
**fortune-pod-https.yaml**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: fortune-https
spec:
  containers:
  - image: khk9346/fortune # html 생성 컨테이너 
    env:
    - name: INTERVAL
      valueFrom: # 고정값이 아닌 컨피그맵 키에서 값을 가져와 초기화
        configMapKeyRef: # fortune-config에서 sleep-interval 키의 값을 가져옴
          key: sleep-interval
          name: fortune-config
    name: html-generator
    volumeMounts:
    - name: html
      mountPath: /var/htdocs
  - image: nginx:alpine # nginx 읽기 전용 마운트
    name: web-server
    volumeMounts:
    - name: config
      mountPath: /etc/someconfig.conf # 파일을 마운트
      subPath: myconfig.conf # 전체 볼륨이 아니라 myconfig.conf항목만 마운트
      readOnly: true
    - name: html
      mountPath: /usr/share/nginx/html
      readOnly: true
    - name: certs
      mountPath: /etc/nginx/certs/
      readOnly: true
    ports:
    - containerPort: 80
    - containerPort: 443
      protocol: TCP
  volumes: # html 단일 emptyDir 볼륨을 위의 컨테이너 두 개에 마운트
  - name: html
    emptyDir: {}
  - name: config
    configMap:
      name: fortune-config
      items:
      - key: my-nginx-config.conf
        path: gzip.conf
  - name: certs # 여기에 시크릿 마운트한다.
    secret:
      secretName: fortune-https
```
```shell
Volumes:
  html:
    Type:       EmptyDir (a temporary directory that shares a pod's lifetime)
    Medium:
    SizeLimit:  <unset>
  config:
    Type:      ConfigMap (a volume populated by a ConfigMap)
    Name:      fortune-config
    Optional:  false
  certs:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  fortune-https
    Optional:    false
```
이렇게 들어가게 된다.

## 이미지를 가져올 때 사용하는 시크릿 이해
애플리케이션에 시크릿을 전달하고 그 안에 있는 데이터를 사용하는 방법은 배웠다. 하지만 쿠버네티스에서 자격증명 전달이 필요한 경우가 있다. (예를들어 프리이빗 컨테이너 이미지 레지스트리를 사용하려는 경우)

- 이때에도 시크릿을 통해 이뤄진다.
- 지금까진 모두 공개 이미지 레지스트리라 필요 없었지만 현업에선 대부분 필요하다.

### Docker Hub에서 프라이빗 이미지 사용
- 도커 레지스트리 자격증명을 가진 시크릿 생성
- 파드 매니페스트 안에 imagePullSecrets 필드에 해당 시크릿 참조

1. **도커 레지스트리 인증을 위한 시크릿 생성**

방금한 generic 시크릿과 다르지 않다. 단지 docker-registry 형식을 가진 시크릿을 만들 뿐이다.
```shell
$ k create secret docker-registry mydockerhubsecret --docker-username=khk9346 --docker-password=rnfjddl12! --docker-email=gnsrl76@naver.com
secret/mydockerhubsecret created
```
```shell
$ k describe secrets mydockerhubsecret
Name:         mydockerhubsecret
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  kubernetes.io/dockerconfigjson

Data
====
.dockerconfigjson:  150 bytes
```

### 모든 파드에서 이미지를 가져올 때 사용할 시크릿을 모두 지정할 필요는 없다.
- 서비스어카운트에 추가해 모든 파드에 자동으로 추가될 수 있는 방법을 나중에 배울 것이다.