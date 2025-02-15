module github.com/g3n/g3nd

go 1.23

toolchain go1.23.4

require (
	github.com/g3n/engine v0.2.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
)

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20211213063430-748e38ca8aec // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/g3n/engine => ../engine
