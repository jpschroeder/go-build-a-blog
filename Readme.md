
# go-build-a-blog

This project is my attempt at building the simplest possible blog engine in go.  You can list, view, add, edit, and delete blog pages as well as the blog frontpage content.

## features

- project compiles to a single executable with all templates and static resources included
- all state is stored in a single file sqlite database
- all pages are written in markdown
- rendered pages include syntax highlighting for code blocks
- editor includes syntax highlighting for markdown and code blocks
- built in support for https and generating certificates from letsencrypt
- minimalistic design

## installation

With go installed:
```shell
go get -u github.com/jpschroeder/go-build-a-blog
```

## usage

```shell
go-build-a-blog.exe -h
  -db string
        the path to the sqlite database file
        it will be created if it does not already exist
         (default "go-build-a-blog.db")
  -httpaddr string
        the address/port to listen on for http
        use :<port> to listen on all addresses
         (default "localhost:8080")
  -httpsaddr string
        the address/port to listen on for https
        use :<port> to listen on all addresses
        this should only be used when listening publicly with proper dns address configured
        this will generate a certificate using letsencrypt
        the server will also listen on the -httpaddr but will redirect to https

  -httpsdomain string
        the domain to use for https
        this flag should be used in conjunction with the -httpsaddr flag
        this should only be used when listening publicly with proper dns address configured

  -reset
        reset the security key used to edit/delete
```

## building and bundling

In order to build the project, just use:
```shell
go build
```

However, if you change any of the files in the `templates/` folder, you will first have to update the `bindata.go` file using the `go-bindata` tool.  This allows these template files to be bundled into the executable.

Install the bundling tool and build the executable as follows:
```shell
go get -u github.com/shuLhan/go-bindata/cmd/go-bindata
go generate
go build
```

You can also use the included build script that combines these steps.

## deploying

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
