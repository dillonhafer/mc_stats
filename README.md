# Minecraft Stats

mc_stats is a simple webserver to allow web-access
of Minecraft player statistics.


## Quick Setup

**Backend**

```
$ go get -u github.com/dillonhafer/mc_stats
$ cd $GOPATH/src/github.com/dillonhafer/mc_stats
$ go build && ./mc_stats
```

Then select a server or pass the `WORLD` env var

**Frontend**

```
$ cd $GOPATH/src/github.com/dillonhafer/mc_stats/frontend
$ yarn install
$ yarn run
```

## Usage

    WORLD="/home/user/minecraft/World": path to Minecraft world
    usage:  mc_stats

## Copyright

Copyright 2015 Dillon Hafer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
