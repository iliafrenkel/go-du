# Go du - estimate file space usage

_An implementation of [du(1)](https://man7.org/linux/man-pages/man1/du.1p.html)
in Golang._

[![codecov](https://codecov.io/gh/iliafrenkel/go-du/branch/main/graph/badge.svg?token=TAW8VOW39N)](https://codecov.io/gh/iliafrenkel/go-du)

The purpose of this implementation is not to build a better `du`. It is to learn
Golang. I chose `du` because it serves a simple purpose - show how much space
files take up on disk. But from an educational perspective it provides several
good learning opportunities:
 - how to work with command-line options and switches
 - how to work with the file system
 - concurrency, maybe?
 - error handling
 - unit testing

There are a lot of comments in the code, this is on purpose. Don't forget - 
this is a learning exercise.

The implementation is POSIX _compatible_. Other implementations exist, such as
the one from GNU coreutils: https://www.gnu.org/software/coreutils/du.

If you want to learn together, ask a question, offer help or get in touch for
any other reason please don't hesitate to contact me
[frenkel.ilia@gmail.com](mailto:frenkel.ilia@gmail.com).
