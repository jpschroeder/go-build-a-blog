
# go-build-a-blog

This project is my attempt at building the simplest possible blog engine in go.  You can list, view, add, edit, and delete blog pages.  Pages are written in markdown.

## building and bundling

In order to build the project, just use:
```
go build
```

However, if you change any of the files in the `templates/` folder, you will first have to update the `bindata.go` file using the `go-bindata` tool.  This allows these template files to be bundled into the executable.

Install the bundling tool and build the executable as follows:
```
go get -u github.com/shuLhan/go-bindata/cmd/go-bindata
go-bindata templates/...
go build
```

You can also use the included build script that combines these steps.

## credits

- [txti](http://txti.es/) Minimalistic design inspiration
- [Blackfriday](https://github.com/russross/blackfriday) Markdown parsing
- [gorilla/mux](https://github.com/gorilla/mux) Routing
- [slugify](https://github.com/avelino/slugify) Create url slugs
- [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt) Cryptographic hashes
