fio
===

    go get github.com/jtacoma/go-fio

This package also provides easy efficiency for high-frequency writing.  For example, the following code will result in a small number of large writes to `sys.Stdout` even though it makes a large number of small writes to `w`:

    func main() {
        w := fio.NewWriter(sys.Stdout)
        for i := 0; i < 1000; i += 1 {
            w.Write([]byte("Hello, World!\n"))
        }
    }

There some abstractions over the standard `"net"` package in the `"fionet"` subdirectory to support the use of this package in the implementation of custom protocols.

As an example, the `"zio1"` subdirectory contains an implementation of the framing parts of [ZMTP/1.0](http://rfc.zeromq.org/spec:13).

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
