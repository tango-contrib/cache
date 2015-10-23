// Copyright 2014 The Macaron Authors
// Copyright 2015 The Tango Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cache

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lunny/tango"
	. "github.com/smartystreets/goconvey/convey"
)

var testStructOption = Options{
	Interval: 2,
}

func Test_Cacher(t *testing.T) {
	Convey("Use cache middleware", t, func() {
		t := tango.Classic()
		t.Use(New())
		t.Get("/", new(testCacheController))

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		t.ServeHTTP(resp, req)
	})

	Convey("Register invalid adapter", t, func() {
		Convey("Adatper not exists", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			t := tango.Classic()
			t.Use(New(Options{
				Adapter: "fake",
			}))
		})

		Convey("Provider value is nil", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Register("fake", nil)
		})

		Convey("Register twice", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Register("memory", &MemoryCacher{})
		})
	})
}

type testCacheController struct {
	Cache
}

func (ts *testCacheController) Get() {
	So(ts.Cache.Put("uname", "unknwon", 1), ShouldBeNil)
	So(ts.Cache.Put("uname2", "unknwon2", 1), ShouldBeNil)
	So(ts.Cache.IsExist("uname"), ShouldBeTrue)

	So(ts.Cache.Get("404"), ShouldBeNil)
	So(ts.Cache.Get("uname").(string), ShouldEqual, "unknwon")

	time.Sleep(1 * time.Second)
	So(ts.Cache.Get("uname"), ShouldBeNil)
	time.Sleep(1 * time.Second)
	So(ts.Cache.Get("uname2"), ShouldBeNil)

	So(ts.Cache.Put("uname", "unknwon", 0), ShouldBeNil)
	So(ts.Cache.Delete("uname"), ShouldBeNil)
	So(ts.Cache.Get("uname"), ShouldBeNil)

	So(ts.Cache.Put("uname", "unknwon", 0), ShouldBeNil)
	So(ts.Cache.Flush(), ShouldBeNil)
	So(ts.Cache.Get("uname"), ShouldBeNil)

	gob.Register(testStructOption)
	So(ts.Cache.Put("struct", testStructOption, 0), ShouldBeNil)
}

type test2CacheController struct {
	Cache
}

func (t2s *test2CacheController) Get() {
	So(t2s.Cache.Incr("404"), ShouldNotBeNil)
	So(t2s.Cache.Decr("404"), ShouldNotBeNil)

	So(t2s.Cache.Put("int", 0, 0), ShouldBeNil)
	So(t2s.Cache.Put("int32", int32(0), 0), ShouldBeNil)
	So(t2s.Cache.Put("int64", int64(0), 0), ShouldBeNil)
	So(t2s.Cache.Put("uint", uint(0), 0), ShouldBeNil)
	So(t2s.Cache.Put("uint32", uint32(0), 0), ShouldBeNil)
	So(t2s.Cache.Put("uint64", uint64(0), 0), ShouldBeNil)
	So(t2s.Cache.Put("string", "hi", 0), ShouldBeNil)

	So(t2s.Cache.Decr("uint"), ShouldNotBeNil)
	So(t2s.Cache.Decr("uint32"), ShouldNotBeNil)
	So(t2s.Cache.Decr("uint64"), ShouldNotBeNil)

	So(t2s.Cache.Incr("int"), ShouldBeNil)
	So(t2s.Cache.Incr("int32"), ShouldBeNil)
	So(t2s.Cache.Incr("int64"), ShouldBeNil)
	So(t2s.Cache.Incr("uint"), ShouldBeNil)
	So(t2s.Cache.Incr("uint32"), ShouldBeNil)
	So(t2s.Cache.Incr("uint64"), ShouldBeNil)

	So(t2s.Cache.Decr("int"), ShouldBeNil)
	So(t2s.Cache.Decr("int32"), ShouldBeNil)
	So(t2s.Cache.Decr("int64"), ShouldBeNil)
	So(t2s.Cache.Decr("uint"), ShouldBeNil)
	So(t2s.Cache.Decr("uint32"), ShouldBeNil)
	So(t2s.Cache.Decr("uint64"), ShouldBeNil)

	So(t2s.Cache.Incr("string"), ShouldNotBeNil)
	So(t2s.Cache.Decr("string"), ShouldNotBeNil)

	So(t2s.Cache.Get("int"), ShouldEqual, 0)
	So(t2s.Cache.Get("int32"), ShouldEqual, 0)
	So(t2s.Cache.Get("int64"), ShouldEqual, 0)
	So(t2s.Cache.Get("uint"), ShouldEqual, 0)
	So(t2s.Cache.Get("uint32"), ShouldEqual, 0)
	So(t2s.Cache.Get("uint64"), ShouldEqual, 0)
}

func testAdapter(opt Options) {
	Convey("Basic operations", func() {
		t := tango.Classic()
		t.Use(New(opt))
		t.Get("/", new(testCacheController))

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		t.ServeHTTP(resp, req)
	})

	Convey("Increase and decrease operations", func() {
		t := tango.Classic()
		t.Use(New(opt))

		t.Get("/", new(test2CacheController))

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		t.ServeHTTP(resp, req)
	})
}
