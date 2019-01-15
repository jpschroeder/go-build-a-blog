
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

# deploying

You can build the project under linux (or Windows Subsystem for Linux) and just copy the executable to your server.

You can then run the program directly or use systemd to install it as a service and keep it running.

Customize the `go-build-a-blog.service` file in the repo for your server and copy it to `/lib/systemd/system/go-build-a-blog.service`

Start the app with: `systemctl start go-build-a-blog`  
Enable it on boot with: `systemctl enable go-build-a-blog`  
Check it's status with: `systemctl status go-build-a-blog`  
See standard output/error with: `journalctl -f -u go-build-a-blog`

## credits

- [txti](http://txti.es/) Minimalistic design inspiration
- [Blackfriday](https://github.com/russross/blackfriday) Markdown parsing
- [gorilla/mux](https://github.com/gorilla/mux) Routing
- [slugify](https://github.com/avelino/slugify) Create url slugs
- [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt) Cryptographic hashes
