# q2go
Queue server with Rest access written in Go.

### publish 100 messages
for ((i=1;i<=100;i++)) ; do curl -d "message=hello this is some more text $i" -X POST "http://localhost:8080/push" ; done

### read 100 messages
for ((i=1;i<=100;i++)); do curl http://localhost:8080/pop ; done
