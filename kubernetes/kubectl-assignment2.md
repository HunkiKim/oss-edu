# Deployment의 데이터
PV, StorageClass, PVC를 따로 정의해야한다. (자동으로 해주지 않기 때문에.)

![image](https://github.com/HunkiKim/Mantech-Edu/assets/66348135/685b9a4a-774b-4519-b721-54b29b51c03d)

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
- 롤링업데이트는 가능하지만, pvc를 각각 파드 개수에 맞춰서 만들어놔야한다.

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