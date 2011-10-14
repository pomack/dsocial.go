package dsocial

type ImageService uint32
type ImageSize uint8

const (
    IMAGE_SERVICE_SELF      ImageService = iota
    IMAGE_SERVICE_AWS_S3    ImageService = iota
    IMAGE_SERVICE_RACKSPACE ImageService = iota
    IMAGE_SERVICE_SMUGMUG   ImageService = iota
    IMAGE_SERVICE_FLICKR    ImageService = iota
    IMAGE_SERVICE_PICASA    ImageService = iota
)

const (
    IMAGE_SIZE_ICON            ImageSize = iota
    IMAGE_SIZE_QUICK_VIEW      ImageSize = iota
    IMAGE_SIZE_MAIN_IN_GALLERY ImageSize = iota
    IMAGE_SIZE_MAIN            ImageSize = iota
    IMAGE_SIZE_FULL_SCREEN     ImageSize = iota
    IMAGE_SIZE_VERY_SMALL      ImageSize = iota
    IMAGE_SIZE_SMALL           ImageSize = iota
    IMAGE_SIZE_MEDIUM          ImageSize = iota
    IMAGE_SIZE_LARGE           ImageSize = iota
    IMAGE_SIZE_VERY_LARGE      ImageSize = iota
    IMAGE_SIZE_ORIGINAL        ImageSize = iota
)

type InImageTag struct {
    PersistableModel
    Top         float64 `json:"top,omitempty"`
    Left        float64 `json:"left,omitempty"`
    Height      float64 `json:"height,omitempty"`
    Width       float64 `json:"width,omitempty"`
    Description string  `json:"description,omitempty"`
    Contact     string  `json:"contact,omitempty"`
}

type ImageStorageMetadata struct {
    AclPersistableModel
    Service  *ImageService `json:"service,omitempty"`
    Uri      string        `json:"uri,omitempty"`
    Width    uint32        `json:"width,omitempty"`
    Height   uint32        `json:"height,omitempty"`
    Size     *ImageSize    `json:"size,omitempty"`
    FileSize uint64        `json:"file_size,omitempty"`
    Location *Location     `json:"location,omitempty"`
}

type ImageMetadata struct {
    AclPersistableModel
    UserId          string                  `json:"user_id,omitempty"`
    Keywords        []string                `json:"keywords,omitempty"`
    People          []*InImageTag           `json:"people,omitempty"`
    Caption         string                  `json:"caption,omitempty"`
    ContentType     string                  `json:"content_type,omitempty"`
    DateTaken       *DateTime               `json:"date_taken,omitempty"`
    StorageMetadata []*ImageStorageMetadata `json:"storage_metadata,omitempty"`
}
