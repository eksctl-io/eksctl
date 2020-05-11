package nodebootstrap

// This file was generated 2020-03-26T16:18:43-07:00 by reserved_generate.go; DO NOT EDIT.

// Data downloaded through the API.
var instanceTypeInfos = map[string]InstanceTypeInfo{
	"a1.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"a1.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"a1.large": {
		CPU:     int64(2),
		Memory:  int64(4096),
		Storage: int64(0),
	},
	"a1.medium": {
		CPU:     int64(1),
		Memory:  int64(2048),
		Storage: int64(0),
	},
	"a1.metal": {
		CPU:     int64(16),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"a1.xlarge": {
		CPU:     int64(4),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"c1.medium": {
		CPU:     int64(2),
		Memory:  int64(1740),
		Storage: int64(350),
	},
	"c1.xlarge": {
		CPU:     int64(8),
		Memory:  int64(7168),
		Storage: int64(1680),
	},
	"c3.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(15360),
		Storage: int64(160),
	},
	"c3.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(30720),
		Storage: int64(320),
	},
	"c3.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(61440),
		Storage: int64(640),
	},
	"c3.large": {
		CPU:     int64(2),
		Memory:  int64(3840),
		Storage: int64(32),
	},
	"c3.xlarge": {
		CPU:     int64(4),
		Memory:  int64(7680),
		Storage: int64(80),
	},
	"c4.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(15360),
		Storage: int64(0),
	},
	"c4.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(30720),
		Storage: int64(0),
	},
	"c4.8xlarge": {
		CPU:     int64(36),
		Memory:  int64(61440),
		Storage: int64(0),
	},
	"c4.large": {
		CPU:     int64(2),
		Memory:  int64(3840),
		Storage: int64(0),
	},
	"c4.xlarge": {
		CPU:     int64(4),
		Memory:  int64(7680),
		Storage: int64(0),
	},
	"c5.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(98304),
		Storage: int64(0),
	},
	"c5.18xlarge": {
		CPU:     int64(72),
		Memory:  int64(147456),
		Storage: int64(0),
	},
	"c5.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"c5.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"c5.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"c5.9xlarge": {
		CPU:     int64(36),
		Memory:  int64(73728),
		Storage: int64(0),
	},
	"c5.large": {
		CPU:     int64(2),
		Memory:  int64(4096),
		Storage: int64(0),
	},
	"c5.metal": {
		CPU:     int64(96),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"c5.xlarge": {
		CPU:     int64(4),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"c5d.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(98304),
		Storage: int64(1800),
	},
	"c5d.18xlarge": {
		CPU:     int64(72),
		Memory:  int64(147456),
		Storage: int64(1800),
	},
	"c5d.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(196608),
		Storage: int64(3600),
	},
	"c5d.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(16384),
		Storage: int64(200),
	},
	"c5d.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(32768),
		Storage: int64(400),
	},
	"c5d.9xlarge": {
		CPU:     int64(36),
		Memory:  int64(73728),
		Storage: int64(900),
	},
	"c5d.large": {
		CPU:     int64(2),
		Memory:  int64(4096),
		Storage: int64(50),
	},
	"c5d.metal": {
		CPU:     int64(96),
		Memory:  int64(196608),
		Storage: int64(3600),
	},
	"c5d.xlarge": {
		CPU:     int64(4),
		Memory:  int64(8192),
		Storage: int64(100),
	},
	"c5n.18xlarge": {
		CPU:     int64(72),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"c5n.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(21504),
		Storage: int64(0),
	},
	"c5n.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(43008),
		Storage: int64(0),
	},
	"c5n.9xlarge": {
		CPU:     int64(36),
		Memory:  int64(98304),
		Storage: int64(0),
	},
	"c5n.large": {
		CPU:     int64(2),
		Memory:  int64(5376),
		Storage: int64(0),
	},
	"c5n.metal": {
		CPU:     int64(72),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"c5n.xlarge": {
		CPU:     int64(4),
		Memory:  int64(10752),
		Storage: int64(0),
	},
	"cc2.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(61952),
		Storage: int64(3360),
	},
	"d2.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(62464),
		Storage: int64(12288),
	},
	"d2.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(124928),
		Storage: int64(24576),
	},
	"d2.8xlarge": {
		CPU:     int64(36),
		Memory:  int64(249856),
		Storage: int64(49152),
	},
	"d2.xlarge": {
		CPU:     int64(4),
		Memory:  int64(31232),
		Storage: int64(6144),
	},
	"f1.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(999424),
		Storage: int64(3760),
	},
	"f1.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(124928),
		Storage: int64(470),
	},
	"f1.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(249856),
		Storage: int64(940),
	},
	"g2.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(15360),
		Storage: int64(60),
	},
	"g2.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(61440),
		Storage: int64(240),
	},
	"g3.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(499712),
		Storage: int64(0),
	},
	"g3.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(124928),
		Storage: int64(0),
	},
	"g3.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(249856),
		Storage: int64(0),
	},
	"g3s.xlarge": {
		CPU:     int64(4),
		Memory:  int64(31232),
		Storage: int64(0),
	},
	"g4dn.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(900),
	},
	"g4dn.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(900),
	},
	"g4dn.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(225),
	},
	"g4dn.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(225),
	},
	"g4dn.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(900),
	},
	"g4dn.metal": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(1800),
	},
	"g4dn.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(125),
	},
	"h1.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(16000),
	},
	"h1.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(2000),
	},
	"h1.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(4000),
	},
	"h1.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(8000),
	},
	"i2.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(62464),
		Storage: int64(1600),
	},
	"i2.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(124928),
		Storage: int64(3200),
	},
	"i2.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(249856),
		Storage: int64(6400),
	},
	"i2.xlarge": {
		CPU:     int64(4),
		Memory:  int64(31232),
		Storage: int64(800),
	},
	"i3.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(499712),
		Storage: int64(15200),
	},
	"i3.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(62464),
		Storage: int64(1900),
	},
	"i3.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(124928),
		Storage: int64(3800),
	},
	"i3.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(249856),
		Storage: int64(7600),
	},
	"i3.large": {
		CPU:     int64(2),
		Memory:  int64(15616),
		Storage: int64(475),
	},
	"i3.metal": {
		CPU:     int64(72),
		Memory:  int64(524288),
		Storage: int64(15200),
	},
	"i3.xlarge": {
		CPU:     int64(4),
		Memory:  int64(31232),
		Storage: int64(950),
	},
	"i3en.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(30000),
	},
	"i3en.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(60000),
	},
	"i3en.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(5000),
	},
	"i3en.3xlarge": {
		CPU:     int64(12),
		Memory:  int64(98304),
		Storage: int64(7500),
	},
	"i3en.6xlarge": {
		CPU:     int64(24),
		Memory:  int64(196608),
		Storage: int64(15000),
	},
	"i3en.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(1250),
	},
	"i3en.metal": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(60000),
	},
	"i3en.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(2500),
	},
	"inf1.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"inf1.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"inf1.6xlarge": {
		CPU:     int64(24),
		Memory:  int64(49152),
		Storage: int64(0),
	},
	"inf1.xlarge": {
		CPU:     int64(4),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"m1.large": {
		CPU:     int64(2),
		Memory:  int64(7680),
		Storage: int64(840),
	},
	"m1.medium": {
		CPU:     int64(1),
		Memory:  int64(3788),
		Storage: int64(410),
	},
	"m1.small": {
		CPU:     int64(1),
		Memory:  int64(1740),
		Storage: int64(160),
	},
	"m1.xlarge": {
		CPU:     int64(4),
		Memory:  int64(15360),
		Storage: int64(1680),
	},
	"m2.2xlarge": {
		CPU:     int64(4),
		Memory:  int64(35020),
		Storage: int64(850),
	},
	"m2.4xlarge": {
		CPU:     int64(8),
		Memory:  int64(70041),
		Storage: int64(1680),
	},
	"m2.xlarge": {
		CPU:     int64(2),
		Memory:  int64(17510),
		Storage: int64(420),
	},
	"m3.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(30720),
		Storage: int64(160),
	},
	"m3.large": {
		CPU:     int64(2),
		Memory:  int64(7680),
		Storage: int64(32),
	},
	"m3.medium": {
		CPU:     int64(1),
		Memory:  int64(3840),
		Storage: int64(4),
	},
	"m3.xlarge": {
		CPU:     int64(4),
		Memory:  int64(15360),
		Storage: int64(80),
	},
	"m4.10xlarge": {
		CPU:     int64(40),
		Memory:  int64(163840),
		Storage: int64(0),
	},
	"m4.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"m4.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"m4.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"m4.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"m4.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"m5.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"m5.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"m5.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"m5.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"m5.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"m5.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(0),
	},
	"m5.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"m5.metal": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"m5.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"m5a.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"m5a.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"m5a.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"m5a.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"m5a.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"m5a.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(0),
	},
	"m5a.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"m5a.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"m5ad.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(1800),
	},
	"m5ad.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(2400),
	},
	"m5ad.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(3600),
	},
	"m5ad.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(300),
	},
	"m5ad.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(600),
	},
	"m5ad.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(1200),
	},
	"m5ad.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(75),
	},
	"m5ad.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(150),
	},
	"m5d.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(1800),
	},
	"m5d.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(2400),
	},
	"m5d.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(3600),
	},
	"m5d.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(300),
	},
	"m5d.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(600),
	},
	"m5d.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(1200),
	},
	"m5d.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(75),
	},
	"m5d.metal": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(3600),
	},
	"m5d.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(150),
	},
	"m5dn.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(1800),
	},
	"m5dn.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(2400),
	},
	"m5dn.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(3600),
	},
	"m5dn.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(300),
	},
	"m5dn.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(600),
	},
	"m5dn.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(1200),
	},
	"m5dn.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(75),
	},
	"m5dn.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(150),
	},
	"m5n.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(196608),
		Storage: int64(0),
	},
	"m5n.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"m5n.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"m5n.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"m5n.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"m5n.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(131072),
		Storage: int64(0),
	},
	"m5n.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"m5n.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"p2.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(749568),
		Storage: int64(0),
	},
	"p2.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(499712),
		Storage: int64(0),
	},
	"p2.xlarge": {
		CPU:     int64(4),
		Memory:  int64(62464),
		Storage: int64(0),
	},
	"p3.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(499712),
		Storage: int64(0),
	},
	"p3.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(62464),
		Storage: int64(0),
	},
	"p3.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(249856),
		Storage: int64(0),
	},
	"p3dn.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(1800),
	},
	"r3.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(62464),
		Storage: int64(160),
	},
	"r3.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(124928),
		Storage: int64(320),
	},
	"r3.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(249856),
		Storage: int64(640),
	},
	"r3.large": {
		CPU:     int64(2),
		Memory:  int64(15360),
		Storage: int64(32),
	},
	"r3.xlarge": {
		CPU:     int64(4),
		Memory:  int64(31232),
		Storage: int64(80),
	},
	"r4.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(499712),
		Storage: int64(0),
	},
	"r4.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(62464),
		Storage: int64(0),
	},
	"r4.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(124928),
		Storage: int64(0),
	},
	"r4.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(249856),
		Storage: int64(0),
	},
	"r4.large": {
		CPU:     int64(2),
		Memory:  int64(15616),
		Storage: int64(0),
	},
	"r4.xlarge": {
		CPU:     int64(4),
		Memory:  int64(31232),
		Storage: int64(0),
	},
	"r5.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"r5.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(524288),
		Storage: int64(0),
	},
	"r5.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(0),
	},
	"r5.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"r5.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(131072),
		Storage: int64(0),
	},
	"r5.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"r5.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"r5.metal": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(0),
	},
	"r5.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"r5a.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"r5a.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(524288),
		Storage: int64(0),
	},
	"r5a.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(0),
	},
	"r5a.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"r5a.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(131072),
		Storage: int64(0),
	},
	"r5a.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"r5a.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"r5a.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"r5ad.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(1800),
	},
	"r5ad.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(524288),
		Storage: int64(2400),
	},
	"r5ad.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(3600),
	},
	"r5ad.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(300),
	},
	"r5ad.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(131072),
		Storage: int64(600),
	},
	"r5ad.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(262144),
		Storage: int64(1200),
	},
	"r5ad.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(75),
	},
	"r5ad.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(150),
	},
	"r5d.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(1800),
	},
	"r5d.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(524288),
		Storage: int64(2400),
	},
	"r5d.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(3600),
	},
	"r5d.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(300),
	},
	"r5d.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(131072),
		Storage: int64(600),
	},
	"r5d.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(262144),
		Storage: int64(1200),
	},
	"r5d.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(75),
	},
	"r5d.metal": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(3600),
	},
	"r5d.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(150),
	},
	"r5dn.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(1800),
	},
	"r5dn.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(524288),
		Storage: int64(2400),
	},
	"r5dn.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(3600),
	},
	"r5dn.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(300),
	},
	"r5dn.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(131072),
		Storage: int64(600),
	},
	"r5dn.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(262144),
		Storage: int64(1200),
	},
	"r5dn.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(75),
	},
	"r5dn.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(150),
	},
	"r5n.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(0),
	},
	"r5n.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(524288),
		Storage: int64(0),
	},
	"r5n.24xlarge": {
		CPU:     int64(96),
		Memory:  int64(786432),
		Storage: int64(0),
	},
	"r5n.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(0),
	},
	"r5n.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(131072),
		Storage: int64(0),
	},
	"r5n.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(262144),
		Storage: int64(0),
	},
	"r5n.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"r5n.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"t1.micro": {
		CPU:     int64(1),
		Memory:  int64(627),
		Storage: int64(0),
	},
	"t2.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"t2.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"t2.medium": {
		CPU:     int64(2),
		Memory:  int64(4096),
		Storage: int64(0),
	},
	"t2.micro": {
		CPU:     int64(1),
		Memory:  int64(1024),
		Storage: int64(0),
	},
	"t2.nano": {
		CPU:     int64(1),
		Memory:  int64(512),
		Storage: int64(0),
	},
	"t2.small": {
		CPU:     int64(1),
		Memory:  int64(2048),
		Storage: int64(0),
	},
	"t2.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"t3.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"t3.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"t3.medium": {
		CPU:     int64(2),
		Memory:  int64(4096),
		Storage: int64(0),
	},
	"t3.micro": {
		CPU:     int64(2),
		Memory:  int64(1024),
		Storage: int64(0),
	},
	"t3.nano": {
		CPU:     int64(2),
		Memory:  int64(512),
		Storage: int64(0),
	},
	"t3.small": {
		CPU:     int64(2),
		Memory:  int64(2048),
		Storage: int64(0),
	},
	"t3.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"t3a.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(32768),
		Storage: int64(0),
	},
	"t3a.large": {
		CPU:     int64(2),
		Memory:  int64(8192),
		Storage: int64(0),
	},
	"t3a.medium": {
		CPU:     int64(2),
		Memory:  int64(4096),
		Storage: int64(0),
	},
	"t3a.micro": {
		CPU:     int64(2),
		Memory:  int64(1024),
		Storage: int64(0),
	},
	"t3a.nano": {
		CPU:     int64(2),
		Memory:  int64(512),
		Storage: int64(0),
	},
	"t3a.small": {
		CPU:     int64(2),
		Memory:  int64(2048),
		Storage: int64(0),
	},
	"t3a.xlarge": {
		CPU:     int64(4),
		Memory:  int64(16384),
		Storage: int64(0),
	},
	"x1.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(999424),
		Storage: int64(1920),
	},
	"x1.32xlarge": {
		CPU:     int64(128),
		Memory:  int64(1998848),
		Storage: int64(3840),
	},
	"x1e.16xlarge": {
		CPU:     int64(64),
		Memory:  int64(1998848),
		Storage: int64(1920),
	},
	"x1e.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(249856),
		Storage: int64(240),
	},
	"x1e.32xlarge": {
		CPU:     int64(128),
		Memory:  int64(3997696),
		Storage: int64(3840),
	},
	"x1e.4xlarge": {
		CPU:     int64(16),
		Memory:  int64(499712),
		Storage: int64(480),
	},
	"x1e.8xlarge": {
		CPU:     int64(32),
		Memory:  int64(999424),
		Storage: int64(960),
	},
	"x1e.xlarge": {
		CPU:     int64(4),
		Memory:  int64(124928),
		Storage: int64(120),
	},
	"z1d.12xlarge": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(1800),
	},
	"z1d.2xlarge": {
		CPU:     int64(8),
		Memory:  int64(65536),
		Storage: int64(300),
	},
	"z1d.3xlarge": {
		CPU:     int64(12),
		Memory:  int64(98304),
		Storage: int64(450),
	},
	"z1d.6xlarge": {
		CPU:     int64(24),
		Memory:  int64(196608),
		Storage: int64(900),
	},
	"z1d.large": {
		CPU:     int64(2),
		Memory:  int64(16384),
		Storage: int64(75),
	},
	"z1d.metal": {
		CPU:     int64(48),
		Memory:  int64(393216),
		Storage: int64(1800),
	},
	"z1d.xlarge": {
		CPU:     int64(4),
		Memory:  int64(32768),
		Storage: int64(150),
	},
}
