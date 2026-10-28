package filesize

const (
	video    int64 = 2 * 1024 * 1024 * 1024
	layered  int64 = 800 * 1024 * 1024
	category3d int64 = 500 * 1024 * 1024
	audio    int64 = 100 * 1024 * 1024
	archive  int64 = 500 * 1024 * 1024
	document int64 = 25 * 1024 * 1024
	image    int64 = 50 * 1024 * 1024
	fallback int64 = 100 * 1024 * 1024
)

// extensionLimits maps file extensions (without leading dot) to their maximum
// allowed upload size in bytes.
// Keep in sync with frontend/src/features/files/lib/fileCategory.ts FILE_SIZE_LIMITS — same values.
var extensionLimits = map[string]int64{
	// Video
	"mp4": video, "mov": video, "webm": video, "mkv": video, "m4v": video,
	"avi": video, "wmv": video, "flv": video, "ogv": video, "ts": video,

	// Layered / project
	"psd": layered, "clip": layered, "kra": layered, "sai": layered, "sai2": layered,
	"xcf": layered, "procreate": layered, "ai": layered, "afphoto": layered,
	"psb": layered, "csp": layered, "afdesign": layered, "pdn": layered,
	"ase": layered, "aseprite": layered,

	// 3D / rigged model
	"blend": category3d, "fbx": category3d, "obj": category3d, "glb": category3d,
	"gltf": category3d, "vrm": category3d, "moc3": category3d, "cmo3": category3d,
	"vox": category3d, "dae": category3d, "stl": category3d,

	// Audio
	"mp3": audio, "wav": audio, "flac": audio, "m4a": audio, "ogg": audio,
	"aac": audio, "wma": audio, "opus": audio, "aiff": audio, "alac": audio,

	// Archive
	"zip": archive, "rar": archive, "7z": archive, "tar": archive, "gz": archive,
	"tgz": archive, "bz2": archive, "xz": archive, "zst": archive,

	// Document
	"pdf": document, "txt": document, "md": document, "doc": document, "docx": document,
	"rtf": document, "odt": document, "pages": document,

	// Image
	"jpg": image, "jpeg": image, "png": image, "gif": image, "svg": image,
	"webp": image, "tiff": image, "tif": image, "bmp": image, "heic": image,
	"heif": image, "avif": image, "ico": image,
}

// MaxBytesForExtension returns the maximum allowed upload size for the given
// file extension (without leading dot). Unknown extensions return the fallback
// limit of 100 MB.
func MaxBytesForExtension(ext string) int64 {
	if limit, ok := extensionLimits[ext]; ok {
		return limit
	}
	return fallback
}
