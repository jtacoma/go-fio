# zio
--
    import "github.com/jtacoma/go-fio/zio"

Package zio provides functions for ZMTP/1.0 compliant I/O.

ZMTP/1.0 is defined here: http://rfc.zeromq.org/spec:13

Because these functions wrap lower-level operations, unless otherwise
informed clients should not assume they are safe for parallel execution.

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
