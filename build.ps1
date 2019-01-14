write-output "Bundling Templates"
go-bindata templates/...
write-output "Building Source"
go build