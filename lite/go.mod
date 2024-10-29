module CentralizedControl

go 1.21

toolchain go1.21.0

replace (
	github.com/4kills/go-libdeflate/v2 => ./three/go-libdeflate
	github.com/bytedance/sonic => ./three/sonic
	github.com/elliotchance/orderedmap/v2 => ./three/orderedmap/v2@v2.2.0
	github.com/mzz2017/gg => ./three/gg
	github.com/mzz2017/softwind => ./three/softwind
	github.com/pquerna/otp => ./three/otp
	github.com/utahta/go-cronowriter => ./three/go-cronowriter
	github.com/v2rayA/shadowsocksR => ./three/shadowsocksR
)

require (
	fyne.io/fyne/v2 v2.5.1
	github.com/bogdanfinn/fhttp v0.5.28
	github.com/bogdanfinn/utls v1.6.1
	github.com/bytedance/sonic v1.10.0-rc
	github.com/emersion/go-imap v1.2.1
	github.com/emirpasic/gods v1.18.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gorilla/mux v1.8.1
	github.com/jamespfennell/xz v0.1.2
	github.com/klauspost/compress v1.17.9
	github.com/mzz2017/gg v0.0.0-00010101000000-000000000000
	github.com/mzz2017/softwind v0.0.0-00010101000000-000000000000
	github.com/pquerna/otp v0.0.0-00010101000000-000000000000
	github.com/refraction-networking/utls v1.6.2
	github.com/twmb/murmur3 v1.1.8
	github.com/utahta/go-cronowriter v0.0.0-00010101000000-000000000000
	github.com/v2rayA/shadowsocksR v1.0.4
	golang.org/x/net v0.28.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.11
)

require (
	fyne.io/systray v1.11.0 // indirect
	github.com/AlecAivazis/survey/v2 v2.3.2 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/cloudflare/circl v1.3.9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-camellia v0.0.0-20191119043421-69a8a13fb23d // indirect
	github.com/dgryski/go-idea v0.0.0-20170306091226-d2fb45a411fb // indirect
	github.com/dgryski/go-metro v0.0.0-20200812162917-85c65e2d0165 // indirect
	github.com/dgryski/go-rc2 v0.0.0-20150621095337-8a9021637152 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eknkc/basex v1.0.1 // indirect
	github.com/emersion/go-sasl v0.0.0-20200509203442-7bfe0ed36a21 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fredbi/uri v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fyne-io/gl-js v0.0.0-20220119005834-d2da28d9ccfe // indirect
	github.com/fyne-io/glfw-js v0.0.0-20240101223322-6e1efdc71b7a // indirect
	github.com/fyne-io/image v0.0.0-20220602074514-4956b0afb3d2 // indirect
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-text/render v0.1.1-0.20240418202334-dd62631dae9b // indirect
	github.com/go-text/typesetting v0.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20240223122105-ce5225dcaa49 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jsummers/gobmp v0.0.0-20151104160322-e2ba15ffa76e // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/lestrrat-go/strftime v0.0.0-20180220091553-9948d03c6207 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mzz2017/disk-bloom v1.0.1 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.4.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/quic-go v0.40.1 // indirect
	github.com/rymdport/portal v0.2.6 // indirect
	github.com/seiflotfy/cuckoofilter v0.0.0-20201222105146-bc6005554a0c // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/yuin/goldmark v1.7.1 // indirect
	gitlab.com/yawning/chacha20.git v0.0.0-20230427033715-7877545b1b37 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/mobile v0.0.0-20231127183840-76ac6878050a // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71 // indirect
	google.golang.org/grpc v1.49.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
