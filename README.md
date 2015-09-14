# webapp-vim

## What is this?

Web server that can write web application in Vim script.

![](http://go-gyazo.appspot.com/9f3e1755f0ee695b.png)



# Requirements

[webapi-vim](https://github.com/mattn/webapi-vim)

# Install

You need to compile server. To comiple server, you need to install [golang](http://golang.org).
After installing golang, type following.

    $ cd ~/.vim/bundle
    $ git clone https://github.com/mattn/webapp-vim
    $ cd webapp-vim/server
    $ go build
    $ ./server

# Example application

This is application server. So this don't contains example to run webapp.
Check [webapp-foo-vim](https://github.com/mattn/webapp-foo-vim)

# How to register your webapp

You need to make following directory structure.

    +---autoload
    |   |
    |   +--- myapp.vim ... add code for your application
    |
    +--- plugin ... script to register your application
    |
    +--- static ... static files
    
1. Add script to register your application in `plugin/myapp.vim`.

        call webapp#handle("/myapp", function('myapp#handle'))

2. Put html/js/css into `static` directory.

3. Write application

        function! myapp#handle(req)
          if a:req.path == '/foo'
            return {"body", "hello world"}
          else
            let a:req.path = a:req.path[4:]
            return webapp#servefile(a:req, s:basedir)
          endif
        endfunction

# License

MIT

# Author

Yasuhiro Matsumoto
