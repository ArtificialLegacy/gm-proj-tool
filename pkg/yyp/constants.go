package yyp

type Version string

const (
	VERSION_PROJ                         Version = "1.7"
	VERSION_ANIMCURVE                    Version = "1.2"
	VERSION_ANIMCURVECHAN                Version = "1.0"
	VERSION_EXT                          Version = "1.2"
	VERSION_FONT                         Version = "1.0"
	VERSION_NOTE                         Version = "1.1"
	VERSION_OBJ                          Version = "1.0"
	VERSION_OBJEVENT                     Version = "1.0"
	VERSION_PARTICLE                     Version = "1.0"
	VERSION_PARTICLEEMIT                 Version = "1.0"
	VERSION_PATH                         Version = "1.0"
	VERSION_ROOM                         Version = "1.0"
	VERSION_SCRIPT                       Version = "1.0"
	VERSION_SEQ                          Version = "1.4"
	VERSION_SEQEVENT                     Version = "1.0"
	VERSION_SEQMOMENT                    Version = "1.0"
	VERSION_SHADER                       Version = "1.0"
	VERSION_SOUND                        Version = "1.0"
	VERSION_SPRITE                       Version = "1.0"
	VERSION_SPRITEFRAME                  Version = "1.1"
	VERSION_IMAGELAYER                   Version = "1.0"
	VERSION_IMAGEFOLDER                  Version = "1.0"
	VERSION_SPRITETRACK                  Version = "1.0"
	VERSION_TILESET                      Version = "1.0"
	VERSION_TIMELINE                     Version = "1.0"
	VERSION_FOLDER                       Version = "1.0"
	VERSION_INCLUDEDFILE                 Version = "1.0"
	VERSION_AUDIOGROUP                   Version = "1.3"
	VERSION_TEXGROUP                     Version = "1.3"
	VERSION_NINESLICE                    Version = "1.0"
	VERSION_KEYFRAMESTORE_MESSAGEEVENT   Version = "1.0"
	VERSION_KEYFRAME_MESSAGEEVENT        Version = "1.0"
	VERSION_MESSAGEEVENTKEYFRAME         Version = "1.0"
	VERSION_KEYFRAMESTORE_SPRITEKEYFRAME Version = "1.0"
	VERSION_KEYFRAME_SPRITEKEYFRAME      Version = "1.0"
	VERSION_SPRITEKEYFRAME               Version = "1.0"
)

// resource directories
const DIR_ANIMCURVE = "animcurves"
const DIR_EXT = "extensions"
const DIR_FONT = "fonts"
const DIR_NOTE = "notes"
const DIR_OBJ = "objects"
const DIR_PARTICLE = "particles"
const DIR_PATH = "paths"
const DIR_ROOM = "rooms"
const DIR_SCRIPT = "scripts"
const DIR_SEQ = "sequences"
const DIR_SHADER = "shaders"
const DIR_SOUND = "sounds"
const DIR_SPRITE = "sprites"
const DIR_TILESET = "tilesets"
const DIR_TIMELINE = "timelines"
const DIR_DATAFILE = "datafiles"
const DIR_PARTICLELIB = "particlelib"

// project file extensions
const EXT_PROJ = ".yyp"
const EXT_ORDER = ".resource_order"
const EXT_RESOURCE = ".yy"
const EXT_SCRIPT = ".gml"
const EXT_SHADER_VERTEX = ".vsh"
const EXT_SHADER_FRAGMENT = ".fsh"
const EXT_NOTE = ".txt"

// project options directories
const OPT = "options"
const OPT_ANDROID = "android"
const OPT_EXT = "extensions"
const OPT_HTML = "html5"
const OPT_IOS = "ios"
const OPT_LINUX = "linux"
const OPT_MAC = "mac"
const OPT_MAIN = "main"
const OPT_OPERAGX = "operagx"
const OPT_TV = "tvos"
const OPT_WINDOWS = "windows"

type ResourceType string

const (
	RESTYPE_PROJ                         ResourceType = "GMProject"
	RESTYPE_AUDIOGROUP                   ResourceType = "GMAudioGroup"
	RESTYPE_TEXGROUP                     ResourceType = "GMTextureGroup"
	RESTYPE_FOLDER                       ResourceType = "GMFolder"
	RESTYPE_INCLUDEDFILE                 ResourceType = "GMIncludedFile"
	RESTYPE_NOTE                         ResourceType = "GMNotes"
	RESTYPE_SPRITE                       ResourceType = "GMSprite"
	RESTYPE_SPRITEFRAME                  ResourceType = "GMSpriteFrame"
	RESTYPE_IMAGELAYER                   ResourceType = "GMImageLayer"
	RESTYPE_SEQ                          ResourceType = "GMSequence"
	RESTYPE_SPRITEFRAMETRACK             ResourceType = "GMSpriteFramesTrack"
	RESTYPE_IMAGEFOLDERLAYER             ResourceType = "GMImageFolderLayer"
	RESTYPE_NINESLICE                    ResourceType = "GMNineSliceData"
	RESTYPE_AUDIOTRACK                   ResourceType = "GMAudioTrack"
	RESTYPE_REALTRACK                    ResourceType = "GMRealTrack"
	RESTYPE_PARTICLETRACK                ResourceType = "GMParticleTrack"
	RESTYPE_INSTTRACK                    ResourceType = "GMInstanceTrack"
	RESTYPE_GRAPHICTRACK                 ResourceType = "GMGraphicTrack"
	RESTYPE_TEXTTRACK                    ResourceType = "GMTextTrack"
	RESTYPE_ANIMCURVE                    ResourceType = "GMAnimCurve"
	RESTYPE_ANIMCURVECHAN                ResourceType = "GMAnimCurveChannel"
	RESTYPE_GROUPTRACK                   ResourceType = "GMGroupTrack"
	RESTYPE_CLIPMASKTRACK                ResourceType = "GMClipMaskTrack"
	RESTYPE_CLIPMASKMASK                 ResourceType = "GMClipMask_Mask"
	RESTYPE_CLIPMASKSUBJECT              ResourceType = "GMClipMask_Subject"
	RESTYPE_KEYFRAMESTORE_MESSAGEEVENT   ResourceType = "KeyframeStore<MessageEventKeyframe>"
	RESTYPE_KEYFRAME_MESSAGEEVENT        ResourceType = "Keyframe<MessageEventKeyframe>"
	RESTYPE_MESSAGEEVENTKEYFRAME         ResourceType = "MessageEventKeyframe"
	RESTYPE_KEYFRAMESTORE_SPRITEKEYFRAME ResourceType = "KeyframeStore<SpriteFrameKeyframe>"
	RESTYPE_KEYFRAME_SPRITEKEYFRAME      ResourceType = "Keyframe<SpriteFrameKeyframe>"
	RESTYPE_SPRITEKEYFRAME               ResourceType = "SpriteFrameKeyFrame"
)

type ScriptType int

const (
	SCRIPTTYPE_ASK ScriptType = iota
	SCRIPTTYPE_GML
	SCRIPTTYPE_VISUAL
)

type TemplateType string

const (
	TEMPLATETYPE_GAME      TemplateType = "game"
	TEMPLATETYPE_WALLPAPER TemplateType = "live_wallpaper"
	TEMPLATETYPE_STRIP     TemplateType = "game_strip"
)

type TextureGroupType string

const (
	TEXGROUPTYPE_DEFAULT TextureGroupType = "default"
	TEXGROUPTYPE_DYNAMIC TextureGroupType = "dynamicpages"
)

type TextureGroupCompression string

const (
	TEXGROUPCOMPRESS_BZ2 TextureGroupCompression = "bz2"
	TEXGROUPCOMPRESS_QOI TextureGroupCompression = "qoi"
	TEXGROUPCOMPRESS_PNG TextureGroupCompression = "png"
)
