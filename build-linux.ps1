
write-output "Bundling Templates"
go-bindata templates/...

write-output "Building Source"

$env:GOOS = 'linux';
$env:GOARCH = 'amd64';

go build

Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
