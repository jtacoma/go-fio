fio
===

    go get github.com/jtacoma/go-fio

This package provides a framing abstraction over the `"bufio"` package.

    w := fio.NewSender(sys.Stdout, fio.Zio1)
    for i := 0; i < 1000; i += 1 {
        w.Send([]byte("Hello, World!\n"))
    }

An implementation of the framing layer of [ZMTP/1.0](http://rfc.zeromq.org/spec:13) is included as reference implentation of the `fio.Encoding` interface.

License
-------

Copyright 2013 Joshua Tacoma

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
