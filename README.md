# Go du - estimate file space usage

_An implementation of [du(1)](https://man7.org/linux/man-pages/man1/du.1p.html)
in Golang._

The purpose of this implementation is not to build a better `du`. It is to learn
Golang. I chose `du` because it serves a simple purpose - show how much space
files take up on disk. But from an educational perspective it provides several
good learning opportunities:
 - how to work with command-line options and switches
 - how to work with the file system
 - concurrency
 - error handling

The implementation is POSIX _compatible_. Other implementations exist, such as
the one from GNU coreutils: https://www.gnu.org/software/coreutils/du.

If you want to learn together, ask a question, offer help or get in touch for
any other reason please don't hesitate to contact me
[frenkel.ilia@gmail.com](mailto:frenkel.ilia@gmail.com).
