# webapp-vim

## What is this?

Web server that can write web application in Vim script.

# Requirements

[webapi-vim](https://github.com/mattn/webapi-vim)

# Install

You need to compile server. To comiple server, you need to install [golang](http://golang.org).
After installing golang, type following.

    $ cd ~/.vim/bundle
    $ git clone https://github.com/mattn/webapp-vim
    $ cd webapp-vim/server
    $ go build webappvimd.go
    $ ./webappvimd

# Note

This is application server. So this don't contains example to run webapp.
Check [webapp-vim-vim](https://github.com/mattn/webapp-foo-vim)

# License

MIT

# Author

Yasuhiro Matsumoto
