package database

// Announcement is a representation of a Discord announcement channel.
//
// It can be obtained by GetAnnouncement for a given channel on a platform.
type Announcement struct {
	ChannelID string
	Role      string
}

// GetAnnouncement reads all Discord announcement channels from the database for a given channel ID
// on a platform.
// A platform could be "twitch" or "youtube".
//
// If no result matches the given platform and channel ID the returned error will be sql.ErrNoRows.
// Other errors may exist.
func GetAnnouncement(platform, id string) ([]Announcement, error) {
	rows, err := Query("SELECT IFNULL(channel, ''),IFNULL(role, '') FROM announcements WHERE type=? AND id=?", platform, id)
	if err != nil {
		return []Announcement{}, err
	}
	defer rows.Close()
	var announcements []Announcement
	for rows.Next() {
		var channelID, roleID string
		if err := rows.Scan(&channelID, &roleID); err != nil {
			return []Announcement{}, err
		}
		announcements = append(announcements, Announcement{channelID, roleID})
	}
	return announcements, err
}
