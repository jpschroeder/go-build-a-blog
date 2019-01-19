
write-output "Bundling Templates"
go-bindata templates/... static/...

write-output "Building Source"
go build
