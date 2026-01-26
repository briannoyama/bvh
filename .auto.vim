let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/Development/go/src/bvh
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +90 bvh/node.go
badd +45 bvh/node_test.go
badd +113 ~/Development/go/src/bvh/volume/orthotope.go
badd +88 rect/bvh_test.go
badd +34 ~/Development/go/src/bvh/bvh/tree_test.go
badd +130 bvh/tree.go
badd +127 main/main.go
badd +2 .gitignore
badd +12 test.json
badd +2 main/example_test.go
badd +123 main2/main.go
badd +6 bvh/impl.go
badd +12 testBVH.json
badd +1 volume/constants.go
badd +249984 testdataBVH.csv
badd +249981 testdataOrig.csv
badd +33 ~/Development/go/src/bvh/volume/sphere.go
badd +1 go.mod
badd +1 go.sum
badd +42 ~/Development/go/src/bvh/volume/sphere_test.go
argglobal
%argdel
edit bvh/impl.go
argglobal
balt ~/Development/go/src/bvh/volume/sphere_test.go
setlocal foldmethod=manual
setlocal foldexpr=0
setlocal foldmarker={{{,}}}
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldenable
silent! normal! zE
let &fdl = &fdl
let s:l = 6 - ((5 * winheight(0) + 25) / 50)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 6
normal! 0
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
