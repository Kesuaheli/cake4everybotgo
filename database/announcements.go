package database

type Announcement struct {
	ChannelID string
	Role      string
}

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
