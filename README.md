# iwebd -- instant web daemons

iwebd is a monolitic suite of instant web daemons for sharing files.
It began as a Python program,
[iwebd.py](https://github.com/phillipberndt/scripts/blob/master/iwebd/iwebd.py),
and grew to supporting a variety of protocols, most of which I only
added to get familiar with them.

This is a more practically oriented rewrite. It focuses entirely on
up- and download of files, only contains the relevant core protocols,
http(s), ftp, webdav(s) and upnp/dlna, and a web frontend to allow
working on sets of files a bit better than http would allow otherwise.

The advantage of this reimplementation is that it can handle load and
has code closer to production-ready. This of course is because Go ships with
a huge, high quality standard library, and because there's open source
implementations of protocols available of similar high quality as libraries
that can be included in programs.

## Relevant Third-Party Libraries Used
* dlna functionality is provided by [anacrolix/dms](https://github.com/anacrolix/dms) (BSD 3),
* ftp functionality is provided by [goftp](https://gitea.com/goftp/server) by yob et al. (MIT),
* icons are taken from the [Yaru theme from Ubuntu](https://github.com/ubuntu/yaru) (CC BY-SA 4.0),
* and this entire project benefits a lot from the huge stdlib Go brings.

## License
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

A full copy of the license is available in [LICENSE](LICENSE).
