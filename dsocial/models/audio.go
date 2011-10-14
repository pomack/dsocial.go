package dsocial

type AudioMetadata struct {
    AclPersistableModel
    Name          string   `json:"name,omitempty"`
    Artists       []string `json:"artists,omitempty"`
    AlbumArtists  []string `json:"album_artists,omitempty"`
    Album         string   `json:"album,omitempty"`
    Groupings     []string `json:"groupings,omitempty"`
    Composers     []string `json:"composers,omitempty"`
    YearRecorded  int16    `json:"year_recorded,omitempty"`
    CopyrightYear int16    `json:"copyright_year,omitempty"`
    TrackNumber   int16    `json:"track_number,omitempty"`
    TotalTracks   int16    `json:"total_tracks,omitempty"`
    DiscNumber    int16    `json:"disc_number,omitempty"`
    TotalDiscs    int16    `json:"total_discs,omitempty"`
    Genre         []string `json:"genre,omitempty"`
    Rating        int8     `json:"rating,omitempty"`
    AlbumRating   int8     `json:"album_rating,omitempty"`
    TimesPlayed   int32    `json:"times_played,omitempty"`
    LastPlayed    int64    `json:"last_played,omitempty"`
    Tags          []string `json:"tags,omitempty"`
    LengthMsec    int32    `json:"length_msec,omitempty"`
    Comments      string   `json:"comments,omitempty"`
    Lyrics        string   `json:"lyrics,omitempty"`
    License       string   `json:"license,omitempty"`
    LicenseUrl    string   `json:"license_url,omitempty"`
    FromUrl       string   `json:"from_url,omitempty"`
    ArtworkUrl    string   `json:"artwork_url,omitempty"`
}
