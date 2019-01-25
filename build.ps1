
write-output "Bundling Templates"
go generate

write-output "Building Source"
go build
