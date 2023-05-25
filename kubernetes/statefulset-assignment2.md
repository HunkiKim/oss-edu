# Deployment의 데이터

<img width="1093" alt="image" src="https://github.com/HunkiKim/Mantech-Edu/assets/66348135/9dc69082-c7af-425a-aa31-c6ffe7e07167">

```shell
$ k get pvc                 
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
mysql-pvc   Bound    pvc-2e46d161-83a0-4257-9289-5ba1ac8d48f4   1Gi        RWO            local-path     2m18s
$ kubernetes [kubectl-assignment2] ⚡  k get pv 
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
pvc-2e46d161-83a0-4257-9289-5ba1ac8d48f4   1Gi        RWO            Delete           Bound    default/mysql-pvc   local-path              9s
$ kubernetes [kubectl-assignment2] ⚡  k get storageclass
NAME                   PROVISIONER             RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
local-path (default)   rancher.io/local-path   Delete          WaitForFirstConsumer   false                  6d2h
```
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: local-path
  resources:
    requests:
      storage: 1Gi

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:latest
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: "1234"
          ports:
            - containerPort: 3306
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql
      volumes:
        - name: mysql-data
          persistentVolumeClaim:
            claimName: mysql-pvc

```
```mysql
mysql> select * from users;
+------+
| age  |
+------+
|    1 |
+------+
1 row in set (0.00 sec)

```
- 위처럼 pod 내부에서 데이터베이스를 만들어 데이터를 넣어둔다.
- 그러면 갑자기 나머지 컨테이너에서 오류가 발생한다.
  - 각각 하나의 디플로이먼트에 대해 PVC를 각각 하나씩 만들어 줘야하기 때문이다.
- 위의 결과는 POD를 삭제하고 다시 열어도 그대로 파일은 유지된다.
  - 하지만 replica를 여러개로 늘리는 순간 PVC를 Pod마다 유지해줘야한다.
- 즉 파드가 여러개로 되는 순간 Deployment에서 replicas를 여러개로 유지할 수 없다.
- 즉 위의 이미지대로의 결과는 나올 수 없으며, Pod별로 직접 PVC를 할당해 주지 않는 이상 불가능하다.
- pvc같은 경우는 모두 연결되어 있지만, mysql에 특정 파일에 lock이 걸려 다른 파드에서 접근할 수 없어 지속적으로 RESTARTS를 하는 현상이 나타난다.

```
$ k get pods
k get pods
NAME                                READY   STATUS    RESTARTS      AGE
kubia-deployment-5fb8568f6-m9xw2    1/1     Running   0             3h17m
kubia-deployment-5fb8568f6-5dwfh    1/1     Running   0             3h17m
kubia-deployment-5fb8568f6-p744s    1/1     Running   0             3h17m
mysql-deployment-75b74b68db-642sn   1/1     Running   0             2m11s
mysql-deployment-75b74b68db-vmkn4   1/1     Running   1 (19s ago)   2m11s
mysql-deployment-75b74b68db-vkftp   1/1     Running   2 (20s ago)   2m11s

$ k get pods
NAME                                READY   STATUS    RESTARTS      AGE
kubia-deployment-5fb8568f6-m9xw2    1/1     Running   0             3h29m
kubia-deployment-5fb8568f6-5dwfh    1/1     Running   0             3h29m
kubia-deployment-5fb8568f6-p744s    1/1     Running   0             3h29m
mysql-deployment-66bc478755-8tghj   1/1     Running   0             4m10s
mysql-deployment-66bc478755-hznpc   1/1     Running   2 (45s ago)   4m15s
mysql-deployment-66bc478755-rhh5v   1/1     Running   2 (41s ago)   4m8s
```

- Restart 하는 파드의 로그를 확인해보면 mysql 파일을 읽고 사용중이라 락에 걸린걸 확인할 수 있다.
```
k logs mysql-deployment-66bc478755-hznpc
2023-05-24 04:47:59+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.0.33-1.el8 started.
2023-05-24 04:47:59+00:00 [Note] [Entrypoint]: Switching to dedicated user 'mysql'
2023-05-24 04:47:59+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.0.33-1.el8 started.
'/var/lib/mysql/mysql.sock' -> '/var/run/mysqld/mysqld.sock'
2023-05-24T04:48:00.131565Z 0 [Warning] [MY-011068] [Server] The syntax '--skip-host-cache' is deprecated and will be removed in a future release. Please use SET GLOBAL host_cache_size=0 instead.
2023-05-24T04:48:00.134585Z 0 [System] [MY-010116] [Server] /usr/sbin/mysqld (mysqld 8.0.33) starting as process 1
2023-05-24T04:48:00.141674Z 1 [System] [MY-013576] [InnoDB] InnoDB initialization has started.
2023-05-24T04:48:00.157869Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:01.158065Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:02.158399Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:03.159021Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:04.159571Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:05.159944Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:06.160524Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:07.161039Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
2023-05-24T04:48:08.162173Z 1 [ERROR] [MY-012574] [InnoDB] Unable to lock ./ibdata1 error: 11
```

# StatefulSet 데이터
- yaml은 단순히 statefulset으로 바꾸고 진행한다.
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql-statefulset
spec:
  serviceName: mysql
  replicas: 3
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      partition: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:5.7
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: "1234"
          ports:
            - containerPort: 3306
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql
  volumeClaimTemplates:
    - metadata:
        name: mysql-data
      spec:
        accessModes:
          - ReadWriteOnce
        storageClassName: local-path
        resources:
          requests:
            storage: 1Gi

```
```shell
$ k get pvc 
NAME                             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
mysql-pvc                        Bound    pvc-2e46d161-83a0-4257-9289-5ba1ac8d48f4   1Gi        RWO            local-path     27m
mysql-data-mysql-statefulset-0   Bound    pvc-c63c9518-318d-4ff0-9307-bce8d4d6577f   1Gi        RWO            local-path     68s
mysql-data-mysql-statefulset-1   Bound    pvc-92da75ff-c6ea-4eda-9aad-809ffca2ad22   1Gi        RWO            local-path     56s
mysql-data-mysql-statefulset-2   Bound    pvc-bbef0545-0ea4-45ed-94d5-50323656ed5c   1Gi        RWO            local-path     45s
```
기존 deployment와 다르게 각각의 파드마다 pvc를 만들어 준 것을 알 수 있다.
```shell
$ k get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                    STORAGECLASS   REASON   AGE
pvc-2e46d161-83a0-4257-9289-5ba1ac8d48f4   1Gi        RWO            Delete           Bound    default/mysql-pvc                        local-path              26m
pvc-c63c9518-318d-4ff0-9307-bce8d4d6577f   1Gi        RWO            Delete           Bound    default/mysql-data-mysql-statefulset-0   local-path              113s
pvc-92da75ff-c6ea-4eda-9aad-809ffca2ad22   1Gi        RWO            Delete           Bound    default/mysql-data-mysql-statefulset-1   local-path              103s
pvc-bbef0545-0ea4-45ed-94d5-50323656ed5c   1Gi        RWO            Delete           Bound    default/mysql-data-mysql-statefulset-2   local-path              92s
```
pv도 마찬가지이다.

그러면 각각의 mysql 컨테이너에 데이터를 넣고 롤링업데이트를 해보자. 
```shell
k edit statefulsets.apps mysql-statefulset            
statefulset.apps/mysql-statefulset edited
☁  kubernetes [kubectl-assignment2] ⚡  k get pods                                
NAME                  READY   STATUS              RESTARTS         AGE
liveness-pod          1/1     Terminating         543 (3d2h ago)   6d
mysql-statefulset-0   1/1     Running             0                3m53s
mysql-statefulset-1   1/1     Running             0                3m47s
mysql-statefulset-2   0/1     ContainerCreating   0                0s
```
image쪽에서 계속 오류가 발생해 파드를 모두 삭제하고 재실행했다.
```shell
k delete pod --all                
pod "liveness-pod" deleted
pod "mysql-statefulset-0" deleted
pod "mysql-statefulset-1" deleted
pod "mysql-statefulset-2" deleted

$ k get pods        
NAME                  READY   STATUS              RESTARTS         AGE
liveness-pod          1/1     Terminating         543 (3d2h ago)   6d
mysql-statefulset-0   0/1     ContainerCreating   0                3s

$ k get pods
NAME                  READY   STATUS        RESTARTS         AGE
liveness-pod          1/1     Terminating   543 (3d2h ago)   6d
mysql-statefulset-0   1/1     Running       0                53s
mysql-statefulset-1   1/1     Running       0                46s
mysql-statefulset-2   1/1     Running       0                39s
```

- 0번과 1번 테이블에 각각 데이터 1,2를 넣었다.
- 다시 실행된 파드들에 대해 결과를 확인해보면 아래와 같다.
```shell
mysql> select * from users;
+------+
| age  |
+------+
|    1 |
+------+
1 row in set (0.00 sec)

mysql> select * from users;
+------+
| age  |
+------+
|    2 |
+------+
1 row in set (0.01 sec)
```
- 결과적으로 다시 업데이트 해도 같은 결과가 나오는 것을 알 수 있다.
  - 파드를 삭제해도 pvc는 삭제되지 않는다.
  - 0,1,2의 index가 파드의 순번을 나타내는데 이런 고유한 식별자로 순서대로 큰 숫자부터 줄어들고  작은숫자부터 생성되면서 기존 PVC와 다시 연동되는 것을 확인했다.
