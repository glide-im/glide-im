protoc \
--proto_path=./google/ \
--proto_path=./rpc/ \
--proto_path=./im/ \
--go_out=./gen/ \
./google/*.proto \
./rpc/*.proto \
./im/*.proto