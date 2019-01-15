
echo "Bundling Templates"
go-bindata templates/...

echo "Building Source"
go build
