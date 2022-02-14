# q2go
minimalistic queue server with REST access written in Golang

### create a queue with name "queue123"
curl -d "qname=queue123" -X POST "http://localhost:8080/queue"

### create some message in the above created queue
curl -d "message=hello this is some more text" -X POST "http://localhost:8080/queue/queue123/message"

### read the first message (fifo) from the above create queue
curl http://localhost:8080/queue/queue123/message

### delete a queue
curl -X DELETE "http://localhost:8080/queue/queue123"

### publish 100 messages (wget as alternative)
for ((i=1;i<=100;i++)) ; do curl -d "message=hello this is some more text $i" -X POST "http://localhost:8080/queue/queue123/message" ; done
for ((i=1;i<=100;i++)) ; do wget http://localhost:8080/queue/queue123/message --post-data="message=hello this is some more text $i" ; done

### read 100 messages
for ((i=1;i<=100;i++)); do curl http://localhost:8080/queue/queue123/message ; done




