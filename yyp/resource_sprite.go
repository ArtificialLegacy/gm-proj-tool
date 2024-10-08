package yyp

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
)

type SpriteLayer struct {
	Layers []SpriteLayer
	Frames []*image.NRGBA
	Name   string
}

type Sprite struct {
	Name     string
	Layers   []SpriteLayer
	Resource *ResourceSprite
}

func NewResourceSprite(name string, parent, textureGroup ProjectResourceNode, width, height int, layers []SpriteLayer) (*ResourceSprite, error) {
	var first SpriteLayer
	firstFound := false
	frameCount := 0
	for _, l := range layers {
		if len(l.Frames) > 0 && len(l.Layers) > 0 {
			return nil, fmt.Errorf("sprite layer cannot have both frames and layers")
		}

		if len(l.Frames) == 0 && len(l.Layers) == 0 {
			return nil, fmt.Errorf("sprite layer must have either frames or layers")
		}

		if len(l.Frames) > 0 {
			first = l
			firstFound = true

			if frameCount == 0 {
				frameCount = len(l.Frames)
			} else if len(l.Frames) != frameCount {
				return nil, fmt.Errorf("all sprite layers must have the same amount of frames")
			}
		}
	}

	if !firstFound {
		return nil, fmt.Errorf("no valid layer containing frames was found")
	}

	frames := make([]ResourceSpriteFrame, len(first.Frames))
	for i := range first.Frames {
		frames[i] = NewResourceSpriteFrame()
	}

	frameIds := make([]ProjectResourceNode, len(frames))
	for i, f := range frames {
		frameIds[i] = ProjectResourceNode{
			Name: f.Name,
			Path: path.Join(DIR_SPRITE, name, name+EXT_RESOURCE),
		}
	}

	resLayers := make([]ResourceImageLayer, len(layers))
	for i, l := range layers {
		resLayers[i] = makeLayers(l)
	}

	return &ResourceSprite{
		ResourceType:    RESTYPE_SPRITE,
		ResourceVersion: VERSION_SPRITE,
		Name:            name,
		Parent:          parent,
		TexGroupID:      textureGroup,

		Width:  width,
		Height: height,

		BBOX_Right:  width - 1,
		BBOX_Bottom: height - 1,

		Frames: frames,
		Layers: resLayers,

		Sequence: NewResourceSpriteSequence(name, frameIds),
	}, nil
}

func makeLayers(layer SpriteLayer) ResourceImageLayer {
	if len(layer.Frames) > 0 {
		return NewResourceImageLayer(layer.Name)
	}

	ls := make([]ResourceImageLayer, len(layer.Layers))
	for i, l := range layer.Layers {
		ls[i] = makeLayers(l)
	}
	return NewResourceImageLayerFolder(layer.Name, ls)
}

func NewSprite(name string, parent, textureGroup ProjectResourceNode, width, height int, layers []SpriteLayer) (*Sprite, error) {
	res, err := NewResourceSprite(name, parent, textureGroup, width, height, layers)

	return &Sprite{
		Name:     name,
		Layers:   layers,
		Resource: res,
	}, err
}

func (s *Sprite) Save(pdir string) (string, string, *ProjectResourceNode, error) {
	d := path.Join(pdir, DIR_SPRITE, s.Name)

	f, err := os.Stat(d)
	if err != nil {
		err := os.Mkdir(d, 0o777)
		if err != nil {
			return "", "", nil, err
		}
	} else {
		if !f.IsDir() {
			return "", "", nil, fmt.Errorf("path for sprite already exists, and it's not a directory")
		}
	}

	if _, err := os.Stat(path.Join(d, DIR_IMAGELAYERS)); err == nil {
		os.RemoveAll(path.Join(d, DIR_IMAGELAYERS))
	}

	fl, _ := os.ReadDir(d)
	for _, fli := range fl {
		if fli.IsDir() {
			continue
		}

		if path.Ext(fli.Name()) == EXT_SPRITE {
			os.Remove(path.Join(d, fli.Name()))
		}
	}

	bases := []*image.NRGBA{}
	layers, layerNames := flattenLayers(s.Layers, s.Resource.Layers)

	for li, layer := range layers {
		layerName := layerNames[li]
		for fi, frame := range layer.Frames {
			frameResource := s.Resource.Frames[fi]
			layerDir := path.Join(d, DIR_IMAGELAYERS, frameResource.Name)
			err := os.MkdirAll(layerDir, 0o777)
			if err != nil {
				return "", "", nil, err
			}

			ff, err := os.OpenFile(path.Join(layerDir, layerName+EXT_SPRITE), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
			if err != nil {
				return "", "", nil, err
			}
			defer ff.Close()
			err = png.Encode(ff, frame)
			if err != nil {
				return "", "", nil, err
			}

			if fi == len(bases) {
				bases = append(bases, frame)
			} else if fi < len(bases) {
				base := bases[fi]
				draw.Draw(frame, frame.Bounds(), base, base.Bounds().Min, draw.Over)
				bases[fi] = frame
			}
		}
	}

	for bi, base := range bases {
		frameResource := s.Resource.Frames[bi]

		ff, err := os.OpenFile(path.Join(d, frameResource.Name+EXT_SPRITE), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
		if err != nil {
			return "", "", nil, err
		}
		defer ff.Close()
		err = png.Encode(ff, base)
		if err != nil {
			return "", "", nil, err
		}
	}

	err = saveJSON(path.Join(d, s.Name+EXT_RESOURCE), s.Resource)
	if err != nil {
		return "", "", nil, err
	}

	return s.Name, path.Join(DIR_SPRITE, s.Name, s.Name+EXT_RESOURCE), &s.Resource.Parent, nil
}

func flattenLayers(layers []SpriteLayer, resourceLayers []ResourceImageLayer) ([]SpriteLayer, []string) {
	out := []SpriteLayer{}
	names := []string{}

	for li, layer := range layers {
		if len(layer.Frames) > 0 {
			out = append(out, layer)
			names = append(names, resourceLayers[li].Name)
		}

		if len(layer.Layers) > 0 {
			fOut, fNames := flattenLayers(layer.Layers, resourceLayers[li].Layers)
			out = append(out, fOut...)
			names = append(names, fNames...)
		}
	}

	return out, names
}

func (p *Project) SpriteExists(name string) bool {
	pth := path.Join(p.Path, DIR_SPRITE, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return false
	}
	if !fs.IsDir() {
		return false
	}

	return true
}

func (p *Project) SpriteLoad(name string) (*Sprite, error) {
	pth := path.Join(p.Path, DIR_SPRITE, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return nil, err
	}
	if !fs.IsDir() {
		return nil, fmt.Errorf("resource %s is not a directory", fs.Name())
	}

	data := &ResourceSprite{}
	err = loadJSON(path.Join(pth, name+EXT_RESOURCE), data)
	if err != nil {
		return nil, fmt.Errorf("error parsing sprite json: %s", err)
	}

	if data.ResourceVersion != VERSION_SPRITE {
		return nil, fmt.Errorf("resource version is unsupported: %s, expected %s", data.ResourceVersion, VERSION_SPRITE)
	}
	if data.ResourceType != RESTYPE_SPRITE {
		return nil, fmt.Errorf("resource is of incorrect type: %s, expected %s", data.ResourceType, RESTYPE_SPRITE)
	}

	if data.Type != SPRITETYPE_BITMAP {
		return nil, fmt.Errorf("sprite is of unsupported type: %d, expected %d", data.Type, SPRITETYPE_BITMAP)
	}

	frameList := make([]string, len(data.Frames))

	for i, f := range data.Frames {
		if f.ResourceVersion != VERSION_SPRITEFRAME {
			return nil, fmt.Errorf("resource version is unsupported: %s, expected %s", f.ResourceVersion, VERSION_SPRITEFRAME)
		}
		if f.ResourceType != RESTYPE_SPRITEFRAME {
			return nil, fmt.Errorf("resource is of incorrect type: %s, expected %s", f.ResourceType, RESTYPE_SPRITEFRAME)
		}

		frameList[i] = f.Name
	}

	sprite := Sprite{
		Name:     data.Name,
		Resource: data,
	}

	for _, l := range data.Layers {
		err := verifyLayers(l, &sprite.Layers, pth, frameList)
		if err != nil {
			return nil, err
		}
	}

	return &sprite, nil
}

func verifyLayers(layer ResourceImageLayer, layerList *[]SpriteLayer, pth string, frames []string) error {
	if layer.ResourceType == RESTYPE_IMAGEFOLDERLAYER {
		if layer.ResourceVersion != VERSION_IMAGEFOLDER {
			return fmt.Errorf("resource version is unsupported: %s, expected %s", layer.ResourceVersion, VERSION_IMAGEFOLDER)
		}

		folderLayer := SpriteLayer{
			Name:   layer.DisplayName,
			Layers: []SpriteLayer{},
		}

		for _, l := range layer.Layers {
			err := verifyLayers(l, &folderLayer.Layers, pth, frames)
			if err != nil {
				return err
			}
		}

		*layerList = append(*layerList, folderLayer)

		return nil
	}

	if layer.ResourceVersion != VERSION_IMAGELAYER {
		return fmt.Errorf("resource version is unsupported: %s, expected %s", layer.ResourceVersion, VERSION_IMAGELAYER)
	}
	if layer.ResourceType != RESTYPE_IMAGELAYER {
		return fmt.Errorf("resource is of incorrect type: %s, expected %s", layer.ResourceType, RESTYPE_IMAGELAYER)
	}

	layerFrame := SpriteLayer{
		Name:   layer.DisplayName,
		Frames: make([]*image.NRGBA, len(frames)),
	}

	for i, f := range frames {
		fdir := path.Join(pth, DIR_IMAGELAYERS, f)

		fd, err := os.Stat(fdir)
		if err != nil {
			return fmt.Errorf("failed to find sprite layer images, %s", err)
		}
		if !fd.IsDir() {
			return fmt.Errorf("frame layer path is not a directory, %s", fd.Name())
		}

		b, err := os.ReadFile(path.Join(fdir, layer.Name+EXT_SPRITE))
		if err != nil {
			return fmt.Errorf("failed to read sprite layer image, %s", err)
		}

		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("failed to decode sprite layer image, %s", err)
		}

		bounds := img.Bounds()
		m := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(m, m.Bounds(), img, bounds.Min, draw.Src)
		layerFrame.Frames[i] = m
	}

	*layerList = append(*layerList, layerFrame)
	return nil
}

func (p *Project) SpriteDelete(name string) error {
	pth := path.Join(p.Path, DIR_SPRITE, name)

	fs, err := os.Stat(pth)
	if err != nil {
		return nil // Trying to delete smth that doesn't exist shouldn't be an error
	}
	if !fs.IsDir() {
		return fmt.Errorf("resource %s is not a directory", fs.Name())
	}

	err = os.RemoveAll(pth)
	if err != nil {
		return fmt.Errorf("failed to delete resource directory: %s", err)
	}

	idPath := path.Join(DIR_SPRITE, name, name+EXT_RESOURCE)
	for i, r := range p.Data.Resources {
		if r.ID.Name == name && r.ID.Path == idPath {
			p.Data.Resources = append(p.Data.Resources[:i], p.Data.Resources[i+1:]...)
			break
		}
	}

	return nil
}

type ResourceSprite struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`

	GridX int `json:"gridX"`
	GridY int `json:"gridY"`

	HTile bool `json:"HTile"`
	VTile bool `json:"VTile"`

	Width  int `json:"width"`
	Height int `json:"height"`

	BBOX_Mode   BBOXMode `json:"bbox_mode"`
	BBOX_Bottom int      `json:"bbox_bottom"`
	BBOX_Left   int      `json:"bbox_left"`
	BBOX_Right  int      `json:"bbox_right"`
	BBOX_Top    int      `json:"bbox_top"`

	CollisionKind      CollMask `json:"collisionKind"`
	CollisionTolerance int      `json:"collisionTolerance"`

	DynamicTexturePage bool `json:"DynamicTexturePage"`
	EdgeFiltering      bool `json:"edgeFiltering"`
	For3D              bool `json:"For3D"`

	Sequence ResourceSpriteSequence `json:"sequence"`

	Frames []ResourceSpriteFrame `json:"frames"`

	Layers []ResourceImageLayer `json:"layers"`

	NineSlice *ResourceNineSlice `json:"nineSlice"`

	Origin SpriteOrigin `json:"origin"`

	Parent ProjectResourceNode `json:"parent"`

	Type             SpriteType `json:"type"`
	PreMultiplyAlpha bool       `json:"preMultiplyAlpha"`
	SwatchColors     []int      `json:"swatchColours"`
	SWFPrecision     float64    `json:"swfPrecision"`

	TexGroupID ProjectResourceNode `json:"textureGroupId"`

	Tags []string `json:"tags,omitempty"`
}

type ResourceSpriteFrame struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
}

func NewResourceSpriteFrame() ResourceSpriteFrame {
	return ResourceSpriteFrame{
		ResourceType:    RESTYPE_SPRITEFRAME,
		ResourceVersion: VERSION_SPRITEFRAME,

		Name: uuid.NewString(),
	}
}

type ResourceImageLayer struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`

	BlendMode   BlendMode            `json:"blendMode"`
	DisplayName string               `json:"displayName"`
	IsLocked    bool                 `json:"isLocked"`
	Opacity     float64              `json:"opacity"`
	Visible     bool                 `json:"visible"`
	Layers      []ResourceImageLayer `json:"layers,omitempty"`
}

const SPRITELAYER_DEFAULTNAME = "default"

func NewResourceImageLayer(name string) ResourceImageLayer {
	return ResourceImageLayer{
		ResourceType:    RESTYPE_IMAGELAYER,
		ResourceVersion: VERSION_IMAGELAYER,

		DisplayName: name,

		Name:    uuid.NewString(),
		Opacity: 100.0,
		Visible: true,
	}
}

func NewResourceImageLayerFolder(name string, layers []ResourceImageLayer) ResourceImageLayer {
	return ResourceImageLayer{
		ResourceType:    RESTYPE_IMAGEFOLDERLAYER,
		ResourceVersion: VERSION_IMAGEFOLDER,

		DisplayName: name,
		Layers:      layers,

		Name:    uuid.NewString(),
		Opacity: 100.0,
		Visible: true,
	}
}

type ResourceNineSlice struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Enabled bool `json:"enabled"`

	Top    int `json:"top"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
	Right  int `json:"right"`

	GuideColor     []int              `json:"guideColour"`
	HighlightColor int                `json:"highlightColour"`
	HighlightStyle NinesliceHighlight `json:"highlightStyle"`

	TileMode []NineSliceTile `json:"tileMode"`
}

const NINESLICECOLOR_GUIDE int = 4294902015
const NINESLICECOLOR_HIGHLIGHT int = 1728023040

func NewResourceNineSlice() *ResourceNineSlice {
	return &ResourceNineSlice{
		ResourceType:    RESTYPE_NINESLICE,
		ResourceVersion: VERSION_NINESLICE,

		Enabled: true,

		GuideColor: []int{
			NINESLICECOLOR_GUIDE,
			NINESLICECOLOR_GUIDE,
			NINESLICECOLOR_GUIDE,
			NINESLICECOLOR_GUIDE,
		},

		HighlightColor: NINESLICECOLOR_HIGHLIGHT,

		TileMode: []NineSliceTile{
			NINESLICETILE_STRETCH,
			NINESLICETILE_STRETCH,
			NINESLICETILE_STRETCH,
			NINESLICETILE_STRETCH,
			NINESLICETILE_STRETCH,
		},
	}
}

type ResourceSpriteSequence struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`

	AutoRecord bool `json:"autoRecord"`

	BackdropHeight       int     `json:"backdropHeight"`
	BackdropImageOpacity float64 `json:"backdropImageOpacity"`
	BackdropImagePath    string  `json:"backdropImagePath"`
	BackdropWidth        int     `json:"backdropWidth"`
	BackdropXOffset      float64 `json:"backdropXOffset"`
	BackdropYOffset      float64 `json:"backdropYOffset"`
	ShowBackdrop         bool    `json:"showBackdrop"`
	ShowBackdropImage    bool    `json:"showBackdropImage"`

	Events  ResourceSpriteSequenceEvent `json:"events"`
	Moments ResourceSpriteSequenceEvent `json:"moments"`

	EventStubScript any          `json:"eventStubScript"`
	EventToFunction ProjectEmpty `json:"eventToFunction"`

	LockOrigin bool `json:"lockOrigin"`
	XOrigin    int  `json:"xorigin"`
	YOrigin    int  `json:"yorigin"`

	Length            float64         `json:"length"`
	Playback          int             `json:"playback"`
	PlaybackSpeed     float64         `json:"playbackSpeed"`
	PlaybackSpeedType SeqPlaybackType `json:"playbackSpeedType"`
	TimeUnits         SeqTimeUnits    `json:"timeUnits"`

	VisibleRange any     `json:"visibleRange"`
	Volume       float64 `json:"volume"`

	Tracks []ResourceSpriteSequenceTrack `json:"tracks"`
}

func NewResourceSpriteSequence(name string, frames []ProjectResourceNode) ResourceSpriteSequence {
	return ResourceSpriteSequence{
		ResourceType:    RESTYPE_SEQ,
		ResourceVersion: VERSION_SEQ,

		Name: name,

		Tracks: []ResourceSpriteSequenceTrack{
			NewResourceSpriteSequenceTrack(frames),
		},

		AutoRecord: true,

		Events:  NewResourceSpriteSequenceEvent(),
		Moments: NewResourceSpriteSequenceEvent(),

		ShowBackdrop:  true,
		Volume:        1.0,
		PlaybackSpeed: 30.0,
	}
}

type ResourceSpriteSequenceTrack struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`
	Name            string       `json:"name"`
	BuiltinName     int          `json:"builtinName"`

	Events             []any    `json:"events"`
	InheritsTrackColor bool     `json:"inheritsTrackColour"`
	Interpolation      SeqInter `json:"interpolation"`
	IsCreationTrack    bool     `json:"isCreationTrack"`

	Keyframes ResourceSpriteSequenceKeyframeStore `json:"keyframes"`

	SpriteID   any `json:"spriteId"`
	TrackColor int `json:"trackColour"`
	Traits     int `json:"traits"`

	Modifiers []ProjectResourceBasic        `json:"modifiers"`
	Tracks    []ResourceSpriteSequenceTrack `json:"tracks"`
}

const SPRITESEQ_DEFAULTNAME = "frames"

func NewResourceSpriteSequenceTrack(frames []ProjectResourceNode) ResourceSpriteSequenceTrack {
	return ResourceSpriteSequenceTrack{
		ResourceType:    RESTYPE_SPRITEFRAMETRACK,
		ResourceVersion: VERSION_SPRITETRACK,

		Name:               SPRITESEQ_DEFAULTNAME,
		InheritsTrackColor: true,
		Interpolation:      SEQINTER_LERP,

		Modifiers: []ProjectResourceBasic{},
		Tracks:    []ResourceSpriteSequenceTrack{},
		Events:    []any{},

		Keyframes: NewResourceSpriteSequenceKeyframeStore(frames),
	}
}

type ResourceSpriteSequenceKeyframeStore struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Keyframes []ResourceSpriteSequenceKeyframe `json:"Keyframes"`
}

func NewResourceSpriteSequenceKeyframeStore(frames []ProjectResourceNode) ResourceSpriteSequenceKeyframeStore {
	keyframes := []ResourceSpriteSequenceKeyframe{}

	for i, f := range frames {
		keyframes = append(keyframes, NewResourceSpriteSequenceKeyframe(i, f))
	}

	return ResourceSpriteSequenceKeyframeStore{
		ResourceType:    RESTYPE_KEYFRAMESTORE_SPRITEKEYFRAME,
		ResourceVersion: VERSION_KEYFRAMESTORE_SPRITEKEYFRAME,

		Keyframes: keyframes,
	}
}

type ResourceSpriteSequenceKeyframe struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Channels map[string]ResourceSpriteSequenceKeyframeChannel `json:"Channels"`

	Disabled      bool    `json:"Disabled"`
	ID            string  `json:"id"`
	IsCreationKey bool    `json:"IsCreationKey"`
	Key           float64 `json:"Key"`
	Length        float64 `json:"Length"`
	Stretch       bool    `json:"Stretch"`
}

func NewResourceSpriteSequenceKeyframe(key int, id ProjectResourceNode) ResourceSpriteSequenceKeyframe {
	return ResourceSpriteSequenceKeyframe{
		ResourceType:    RESTYPE_KEYFRAME_SPRITEKEYFRAME,
		ResourceVersion: VERSION_KEYFRAME_SPRITEKEYFRAME,

		Key: float64(key),

		Channels: map[string]ResourceSpriteSequenceKeyframeChannel{
			"0": NewResourceSpriteSequenceKeyframeChannel(id),
		},

		Length: 1.0,

		ID: uuid.NewString(),
	}
}

type ResourceSpriteSequenceKeyframeChannel struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	ID ProjectResourceNode `json:"Id"`
}

func NewResourceSpriteSequenceKeyframeChannel(id ProjectResourceNode) ResourceSpriteSequenceKeyframeChannel {
	return ResourceSpriteSequenceKeyframeChannel{
		ResourceType:    RESTYPE_SPRITEKEYFRAME,
		ResourceVersion: VERSION_SPRITEKEYFRAME,

		ID: id,
	}
}

type ResourceSpriteSequenceEvent struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Keyframes []ResourceSpriteSequenceEventKeyframe `json:"Keyframes"`
}

func NewResourceSpriteSequenceEvent() ResourceSpriteSequenceEvent {
	return ResourceSpriteSequenceEvent{
		ResourceType:    RESTYPE_KEYFRAMESTORE_MESSAGEEVENT,
		ResourceVersion: VERSION_KEYFRAMESTORE_MESSAGEEVENT,

		Keyframes: []ResourceSpriteSequenceEventKeyframe{},
	}
}

type ResourceSpriteSequenceEventKeyframe struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Channels map[string]ResourceSpriteSequenceEventChannel `json:"Channels"`

	Disabled      bool    `json:"Disabled"`
	ID            string  `json:"id"`
	IsCreationKey bool    `json:"IsCreationKey"`
	Key           float64 `json:"Key"`
	Length        float64 `json:"Length"`
	Stretch       bool    `json:"Stretch"`
}

func NewResourceSpriteSequenceEventKeyframe(events [][]string, key float64) ResourceSpriteSequenceEventKeyframe {
	channels := map[string]ResourceSpriteSequenceEventChannel{}

	for i, ev := range events {
		channels[strconv.Itoa(i)] = NewResourceSpriteSequenceEventChannel(ev)
	}

	return ResourceSpriteSequenceEventKeyframe{
		ResourceType:    RESTYPE_KEYFRAME_MESSAGEEVENT,
		ResourceVersion: VERSION_KEYFRAME_MESSAGEEVENT,

		Channels: channels,

		Key:    key,
		Length: 1.0,

		ID: uuid.NewString(),
	}
}

type ResourceSpriteSequenceEventChannel struct {
	ResourceType    ResourceType `json:"resourceType"`
	ResourceVersion Version      `json:"resourceVersion"`

	Events []string `json:"Events"`
}

func NewResourceSpriteSequenceEventChannel(events []string) ResourceSpriteSequenceEventChannel {
	return ResourceSpriteSequenceEventChannel{
		ResourceType:    RESTYPE_MESSAGEEVENTKEYFRAME,
		ResourceVersion: VERSION_MESSAGEEVENTKEYFRAME,

		/// Wrong
		Events: events,
	}
}

type BBOXMode int

const (
	BBOXMODE_AUTO BBOXMode = iota
	BBOXMODE_FULL
	BBOXMODE_MANUAL
)

type CollMask int

const (
	COLLMASK_PRECISE CollMask = iota
	COLLMASK_RECT
	COLLMASK_ELLIPSE
	COLLMASK_DIAMOND
	COLLMASK_PRECISEFRAME
	COLLMASK_RECTROT
	COLLMASK_SPINE
)

type BlendMode int

const (
	BLENDMODE_NORMAL BlendMode = iota
	BLENDMODE_ADD
	BLENDMODE_SUBTRACT
	BLENDMODE_MULTIPLY
)

type NinesliceHighlight int

const (
	NINESLICEHIGHLIGHT_INVERTED NinesliceHighlight = iota
	NINESLICEHIGHLIGHT_OVERLAY
)

type NineSliceTile int

const (
	NINESLICETILE_STRETCH NineSliceTile = iota
	NINESLICETILE_REPEAT
	NINESLICETILE_MIRROR
	NINESLICETILE_BLANKREPEAT
	NINESLICETILE_HIDE
)

const (
	NINESLICESLICE_LEFT int = iota
	NINESLICESLICE_TOP
	NINESLICESLICE_RIGHT
	NINESLICESLICE_BOTTOM
	NINESLICESLICE_CENTER
)

type SpriteOrigin int

const (
	SPRITEORIGIN_TOPLEFT SpriteOrigin = iota
	SPRITEORIGIN_TOPCENTER
	SPRITEORIGIN_TOPRIGHT
	SPRITEORIGIN_MIDDLELEFT
	SPRITEORIGIN_MIDDLECENTER
	SPRITEORIGIN_MIDDLERIGHT
	SPRITEORIGIN_BOTTOMLEFT
	SPRITEORIGIN_BOTTOMCENTER
	SPRITEORIGIN_BOTTOMRIGHT
	SPRITEORIGIN_CUSTOM
)

type SpriteType int

const (
	SPRITETYPE_BITMAP SpriteType = iota
	SPRITETYPE_SWF
	SPRITETYPE_SPINE
)

type SeqPlaybackType int

const (
	SEQPLAYBACK_NORMAL SeqPlaybackType = iota
	SEQPLAYBACK_LOOP
	SEQPLAYBACK_PING
)

type SeqTimeUnits int

const (
	SEQUNITS_TIME SeqTimeUnits = iota
	SEQUNITS_FRAME
)

type SeqInter int

const (
	SEQINTER_ASSIGN SeqInter = iota
	SEQINTER_LERP
)

const TRACKMOD_LOCK = "LockedModifier"
const TRACKMOD_INVIS = "InvisibleModifier"
