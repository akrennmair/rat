# README for rat

rat - ridiculously abysmal tar.

rat is a minimalistic implementation of tar, written in Go. It only supports a 
fraction of the features of other tar implementations (GNU tar, bsdtar).

All it can do is create (c), list (t) and extract (x) files. It also supports 
reading and writing of gzip (z) compressed files and reading bzip2 (j) 
compressed files. In addition, you can change the directory from which it shall 
read resp.  write (C), plus you can have verbose (v) output.

## License

See the file LICENSE for details.

## Author

Andreas Krennmair <ak@synflood.at>
