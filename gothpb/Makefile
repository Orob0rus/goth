GOPATH="/home/vishnu/go"
GOBIN="/home/vishnu/go/bin"
GOTH="/home/vishnu/projects"
SRC_DIR="${GOTH}/goth/gothpb"
DST_DIR="${GOTH}/goth/goth-proto-go"
compile:
	protoc -I=${SRC_DIR} --go_out=${DST_DIR} ${SRC_DIR}/goth.proto
