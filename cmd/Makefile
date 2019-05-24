proto: proto/sign/sign.pb.go
	echo "finish!"

clean:
	rm proto/*/*.pb.go

proto/sign/sign.pb.go: proto/sign/sign.proto
	protoc --go_out=. proto/sign/sign.proto

