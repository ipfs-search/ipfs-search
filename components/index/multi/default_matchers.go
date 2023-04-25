package multi

var DefaultMatchers = []*RegexpMatcher{
	documentMatcher,
	imageMatcher,
	dataMatcher,
	archiveMatcher,
	unknownMatcher,
	videoMatcher,
	audioMatcher,
}

var documentMatcher = NewRegexpMatcher(
	"ipfs_documents",
	[]string{
		"^text/x-web-markdown",
		"^text/x-rst",
		"^text/x-log",
		"^text/x-asciidoc",
		"^text/troff",
		"^text/plain",
		"^text/html",
		"^message/rfc822",
		"^message/news",
		"^image/vnd.djvu", // This requires documents to be matched before images.
		"^application/xhtml+xml",
		"^application/x-tika-ooxml",
		"^application/x-tika-msoffice",
		"^application/x-tex",
		"^application/x-mobipocket-ebook",
		"^application/x-fictionbook+xml",
		"^application/x-dvi",
		"^application/vnd.sun.xml.writer.global",
		"^application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"^application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"^application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"^application/vnd.oasis.opendocument.text",
		"^application/vnd.ms-powerpoint",
		"^application/vnd.ms-htmlhelp",
		"^application/vnd.ms-excel",
		"^application/vnd.sun.xml.draw",
		"^application/rtf",
		"^application/postscript",
		"^application/pdf",
		"^application/msword5",
		"^application/msword2",
		"^application/msword",
		"^application/epub+zip",
	},
)

var imageMatcher = NewRegexpMatcher(
	"ipfs_images",
	[]string{
		"^image",
		"^application/dicom",
	},
)

var dataMatcher = NewRegexpMatcher(
	"ipfs_data",
	[]string{
		"^text/x-yaml",
		"^text/x-java-properties",
		"^text/x-ini",
		"^text/tsv",
		"^text/tab-separated-values",
		"^text/csv",
		"^multipart/appledouble",
		"^application/xml",
		"^application/x-sqlite3",
		"^application/x-plist",
		"^application/rdf+xml",
		"^application/json",
	},
)

var archiveMatcher = NewRegexpMatcher(
	"ipfs_archives",
	[]string{
		"^application/gzip",
		"^application/zlib",
		"^application/x-xz",
		"^application/x-bzip2",
		"^application/x-brotli",
		"^application/x-lzma",
		"^application/x-compress",
		"^application/x-lz4",
		"^application/x-snappy",
		"^application/x-java-pack200",
		"^application/zip",
		"^application/x-rpm",
		"^application/zstd",
		"^application/x-archive",
		"^application/x-tar",
		"^application/x-rar-compressed",
		"^application/x-7z-compressed",
		"^application/java-archive",
		"^application/x-iso9660-image",
		"^application/vnd.android.package-archive",
		"^application/x-lha",
		"^application/x-apple-diskimage",
		"^application/x-gtar",
		"^application/x-lharc",
		"^application/vnd.ms-cab-compressed",
		"^application/x-tika-ooxml",
		"^application/x-arj",
		"^application/x-tika-java-web-archive",
		"^application/x-cpio",
		"^application/x-itunes-ipa",
		"^application/vnd.google-earth.kmz",
		"^application/x-xmind",
		"^application/vnd.adobe.air-application-installer-package+zip",
		"^application/vnd.etsi.asic-e+zip",
	},
)

var unknownMatcher = NewRegexpMatcher(
	"ipfs_unknown",
	[]string{
		"^$",
	},
)

var videoMatcher = NewRegexpMatcher(
	"ipfs_videos",
	[]string{
		"^video",
		"^application/x-matroska",
		"^application/mp4",
	},
)

var audioMatcher = NewRegexpMatcher(
	"ipfs_audio",
	[]string{
		"^audio",
		"^application/ogg",
	},
)
