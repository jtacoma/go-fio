fio
===

    go get github.com/jtacoma/go-fio

This is an experiment in easing the development of custom protocols through a framing abstraction over the standard "io" package.

Coincidentally, this package also provides easy efficiency for high-frequency writing.  For example, the following code will result in a small number of large writes to `sys.Stdout` even though it makes a large number of small writes to `w`:

    func main() {
        w := fio.NewWriter(sys.Stdout)
        for i := 0; i < 1000; i += 1 {
            w.Write([]byte("Hello, World!\n"))
        }
    }

There are also some abstractions over the standard `"net"` package in the `"fionet"` subdirectory.

Plan
----

* Use this package to develop some custom protocols (or implementations of published protocols).
* Increase test coverate with `gocov`.  It should be possible to hit 100% or very close to it.

License
-------

See the NOTICE file distributed with this work for additional information regarding copyright ownership.  Joshua Tacoma licenses this file to you under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the License for the specific language governing permissions and limitations under the License.
