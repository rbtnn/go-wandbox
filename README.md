
# go-wandbox

[wandbox](http://melpon.org/wandbox) CLI

## usage

    $go get github.com/rbtnn/go-wandbox/cmd/wandbox
    $cat ~/a.vim
    echo "hi"
    $wandbox -compiler vim-head -source ~/a.vim
    hi
    $wandbox -compiler vim-head -code "echo 'hi'"
    hi
    $wandbox -help
    Usage of wandbox:
      -code string
            code
      -compiler string
            compiler
      -list
            compiler list
      -source string
            source file

    $wandbox -list
    [Elixir]
      elixir-head
      elixir-1.4.1
      elixir-1.3.4
    [CoffeeScript]
      coffeescript-head
      coffeescript-1.12.3
      coffeescript-1.11.1
      coffeescript-1.10.0
    [C]
      gcc-head-c
      gcc-6.3.0-c
      gcc-6.2.0-c
      gcc-6.1.0-c
      gcc-5.4.0-c
      gcc-5.3.0-c
      gcc-5.2.0-c
      gcc-5.1.0-c
      gcc-4.9.3-c
    ...
