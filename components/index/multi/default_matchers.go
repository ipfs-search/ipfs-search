package multi

import "regexp"

var DefaultMatchers = []*RegexpMatcher{
	documentMatcher,
	imageMatcher,
	dataMatcher,
	archiveMatcher,
	unknownMatcher,
	videoMatcher,
	audioMatcher,
}

func PrefixRegex(prefix string) string {
	return "^" + regexp.QuoteMeta(prefix)
}

var documentMatcher = NewRegexpMatcher(
	"ipfs_documents",
	[]string{
		PrefixRegex("text/x-web-markdown"),
		PrefixRegex("text/x-rst"),
		PrefixRegex("text/x-log"),
		PrefixRegex("text/x-asciidoc"),
		PrefixRegex("text/troff"),
		PrefixRegex("text/plain"),
		PrefixRegex("text/html"),
		PrefixRegex("message/rfc822"),
		PrefixRegex("message/news"),
		PrefixRegex("image/vnd.djvu"), // This requires documents to be matched before images.
		PrefixRegex("application/xhtml+xml"),
		PrefixRegex("application/x-tika-ooxml"),
		PrefixRegex("application/x-tika-msoffice"),
		PrefixRegex("application/x-tex"),
		PrefixRegex("application/x-mobipocket-ebook"),
		PrefixRegex("application/x-fictionbook+xml"),
		PrefixRegex("application/x-dvi"),
		PrefixRegex("application/vnd.sun.xml.writer.global"),
		PrefixRegex("application/vnd.openxmlformats-officedocument.wordprocessingml.document"),
		PrefixRegex("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"),
		PrefixRegex("application/vnd.openxmlformats-officedocument.presentationml.presentation"),
		PrefixRegex("application/vnd.oasis.opendocument.text"),
		PrefixRegex("application/vnd.ms-powerpoint"),
		PrefixRegex("application/vnd.ms-htmlhelp"),
		PrefixRegex("application/vnd.ms-excel"),
		PrefixRegex("application/vnd.sun.xml.draw"),
		PrefixRegex("application/rtf"),
		PrefixRegex("application/postscript"),
		PrefixRegex("application/pdf"),
		PrefixRegex("application/msword5"),
		PrefixRegex("application/msword2"),
		PrefixRegex("application/msword"),
		PrefixRegex("application/epub+zip"),
	},
)

var imageMatcher = NewRegexpMatcher(
	"ipfs_images",
	[]string{
		PrefixRegex("image"),
		PrefixRegex("application/dicom"),
	},
)

var dataMatcher = NewRegexpMatcher(
	"ipfs_data",
	[]string{
		PrefixRegex("text/x-yaml"),
		PrefixRegex("text/x-java-properties"),
		PrefixRegex("text/x-ini"),
		PrefixRegex("text/tsv"),
		PrefixRegex("text/tab-separated-values"),
		PrefixRegex("text/csv"),
		PrefixRegex("multipart/appledouble"),
		PrefixRegex("application/xml"),
		PrefixRegex("application/x-sqlite3"),
		PrefixRegex("application/x-plist"),
		PrefixRegex("application/rdf+xml"),
		PrefixRegex("application/json"),
	},
)

var archiveMatcher = NewRegexpMatcher(
	"ipfs_archives",
	[]string{
		PrefixRegex("application/gzip"),
		PrefixRegex("application/zlib"),
		PrefixRegex("application/x-xz"),
		PrefixRegex("application/x-bzip2"),
		PrefixRegex("application/x-brotli"),
		PrefixRegex("application/x-lzma"),
		PrefixRegex("application/x-compress"),
		PrefixRegex("application/x-lz4"),
		PrefixRegex("application/x-snappy"),
		PrefixRegex("application/x-java-pack200"),
		PrefixRegex("application/zip"),
		PrefixRegex("application/x-rpm"),
		PrefixRegex("application/zstd"),
		PrefixRegex("application/x-archive"),
		PrefixRegex("application/x-tar"),
		PrefixRegex("application/x-rar-compressed"),
		PrefixRegex("application/x-7z-compressed"),
		PrefixRegex("application/java-archive"),
		PrefixRegex("application/x-iso9660-image"),
		PrefixRegex("application/vnd.android.package-archive"),
		PrefixRegex("application/x-lha"),
		PrefixRegex("application/x-apple-diskimage"),
		PrefixRegex("application/x-gtar"),
		PrefixRegex("application/x-lharc"),
		PrefixRegex("application/vnd.ms-cab-compressed"),
		PrefixRegex("application/x-tika-ooxml"),
		PrefixRegex("application/x-arj"),
		PrefixRegex("application/x-tika-java-web-archive"),
		PrefixRegex("application/x-cpio"),
		PrefixRegex("application/x-itunes-ipa"),
		PrefixRegex("application/vnd.google-earth.kmz"),
		PrefixRegex("application/x-xmind"),
		PrefixRegex("application/vnd.adobe.air-application-installer-package+zip"),
		PrefixRegex("application/vnd.etsi.asic-e+zip"),
	},
)

var unknownMatcher = NewRegexpMatcher(
	"ipfs_unknown",
	[]string{
		"^$", // Unknown matches absent content type.
	},
)

var videoMatcher = NewRegexpMatcher(
	"ipfs_videos",
	[]string{
		PrefixRegex("video"),
		PrefixRegex("application/x-matroska"),
		PrefixRegex("application/mp4"),
	},
)

var audioMatcher = NewRegexpMatcher(
	"ipfs_audio",
	[]string{
		PrefixRegex("audio"),
		PrefixRegex("application/ogg"),
	},
)
