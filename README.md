# Vimpack

Vim comes with a default package manager which is able to load packages on start or on demand. Unfortunately vim's package manager is not able to clone/pull packages from GitHub automatically.

Vimpack defines a file format respecting vim's package manager structure and letting the user defines how to group plugins.

Here an example:

```
// ~/.vim/Packfile
namespace "musthave"

onstart "tpope/vim-commentary"
onstart "tpope/vim-repeat"
onstart "tpope/vim-sensible"
onstart "tpope/vim-surround"
onstart "tpope/vim-unimpaired"

namespace "textobjects"

ondemand "andyl/vim-textobj-elixir"
onstart "chaoren/vim-wordmotion"
onstart "christoomey/vim-sort-motion"
onstart "kana/vim-textobj-user"
onstart "wellle/targets.vim"

namespace "snippets"

onstart "SirVer/ultisnips"
onstart "honza/vim-snippets"
```

This file will produce:

```
~/.vim/pack/
  musthave/
    start/
      vim-commentary/
      vim-repeat/
      vim-sensible/
      vim-surround/
      vim-unimpaired/
  textobjects/
    opt/
      vim-textobj-elixir
    start/
      vim-wordmotion/
      vim-sort-motion/
      vim-textobj-user/
      targets.vim/
  snippets/
    start/
      ultisnips/
      vim-snippets/
```

The only configuration required to vim, is to add the following command in your `~/.vimrc`:

```
packloadall
```

## Usage

```
Usage of vimpack:
  -dest string
    	path where to download plugins (default "~/.vim/pack")
  -file string
    	path to the Vimpackfile (default "~/.vim/Vimpackfile")
```

## Development

* `make refresh-mocks`: rebuild the mocks
* `make test-style`: run the style suite (`gofmt`, `golint`)
* `make test-unit`: run the test suite

###  Go Versions

We use `asdf` to manage my Go versions. Installation:

```
asdf install
```

Each time we bump the Go version, we should regenerate the Circle-CI config file:

```
make refresh-templates
```
