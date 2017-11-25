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

Response body:  copy of *Post body*

## tail

`tail -f` a tmp file for 60 seconds

```
http://<host>:8080/tail
```

Method: `GET`
