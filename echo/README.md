# Echo

`echo` is a REST service for demo, which has the following API:

## health check

```
http://<host>:8080/healthz
```

Method: `GET`

Response body:

```
OK
```

## echo

```
http://<host>:8080/echo
```

Method: `POST`

Response body:  copy of *Post body* appending with env var `echo`. Responsing with env is supported after 1.1 .

```
$ curl -X POST http://localhost:8080/echo -d 'test'

body=[test] echo=[]
```

## tail

`tail -f` a tmp file for 60 seconds

```
http://<host>:8080/tail
```

Method: `GET`
