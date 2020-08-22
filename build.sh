protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/digitalwantlist.proto
mv proto/github.com/brotherlogic/digitalwantlist/proto/* ./proto
