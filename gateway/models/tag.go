package models

type Tag struct {
	Name  string `gorm:"primaryKey;type:varchar(255);not null;unique"`
	Type  string `gorm:"type:varchar(100);not null"`
	Files []File `gorm:"many2many:datafile_tags;foreignKey:Name;joinForeignKey:tag_name;inverseJoinForeignKey:data_file_c_id"`
	// Files []File `gorm:"many2many:datafile_tags;foreignKey:Name;joinForeignKey:tag_name;inverseJoinForeignKey:data_file_id"
}
