let s:basedir = expand('<sfile>:h:h') . '/static'
let s:count = 1

let s:mimetypes = {
\ 'ico':  'image/x-icon',
\ 'html': 'text/html; charset=UTF-8',
\ 'js':   'application/javascript; charset=UTF-8',
\ 'txt':  'text/plain; charset=UTF-8',
\ 'jpg':  'image/jpeg',
\ 'gif':  'image/gif',
\ 'png':  'image/png',
\}

function! webapp#path2slash(path)
  return substitute(a:path, '\\', '/', 'g')
endfunction

function! webapp#fname2mimetype(fname)
  let ext = fnamemodify(a:fname, ':e')
  if has_key(s:mimetypes, ext)
    return s:mimetypes[ext]
  else
    return 'application/octet-stream'
  endif
endfunction

if !exists('s:handlers')
  let s:handlers = {}
endif

function! webapp#params(req)
  let params = {}
  for q in split(a:req.query, '&') 
    let pos = stridx(q, '=')
    if pos > 0
      let params[q[:pos-1]] = q[pos+1:]
    endif
  endfor
  return params
endfunction

function! webapp#handle(path, Func)
  let s:handlers[a:path] = a:Func
endfunction

function! webapp#json(req, obj, ...)
  let res = webapi#json#encode(a:obj)
  let cb = get(a:000, 0, '')
  if len(cb) != 0
    let res = cb . '(' . res . ')'
  endif
  return {"header": ["Content-Type: application/json"], "body": res}
endfunction

function! webapp#redirect(req, to)
  return {"header": ["Location: " . a:to], "status": 302}
endfunction

function! webapp#servefile(req, basedir)
  let res = {"header": [], "body": "", "status": 200}
  let fname = a:basedir . a:req.path
  if isdirectory(fname)
    if filereadable(fname . '/index.html')
      let fname .= '/index.html'
      let mimetype = webapp#fname2mimetype(fname)
      call add(res.header, "Content-Type: " . mimetype)
      if mimetype =~ '^text/'
        let res.body = iconv(join(readfile(fname, 'b'), "\n"), "UTF-8", &encoding)
      else
        let res.body = map(split(substitute(system("xxd -ps " . fname), "[\r\n]", "", "g"), '..\zs'), '"0x".v:val+0')
      endif
    else
      call add(res.header, "Content-Type: text/plain; charset=UTF-8")
      let res.body = join(map(map(split(glob(fname . '/*'), "\n"), 'a:req.path . webapp#path2slash(v:val[len(fname):])'), '"<a href=\"".webapi#http#encodeURIComponent(v:val)."\">".webapi#html#encodeEntityReference(v:val)."</a><br>"'), "\n")
    endif
  elseif filereadable(fname)
    let mimetype = webapp#fname2mimetype(fname)
    call add(res.header, "Content-Type: " . mimetype)
    if mimetype =~ '^text/'
      let res.body = iconv(join(readfile(fname, 'b'), "\n"), "UTF-8", &encoding)
    else
      let res.body = map(split(substitute(system("xxd -ps " . fname), "[\r\n]", "", "g"), '..\zs'), '"0x".v:val+0')
    endif
  else
    let res.status = 404
    let res.body = "Not Found"
  endif
  return res
endfunction

function! webapp#serve(req)
  try
    for path in reverse(sort(keys(s:handlers)))
      if stridx(a:req.path, path) == 0
        return s:handlers[path](a:req)
      endif
    endfor
    let res = webapp#servefile(a:req, s:basedir)
  catch
    let res = {"header": [], "body": "Internal Server Error: " . v:exception, "status": 500}
  endtry
  return res
endfunction
