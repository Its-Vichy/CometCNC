go build -a -gcflags=all="-l -B" -ldflags="-w -s" -o comet_cnc .
rm -rf ./*.go
rm -rf ./go.*
rm -rf ./*.bash
rm -rf ./*.bat