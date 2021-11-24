package resources

// Load loads gg's assets.
func Load() error {
	if err := loadImages(); err != nil {
		return err
	}

	if err := loadMaps(); err != nil {
		return err
	}

	return nil
}
