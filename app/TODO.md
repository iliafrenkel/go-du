# TODO

 - [ ] More robust unit tests with higher coverage
    * Test with -a flag
    * ~Test with -k flag~
    * Test with -s flag
    * Test flag conflicts
    * ~Test single file~
    * ~Test single directory~
    * Test multiple files
    * ~Test multiple directories~
    * Test multiple files and directories
    * Test different sizes
        - empty directory
        - empty file
        - file under 4k
        - file over 4k
        - big file, over 1M
    * Test access denied

 - [ ] CI/CD pipeline that automatically tests, builds and creates releases for
       multiple architectures

 - [ ] Cross platform detection of the silesystem block size. At the moment it
       is hardcoded to be 4096 bytes.

 - [ ] Generate coverage report:
       ```
       go test --cover -coverprofile=coverage.out ./...
       go tool cover -html=coverage.out -o covearge.html
       ```
 - [ ] Add `--version` flag to show license, contact and version info.
 - [ ] Cross-platform compilation
