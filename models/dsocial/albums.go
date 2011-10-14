package dsocial

type AlbumMetadata struct {
    AclPersistableModel `json:"model,omitempty,collapse"`
    UserId              string    `json:"user_id,omitempty"`
    Name                string    `json:"name,omitempty"`
    UrlName             string    `json:"url_name,omitempty"`
    Location            *Location `json:"location,omitempty"`
    LocationDescription string    `json:"location_description,omitempty"`
    Description         string    `json:"description,omitempty"`
    IconImageId         string    `json:"icon_image_id,omitempty"`
    Theme               string    `json:"theme,omitempty"`
}

type Album struct {
    AlbumMetadata `json:"metadata,omitempty,collapse"`
    ImageIds      []string `json:"image_ids,omitempty"`
}
