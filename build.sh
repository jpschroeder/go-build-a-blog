
echo "Bundling Templates"
go-bindata templates/... static/...

echo "Building Source"
go build
